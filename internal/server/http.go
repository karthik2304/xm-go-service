package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/getkin/kin-openapi/openapi3filter"
	Openapi "github.com/karthik2304/xm-go-service/api/v1/go"
	"github.com/karthik2304/xm-go-service/internal/auth"
	"github.com/karthik2304/xm-go-service/internal/repository"
	"github.com/karthik2304/xm-go-service/pkg"
	"github.com/labstack/echo/v4"
	middleware "github.com/oapi-codegen/echo-middleware"
	"github.com/segmentio/kafka-go"
)

type Server struct {
	Server *echo.Echo
	db     repository.Repository
	writer pkg.KafkaWriter
	reader pkg.KafkaReader
}

func New(e *echo.Echo, db repository.Repository, writer pkg.KafkaWriter, reader pkg.KafkaReader) *Server {
	s := Server{
		Server: e,
		db:     db,
		writer: writer,
		reader: reader,
	}

	go ConsumeMessage(db, reader)

	swagger, err := Openapi.GetSwagger()
	if err != nil {
		fmt.Println(err)
	}
	swagger.Servers = nil

	validatorOptions := &middleware.Options{
		Options: openapi3filter.Options{
			AuthenticationFunc: auth.OpenAPIAuthFunc,
		},
	}
	validatorOptions.ErrorHandler = errorHandler

	e.Use(middleware.OapiRequestValidatorWithOptions(swagger, validatorOptions))

	Openapi.RegisterHandlers(e, &s)
	return &s
}
func errorHandler(c echo.Context, err *echo.HTTPError) error {
	errText := err.Error()
	var msg Openapi.BadRequest
	msg.Message = errText
	return c.JSON(err.Code, &msg)
}

func ConsumeMessage(db repository.Repository, reader pkg.KafkaReader) {
	const maxRetries = 5

	for {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		msg, err := reader.ReadMessage(ctx)
		cancel()

		if err != nil {
			if err == kafka.LeaderNotAvailable {
				log.Printf("Leader Not Available. Retrying...")
				retryWithBackoff()
				continue
			}

			log.Printf("Error reading message: %v", err)
			time.Sleep(2 * time.Second)
			continue
		}

		var receivedMsg Openapi.Event
		err = json.Unmarshal(msg.Value, &receivedMsg)
		if err != nil {
			log.Printf("Error unmarshaling JSON: %v", err)
			continue
		}

		for attempt := 1; attempt <= maxRetries; attempt++ {
			err = db.Send(context.Background(), receivedMsg)
			if err == nil {
				break
			}

			log.Printf("Error sending to DB (attempt %d/%d): %v", attempt, maxRetries, err)
			time.Sleep(time.Duration(attempt) * time.Second)
		}

		if err != nil {
			log.Printf("Dropping message after %d retries", maxRetries)
		}
	}
}

func retryWithBackoff() {
	delay := time.Duration(2+rand.Intn(3)) * time.Second
	time.Sleep(delay)
}

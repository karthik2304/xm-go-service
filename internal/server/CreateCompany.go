package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	Openapi "github.com/karthik2304/xm-go-service/api/v1/go"
	"github.com/karthik2304/xm-go-service/internal/utils"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/mongo"
)

func (s *Server) CreateCompany(ctx echo.Context) error {

	// get user details
	userName, ok := ctx.Get("username").(string)
	if !ok {
		return ctx.JSON(http.StatusUnauthorized, Openapi.Error{Message: "No UserName in context"})
	}

	// get payload
	body := Openapi.CompanyRequestBody{}
	err := ctx.Bind(&body)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, Openapi.Error{Message: fmt.Sprintf("%v", err)})
	}

	// validate existing user
	data, err := s.db.GetOne(context.Background(), map[string]interface{}{
		"companyname": body.CompanyName,
	})
	if err != nil && err != mongo.ErrNoDocuments {
		return ctx.JSON(http.StatusInternalServerError, Openapi.Error{Message: fmt.Sprintf("Database error: %v", err)})
	}
	if data != nil && data.CompanyName == body.CompanyName {
		return ctx.JSON(http.StatusBadRequest, Openapi.Error{Message: "Company Name should be unique"})
	}

	// db create record
	err = s.db.Create(context.Background(), Openapi.CreateCompanyJSONBody(body))
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, Openapi.Error{Message: fmt.Sprintf("%v", err)})
	}

	// Event Payload
	err = s.WriteMessage(Openapi.Event{
		EventDetails: &body,
		EventType:    "CREATE-COMPANY-EVENT",
		Id:           uuid.New().String(),
		Timestamp:    utils.GetCurrentTime(),
		UserName:     userName,
	})
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, Openapi.Error{Message: fmt.Sprintf("%v", err)})
	}

	return ctx.JSON(http.StatusCreated, Openapi.Success{Message: "Company Created Successfully"})

}

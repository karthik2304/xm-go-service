package server

import (
	"encoding/json"
	"net/http"
	"testing"

	Openapi "github.com/karthik2304/xm-go-service/api/v1/go"
	"github.com/karthik2304/xm-go-service/internal/auth"
	"github.com/karthik2304/xm-go-service/internal/mocks"
	"github.com/labstack/echo/v4"
	"github.com/segmentio/kafka-go"
	"github.com/steinfletcher/apitest"
	"github.com/stretchr/testify/mock"
)

func TestLogin_Success(t *testing.T) {
	e := echo.New()

	mockRepo := new(mocks.MockRepository)
	mockKafkaWriter := new(mocks.MockKafkaWriter)
	mockKafkaReader := new(mocks.MockKafkaReader)

	pass, _ := auth.HashPassword("securepassword")

	mockRepo.On("GetAuth", mock.Anything, "testuser").Return(&Openapi.SignUpJSONRequestBody{
		Username: "testuser",
		Password: pass,
	}, nil)
	mockRepo.On("Send", mock.Anything, mock.Anything).Return(nil)
	mockKafkaWriter.On("WriteMessages", mock.Anything, mock.Anything).Return(nil)
	mockKafkaReader.On("ReadMessage", mock.Anything).Return(kafka.Message{
		Value: []byte(`{"event": "LoginEvent"}`),
	}, nil)

	srv := New(e, mockRepo, mockKafkaWriter, mockKafkaReader)

	loginRequest := Openapi.LoginJSONRequestBody{
		Username: "testuser",
		Password: "securepassword",
	}
	requestBody, _ := json.Marshal(loginRequest)

	apitest.New().
		Handler(srv.Server).
		Post("/v1/auth/login").
		JSON(string(requestBody)).
		Expect(t).
		Status(http.StatusOK).
		End()

	mockRepo.AssertExpectations(t)
	mockKafkaWriter.AssertExpectations(t)
}

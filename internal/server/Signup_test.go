package server

import (
	"encoding/json"
	"net/http"
	"testing"

	Openapi "github.com/karthik2304/xm-go-service/api/v1/go"
	"github.com/karthik2304/xm-go-service/internal/mocks"
	"github.com/labstack/echo/v4"
	"github.com/segmentio/kafka-go"
	"github.com/steinfletcher/apitest"
	"github.com/stretchr/testify/mock"
)

func TestSignUp(t *testing.T) {
	e := echo.New()

	mockRepo := new(mocks.MockRepository)
	mockKafkaWriter := new(mocks.MockKafkaWriter)
	mockKafkaReader := new(mocks.MockKafkaReader)

	mockRepo.On("GetAuth", mock.Anything, "newuser").Return(nil, nil)
	mockRepo.On("CreateAuth", mock.Anything, mock.Anything).Return(nil)
	mockKafkaWriter.On("WriteMessages", mock.Anything, mock.Anything).Return(nil)
	mockRepo.On("Send", mock.Anything, mock.Anything).Return(nil)

	mockKafkaReader.On("ReadMessage", mock.Anything).Return(kafka.Message{
		Value: []byte(`{"event": "Signup-Event"}`),
	}, nil)
	srv := New(e, mockRepo, mockKafkaWriter, mockKafkaReader)

	signUpRequest := Openapi.LoginJSONRequestBody{
		Username: "newuser",
		Password: "securepassword",
	}
	requestBody, _ := json.Marshal(signUpRequest)

	apitest.New().
		Handler(srv.Server).
		Post("/v1/auth/signup").
		JSON(string(requestBody)).
		Expect(t).
		Status(http.StatusOK).
		Body(`{"message":"User Created Successfully"}`).
		End()

	mockRepo.AssertExpectations(t)
	mockKafkaWriter.AssertExpectations(t)
}

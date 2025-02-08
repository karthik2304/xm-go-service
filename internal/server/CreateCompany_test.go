package server_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	Openapi "github.com/karthik2304/xm-go-service/api/v1/go"
	constants "github.com/karthik2304/xm-go-service/internal/constant"
	"github.com/karthik2304/xm-go-service/internal/mocks"
	"github.com/karthik2304/xm-go-service/internal/server"
	"github.com/labstack/echo/v4"
	"github.com/segmentio/kafka-go"
	"github.com/steinfletcher/apitest"
	"github.com/stretchr/testify/mock"
)

func TestCreateCompany_Success(t *testing.T) {
	e := echo.New()

	// mocking db & kafka
	mockRepo := new(mocks.MockRepository)
	mockKafkaWriter := new(mocks.MockKafkaWriter)
	mockKafkaReader := new(mocks.MockKafkaReader)

	mockRepo.On("GetOne", mock.Anything, mock.Anything).Return(nil, nil)
	mockRepo.On("Create", mock.Anything, mock.Anything).Return(nil)
	mockRepo.On("Send", mock.Anything, mock.Anything).Return(nil)
	mockKafkaWriter.On("WriteMessages", mock.Anything, mock.Anything).Return(nil)
	mockKafkaReader.On("ReadMessage", mock.Anything).Return(kafka.Message{
		Value: []byte(`{"event": "CompanyCreated"}`),
	}, nil)

	server := server.New(e, mockRepo, mockKafkaWriter, mockKafkaReader)

	companyRequest := Openapi.CreateCompanyJSONBody{
		CompanyName:    "tcs",
		CompanyUUID:    "tcs",
		Registered:     false,
		TotalEmployees: 10,
		Type:           Openapi.Cooperative,
	}

	requestBody, _ := json.Marshal(companyRequest)

	apitest.New().
		Handler(server.Server).
		Post("/v1/create-company").
		Header("Authorization", fmt.Sprintf("Bearer %s", constants.Token)).
		JSON(string(requestBody)).
		Expect(t).
		Status(http.StatusCreated).
		Body(`{"message":"Company Created Successfully"}`).
		End()

	mockRepo.AssertExpectations(t)
	mockKafkaWriter.AssertExpectations(t)
	mockKafkaReader.AssertExpectations(t)
}

func TestCreateCompany_WithSameUser(t *testing.T) {
	e := echo.New()

	// mocking db & kafka
	mockRepo := new(mocks.MockRepository)
	mockKafkaWriter := new(mocks.MockKafkaWriter)
	mockKafkaReader := new(mocks.MockKafkaReader)

	companyRequest := Openapi.CreateCompanyJSONBody{
		CompanyName:    "tcs",
		CompanyUUID:    "tcs",
		Registered:     false,
		TotalEmployees: 10,
		Type:           Openapi.Cooperative,
	}

	mockRepo.On("GetOne", mock.Anything, mock.Anything).Return(&companyRequest, nil)
	mockRepo.On("Send", mock.Anything, mock.Anything).Return(nil)
	mockKafkaReader.On("ReadMessage", mock.Anything).Return(kafka.Message{
		Value: []byte(`{"event": "CompanyCreateFailWithSameCompanyDetails"}`),
	}, nil)

	server := server.New(e, mockRepo, mockKafkaWriter, mockKafkaReader)

	requestBody, _ := json.Marshal(companyRequest)

	apitest.New().
		Handler(server.Server).
		Post("/v1/create-company").
		Header("Authorization", fmt.Sprintf("Bearer %s", constants.Token)).
		JSON(string(requestBody)).
		Expect(t).
		Status(http.StatusBadRequest).
		Body(`{"message":"Company Name should be unique"}`).
		End()

}

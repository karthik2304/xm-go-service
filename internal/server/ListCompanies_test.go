package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	Openapi "github.com/karthik2304/xm-go-service/api/v1/go"
	constants "github.com/karthik2304/xm-go-service/internal/constant"
	"github.com/karthik2304/xm-go-service/internal/mocks"
	"github.com/labstack/echo/v4"
	"github.com/segmentio/kafka-go"
	"github.com/steinfletcher/apitest"
	"github.com/stretchr/testify/mock"
)

func TestListCompanies(t *testing.T) {
	e := echo.New()

	mockRepo := new(mocks.MockRepository)
	mockKafkaWriter := new(mocks.MockKafkaWriter)
	mockKafkaReader := new(mocks.MockKafkaReader)

	companyData := []Openapi.SuccessResponse{
		{
			CompanyName:    "tcs",
			CompanyUUID:    "tcs-uuid",
			Registered:     true,
			TotalEmployees: 100,
			Type:           Openapi.NonProfit,
		},
		{
			CompanyName:    "XM",
			CompanyUUID:    "XM-IT",
			Registered:     true,
			TotalEmployees: 500,
			Type:           Openapi.SoleProprietorship,
		},
	}

	mockRepo.On("GetAll", mock.Anything, mock.Anything).Return(&companyData, nil)
	mockRepo.On("Send", mock.Anything, mock.Anything).Return(nil)
	mockKafkaWriter.On("WriteMessages", mock.Anything, mock.Anything).Return(nil)
	mockKafkaReader.On("ReadMessage", mock.Anything).Return(kafka.Message{
		Value: []byte(`{"event": "CompanyGetAll"}`),
	}, nil)

	server := New(e, mockRepo, mockKafkaWriter, mockKafkaReader)

	expectedResponse, _ := json.Marshal(companyData)

	apitest.New().
		Handler(server.Server).
		Get("/v1/list-companies").
		Header("Authorization", fmt.Sprintf("Bearer %s", constants.Token)).
		Expect(t).
		Status(http.StatusOK).
		Body(string(expectedResponse)).
		End()

	mockRepo.AssertExpectations(t)
}

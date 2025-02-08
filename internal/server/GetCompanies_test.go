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

func TestGetCompany(t *testing.T) {
	e := echo.New()

	mockRepo := new(mocks.MockRepository)
	mockKafkaWriter := new(mocks.MockKafkaWriter)
	mockKafkaReader := new(mocks.MockKafkaReader)

	companyData := Openapi.CreateCompanyJSONBody{
		CompanyName:    "tcs",
		CompanyUUID:    "tcs-uuid",
		Registered:     true,
		TotalEmployees: 100,
		Type:           Openapi.NonProfit,
	}

	mockRepo.On("GetOne", mock.Anything, mock.Anything).Return(&companyData, nil)
	mockRepo.On("Send", mock.Anything, mock.Anything).Return(nil)
	mockKafkaWriter.On("WriteMessages", mock.Anything, mock.Anything).Return(nil)
	mockKafkaReader.On("ReadMessage", mock.Anything).Return(kafka.Message{
		Value: []byte(`{"event": "CompanyGetOne"}`),
	}, nil)

	server := New(e, mockRepo, mockKafkaWriter, mockKafkaReader)

	expectedResponse, _ := json.Marshal(companyData)

	apitest.New().
		Handler(server.Server).
		Get("/v1/company-details/tcs-uuid").
		Header("Authorization", fmt.Sprintf("Bearer %s", constants.Token)).
		Expect(t).
		Status(http.StatusOK).
		Body(string(expectedResponse)).
		End()

	mockRepo.AssertExpectations(t)
}

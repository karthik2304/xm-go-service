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

func TestUpdateCompany(t *testing.T) {
	e := echo.New()

	mockRepo := new(mocks.MockRepository)
	mockKafkaWriter := new(mocks.MockKafkaWriter)
	mockKafkaReader := new(mocks.MockKafkaReader)

	mockRepo.On("Update", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	mockRepo.On("Send", mock.Anything, mock.Anything).Return(nil) // ✅ Mocking Send

	mockKafkaWriter.On("WriteMessages", mock.Anything, mock.Anything).Return(nil)
	mockKafkaReader.On("ReadMessage", mock.Anything).Return(kafka.Message{
		Value: []byte(`{"event": "CompanyUpdated"}`),
	}, nil)
	server := New(e, mockRepo, mockKafkaWriter, mockKafkaReader)

	updateRequest := Openapi.UpdateCompanyDetailsJSONRequestBody{
		CompanyName:    "TCS UUID",
		Registered:     true,
		TotalEmployees: 300,
		Type:           Openapi.NonProfit,
	}

	requestBody, _ := json.Marshal(updateRequest)

	apitest.New().
		Handler(server.Server).
		Patch("/v1/company-details/tcs-uuid").
		Header("Authorization", fmt.Sprintf("Bearer %s", constants.Token)).
		JSON(string(requestBody)).
		Expect(t).
		Status(http.StatusOK).
		Body(`{"message":"Company updated successfully"}`).
		End()

	mockRepo.AssertExpectations(t)
}

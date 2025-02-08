package server

import (
	"fmt"
	"net/http"
	"testing"

	constants "github.com/karthik2304/xm-go-service/internal/constant"
	"github.com/karthik2304/xm-go-service/internal/mocks"
	"github.com/labstack/echo/v4"
	"github.com/segmentio/kafka-go"
	"github.com/steinfletcher/apitest"
	"github.com/stretchr/testify/mock"
)

func TestDeleteCompany(t *testing.T) {
	e := echo.New()

	mockRepo := new(mocks.MockRepository)
	mockKafkaWriter := new(mocks.MockKafkaWriter)
	mockKafkaReader := new(mocks.MockKafkaReader)

	mockRepo.On("Delete", mock.Anything, mock.Anything).Return(nil)
	mockRepo.On("Send", mock.Anything, mock.Anything).Return(nil)
	mockKafkaWriter.On("WriteMessages", mock.Anything, mock.Anything).Return(nil)
	mockKafkaReader.On("ReadMessage", mock.Anything).Return(kafka.Message{
		Value: []byte(`{"event": "CompanyDeleteEvent"}`),
	}, nil)

	server := New(e, mockRepo, mockKafkaWriter, mockKafkaReader)

	apitest.New().
		Handler(server.Server).
		Delete("/v1/company-details/tcs-uuid").
		Header("Authorization", fmt.Sprintf("Bearer %s", constants.Token)).
		Expect(t).
		Status(http.StatusOK).
		Body(`{"message":"Company deleted successfully"}`).
		End()

	mockRepo.AssertExpectations(t)
}

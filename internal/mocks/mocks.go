package mocks

import (
	"context"

	Openapi "github.com/karthik2304/xm-go-service/api/v1/go"
	"github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/mock"
)

type MockRepository struct {
	mock.Mock
}

type MockKafkaReader struct {
	mock.Mock
}

type MockKafkaWriter struct {
	mock.Mock
}

func (m *MockRepository) CreateAuth(ctx context.Context, data Openapi.SignUpJSONRequestBody) error {
	args := m.Called(ctx, data)
	return args.Error(0)
}

func (m *MockRepository) Delete(ctx context.Context, filter map[string]interface{}) error {
	args := m.Called(ctx, filter)
	return args.Error(0)
}

func (m *MockRepository) GetAll(ctx context.Context) (*[]Openapi.SuccessResponse, error) {
	args := m.Called(ctx)
	result := args.Get(0)
	if result == nil {
		return nil, args.Error(1)
	}
	return result.(*[]Openapi.SuccessResponse), args.Error(1)
}

func (m *MockRepository) GetAllEvents(ctx context.Context) (*[]Openapi.Event, error) {
	args := m.Called(ctx)
	result := args.Get(0)
	if result == nil {
		return nil, args.Error(1)
	}
	return result.(*[]Openapi.Event), args.Error(1)
}

func (m *MockRepository) GetAuth(ctx context.Context, email string) (*Openapi.SignUpJSONRequestBody, error) {
	args := m.Called(ctx, email)
	result := args.Get(0)
	if result == nil {
		return nil, args.Error(1)
	}
	return result.(*Openapi.SignUpJSONRequestBody), args.Error(1)
}

func (m *MockRepository) Send(ctx context.Context, data Openapi.Event) error {
	args := m.Called(ctx, data)
	return args.Error(0)
}

func (m *MockRepository) Update(ctx context.Context, filter map[string]interface{}, details Openapi.UpdateCompanyDetailsJSONRequestBody) error {
	args := m.Called(ctx, filter, details)
	return args.Error(0)
}

func (m *MockRepository) Create(ctx context.Context, data Openapi.CreateCompanyJSONBody) error {
	args := m.Called(ctx, data)
	return args.Error(0)
}

func (m *MockRepository) GetOne(ctx context.Context, filter map[string]interface{}) (*Openapi.CreateCompanyJSONRequestBody, error) {
	args := m.Called(ctx, filter)
	result := args.Get(0)
	if result == nil {
		return nil, args.Error(1)
	}
	return result.(*Openapi.CreateCompanyJSONRequestBody), args.Error(1)
}

func (m *MockKafkaReader) ReadMessage(ctx context.Context) (kafka.Message, error) {
	args := m.Called(ctx)
	return args.Get(0).(kafka.Message), args.Error(1)
}

func (m *MockKafkaReader) Close() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockKafkaWriter) WriteMessages(ctx context.Context, msgs ...kafka.Message) error {
	args := m.Called(ctx, msgs)
	return args.Error(0)
}

func (m *MockKafkaWriter) Close() error {
	args := m.Called()
	return args.Error(0)
}

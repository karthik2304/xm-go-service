package utils

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	Openapi "github.com/karthik2304/xm-go-service/api/v1/go"
)

func ConvertPayload(data Openapi.Event) ([]byte, error) {
	msgBytes, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	return msgBytes, nil

}

func GetCurrentTime() *string {
	currentTime := time.Now().String()
	return &currentTime
}

func GetUniqueId() string {
	return uuid.New().String()
}

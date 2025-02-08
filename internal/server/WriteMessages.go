package server

import (
	"context"
	"fmt"

	Openapi "github.com/karthik2304/xm-go-service/api/v1/go"
	"github.com/karthik2304/xm-go-service/internal/utils"
	"github.com/segmentio/kafka-go"
)

func (s *Server) WriteMessage(data Openapi.Event) error {
	decodeMsg, _ := utils.ConvertPayload(data)
	err := s.writer.WriteMessages(context.Background(), kafka.Message{
		Key:   []byte(fmt.Sprintf("Key-%s", data.Id)),
		Value: decodeMsg,
	})
	if err != nil {
		return err
	}
	return nil
}

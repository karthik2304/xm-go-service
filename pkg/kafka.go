package pkg

import (
	"context"
	"log"

	"github.com/karthik2304/xm-go-service/configs"
	"github.com/segmentio/kafka-go"
)

type KafkaWriter interface {
	WriteMessages(ctx context.Context, msgs ...kafka.Message) error
	Close() error
}

type KafkaReader interface {
	ReadMessage(ctx context.Context) (kafka.Message, error)
	Close() error
}

func ConnectKafka() (*kafka.Reader, *kafka.Writer) {
	// kafka
	w := kafka.Writer{
		Addr:     kafka.TCP(configs.Settings.APP_KAFKA_ADDR),
		Topic:    configs.Settings.APP_KAFKATOPIC,
		Balancer: &kafka.LeastBytes{},
	}
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{configs.Settings.APP_KAFKA_ADDR},
		Topic:    configs.Settings.APP_KAFKATOPIC,
		GroupID:  configs.Settings.APP_KAFKA_GROUPID,
		MaxBytes: 10e6,
	})
	log.Println("Kafka Connected")
	return r, &w

}

package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/Neokrid/order_service/internal/entity"
	"github.com/segmentio/kafka-go"
)

type KafkaPublisher struct {
	writer *kafka.Writer
}

func NewKafkaPublisher(brokers []string, topic string) *KafkaPublisher {
	fmt.Printf("KAFKA DEBUG: Подключаемся к %v, топик: '%s'\n", brokers, topic)
	return &KafkaPublisher{
		writer: &kafka.Writer{
			Addr:                   kafka.TCP(brokers...),
			Topic:                  topic,
			Balancer:               &kafka.LeastBytes{},
			Async:                  false,
			Logger:                 kafka.LoggerFunc(log.Printf),
			ErrorLogger:            kafka.LoggerFunc(log.Printf),
			AllowAutoTopicCreation: true,
		},
	}
}

func (p *KafkaPublisher) Send(ctx context.Context, order *entity.Order, eventType string) error {
	event := entity.OrderEvent{
		EventType: eventType,
		Payload:   order,
		SentAt:    time.Now(),
	}
	msgBytes, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("marshal event: %w", err)
	}

	err = p.writer.WriteMessages(ctx, kafka.Message{
		Key:   []byte(order.ID.String()),
		Value: msgBytes,
	})

	if err != nil {
		return fmt.Errorf("kafka write: %w", err)
	}

	return nil
}

func (p *KafkaPublisher) Close() error {
	return p.writer.Close()
}

package events

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/segmentio/kafka-go"
)

type EventBus interface {
	Publish(ctx context.Context, topic string, event *Event) error
	Subscribe(ctx context.Context, topic string, handler EventHandler) error
	Close() error
}

type EventHandler func(ctx context.Context, event *Event) error

type EventPublisher interface {
	PublishUserRegistered(ctx context.Context, user interface{}) error
	PublishUserTierUpgraded(ctx context.Context, userID string, oldTier, newTier interface{}) error
	PublishUserQuotaUpdated(ctx context.Context, userID string, quotas interface{}) error
}

type KafkaEventBus struct {
	brokers []string
	writer  *kafka.Writer
	readers map[string]*kafka.Reader
}

func NewKafkaEventBus(brokers []string) *KafkaEventBus {
	writer := &kafka.Writer{
		Addr:                   kafka.TCP(brokers...),
		Balancer:               &kafka.LeastBytes{},
		AllowAutoTopicCreation: true,
	}

	return &KafkaEventBus{
		brokers: brokers,
		writer:  writer,
		readers: make(map[string]*kafka.Reader),
	}
}

func (k *KafkaEventBus) Publish(ctx context.Context, topic string, event *Event) error {
	eventBytes, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	err = k.writer.WriteMessages(ctx, kafka.Message{
		Topic: topic,
		Key:   []byte(event.ID),
		Value: eventBytes,
	})
	if err != nil {
		return fmt.Errorf("failed to publish event: %w", err)
	}

	log.Printf("Event published: %s to topic: %s", event.Type, topic)
	return nil
}

func (k *KafkaEventBus) Subscribe(ctx context.Context, topic string, handler EventHandler) error {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: k.brokers,
		Topic:   topic,
		GroupID: "smm-platform",
	})

	k.readers[topic] = reader

	go k.consumeMessages(ctx, reader, handler)
	return nil
}

func (k *KafkaEventBus) consumeMessages(ctx context.Context, reader *kafka.Reader, handler EventHandler) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			msg, err := reader.ReadMessage(ctx)
			if err != nil {
				log.Printf("Error reading message: %v", err)
				continue
			}

			var event Event
			if err := json.Unmarshal(msg.Value, &event); err != nil {
				log.Printf("Error unmarshaling event: %v", err)
				continue
			}

			if err := handler(ctx, &event); err != nil {
				log.Printf("Error handling event: %v", err)
			}
		}
	}
}

func (k *KafkaEventBus) Close() error {
	if err := k.writer.Close(); err != nil {
		return err
	}

	for _, reader := range k.readers {
		if err := reader.Close(); err != nil {
			return err
		}
	}

	return nil
}

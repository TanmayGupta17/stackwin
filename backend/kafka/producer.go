package kafka

import (
	"context"
	"encoding/json"
	"log"

	"github.com/segmentio/kafka-go"
)

type KafkaProducer struct {
	writer *kafka.Writer
}

type GameEvent struct {
	EventType string      `json:"event_type"` // "move", "game_end", "game_start"
	GameID    string      `json:"game_id"`
	PlayerID  string      `json:"player_id"`
	Timestamp int64       `json:"timestamp"`
	Data      interface{} `json:"data,omitempty"`
}

func NewKafkaProducer(brokers []string, topic string) *KafkaProducer {
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers: brokers,
		Topic:   topic,
	})
	return &KafkaProducer{writer: writer}
}

func (kp *KafkaProducer) SendEvent(event GameEvent) error {
	eventJSON, _ := json.Marshal(event)

	err := kp.writer.WriteMessages(context.Background(), kafka.Message{
		Value: eventJSON,
	})

	if err != nil {
		log.Printf("Kafka write error: %v\n", err)
	}
	return err
}

func (kp *KafkaProducer) SendEventAsync(event GameEvent) {
	// Non-blocking send
	go func() {
		if err := kp.SendEvent(event); err != nil {
			log.Printf("Async Kafka error: %v\n", err)
		}
	}()
}

func (kp *KafkaProducer) Close() error {
	return kp.writer.Close()
}

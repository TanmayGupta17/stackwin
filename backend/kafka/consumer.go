package kafka

import (
	"context"
	"encoding/json"
	"log"

	"github.com/segmentio/kafka-go"
)

type KafkaConsumer struct {
	reader *kafka.Reader
}

func NewKafkaConsumer(brokers []string, topic string, groupID string) *KafkaConsumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: brokers,
		Topic:   topic,
		GroupID: groupID,
	})
	return &KafkaConsumer{reader: reader}
}

// ConsumeMessages starts consuming messages
func (kc *KafkaConsumer) ConsumeMessages(handler func(event GameEvent)) error {
	for {
		msg, err := kc.reader.ReadMessage(context.Background())
		if err != nil {
			log.Printf("Kafka read error: %v\n", err)
			continue
		}

		var event GameEvent
		err = json.Unmarshal(msg.Value, &event)
		if err != nil {
			log.Printf("JSON unmarshal error: %v\n", err)
			continue
		}

		handler(event)
	}
}

func (kc *KafkaConsumer) Close() error {
	return kc.reader.Close()
}


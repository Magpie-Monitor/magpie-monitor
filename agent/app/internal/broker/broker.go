package broker

import (
	"context"
	"log"

	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/plain"
)

type Broker interface {
	Publish(content string)
}

type KafkaWriter struct {
	address   string
	topic     string
	batchSize int
	writer    *kafka.Writer
}

func NewStreamWriter(address, topic, username, password string, batchSize int) Broker {
	writer := kafka.Writer{
		Addr:                   kafka.TCP(address),
		Topic:                  topic,
		AllowAutoTopicCreation: true,
		Transport:              &kafka.Transport{SASL: plain.Mechanism{Username: username, Password: password}},
	}
	return &KafkaWriter{address: address, topic: topic, batchSize: batchSize, writer: &writer}
}

func (s *KafkaWriter) Publish(content string) {
	msg := kafka.Message{Value: []byte(content)}

	err := s.writer.WriteMessages(context.Background(), msg)
	if err != nil {
		log.Printf("Error writing message: %v.", err)
	}
}

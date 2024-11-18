package remote_write

import (
	"context"
	"log"

	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/plain"
)

type RemoteWriter interface {
	Write(content string)
}

type StreamWriter struct {
	address   string
	topic     string
	batchSize int
	buffer    []kafka.Message
	writer    *kafka.Writer
}

func NewStreamWriter(address, topic, username, password string, batchSize int) RemoteWriter {
	writer := kafka.Writer{
		Addr:                   kafka.TCP(address),
		Topic:                  topic,
		AllowAutoTopicCreation: true,
		Transport:              &kafka.Transport{SASL: plain.Mechanism{Username: username, Password: password}},
	}
	return &StreamWriter{address: address, topic: topic, batchSize: batchSize, buffer: make([]kafka.Message, 0), writer: &writer}
}

func (s *StreamWriter) Write(content string) {
	msg := kafka.Message{Value: []byte(content)}

	err := s.writer.WriteMessages(context.Background(), msg)
	if err != nil {
		log.Printf("Error writing message: %v.", err)
	}
}

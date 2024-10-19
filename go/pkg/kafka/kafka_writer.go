package remote_write

import (
	"context"
	"log"
	"os"

	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/plain"
)

type StreamWriter struct {
	topic     string
	batchSize int
	buffer    []kafka.Message
	writer    *kafka.Writer
}

type KafkaCredentials struct {
	address  string
	username string
	password string
}

func NewKafkaCredentials() *KafkaCredentials {
	kafkaAddress, ok := os.LookupEnv("KAFKA_BROKER_URL")
	if !ok {
		panic("KAFKA_ADDRESS env variable not provided")
	}

	kafkaUsername, ok := os.LookupEnv("KAFKA_CLIENT_USERNAME")
	if !ok {
		panic("KAFKA_USERNAME env variable not provided")
	}

	kafkaPassword, ok := os.LookupEnv("KAFKA_CLIENT_PASSWORD")
	if !ok {
		panic("KAFKA_PASSWORD env variable not provided")
	}

	return &KafkaCredentials{address: kafkaAddress, username: kafkaUsername, password: kafkaPassword}
}

func NewStreamWriter(credentials *KafkaCredentials, topic string, batchSize int) *StreamWriter {
	writer := kafka.Writer{
		Addr:                   kafka.TCP(credentials.address),
		Topic:                  topic,
		AllowAutoTopicCreation: true,
		Transport:              &kafka.Transport{SASL: plain.Mechanism{Username: credentials.username, Password: credentials.password}},
	}
	return &StreamWriter{topic: topic, batchSize: batchSize, buffer: make([]kafka.Message, 0), writer: &writer}
}

func (s *StreamWriter) Write(content string) {
	msg := kafka.Message{Value: []byte(content)}
	s.buffer = append(s.buffer, msg)

	if len(s.buffer) > s.batchSize {
		log.Printf("Buffer reached the batch size of: %d, sending messages.", s.batchSize)

		err := s.writer.WriteMessages(context.Background(), s.buffer...)
		if err != nil {
			log.Printf("Error writing message: %v. Buffering, buffer size: %d.", err, len(s.buffer))
		}
	}
}

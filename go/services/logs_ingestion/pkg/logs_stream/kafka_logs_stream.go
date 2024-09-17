package logsstream

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

func NewKafkaLogsStream[T any](host string, port string, topic string, logger *zap.Logger) KafkaLogsStreamReader[T] {

	brokers := []string{fmt.Sprintf("%s:%s", host, port)}
	reader := kafka.NewReader(
		kafka.ReaderConfig{
			Brokers:   brokers,
			Topic:     topic,
			Partition: 0,
			MaxBytes:  10e8,
		},
	)

	return KafkaLogsStreamReader[T]{
		topic:   topic,
		brokers: brokers,
		reader:  reader,
		stream:  make(chan T),
		logger:  logger,
	}
}

func (s *KafkaLogsStreamReader[T]) handleMessage(message []byte) {

	s.logger.Debug("Read message", zap.String("msg", string(message)))

	var a T
	err := json.NewDecoder(bytes.NewReader(message)).Decode(&a)
	if err != nil {
		s.logger.Error("Failed to decode message from Kafka", zap.Error(err))
		return
	}

	if s.handler != nil {
		s.handler(a)
	}

	s.stream <- a
}
func (s *KafkaLogsStreamReader[T]) Listen() {
	for {
		m, err := s.reader.ReadMessage(context.Background())
		if err != nil {
			s.logger.Error("Failed to read message from Kafka", zap.Error(err))
			continue
		}

		go s.handleMessage(m.Value)

	}
}

type KafkaLogsStreamReader[T any] struct {
	topic   string
	brokers []string
	reader  *kafka.Reader
	stream  chan T
	logger  *zap.Logger
	handler func(T)
}

func (s *KafkaLogsStreamReader[T]) Stream() chan T {
	return s.stream
}

func (s *KafkaLogsStreamReader[T]) SetHandler(f func(T)) {

	s.handler = f
}

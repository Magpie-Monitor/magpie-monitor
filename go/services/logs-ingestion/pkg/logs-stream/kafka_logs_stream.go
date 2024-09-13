package logsstream

import (
	// "bytes"
	"context"
	// "encoding/json"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

func NewKafkaLogsStream[T any](brokers []string, topic string, logger *zap.Logger) KafkaLogsStreamReader[T] {

	reader := kafka.NewReader(
		kafka.ReaderConfig{
			Brokers:   brokers,
			Topic:     topic,
			Partition: 0,
			MaxBytes:  10e8,
			// GroupID:   "consumer-group-1",
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

func (s *KafkaLogsStreamReader[T]) Listen() {
	for {
		// var a T
		m, err := s.reader.ReadMessage(context.Background())
		if err != nil {
			s.logger.Error("Failed to read message from Kafka", zap.Error(err))
			// break
		}

		s.logger.Info("Read message", zap.String("msg", string(m.Value)))

		// err = json.NewDecoder(bytes.NewReader(m.Value)).Decode(&a)
		if err != nil {
			s.logger.Error("Failed to decode message from Kafka", zap.Error(err))
			// break
		}

		if s.handler != nil {
			// s.handler(a)
		}

		// s.stream <- a
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

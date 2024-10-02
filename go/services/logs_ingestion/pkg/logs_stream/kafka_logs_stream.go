package logsstream

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/scram"
	"go.uber.org/zap"
)

type KafkaLogsStreamParams struct {
	Host     string
	Port     string
	Topic    string
	Logger   *zap.Logger
	Username string
	Password string
}

func NewKafkaLogsStream[T any](params *KafkaLogsStreamParams) KafkaLogsStreamReader[T] {

	brokers := []string{fmt.Sprintf("%s:%s", params.Host, params.Port)}

	mechanism, err := scram.Mechanism(scram.SHA512, params.Username, params.Password)
	if err != nil {
		panic("Failed to set sasl mechanism for logs ingestion kafka queue")
	}

	dialer := &kafka.Dialer{
		SASLMechanism: mechanism,
	}

	reader := kafka.NewReader(
		kafka.ReaderConfig{
			Brokers:   brokers,
			Topic:     params.Topic,
			Partition: 0,
			MaxBytes:  10e8,
			Dialer:    dialer,
		},
	)

	return KafkaLogsStreamReader[T]{
		topic:   params.Topic,
		brokers: brokers,
		reader:  reader,
		stream:  make(chan T),
		logger:  params.Logger,
	}
}

func (s *KafkaLogsStreamReader[T]) handleMessage(message []byte) {

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

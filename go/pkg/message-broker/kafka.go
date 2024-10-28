package messagebroker

import (
	"context"
	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/plain"
	"github.com/segmentio/kafka-go/sasl/scram"
	"go.uber.org/zap"
)

type KafkaMessageBroker struct {
	logger *zap.Logger
	writer *kafka.Writer
	reader *kafka.Reader
}

func NewKafkaMessageBroker(addr string, topic string, username string, password string, logger *zap.Logger) *KafkaMessageBroker {

	writer := &kafka.Writer{
		Addr:                   kafka.TCP(addr),
		Topic:                  topic,
		AllowAutoTopicCreation: true,
		Transport:              &kafka.Transport{SASL: plain.Mechanism{Username: username, Password: password}},
	}

	mechanism, err := scram.Mechanism(scram.SHA512, username, password)
	if err != nil {
		panic("Failed to set sasl mechanism for logs ingestion kafka queue")
	}

	dialer := &kafka.Dialer{
		SASLMechanism: mechanism,
	}

	reader := kafka.NewReader(
		kafka.ReaderConfig{
			Brokers:   []string{addr},
			Topic:     topic,
			Partition: 0,
			MaxBytes:  10e8,
			Dialer:    dialer,
		},
	)

	return &KafkaMessageBroker{
		writer: writer,
		reader: reader,
		logger: logger,
	}
}

func (b *KafkaMessageBroker) Publish(ctx context.Context, key []byte, value []byte) error {
	err := b.writer.WriteMessages(context.Background(), kafka.Message{
		Key:   key,
		Value: value,
	})

	if err != nil {
		b.logger.Error("Failed to publish a message", zap.Error(err),
			zap.String("key", string(key)),
			zap.String("value", string(value)))
		return err
	}

	return nil
}

func (b *KafkaMessageBroker) Subscribe(ctx context.Context, channel chan<- []byte, errChannel chan<- error) {
	for {
		msg, err := b.reader.ReadMessage(ctx)

		if err != nil {
			b.logger.Error("Failed to read message", zap.Error(err))
			errChannel <- err
		}

		channel <- msg.Value
	}
}

var _ MessageBroker[any] = &KafkaJsonMessageBroker[any]{}

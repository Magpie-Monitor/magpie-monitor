package messagebroker

import (
	"context"
	"encoding/json"

	"go.uber.org/zap"
)

type KafkaJsonMessageBroker[T any] struct {
	broker *KafkaMessageBroker
	logger *zap.Logger
}

func NewKafkaJsonMessageBroker[T any](logger *zap.Logger, addr string, topic string, username string, password string) *KafkaJsonMessageBroker[T] {
	return &KafkaJsonMessageBroker[T]{
		logger: logger,
		broker: NewKafkaMessageBroker(addr, topic, username, password, logger),
	}
}

func (b *KafkaJsonMessageBroker[T]) Publish(key string, message T) error {

	encodedMessage, err := json.Marshal(message)
	if err != nil {
		b.logger.Error("Failed to encode broker json message", zap.Error(err), zap.Any("message", message))
		return err
	}

	return b.broker.Publish(context.Background(), []byte(key), encodedMessage)
}

func (b *KafkaJsonMessageBroker[T]) Subscribe(ctx context.Context, messages chan<- T, errors chan<- error) {

	msgChannel := make(chan []byte)

	go b.broker.Subscribe(ctx, msgChannel, errors)

	for {
		select {
		case message := <-msgChannel:
			var decodedMessage T

			err := json.Unmarshal(message, &decodedMessage)
			if err != nil {
				b.logger.Error("Failed to decode broker json message", zap.Error(err),
					zap.Any("messsage", message))
				errors <- err
			}

			messages <- decodedMessage

		case <-ctx.Done():
			b.logger.Info("KafkaJsonMessageBroker conext was cancelled")
			break
		}

	}

}
func (b *KafkaJsonMessageBroker[T]) CloseReader() error {
	return b.broker.CloseReader()

}

var _ MessageBroker[any] = &KafkaJsonMessageBroker[any]{}

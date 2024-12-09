package tests

import (
	"context"
	"time"

	messagebroker "github.com/Magpie-Monitor/magpie-monitor/pkg/message-broker"
	"github.com/Magpie-Monitor/magpie-monitor/services/cluster_metadata/pkg/repositories"
	"github.com/Magpie-Monitor/magpie-monitor/services/cluster_metadata/pkg/services"
	"go.uber.org/zap"
)

func NewMockApplicationMetadataBroker(logger *zap.Logger) messagebroker.MessageBroker[repositories.ApplicationState] {
	return NewKafkaJsonMessageBroker[repositories.ApplicationState](logger)
}

func NewMockNodeMetadataBroker(logger *zap.Logger) messagebroker.MessageBroker[repositories.NodeState] {
	return NewKafkaJsonMessageBroker[repositories.NodeState](logger)
}

func NewMockNodeMetadataUpdatedBroker(logger *zap.Logger) messagebroker.MessageBroker[services.NodeMetadataUpdated] {
	return NewKafkaJsonMessageBroker[services.NodeMetadataUpdated](logger)
}

func NewMockApplicationMetadataUpdatedBroker(logger *zap.Logger) messagebroker.MessageBroker[services.ApplicationMetadataUpdated] {
	return NewKafkaJsonMessageBroker[services.ApplicationMetadataUpdated](logger)
}

func NewMockClusterMetadataUpdatedBroker(logger *zap.Logger) messagebroker.MessageBroker[services.ClusterMetadataUpdated] {
	return NewKafkaJsonMessageBroker[services.ClusterMetadataUpdated](logger)
}

type MockKafkaJsonMessageBroker[T any] struct {
	messages []T
	logger   *zap.Logger
}

func NewKafkaJsonMessageBroker[T any](logger *zap.Logger) messagebroker.MessageBroker[T] {
	return &MockKafkaJsonMessageBroker[T]{
		logger:   logger,
		messages: make([]T, 0),
	}
}

func (b *MockKafkaJsonMessageBroker[T]) Publish(key string, message T) error {
	b.messages = append(b.messages, message)
	b.logger.Info("Received message", zap.Any("messages", b.messages))
	return nil
}

func (b *MockKafkaJsonMessageBroker[T]) Subscribe(ctx context.Context, messages chan<- T, errors chan<- error) {
	for {
		for _, msg := range b.messages {
			messages <- msg
			b.logger.Info("Published message", zap.Any("message", msg))
		}
		b.messages = make([]T, 0)

		time.Sleep(5 * time.Second)
	}
}

func (b *MockKafkaJsonMessageBroker[T]) CloseReader() error {
	return nil
}

var _ messagebroker.MessageBroker[any] = &MockKafkaJsonMessageBroker[any]{}

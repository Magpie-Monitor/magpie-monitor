package messagebroker

import "context"

type MessageBroker[T any] interface {
	Publish(key string, message T) error
	Subscribe(ctx context.Context, messages chan<- T, errors chan<- error)
}

package messagebroker

type MessageBroker[T any] interface {
	Publish(key string, message T) error
	Subscribe(ch chan<- T, errCh chan<- error)
}

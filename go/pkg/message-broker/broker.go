package messagebroker

type MessageBroker[T any] interface {
	Publish(key string, message T) error
	Subscribe(messages chan<- T, errors chan<- error)
}

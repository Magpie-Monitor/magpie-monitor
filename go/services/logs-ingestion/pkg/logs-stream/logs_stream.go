package logsstream

type LogsStreamReader[T any] interface {
	Stream() chan T
	SetHandler(func(T))
}

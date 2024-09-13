package logsstream

import (
	"github.com/Magpie-Monitor/magpie-monitor/pkg/repositories"
	"go.uber.org/zap"
)

type ApplicationLogsStreamReader struct {
	kafkaReader *KafkaLogsStreamReader[repositories.ApplicationLogs]
}

func NewApplicationLogsStreamReader(logger *zap.Logger) *ApplicationLogsStreamReader {

	kafkaReader := NewKafkaLogsStream[repositories.ApplicationLogs](
		[]string{"kafka:9092"},
		"applications",
		logger,
	)

	return &ApplicationLogsStreamReader{
		kafkaReader: &kafkaReader,
	}
}

func (a *ApplicationLogsStreamReader) Listen() {
	a.kafkaReader.Listen()
}

func (a *ApplicationLogsStreamReader) Stream() chan repositories.ApplicationLogs {
	return a.kafkaReader.Stream()
}

func (s *ApplicationLogsStreamReader) SetHandler(f func(repositories.ApplicationLogs)) {
	s.kafkaReader.SetHandler(f)
}

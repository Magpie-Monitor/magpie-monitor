package logsstream

import (
	"github.com/Magpie-Monitor/magpie-monitor/pkg/repositories"
	"go.uber.org/zap"
)

type NodeLogsStreamReader struct {
	kafkaReader *KafkaLogsStreamReader[repositories.NodeLogs]
}

func NewNodeLogsStreamReader(logger *zap.Logger) *NodeLogsStreamReader {

	kafkaReader := NewKafkaLogsStream[repositories.NodeLogs](
		[]string{"kafka:9094"},
		"nodes",
		logger,
	)

	return &NodeLogsStreamReader{
		kafkaReader: &kafkaReader,
	}
}

func (a *NodeLogsStreamReader) Listen() {
	a.kafkaReader.logger.Info("Starting to listen")
	a.kafkaReader.Listen()
}

func (a *NodeLogsStreamReader) Stream() chan repositories.NodeLogs {
	return a.kafkaReader.Stream()
}

func (s *NodeLogsStreamReader) SetHandler(f func(repositories.NodeLogs)) {
	s.kafkaReader.SetHandler(f)
}

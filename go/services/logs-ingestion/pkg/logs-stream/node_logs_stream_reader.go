package logsstream

import (
	"context"
	"github.com/Magpie-Monitor/magpie-monitor/pkg/repositories"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type NodeLogsStreamReader interface {
	Handle(ctx context.Context, nodeLogs *repositories.NodeLogs) error
	Listen()
	AddHandler(func(*repositories.NodeLogs))
}

type KafkaNodeLogsStreamReader struct {
	kafkaReader         *KafkaLogsStreamReader[*repositories.NodeLogs]
	nodesLogsRepository repositories.NodeLogsRepository
	logger              *zap.Logger
}

type NodeLogsStreamReaderParams struct {
	fx.In
	NodesLogsRepository repositories.NodeLogsRepository
	Logger              *zap.Logger
}

func NewKafkaNodeLogsStreamReader(params NodeLogsStreamReaderParams) *KafkaNodeLogsStreamReader {

	kafkaReader := NewKafkaLogsStream[*repositories.NodeLogs](
		[]string{"kafka:9094"},
		"nodes",
		params.Logger,
	)

	return &KafkaNodeLogsStreamReader{
		kafkaReader:         &kafkaReader,
		logger:              params.Logger,
		nodesLogsRepository: params.NodesLogsRepository,
	}
}

func (r *KafkaNodeLogsStreamReader) Handle(
	ctx context.Context,
	nodeLogs *repositories.NodeLogs) error {

	err := r.nodesLogsRepository.InsertLogs(ctx, nodeLogs)
	if err != nil {
		r.logger.Error("Failed to index node logs", zap.Error(err))
	}

	return err
}

func (r *KafkaNodeLogsStreamReader) Listen() {
	r.kafkaReader.logger.Info("Starting to listen for node logs")
	go r.kafkaReader.Listen()
	nodeStream := r.kafkaReader.Stream()

	for {
		log := <-nodeStream
		r.kafkaReader.logger.Debug("Got message from stream", zap.Any("log", log))
		r.Handle(context.Background(), log)
	}
}

func (s *KafkaNodeLogsStreamReader) AddHandler(f func(*repositories.NodeLogs)) {
	s.kafkaReader.SetHandler(f)
}

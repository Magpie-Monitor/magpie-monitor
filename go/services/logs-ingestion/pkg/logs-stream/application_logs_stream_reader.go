package logsstream

import (
	"context"
	"github.com/Magpie-Monitor/magpie-monitor/pkg/repositories"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type ApplicationLogsStreamReader interface {
	Handle(ctx context.Context, nodeLogs *repositories.ApplicationLogs) error
	Listen()
	AddHandler(func(*repositories.ApplicationLogs))
}

type KafkaApplicationLogsStreamReader struct {
	kafkaReader               *KafkaLogsStreamReader[*repositories.ApplicationLogs]
	applicationLogsRepository repositories.ApplicationLogsRepository
	logger                    *zap.Logger
}

type ApplicationLogsStreamReaderParams struct {
	fx.In
	ApplicationsLogsRepository repositories.ApplicationLogsRepository
	Logger                     *zap.Logger
}

func NewKafkaApplicationLogsStreamReader(params ApplicationLogsStreamReaderParams) *KafkaApplicationLogsStreamReader {

	kafkaReader := NewKafkaLogsStream[*repositories.ApplicationLogs](
		[]string{"kafka:9094"},
		"nodes",
		params.Logger,
	)

	return &KafkaApplicationLogsStreamReader{
		kafkaReader:               &kafkaReader,
		logger:                    params.Logger,
		applicationLogsRepository: params.ApplicationsLogsRepository,
	}
}

func (r *KafkaApplicationLogsStreamReader) Handle(
	ctx context.Context,
	applicationLogs *repositories.ApplicationLogs) error {

	err := r.applicationLogsRepository.InsertLogs(ctx, applicationLogs)
	if err != nil {
		r.logger.Error("Failed to index node logs", zap.Error(err))
	}

	return err
}

func (r *KafkaApplicationLogsStreamReader) Listen() {
	r.kafkaReader.logger.Info("Starting to listen for node logs")
	go r.kafkaReader.Listen()
	nodeStream := r.kafkaReader.Stream()

	for {
		log := <-nodeStream
		r.kafkaReader.logger.Debug("Got message from stream", zap.Any("log", log))
		r.Handle(context.Background(), log)
	}
}

func (s *KafkaApplicationLogsStreamReader) AddHandler(f func(*repositories.ApplicationLogs)) {
	s.kafkaReader.SetHandler(f)
}

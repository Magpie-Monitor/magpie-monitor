package logsstream

import (
	"context"
	"github.com/Magpie-Monitor/magpie-monitor/pkg/envs"
	"github.com/Magpie-Monitor/magpie-monitor/pkg/repositories"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"os"
)

var APPLICATION_LOGS_QUQUE_HOST_KEY = "LOGS_INGESTION_QUEUE_HOST"
var APPLICATION_LOGS_QUEUE_PORT_KEY = "LOGS_INGESTION_QUEUE_PORT"
var APPLICATION_LOGS_TOPIC = "applications"

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

	envs.ValidateEnvs("Failed to connect to Kafka for application logs",
		[]string{APPLICATION_LOGS_QUQUE_HOST_KEY, APPLICATION_LOGS_QUEUE_PORT_KEY})

	kafkaHost := os.Getenv(APPLICATION_LOGS_QUQUE_HOST_KEY)
	kafkaPort := os.Getenv(APPLICATION_LOGS_QUEUE_PORT_KEY)

	kafkaReader := NewKafkaLogsStream[*repositories.ApplicationLogs](
		kafkaHost,
		kafkaPort,
		APPLICATION_LOGS_TOPIC,
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

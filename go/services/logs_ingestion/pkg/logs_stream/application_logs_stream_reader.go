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
var APPLICATION_LOGS_QUEUE_USERNAME_KEY = "LOGS_INGESTION_QUEUE_USERNAME"
var APPLICATION_LOGS_QUEUE_PASSWORD_KEY = "LOGS_INGESTION_QUEUE_PASSWORD"
var APPLICATION_LOGS_TOPIC_KEY = "LOGS_INGESTION_APPLICATION_LOGS_TOPIC"

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
		[]string{APPLICATION_LOGS_QUQUE_HOST_KEY,
			APPLICATION_LOGS_QUEUE_PORT_KEY,
			APPLICATION_LOGS_QUEUE_PASSWORD_KEY,
			APPLICATION_LOGS_QUEUE_USERNAME_KEY,
			APPLICATION_LOGS_TOPIC_KEY,
		},
	)

	kafkaHost := os.Getenv(APPLICATION_LOGS_QUQUE_HOST_KEY)
	kafkaPort := os.Getenv(APPLICATION_LOGS_QUEUE_PORT_KEY)
	kafkaUsername := os.Getenv(APPLICATION_LOGS_QUEUE_USERNAME_KEY)
	kafkaPassword := os.Getenv(APPLICATION_LOGS_QUEUE_PASSWORD_KEY)
	kafkaTopic := os.Getenv(APPLICATION_LOGS_TOPIC_KEY)

	kafkaReader := NewKafkaLogsStream[*repositories.ApplicationLogs](&KafkaLogsStreamParams{
		Host:     kafkaHost,
		Port:     kafkaPort,
		Topic:    kafkaTopic,
		Username: kafkaUsername,
		Password: kafkaPassword,
		Logger:   params.Logger,
	})

	return &KafkaApplicationLogsStreamReader{
		kafkaReader:               &kafkaReader,
		logger:                    params.Logger,
		applicationLogsRepository: params.ApplicationsLogsRepository,
	}
}

func (r *KafkaApplicationLogsStreamReader) Handle(
	ctx context.Context,
	applicationLogs *repositories.ApplicationLogs) error {

	r.kafkaReader.logger.Debug("Got message from stream", zap.Any("log", applicationLogs))

	_, err := r.applicationLogsRepository.InsertLogs(ctx, applicationLogs)
	if err != nil {
		r.logger.Error("Failed to index application logs", zap.Error(err))
	}

	return err
}

func (r *KafkaApplicationLogsStreamReader) Listen() {

	r.kafkaReader.logger.Info("Starting to listen for application logs at", zap.String(
		"addr", r.kafkaReader.brokers[0],
	))
	go r.kafkaReader.Listen()
	applicationLogs := r.kafkaReader.Stream()

	for {
		log := <-applicationLogs
		go r.Handle(context.Background(), log)
	}
}

func (s *KafkaApplicationLogsStreamReader) AddHandler(f func(*repositories.ApplicationLogs)) {
	s.kafkaReader.SetHandler(f)
}

package logsstream

import (
	"context"
	"github.com/Magpie-Monitor/magpie-monitor/pkg/envs"
	"github.com/Magpie-Monitor/magpie-monitor/pkg/repositories"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"os"
)

var NODE_LOGS_QUEUE_HOST_KEY = "LOGS_INGESTION_QUEUE_HOST"
var NODE_LOGS_QUEUE_PORT_KEY = "LOGS_INGESTION_QUEUE_PORT"
var NODE_LOGS_QUEUE_USERNAME_KEY = "LOGS_INGESTION_QUEUE_USERNAME"
var NODE_LOGS_QUEUE_PASSWORD_KEY = "LOGS_INGESTION_QUEUE_PASSWORD"
var NODE_LOGS_TOPIC_KEY = "LOGS_INGESTION_NODE_LOGS_TOPIC"

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

	envs.ValidateEnvs("Failed to connect to Kafka for node logs",
		[]string{NODE_LOGS_QUEUE_HOST_KEY,
			NODE_LOGS_QUEUE_PORT_KEY,
			NODE_LOGS_QUEUE_USERNAME_KEY,
			NODE_LOGS_QUEUE_PASSWORD_KEY,
			NODE_LOGS_TOPIC_KEY})

	kafkaHost := os.Getenv(NODE_LOGS_QUEUE_HOST_KEY)
	kafkaPort := os.Getenv(NODE_LOGS_QUEUE_PORT_KEY)
	kafkaUsername := os.Getenv(NODE_LOGS_QUEUE_USERNAME_KEY)
	kafkaPassword := os.Getenv(NODE_LOGS_QUEUE_PASSWORD_KEY)
	kafkaTopic := os.Getenv(NODE_LOGS_TOPIC_KEY)

	kafkaReader := NewKafkaLogsStream[*repositories.NodeLogs](
		&KafkaLogsStreamParams{
			Host:     kafkaHost,
			Port:     kafkaPort,
			Username: kafkaUsername,
			Password: kafkaPassword,
			Topic:    kafkaTopic,
			Logger:   params.Logger,
		})

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
	r.kafkaReader.logger.Info("Starting to listen for node logs at", zap.String(
		"addr", r.kafkaReader.brokers[0],
	))
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

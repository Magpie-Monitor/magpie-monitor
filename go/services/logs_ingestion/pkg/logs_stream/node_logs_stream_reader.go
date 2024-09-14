package logsstream

import (
	"context"
	"github.com/Magpie-Monitor/magpie-monitor/pkg/envs"
	"github.com/Magpie-Monitor/magpie-monitor/pkg/repositories"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"os"
)

var NODE_LOGS_QUQUE_HOST_KEY = "LOGS_INGESTION_QUEUE_HOST"
var NODE_LOGS_QUEUE_PORT_KEY = "LOGS_INGESTION_QUEUE_PORT"
var NODE_LOGS_TOPIC = "nodes"

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
		[]string{NODE_LOGS_QUQUE_HOST_KEY, NODE_LOGS_QUEUE_PORT_KEY})

	kafkaHost := os.Getenv(NODE_LOGS_QUQUE_HOST_KEY)
	kafkaPort := os.Getenv(NODE_LOGS_QUEUE_PORT_KEY)

	kafkaReader := NewKafkaLogsStream[*repositories.NodeLogs](
		kafkaHost,
		kafkaPort,
		NODE_LOGS_TOPIC,
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

package logsgeneration

import (
	"context"
	"fmt"
	"github.com/Magpie-Monitor/magpie-monitor/pkg/envs"
	"github.com/segmentio/kafka-go"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"os"
)

var APPLICATION_LOGS_QUQUE_HOST_KEY = "LOGS_INGESTION_QUEUE_HOST"
var APPLICATION_LOGS_QUEUE_PORT_KEY = "LOGS_INGESTION_QUEUE_PORT"
var APPLICATION_LOGS_TOPIC = "applications"

var NODE_LOGS_QUQUE_HOST_KEY = "LOGS_INGESTION_QUEUE_HOST"
var NODE_LOGS_QUEUE_PORT_KEY = "LOGS_INGESTION_QUEUE_PORT"
var NODE_LOGS_TOPIC = "nodes"

type KafkaLogsStreamWriter struct {
	nodeLogsWriter        *kafka.Writer
	applicationLogsWriter *kafka.Writer
	logger                *zap.Logger
}

func NewKafkaLogsStreamWriter(lc fx.Lifecycle, logger *zap.Logger) *KafkaLogsStreamWriter {

	envs.ValidateEnvs(
		"Host and port for log ingestion queue are not set",
		[]string{
			APPLICATION_LOGS_QUEUE_PORT_KEY,
			APPLICATION_LOGS_QUQUE_HOST_KEY,
			NODE_LOGS_QUEUE_PORT_KEY,
			NODE_LOGS_QUQUE_HOST_KEY,
		},
	)

	nodeLogsWriter := &kafka.Writer{
		Addr: kafka.TCP(fmt.Sprintf("%s:%s",
			os.Getenv(NODE_LOGS_QUQUE_HOST_KEY),
			os.Getenv(NODE_LOGS_QUEUE_PORT_KEY))),
		Topic:                  NODE_LOGS_TOPIC,
		AllowAutoTopicCreation: true,
	}

	applicationLogsWriter := &kafka.Writer{
		Addr: kafka.TCP(fmt.Sprintf(
			fmt.Sprintf("%s:%s",
				os.Getenv(APPLICATION_LOGS_QUQUE_HOST_KEY),
				os.Getenv(APPLICATION_LOGS_QUEUE_PORT_KEY)),
		)),
		Topic:                  APPLICATION_LOGS_TOPIC,
		AllowAutoTopicCreation: true,
	}

	kafkaWriter := &KafkaLogsStreamWriter{
		nodeLogsWriter:        nodeLogsWriter,
		applicationLogsWriter: applicationLogsWriter,
		logger:                logger,
	}

	return kafkaWriter

}

func (w *KafkaLogsStreamWriter) WriteApplicationLogs(ctx context.Context, key string, value string) error {
	err := w.applicationLogsWriter.WriteMessages(context.Background(), kafka.Message{
		Key:   []byte(key),
		Value: []byte(value),
	})

	if err != nil {
		w.logger.Error("Failed to write application logs to kafka", zap.Error(err))
		return err
	}

	return nil
}

func (w *KafkaLogsStreamWriter) WriteNodeLogs(ctx context.Context, key string, value string) error {
	err := w.nodeLogsWriter.WriteMessages(context.Background(), kafka.Message{
		Key:   []byte(key),
		Value: []byte(value),
	})

	if err != nil {
		w.logger.Error("Failed to write a node logs to kafka", zap.Error(err))
		return err
	}

	return nil
}

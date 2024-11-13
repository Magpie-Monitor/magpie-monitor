package tests

import (
	"context"
	"fmt"
	"github.com/Magpie-Monitor/magpie-monitor/pkg/envs"
	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/plain"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"os"
)

var LOGS_QUEUE_USERNAME_KEY = "LOGS_INGESTION_QUEUE_USERNAME"
var LOGS_QUEUE_PASSWORD_KEY = "LOGS_INGESTION_QUEUE_PASSWORD"

var APPLICATION_LOGS_QUQUE_HOST_KEY = "LOGS_INGESTION_QUEUE_HOST"
var APPLICATION_LOGS_QUEUE_PORT_KEY = "LOGS_INGESTION_QUEUE_PORT"
var APPLICATION_LOGS_TOPIC_KEY = "LOGS_INGESTION_APPLICATION_LOGS_TOPIC"

var NODE_LOGS_QUQUE_HOST_KEY = "LOGS_INGESTION_QUEUE_HOST"
var NODE_LOGS_QUEUE_PORT_KEY = "LOGS_INGESTION_QUEUE_PORT"
var NODE_LOGS_TOPIC_KEY = "LOGS_INGESTION_NODE_LOGS_TOPIC"

type KafkaLogsStreamWriter struct {
	nodeLogsWriter        *kafka.Writer
	applicationLogsWriter *kafka.Writer
	logger                *zap.Logger
}

func NewKafkaLogsStreamWriter(lc fx.Lifecycle, logger *zap.Logger) *KafkaLogsStreamWriter {

	envs.ValidateEnvs(
		"Host/port/username/password for log ingestion queue are not set",
		[]string{
			APPLICATION_LOGS_QUEUE_PORT_KEY,
			APPLICATION_LOGS_QUQUE_HOST_KEY,
			NODE_LOGS_QUEUE_PORT_KEY,
			NODE_LOGS_QUQUE_HOST_KEY,
			LOGS_QUEUE_USERNAME_KEY,
			LOGS_QUEUE_PASSWORD_KEY,
			APPLICATION_LOGS_TOPIC_KEY,
			NODE_LOGS_TOPIC_KEY,
		},
	)

	username := os.Getenv(LOGS_QUEUE_USERNAME_KEY)
	password := os.Getenv(LOGS_QUEUE_PASSWORD_KEY)
	nodeTopic := os.Getenv(NODE_LOGS_TOPIC_KEY)
	applicationTopic := os.Getenv(APPLICATION_LOGS_TOPIC_KEY)

	nodeLogsWriter := &kafka.Writer{
		Addr: kafka.TCP(fmt.Sprintf("%s:%s",
			os.Getenv(NODE_LOGS_QUQUE_HOST_KEY),
			os.Getenv(NODE_LOGS_QUEUE_PORT_KEY))),
		Topic:                  nodeTopic,
		AllowAutoTopicCreation: true,
		Transport:              &kafka.Transport{SASL: plain.Mechanism{Username: username, Password: password}},
	}

	applicationLogsWriter := &kafka.Writer{
		Addr: kafka.TCP(fmt.Sprintf(
			fmt.Sprintf("%s:%s",
				os.Getenv(APPLICATION_LOGS_QUQUE_HOST_KEY),
				os.Getenv(APPLICATION_LOGS_QUEUE_PORT_KEY)),
		)),
		Topic:                  applicationTopic,
		AllowAutoTopicCreation: true,
		Transport:              &kafka.Transport{SASL: plain.Mechanism{Username: username, Password: password}},
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

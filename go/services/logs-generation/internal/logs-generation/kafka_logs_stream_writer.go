package logsgeneration

import (
	"context"
	"github.com/segmentio/kafka-go"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type KafkaLogsStreamWriter struct {
	writer *kafka.Writer
	logger *zap.Logger
}

func NewKafkaLogsStreamWriter(lc fx.Lifecycle, logger *zap.Logger) *KafkaLogsStreamWriter {

	writer := &kafka.Writer{
		Addr:                   kafka.TCP("kafka:9094"),
		Topic:                  "nodes",
		AllowAutoTopicCreation: true,
	}

	kafkaWriter := &KafkaLogsStreamWriter{
		writer: writer,
		logger: logger,
	}

	return kafkaWriter

}

func (w *KafkaLogsStreamWriter) Write(ctx context.Context, key string, value string) error {
	err := w.writer.WriteMessages(ctx, kafka.Message{
		Key:   []byte(key),
		Value: []byte(value),
	})

	if err != nil {
		w.logger.Error("Failed to write a message to kafka", zap.Error(err))
		return err
	}

	return nil
}

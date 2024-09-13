package main

import (
	"context"
	"github.com/Magpie-Monitor/magpie-monitor/services/logs-generation/internal/logs-generation"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
	"time"
)

type LogsGenerator struct {
	writer *logsgeneration.KafkaLogsStreamWriter
	logger *zap.Logger
}

func (g *LogsGenerator) Write(ctx context.Context) {

	for {
		err := g.writer.Write(ctx, "test", "message")
		if err != nil {
			g.logger.Error("Failed to write a message", zap.Error(err))
		} else {
			g.logger.Info("Sent a message")
		}
		time.Sleep(time.Second * 2)
	}
}

func NewLogsGenerator(lc fx.Lifecycle, logger *zap.Logger, writer *logsgeneration.KafkaLogsStreamWriter) *LogsGenerator {

	generator := &LogsGenerator{
		logger: logger,
		writer: writer,
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			logger.Info("Starting sending logs to the queue")
			go generator.Write(ctx)
			return nil
		},
	},
	)

	return generator
}

func main() {
	fx.New(
		fx.WithLogger(func(log *zap.Logger) fxevent.Logger {
			return &fxevent.ZapLogger{Logger: log}
		}),
		fx.Provide(
			logsgeneration.NewKafkaLogsStreamWriter,
			NewLogsGenerator,
			zap.NewExample),
		fx.Invoke(func(*LogsGenerator) {}),
	).Run()
}

package main

import (
	"context"
	"encoding/json"
	"time"

	"github.com/Magpie-Monitor/magpie-monitor/pkg/repositories"
	"github.com/Magpie-Monitor/magpie-monitor/services/logs-generation/internal/logs-generation"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
)

type LogsGenerator struct {
	writer *logsgeneration.KafkaLogsStreamWriter
	logger *zap.Logger
}

func (g *LogsGenerator) WriteNodeLogs(ctx context.Context) {

	nodeLogs := repositories.NodeLogs{
		Cluster:   "cluster-1",
		Kind:      "node",
		Timestamp: time.Now().Unix(),
		Name:      "host-1",
		Namespace: "nms-1",
		Content:   "log-contentkldj",
	}

	jsonNodeLogs, err := json.Marshal(nodeLogs)
	if err != nil {
		g.logger.Error("Failed to encode node logs", zap.Error(err))
	}

	for {

		err = g.writer.Write(ctx, time.Now().String(), string(jsonNodeLogs))
		if err != nil {
			g.logger.Error("Failed to write a message", zap.Error(err))
		} else {
			g.logger.Info("Sent a node log")
		}

		time.Sleep(time.Second * 2)
	}
}

func (g *LogsGenerator) WriteApplicationLogs(ctx context.Context) {

	applicationLogs := repositories.ApplicationLogs{
		Cluster:   "cluster-1",
		Kind:      "application",
		Timestamp: time.Now().Unix(),
		Name:      "my-cool-app",
		Pods: []*repositories.PodLogs{
			{
				Name: "pod-1",
				Containers: []*repositories.ContainerLogs{
					{
						Name:    "container-1",
						Image:   "container-1-image",
						Content: "container-logs-content",
					},
				},
			},
		},
	}

	jsonApplicationLogs, err := json.Marshal(applicationLogs)
	if err != nil {
		g.logger.Error("Failed to encode applicatation logs", zap.Error(err))
	}
	for {
		err := g.writer.Write(ctx, time.Now().String(), string(jsonApplicationLogs))
		if err != nil {
			g.logger.Error("Failed to write a message", zap.Error(err))
		} else {
			g.logger.Info("Sent an applicatoion log")
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
			go generator.WriteNodeLogs(ctx)
			go generator.WriteApplicationLogs(ctx)
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

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Magpie-Monitor/magpie-monitor/pkg/repositories"
	"github.com/Magpie-Monitor/magpie-monitor/services/logs_generation/internal/logs_generation"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
	"math/rand"
	"time"
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

		err = g.writer.WriteNodeLogs(ctx, time.Now().String(), string(jsonNodeLogs))

		if err != nil {
			g.logger.Error("Failed to write a message", zap.Error(err))
		} else {
			g.logger.Debug("Sent a node log")

		}

		time.Sleep(time.Second * 1000)
	}
}

func (g *LogsGenerator) WriteApplicationLogs(ctx context.Context) {

	apps := []string{"app-1", "app-2", "app-3", "app-3"}

	for {

		nodeLogs := repositories.NodeLogs{
			Cluster:   "testcluster",
			Kind:      "node",
			Timestamp: 1728313197000000020,
			Name:      "tools",
			Namespace: "nms",
			Content:   "Failed to save new nginx configuration. Out of disk space.",
		}

		applicationLogs := repositories.ApplicationLogs{
			Cluster:   "testcluster",
			Kind:      "application",
			Timestamp: 1728313197000000010,
			Name:      apps[rand.Intn(len(apps))],
			Pods: []*repositories.PodLogs{
				{
					Name: "pod-1",
					Containers: []*repositories.ContainerLogs{
						{
							Name:    "container-x",
							Image:   "container-x-image",
							Content: "Failed to connect the psql database",
						},
						{
							Name:    "container-2",
							Image:   "container-2-image",
							Content: "Failed to connect the mysql database",
						},
					},
				},
			},
		}

		jsonNodeLogs, _ := json.Marshal(nodeLogs)
		jsonApplicationLogs, _ := json.Marshal(applicationLogs)

		g.handleApplicationLogs(ctx, string(jsonApplicationLogs))
		g.handleNodeLogs(ctx, string(jsonNodeLogs))
		time.Sleep(time.Second * 3609)

	}

}

func (g *LogsGenerator) handleApplicationLogs(ctx context.Context, json string) {

	err := g.writer.WriteApplicationLogs(ctx, fmt.Sprintf("%s%s", time.Now().String(), "2"), string(json))
	if err != nil {
		g.logger.Error("Failed to write a message", zap.Error(err))
	} else {
		g.logger.Info("Sent application logs")
	}
}

func (g *LogsGenerator) handleNodeLogs(ctx context.Context, json string) {

	err := g.writer.WriteNodeLogs(ctx, fmt.Sprintf("%s%s", time.Now().String(), "2"), string(json))
	if err != nil {
		g.logger.Error("Failed to write a message", zap.Error(err))
	} else {
		g.logger.Info("Sent node logs")
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

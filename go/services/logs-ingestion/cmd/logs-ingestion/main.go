package main

import (
	"context"

	"github.com/Magpie-Monitor/magpie-monitor/pkg/elasticsearch"
	"github.com/Magpie-Monitor/magpie-monitor/pkg/repositories"
	logsstream "github.com/Magpie-Monitor/magpie-monitor/services/logs-ingestion/pkg/logs-stream"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
)

type LogsStreamListener struct {
	applicationLogsReader *logsstream.ApplicationLogsStreamReader
	nodeLogsReader        *logsstream.NodeLogsStreamReader
	logger                *zap.Logger
}

func NewLogsStreamListener(
	lc fx.Lifecycle,
	logger *zap.Logger,
	applicationLogsReader *logsstream.ApplicationLogsStreamReader,
	nodeLogsReader *logsstream.NodeLogsStreamReader,
) *LogsStreamListener {

	listener := LogsStreamListener{
		logger:                logger,
		applicationLogsReader: applicationLogsReader,
		nodeLogsReader:        nodeLogsReader,
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {

			listener.nodeLogsReader.SetHandler(func(logs repositories.NodeLogs) {
				logger.Info("Recievied", zap.Any("logs", logs))
			})

			logger.Info("Starting listening for logs from", zap.String("addr", "kafka:9094"))

			go listener.nodeLogsReader.Listen()
			return nil
		},
	})

	return &listener
}

func main() {
	fx.New(

		fx.WithLogger(func(log *zap.Logger) fxevent.Logger {
			return &fxevent.ZapLogger{Logger: log}
		}),
		fx.Provide(

			NewLogsStreamListener,
			logsstream.NewApplicationLogsStreamReader,
			logsstream.NewNodeLogsStreamReader,

			elasticsearch.NewElasticSearchLogsDbClient,
			repositories.ProvideAsApplicationLogsRepository(
				repositories.NewElasticSearchApplicationLogsRepository,
			),

			repositories.ProvideAsNodeLogsRepository(
				repositories.NewElasticSearchNodeLogsRepository,
			),

			zap.NewExample),
		fx.Invoke(func(*LogsStreamListener) {}),
	).Run()
}

package config

import (
	"context"
	"os"

	"github.com/Magpie-Monitor/magpie-monitor/pkg/elasticsearch"
	"github.com/Magpie-Monitor/magpie-monitor/pkg/repositories"
	logsstream "github.com/Magpie-Monitor/magpie-monitor/services/logs_ingestion/pkg/logs_stream"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
)

var AppModule fx.Option

const (
	PRODUCTION_ENVIRONMENT = "prod"
	TEST_ENVIRONMENT       = "test"
)

type LogsStreamListener struct {
	applicationLogsReader logsstream.ApplicationLogsStreamReader
	nodeLogsReader        logsstream.NodeLogsStreamReader
	logger                *zap.Logger
}

func NewLogsStreamListener(
	lc fx.Lifecycle,
	logger *zap.Logger,
	applicationLogsReader logsstream.ApplicationLogsStreamReader,
	nodeLogsReader logsstream.NodeLogsStreamReader,
) *LogsStreamListener {

	listener := LogsStreamListener{
		logger:                logger,
		applicationLogsReader: applicationLogsReader,
		nodeLogsReader:        nodeLogsReader,
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {

			go listener.nodeLogsReader.Listen()
			go listener.applicationLogsReader.Listen()
			return nil
		},
	})

	return &listener
}

func init() {
	env := os.Getenv("APP_ENV")
	if env == TEST_ENVIRONMENT {
		AppModule = fx.Options(
			fx.WithLogger(func(log *zap.Logger) fxevent.Logger {
				return &fxevent.ZapLogger{Logger: log}
			}),
			fx.Provide(
				NewLogsStreamListener,
				fx.Annotate(
					logsstream.NewKafkaApplicationLogsStreamReader,
					fx.As(new(logsstream.ApplicationLogsStreamReader)),
				),

				fx.Annotate(
					logsstream.NewKafkaNodeLogsStreamReader,
					fx.As(new(logsstream.NodeLogsStreamReader)),
				),

				elasticsearch.NewElasticSearchLogsDbClient,
				repositories.ProvideAsApplicationLogsRepository(
					repositories.NewElasticSearchApplicationLogsRepository,
				),

				repositories.ProvideAsNodeLogsRepository(
					repositories.NewElasticSearchNodeLogsRepository,
				),

				zap.NewExample),
		)
	} else {
		AppModule = fx.Options(
			fx.WithLogger(func(log *zap.Logger) fxevent.Logger {
				return &fxevent.ZapLogger{Logger: log}
			}),
			fx.Provide(
				NewLogsStreamListener,
				fx.Annotate(
					logsstream.NewKafkaApplicationLogsStreamReader,
					fx.As(new(logsstream.ApplicationLogsStreamReader)),
				),

				fx.Annotate(
					logsstream.NewKafkaNodeLogsStreamReader,
					fx.As(new(logsstream.NodeLogsStreamReader)),
				),

				elasticsearch.NewElasticSearchLogsDbClient,
				repositories.ProvideAsApplicationLogsRepository(
					repositories.NewElasticSearchApplicationLogsRepository,
				),

				repositories.ProvideAsNodeLogsRepository(
					repositories.NewElasticSearchNodeLogsRepository,
				),
				zap.NewExample),
		)
	}
}

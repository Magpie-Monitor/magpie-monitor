package config

import (
	"fmt"
	"os"

	"github.com/Magpie-Monitor/magpie-monitor/pkg/elasticsearch"
	"github.com/Magpie-Monitor/magpie-monitor/pkg/repositories"
	"github.com/Magpie-Monitor/magpie-monitor/pkg/tests"
	logsstream "github.com/Magpie-Monitor/magpie-monitor/services/logs_ingestion/pkg/logs_stream"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
)

var AppModule fx.Option

func init() {
	env := os.Getenv("APP_ENV")
	fmt.Printf("Starting the app in %s environment", env)

	if env == tests.TEST_ENVIRONMENT {
		AppModule = fx.Options(
			fx.WithLogger(func(log *zap.Logger) fxevent.Logger {
				return &fxevent.ZapLogger{Logger: log}
			}),
			fx.Provide(
				logsstream.NewLogsStreamListener,
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

				//Directly providing implementation for tests
				logsstream.NewKafkaApplicationLogsStreamReader,

				//Directly providing implementation for tests
				logsstream.NewKafkaNodeLogsStreamReader,

				// Mock logs writer used for tests
				tests.NewKafkaLogsStreamWriter,

				zap.NewExample),
		)
	} else {
		AppModule = fx.Options(
			fx.WithLogger(func(log *zap.Logger) fxevent.Logger {
				return &fxevent.ZapLogger{Logger: log}
			}),
			fx.Provide(
				logsstream.NewLogsStreamListener,
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

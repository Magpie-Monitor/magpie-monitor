package logsstream_test

import (
	"os"
	"testing"

	"github.com/Magpie-Monitor/magpie-monitor/services/logs_ingestion/pkg/config"
	logsstream "github.com/Magpie-Monitor/magpie-monitor/services/logs_ingestion/pkg/logs_stream"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
	"go.uber.org/zap"
)

func TestApplicationLogsStreamReader(t *testing.T) {
	os.Setenv("APP_ENV", "test")
	defer os.Unsetenv("APP_ENV")

	app := fxtest.New(
		t,
		config.AppModule,
		fx.Invoke(func(reader logsstream.ApplicationLogsStreamReader, logger *zap.Logger) {
			if reader == nil {
				t.Fatal("Failed to load reader")
			}

			logger.Info("Success", zap.Any("reader", reader))

		}),
	)

	app.Run()
}

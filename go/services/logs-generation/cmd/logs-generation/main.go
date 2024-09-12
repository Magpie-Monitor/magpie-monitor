package main

import (
	"github.com/Magpie-Monitor/magpie-monitor/pkg/elasticsearch"
	"github.com/Magpie-Monitor/magpie-monitor/pkg/repositories"
	"github.com/Magpie-Monitor/magpie-monitor/pkg/routing"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
	"net/http"
)

func main() {
	fx.New(

		fx.WithLogger(func(log *zap.Logger) fxevent.Logger {
			return &fxevent.ZapLogger{Logger: log}
		}),
		fx.Provide(
			routing.ProvideAsRootServeMux(
				routing.NewServeMux,
			),

			NewHTTPServer,
			routing.ProvideAsRoute(handlers.NewLogsIngestionRouter),
			handlers.NewLogsIngestionHandler,
			elasticsearch.NewElasticSearchLogsDbClient,
			repositories.ProvideAsApplicationLogsRepository(
				repositories.NewElasticSearchApplicationLogsRepository,
			),

			repositories.ProvideAsNodeLogsRepository(
				repositories.NewElasticSearchNodeLogsRepository,
			),

			zap.NewExample),
		fx.Invoke(func(*http.Server) {}),
	).Run()
}

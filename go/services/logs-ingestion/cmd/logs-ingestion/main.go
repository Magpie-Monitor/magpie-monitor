package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"

	"github.com/Magpie-Monitor/magpie-monitor/pkg/elasticsearch"
	"github.com/Magpie-Monitor/magpie-monitor/pkg/repositories"
	"github.com/Magpie-Monitor/magpie-monitor/pkg/routing"
	"github.com/Magpie-Monitor/magpie-monitor/services/logs-ingestion/internal/handlers"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
)

func NewHTTPServer(lc fx.Lifecycle, mux *http.ServeMux, log *zap.Logger) *http.Server {
	port := os.Getenv("LOGS_INGESTION_SERVICE_HTTP_PORT")

	srv := &http.Server{Addr: fmt.Sprintf(":%s", port), Handler: mux}
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			ln, err := net.Listen("tcp", srv.Addr)

			if err != nil {
				return err
			}

			log.Info("Starting HTTP server at", zap.String("addr", srv.Addr))
			go srv.Serve(ln)
			return nil
		},
		OnStop: func(ctx context.Context) error {

			log.Info("Shutting down the HTTP server at", zap.String("addr", srv.Addr))
			return srv.Shutdown(ctx)
		},
	})
	return srv
}

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

package main

import (
	"context"
	"fmt"
	elasticsearch "github.com/Magpie-Monitor/magpie-monitor/pkg/elasticsearch"
	sharedrepositories "github.com/Magpie-Monitor/magpie-monitor/pkg/repositories"
	"github.com/Magpie-Monitor/magpie-monitor/pkg/routing"
	"github.com/Magpie-Monitor/magpie-monitor/services/reports/internal/database"
	"github.com/Magpie-Monitor/magpie-monitor/services/reports/internal/handlers"
	"github.com/Magpie-Monitor/magpie-monitor/services/reports/pkg/repositories"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"net"
	"net/http"
	"os"
)

func NewHTTPServer(lc fx.Lifecycle, mux *http.ServeMux, log *zap.Logger) *http.Server {
	port := os.Getenv("REPORTS_SERVICE_HTTP_PORT")

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
		fx.Provide(

			NewHTTPServer,

			routing.ProvideAsRootServeMux(routing.NewServeMux),

			routing.ProvideAsRoute(handlers.NewReportsRouter),
			handlers.NewReportsHandler,

			database.NewReportsDbMongoClient,
			repositories.ProvideAsReportRepository(
				repositories.NewMongoDbReportRepository,
			),

			elasticsearch.NewElasticSearchLogsDbClient,
			sharedrepositories.ProvideAsNodeLogsRepository(
				sharedrepositories.NewElasticSearchNodeLogsRepository,
			),

			sharedrepositories.ProvideAsApplicationLogsRepository(
				sharedrepositories.NewElasticSearchApplicationLogsRepository,
			),

			zap.NewExample),
		fx.Invoke(func(*http.Server) {}),
	).Run()
}

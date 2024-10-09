package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"

	"github.com/Magpie-Monitor/magpie-monitor/pkg/mongodb"
	"github.com/Magpie-Monitor/magpie-monitor/pkg/routing"
	"github.com/Magpie-Monitor/magpie-monitor/pkg/swagger"
	"github.com/Magpie-Monitor/magpie-monitor/services/cluster_metadata/internal/database"
	"github.com/Magpie-Monitor/magpie-monitor/services/cluster_metadata/internal/handlers"
	"github.com/Magpie-Monitor/magpie-monitor/services/cluster_metadata/pkg/repositories"
	"github.com/Magpie-Monitor/magpie-monitor/services/cluster_metadata/pkg/services"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
)

type ServerParams struct {
	fx.In
	Lc             fx.Lifecycle
	Logger         *zap.Logger
	MetadataRouter *handlers.MetadataRouter
	SwaggerRouter  *swagger.SwaggerRouter
}

func NewHTTPServer(ServerParams ServerParams) *http.Server {
	port := os.Getenv("HTTP_PORT")

	srv := &http.Server{Addr: fmt.Sprintf(":%s", port), Handler: ServerParams.MetadataRouter}
	ServerParams.Lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			ln, err := net.Listen("tcp", srv.Addr)
			if err != nil {
				return err
			}

			ServerParams.Logger.Info("Starting HTTP server at", zap.String("addr", srv.Addr))
			go srv.Serve(ln)
			return nil
		},
		OnStop: func(ctx context.Context) error {
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
			NewHTTPServer,

			routing.NewRootRouter,

			// TODO - fix
			swagger.NewSwaggerRouter,
			swagger.NewSwaggerHandler,
			swagger.ProvideSwaggerConfig(),

			handlers.NewMetadataRouter,
			handlers.NewMetadataHandler,

			mongodb.NewMongoDbClient,
			database.NewMongoDbConnectionDetails,

			services.NewMetadataService,

			repositories.NewClusterMetadataCollection,
			repositories.NewNodeMetadataCollection,

			zap.NewProduction,
		),
		fx.Invoke(func(*http.Server) {}),
	).Run()
}

package main

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"github.com/Magpie-Monitor/magpie-monitor/pkg/routing"
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
}

func NewHTTPServer(ServerParams ServerParams) *http.Server {
	// TODO - add variable
	port := "8080"

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
			routing.NewRootRouter,
			NewHTTPServer,
			handlers.NewMetadataRouter,
			handlers.NewMetadataHandler,
			database.NewMetadataDbMongoClient,
			services.NewMetadataService,
			repositories.NewClusterMetadataCollection,
			repositories.NewNodeMetadataCollection,
			zap.NewProduction,
		),
		fx.Invoke(func(*http.Server) {}),
	).Run()
}

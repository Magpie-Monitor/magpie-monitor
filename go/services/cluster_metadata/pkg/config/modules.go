package config

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"

	messagebroker "github.com/Magpie-Monitor/magpie-monitor/pkg/message-broker"
	"github.com/Magpie-Monitor/magpie-monitor/pkg/mongodb"
	"github.com/Magpie-Monitor/magpie-monitor/pkg/routing"
	"github.com/Magpie-Monitor/magpie-monitor/pkg/swagger"
	"github.com/Magpie-Monitor/magpie-monitor/pkg/tests"
	"github.com/Magpie-Monitor/magpie-monitor/services/cluster_metadata/internal/database"
	"github.com/Magpie-Monitor/magpie-monitor/services/cluster_metadata/internal/handlers"
	"github.com/Magpie-Monitor/magpie-monitor/services/cluster_metadata/pkg/repositories"
	"github.com/Magpie-Monitor/magpie-monitor/services/cluster_metadata/pkg/services"
	"github.com/gorilla/mux"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type ServerParams struct {
	fx.In
	Lc             fx.Lifecycle
	Logger         *zap.Logger
	MetadataRouter *handlers.MetadataRouter
	RootRouter     *mux.Router
	SwaggerRouter  *swagger.SwaggerRouter
}

func NewHTTPServer(ServerParams ServerParams) *http.Server {
	port := os.Getenv("HTTP_PORT")

	srv := &http.Server{Addr: fmt.Sprintf(":%s", port), Handler: ServerParams.RootRouter}
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

var AppModule fx.Option

func init() {
	env := os.Getenv("APP_ENV")
	fmt.Printf("Starting the app in %s environment", env)

	if env == tests.TEST_ENVIRONMENT {
		AppModule = fx.Options(
			// fx.WithLogger(func(log *zap.Logger) fxevent.Logger {
			// 	return &fxevent.ZapLogger{Logger: log}
			// }),
			fx.Provide(
				NewHTTPServer,

				routing.NewRootRouter,

				swagger.NewSwaggerRouter,
				swagger.NewSwaggerHandler,
				swagger.ProvideSwaggerConfig(),

				handlers.NewMetadataRouter,
				handlers.NewMetadataHandler,

				mongodb.NewMongoDbClient,
				database.NewMongoDbConnectionDetails,
				messagebroker.NewKafkaCredentials,

				services.NewMetadataService,
				services.NewMetadataEventPublisher,

				services.NewApplicationMetadataBroker,
				services.NewNodeMetadataBroker,

				repositories.NewApplicationMetadataCollection,
				repositories.NewNodeMetadataCollection,
				repositories.NewApplicationAggregatedMetadataCollection,
				repositories.NewNodeAggregatedMetadataCollection,
				repositories.NewClusterAggregatedStateCollection,

				services.NewApplicationMetadataUpdatedBroker,
				services.NewNodeMetadataUpdatedBroker,
				services.NewClusterMetadataUpdatedBroker,

				zap.NewProduction,
			),
		)
	} else {
		AppModule = fx.Options(
			// fx.WithLogger(func(log *zap.Logger) fxevent.Logger {
			// 	return &fxevent.ZapLogger{Logger: log}
			// }),
			fx.Provide(
				NewHTTPServer,

				routing.NewRootRouter,

				swagger.NewSwaggerRouter,
				swagger.NewSwaggerHandler,
				swagger.ProvideSwaggerConfig(),

				handlers.NewMetadataRouter,
				handlers.NewMetadataHandler,

				mongodb.NewMongoDbClient,
				database.NewMongoDbConnectionDetails,
				messagebroker.NewKafkaCredentials,

				services.NewMetadataService,
				services.NewMetadataEventPublisher,

				services.NewApplicationMetadataBroker,
				services.NewNodeMetadataBroker,

				repositories.NewApplicationMetadataCollection,
				repositories.NewNodeMetadataCollection,
				repositories.NewApplicationAggregatedMetadataCollection,
				repositories.NewNodeAggregatedMetadataCollection,
				repositories.NewClusterAggregatedStateCollection,

				zap.NewProduction,
			),
		)
	}
}

package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"

	elasticsearch "github.com/Magpie-Monitor/magpie-monitor/pkg/elasticsearch"
	sharedrepositories "github.com/Magpie-Monitor/magpie-monitor/pkg/repositories"
	"github.com/Magpie-Monitor/magpie-monitor/pkg/routing"
	"github.com/Magpie-Monitor/magpie-monitor/pkg/swagger"
	"github.com/Magpie-Monitor/magpie-monitor/services/reports/internal/brokers"
	"github.com/Magpie-Monitor/magpie-monitor/services/reports/internal/database"
	"github.com/Magpie-Monitor/magpie-monitor/services/reports/internal/handlers"
	"github.com/Magpie-Monitor/magpie-monitor/services/reports/internal/services"
	"github.com/Magpie-Monitor/magpie-monitor/services/reports/pkg/insights"
	"github.com/Magpie-Monitor/magpie-monitor/services/reports/pkg/openai"
	"github.com/Magpie-Monitor/magpie-monitor/services/reports/pkg/repositories"
	"github.com/gorilla/mux"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
)

type ServerParams struct {
	fx.In
	Lc              fx.Lifecycle
	Logger          *zap.Logger
	RootRouter      *mux.Router
	ReportsRouter   *handlers.ReportsRouter
	SwaggerRouter   *swagger.SwaggerRouter
	IncidentsRouter *handlers.IncidentsRouter
}

func NewHTTPServer(serverParams ServerParams) *http.Server {
	port := os.Getenv("REPORTS_SERVICE_PORT")

	srv := &http.Server{Addr: fmt.Sprintf(":%s", port), Handler: serverParams.RootRouter}
	serverParams.Lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			ln, err := net.Listen("tcp", srv.Addr)

			if err != nil {
				return err
			}

			serverParams.Logger.Info("Starting HTTP server at", zap.String("addr", srv.Addr))
			go srv.Serve(ln)
			return nil
		},
		OnStop: func(ctx context.Context) error {

			serverParams.Logger.Info("Shutting down the HTTP server at", zap.String("addr", srv.Addr))
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
			services.NewReportsService,
			handlers.NewReportsRouter,
			handlers.NewReportsHandler,

			handlers.NewIncidentsRouter,
			handlers.NewIncidentHandler,
			services.NewNodeIncidentsService,
			services.NewApplicationIncidentsService,

			swagger.NewSwaggerRouter,
			swagger.NewSwaggerHandler,
			swagger.ProvideSwaggerConfig(),

			database.NewReportsDbMongoClient,
			repositories.NewReportCollection,
			repositories.ProvideAsReportRepository(
				repositories.NewMongoDbReportRepository,
			),

			repositories.NewNodeIncidentsCollection,
			repositories.ProvideAsNodeIncidentRepository(
				repositories.NewMongoDbNodeIncidentRepository,
			),

			openai.NewOpenAiJobsCollection,

			openai.ProvideAsOpenAiJobRepository(
				openai.NewMongoDbOpenAiJobRepository,
			),

			brokers.ProvideAsReportGeneratedBroker(
				brokers.NewReportGeneratedBroker,
			),

			brokers.ProvideAsReportRequestedBroker(
				brokers.NewReportRequestedBroker,
			),

			brokers.ProvideAsReportRequestFailedBroker(
				brokers.NewReportRequestFailedBroker,
			),

			repositories.NewApplicationIncidentsCollection,

			repositories.ProvideAsApplicationIncidentRepository(
				repositories.NewMongoDbApplicationIncidentRepository,
			),

			openai.NewBatchPoller,

			openai.ProvideAsPendingBatchRepository(
				openai.NewRedisPendingBatchRepository,
			),

			elasticsearch.NewElasticSearchLogsDbClient,
			sharedrepositories.ProvideAsNodeLogsRepository(
				sharedrepositories.NewElasticSearchNodeLogsRepository,
			),

			sharedrepositories.ProvideAsApplicationLogsRepository(
				sharedrepositories.NewElasticSearchApplicationLogsRepository,
			),
			openai.NewOpenAiClient,

			fx.Annotate(
				insights.NewOpenAiInsightsGenerator,
				fx.As(new(insights.ApplicationInsightsGenerator)),
				fx.As(new(insights.NodeInsightsGenerator)),
			),

			zap.NewExample),
		fx.Invoke(func(server *http.Server,
			reportsService *services.ReportsService,
			batchPoller *openai.BatchPoller,
			reportsHandler *handlers.ReportsHandler) {

			// Listen for ReportRequested messages
			go reportsHandler.ListenForReportRequests()

			// Poll for pending reports
			go reportsHandler.PollReports()

			// Poll for pending OpenAi batches
			go batchPoller.Start()
		}),
	).Run()
}

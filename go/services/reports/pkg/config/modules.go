package config

import (
	"fmt"
	"github.com/Magpie-Monitor/magpie-monitor/pkg/elasticsearch"
	sharedrepositories "github.com/Magpie-Monitor/magpie-monitor/pkg/repositories"
	"github.com/Magpie-Monitor/magpie-monitor/pkg/tests"
	"github.com/Magpie-Monitor/magpie-monitor/services/reports/internal/brokers"
	"github.com/Magpie-Monitor/magpie-monitor/services/reports/internal/database"
	"github.com/Magpie-Monitor/magpie-monitor/services/reports/internal/handlers"
	"github.com/Magpie-Monitor/magpie-monitor/services/reports/internal/services"
	incidentcorrelation "github.com/Magpie-Monitor/magpie-monitor/services/reports/pkg/incident_correlation"
	"github.com/Magpie-Monitor/magpie-monitor/services/reports/pkg/insights"
	"github.com/Magpie-Monitor/magpie-monitor/services/reports/pkg/openai"
	"github.com/Magpie-Monitor/magpie-monitor/services/reports/pkg/repositories"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
	"os"
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

				services.NewReportsService,
				handlers.NewReportsHandler,

				services.NewNodeIncidentsService,
				services.NewApplicationIncidentsService,

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

				brokers.ProvideAsNodeIncidentSourcedBroker(
					brokers.NewNodeIncidentSourcedBroker,
				),

				brokers.ProvideAsApplicationIncidentSourcedBroker(
					brokers.NewApplicationIncidentSourcedBroker,
				),

				repositories.NewApplicationIncidentsCollection,

				repositories.ProvideAsApplicationIncidentRepository(
					repositories.NewMongoDbApplicationIncidentRepository,
				),

				openai.NewBatchPoller,

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

				fx.Annotate(
					incidentcorrelation.NewOpenAiIncidentMerger,
					fx.As(new(incidentcorrelation.IncidentMerger)),
				),

				// Adding direct implementation for tests
				insights.NewOpenAiInsightsGenerator,

				repositories.NewNodeIncidentSourcesCollection,
				repositories.ProvideAsNodeIncidentSourceRepository(
					repositories.NewMongoDbNodeIncidentSourceRepository,
				),

				repositories.NewApplicationIncidentSourcesCollection,
				repositories.ProvideAsApplicationIncidentSourceRepository(
					repositories.NewMongoDbApplicationIncidentSourceRepository,
				),

				zap.NewExample),
		)
	} else {
		AppModule = fx.Options(
			fx.WithLogger(func(log *zap.Logger) fxevent.Logger {
				return &fxevent.ZapLogger{Logger: log}
			}),
			fx.Provide(

				services.NewReportsService,
				handlers.NewReportsHandler,

				services.NewNodeIncidentsService,
				services.NewApplicationIncidentsService,

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

				brokers.ProvideAsNodeIncidentSourcedBroker(
					brokers.NewNodeIncidentSourcedBroker,
				),

				brokers.ProvideAsApplicationIncidentSourcedBroker(
					brokers.NewApplicationIncidentSourcedBroker,
				),

				repositories.NewApplicationIncidentsCollection,

				repositories.ProvideAsApplicationIncidentRepository(
					repositories.NewMongoDbApplicationIncidentRepository,
				),

				openai.NewBatchPoller,

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

				fx.Annotate(
					incidentcorrelation.NewOpenAiIncidentMerger,
					fx.As(new(incidentcorrelation.IncidentMerger)),
				),

				repositories.NewNodeIncidentSourcesCollection,
				repositories.ProvideAsNodeIncidentSourceRepository(
					repositories.NewMongoDbNodeIncidentSourceRepository,
				),

				repositories.NewApplicationIncidentSourcesCollection,
				repositories.ProvideAsApplicationIncidentSourceRepository(
					repositories.NewMongoDbApplicationIncidentSourceRepository,
				),

				zap.NewExample,
			),
		)
	}
}

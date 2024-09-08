package main

import (
	"context"
	"encoding/json"
	"fmt"
	elasticsearch "github.com/Magpie-Monitor/magpie-monitor/pkg/elasticsearch"
	sharedrepositories "github.com/Magpie-Monitor/magpie-monitor/pkg/repositories"
	"github.com/Magpie-Monitor/magpie-monitor/services/reports/internal/database"
	"github.com/Magpie-Monitor/magpie-monitor/services/reports/pkg/repositories"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"net"
	"net/http"
	"os"
)

type ReportsHandler struct {
	logger                    *zap.Logger
	reportRepository          repositories.ReportRepository
	applicationLogsRepository sharedrepositories.ApplicationLogsRepository
	nodeLogsRepository        sharedrepositories.NodeLogsRepository
}

type ReportsHandlerParams struct {
	fx.In
	Logger                    *zap.Logger
	ReportRepository          repositories.ReportRepository
	ApplicationLogsRepository sharedrepositories.ApplicationLogsRepository
	NodeLogsRepository        sharedrepositories.NodeLogsRepository
}

func NewReportsHandler(p ReportsHandlerParams) *ReportsHandler {
	return &ReportsHandler{
		logger:                    p.Logger,
		reportRepository:          p.ReportRepository,
		applicationLogsRepository: p.ApplicationLogsRepository,
		nodeLogsRepository:        p.NodeLogsRepository,
	}
}

func generateDummyReportsFromLogs(logs []sharedrepositories.NodeLogs) repositories.Report {

	hostReports := make([]*repositories.HostReport, 0, len(logs))
	for _, log := range logs {
		hostReports = append(hostReports, &repositories.HostReport{
			Host:         log.Host,
			CustomPrompt: log.Message,
		})
	}

	report := repositories.Report{
		Title:       "title",
		StartMs:     21,
		EndMs:       43,
		HostReports: hostReports}

	return report

}

func (h *ReportsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	defer r.Body.Close()

	ctx := context.Background()

	logs, err := h.nodeLogsRepository.GetAllLogs(ctx)
	if err != nil {
		h.logger.Error("Failed to get logs", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	report := generateDummyReportsFromLogs(logs)
	err = h.reportRepository.InsertReport(ctx, &report)
	if err != nil {
		h.logger.Error("Failed to generate report", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	reportJson, err := json.Marshal(report)
	if err != nil {
		h.logger.Error("Failed encode report into json", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(reportJson)

	w.WriteHeader(http.StatusOK)
}

func NewServeMux(reportsHandler *ReportsHandler) *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle("/reports", reportsHandler)
	return mux
}

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
			database.NewReportsDbMongoClient,
			fx.Annotate(
				repositories.NewMongoDbReportRepository,
				fx.As(new(repositories.ReportRepository)),
			),
			elasticsearch.NewElasticSearchLogsDbClient,
			fx.Annotate(
				sharedrepositories.NewElasticSearchApplicationLogsRepository,
				fx.As(new(sharedrepositories.ApplicationLogsRepository)),
			),
			fx.Annotate(
				sharedrepositories.NewElasticSearchNodeLogsRepository,
				fx.As(new(sharedrepositories.NodeLogsRepository)),
			),
			NewHTTPServer,
			NewReportsHandler,
			NewServeMux,
			zap.NewExample),
		fx.Invoke(func(*http.Server) {}),
	).Run()
}

package handlers

import (
	"context"
	"encoding/json"
	sharedrepositories "github.com/Magpie-Monitor/magpie-monitor/pkg/repositories"
	"github.com/Magpie-Monitor/magpie-monitor/services/reports/pkg/repositories"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type ReportsRouter struct {
	mux *http.ServeMux
}

func NewReportsRouter(reportsHandler *ReportsHandler) *ReportsRouter {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /", reportsHandler.Post)

	return &ReportsRouter{
		mux: mux,
	}
}

func (r *ReportsRouter) Pattern() string {
	return "/reports"
}

func (router *ReportsRouter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	router.mux.ServeHTTP(w, r)
}

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

type reportsPostParams struct {
	Cluster  string `json:"cluster"`
	FromDate int64  `json:"fromDate"`
	ToDate   int64  `json:"toDate"`
}

func (h *ReportsHandler) Post(w http.ResponseWriter, r *http.Request) {

	defer r.Body.Close()
	var params reportsPostParams
	ctx := context.Background()

	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		h.logger.Error("Failed to parse POST /reports params", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	logs, err := h.nodeLogsRepository.GetLogs(ctx,
		params.Cluster,
		time.Unix(params.FromDate, 0),
		time.Unix(params.ToDate, 0))

	if err != nil {
		h.logger.Error("Failed to get logs", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// TODO: Replace with interface for reports generation based on LLMs response
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
}

func generateDummyReportsFromLogs(logs []*sharedrepositories.NodeLogsDocument) repositories.Report {

	hostReports := make([]*repositories.HostReport, 0, len(logs))
	for _, log := range logs {
		hostReports = append(hostReports, &repositories.HostReport{
			Host:         log.Name,
			CustomPrompt: log.Kind,
		})
	}

	report := repositories.Report{
		Title:       "title",
		StartMs:     21,
		EndMs:       43,
		HostReports: hostReports}

	return report

}

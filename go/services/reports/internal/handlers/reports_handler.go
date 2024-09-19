package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	sharedrepositories "github.com/Magpie-Monitor/magpie-monitor/pkg/repositories"
	"github.com/Magpie-Monitor/magpie-monitor/services/reports/pkg/insights"
	"github.com/Magpie-Monitor/magpie-monitor/services/reports/pkg/repositories"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"math"
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
	logger                       *zap.Logger
	reportRepository             repositories.ReportRepository
	applicationLogsRepository    sharedrepositories.ApplicationLogsRepository
	nodeLogsRepository           sharedrepositories.NodeLogsRepository
	applicationInsightsGenerator insights.ApplicationInsightsGenerator
	nodeInsightsGenerator        insights.NodeInsightsGenerator
}

type ReportsHandlerParams struct {
	fx.In
	Logger                       *zap.Logger
	ReportRepository             repositories.ReportRepository
	ApplicationLogsRepository    sharedrepositories.ApplicationLogsRepository
	NodeLogsRepository           sharedrepositories.NodeLogsRepository
	ApplicationInsightsGenerator insights.ApplicationInsightsGenerator
	NodeInsightsGenerator        insights.NodeInsightsGenerator
}

func NewReportsHandler(p ReportsHandlerParams) *ReportsHandler {
	return &ReportsHandler{
		logger:                       p.Logger,
		reportRepository:             p.ReportRepository,
		applicationLogsRepository:    p.ApplicationLogsRepository,
		nodeLogsRepository:           p.NodeLogsRepository,
		applicationInsightsGenerator: p.ApplicationInsightsGenerator,
		nodeInsightsGenerator:        p.NodeInsightsGenerator,
	}
}

type reportsPostParams struct {
	Cluster   string `json:"cluster"`
	FromDate  int64  `json:"fromDate"`
	ToDate    int64  `json:"toDate"`
	MaxLength int64  `json:"maxLength"`
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

	applicationLogs, err := h.applicationLogsRepository.GetLogs(ctx,
		params.Cluster,
		time.Unix(0, params.FromDate),
		time.Unix(0, params.ToDate))

	h.logger.Sugar().Infof("Number of logs #%d", len(applicationLogs))

	if err != nil {
		h.logger.Error("Failed to get application logs", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	nodeLogs, err := h.nodeLogsRepository.GetLogs(ctx,
		params.Cluster,
		time.Unix(params.FromDate, 0),
		time.Unix(params.ToDate, 0))

	if err != nil {
		h.logger.Error("Failed to get node logs", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// filteredApplicationLogs := applicationLogs[0:params.MaxLength]
	filteredApplicationLogs := applicationLogs[0:int(math.Min(float64(params.MaxLength), float64(len(applicationLogs))))]
	filteredNodeLogs := nodeLogs[0:int(math.Min(float64(params.MaxLength), float64(len(nodeLogs))))]

	applicationInsights, err := h.applicationInsightsGenerator.OnDemandApplicationInsights(filteredApplicationLogs)
	if err != nil {
		h.logger.Error("Failed to generate application insights", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
	}

	nodeInsights, err := h.nodeInsightsGenerator.OnDemandNodeInsights(filteredNodeLogs)
	if err != nil {
		h.logger.Error("Failed to generate application insights", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
	}

	reportJson, err := json.Marshal(nodeInsights)
	if err != nil {
		h.logger.Error("Failed encode report into json", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(reportJson)

	applicationInsightsJson, err := json.Marshal(applicationInsights)
	if err != nil {
		h.logger.Error("Failed encode report into json", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(applicationInsightsJson)
}

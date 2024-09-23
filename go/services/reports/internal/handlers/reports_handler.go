package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/IBM/fp-go/array"
	sharedrepositories "github.com/Magpie-Monitor/magpie-monitor/pkg/repositories"
	"github.com/Magpie-Monitor/magpie-monitor/pkg/routing"
	"github.com/Magpie-Monitor/magpie-monitor/services/reports/pkg/insights"
	"github.com/Magpie-Monitor/magpie-monitor/services/reports/pkg/repositories"
	"github.com/gorilla/mux"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type ReportsRouter struct {
	mux *mux.Router
}

func NewReportsRouter(reportsHandler *ReportsHandler, rootRouter *mux.Router) *ReportsRouter {
	router := rootRouter.PathPrefix("/reports").Subrouter()
	router.Methods(http.MethodGet).Path("/{id}").HandlerFunc(reportsHandler.GetSingle)
	router.Methods(http.MethodGet).HandlerFunc(reportsHandler.GetAll)
	router.Methods(http.MethodPost).HandlerFunc(reportsHandler.Post)

	return &ReportsRouter{
		mux: rootRouter,
	}
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
	Cluster                  *string                                     `json:"cluster"`
	FromDate                 *int64                                      `json:"fromDate"`
	ToDate                   *int64                                      `json:"toDate"`
	ApplicationConfiguration []*insights.ApplicationInsightConfiguration `json:"applicationConfiguration"`
	NodeConfiguration        []*insights.NodeInsightConfiguration        `json:"nodeConfiguration"`
	MaxLength                *int                                        `json:"maxLength"`
}

func (h *ReportsHandler) GetSingle(w http.ResponseWriter, r *http.Request) {

	ctx := context.Background()
	vars := mux.Vars(r)

	id := vars["id"]
	report, repositoryErr := h.reportRepository.GetSingleReport(ctx, id)

	if repositoryErr != nil {
		switch repositoryErr.Kind() {
		case repositories.InternalError:
			w.WriteHeader(http.StatusInternalServerError)
		case repositories.ReportNotFound:
			w.WriteHeader(http.StatusNotFound)
		case repositories.InvalidReportId:
			w.WriteHeader(http.StatusBadRequest)
		}

		h.logger.Error(
			"Failed to fetch single report by id",
			zap.String("id", id),
			zap.Error(repositoryErr))

		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(repositoryErr.Error()))
	}

	encodedReport, err := json.Marshal(report)
	if err != nil {
		h.logger.Error("Failed to encode report", zap.String("id", id), zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write(encodedReport)
}

func (h *ReportsHandler) GetAll(w http.ResponseWriter, r *http.Request) {

	ctx := context.Background()

	query := r.URL.Query()

	cluster, isClusterSet := routing.LookupQueryParam(query, "cluster")
	fromDate, isFromDateSet := routing.LookupQueryParam(query, "fromDate")
	toDate, isToDateSet := routing.LookupQueryParam(query, "toDate")

	filterParams := repositories.FilterParams{}

	if isClusterSet {
		filterParams.Cluster = &cluster
	}

	if isFromDateSet {
		fromDateInt, err := strconv.ParseInt(fromDate, 10, 64)
		if err != nil {
			h.logger.Warn("Invalid fromDate query param", zap.Error(err))
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Invalid fromDate parameter"))
			return
		}
		filterParams.FromDate = &fromDateInt
	}

	if isToDateSet {
		toDateInt, err := strconv.ParseInt(toDate, 10, 64)
		if err != nil {
			h.logger.Warn("Invalid toDate query param", zap.Error(err))
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Invalid toDate parameter"))
			return
		}
		filterParams.ToDate = &toDateInt
	}

	reports, repositoryError := h.reportRepository.GetAllReports(ctx, filterParams)

	if repositoryError != nil {
		h.logger.Error("Failed to fetch all reports", zap.Error(repositoryError))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	encodedReports, err := json.Marshal(reports)
	if err != nil {
		h.logger.Error("Failed to encode all reports", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write(encodedReports)
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

	if params.FromDate == nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Missing fromDate parameter"))
		return
	}

	if params.ToDate == nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Missing toDate parameter"))
		return
	}

	if params.Cluster == nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Missing cluster parameter"))
		return
	}

	if params.MaxLength == nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Missing maxLength parameter"))
		return
	}

	report, err := h.generateCompleteOnDemandReport(ctx, reportGenerationFilters{
		Cluster:                  *params.Cluster,
		FromDate:                 *params.FromDate,
		ToDate:                   *params.ToDate,
		MaxLength:                *params.MaxLength,
		ApplicationConfiguration: params.ApplicationConfiguration,
		NodeConfiguration:        params.NodeConfiguration,
	})
	if err != nil {
		h.logger.Error("Failed to generate report", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	insertedReport, repositoryErr := h.reportRepository.InsertReport(ctx, report)
	if repositoryErr != nil {
		h.logger.Error("Failed to save a report", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	reportJson, err := json.Marshal(insertedReport)
	if err != nil {
		h.logger.Error("Failed encode report into json", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write(reportJson)
}

type reportGenerationFilters struct {
	Cluster                  string                                      `json:"cluster"`
	FromDate                 int64                                       `json:"fromDate"`
	ToDate                   int64                                       `json:"toDate"`
	ApplicationConfiguration []*insights.ApplicationInsightConfiguration `json:"applicationConfiguration"`
	NodeConfiguration        []*insights.NodeInsightConfiguration        `json:"nodeConfiguration"`
	MaxLength                int                                         `json:"maxLength"`
}

func (h *ReportsHandler) generateCompleteOnDemandReport(
	ctx context.Context,
	params reportGenerationFilters) (*repositories.Report, error) {

	fromDate := time.Unix(0, params.FromDate)
	toDate := time.Unix(0, params.ToDate)

	applicationReports, err := h.generateApplicationReports(
		ctx,
		params.Cluster,
		fromDate,
		toDate,
		params.ApplicationConfiguration,
		params.MaxLength,
	)
	if err != nil {
		h.logger.Error("Failed to get application reports", zap.Error(err))
		return nil, err
	}

	nodeReports, err := h.generateNodeReports(
		ctx,
		params.Cluster,
		fromDate,
		toDate,
		params.NodeConfiguration,
		params.MaxLength,
	)
	if err != nil {
		h.logger.Error("Failed to get node reports", zap.Error(err))
		return nil, err
	}

	report := repositories.Report{
		Status:        repositories.ReportState_Generated,
		Cluster:       params.Cluster,
		RequestedAtNs: time.Now().UnixNano(),
		//TODO: Generate the report entity first and assign the insights
		//later once the Batch API is implemented.
		GeneratedAtNs:           time.Now().UnixNano(),
		ScheduledGenerationAtMs: time.Now().UnixNano(),
		Title:                   h.getTitleForReport(params.Cluster, fromDate, toDate),
		FromDateNs:              params.FromDate,
		ToDateNs:                params.ToDate,
		NodeReports:             nodeReports,
		ApplicationReports:      applicationReports,
	}

	return &report, nil
}

func (h *ReportsHandler) getTitleForReport(cluster string, fromDate time.Time, toDate time.Time) string {
	return fmt.Sprintf("On-Demand report for %s (%s - %s)",
		cluster,
		fmt.Sprintf("%d.%d", fromDate.Month(), fromDate.Year()),
		fmt.Sprintf("%d.%d", toDate.Month(), toDate.Year()),
	)
}

func (h *ReportsHandler) generateApplicationReports(
	ctx context.Context,
	cluster string,
	fromDate time.Time,
	toDate time.Time,
	applicationConfiguration []*insights.ApplicationInsightConfiguration,
	maxLength int,
) ([]repositories.ApplicationReport, error) {

	applicationLogs, err := h.applicationLogsRepository.GetLogs(ctx,
		cluster,
		fromDate,
		toDate)

	if err != nil {
		h.logger.Error("Failed to get application logs", zap.Error(err))
		return nil, err
	}

	filteredApplicationLogs := applicationLogs[0:int(math.Min(float64(maxLength), float64(len(applicationLogs))))]

	applicationInsights, err := h.applicationInsightsGenerator.OnDemandApplicationInsights(
		filteredApplicationLogs,
		applicationConfiguration,
	)
	if err != nil {
		h.logger.Error("Failed to generate application insights", zap.Error(err))
		return nil, err
	}

	insightsByApplication := make(map[string][]insights.ApplicationInsightsWithMetadata)

	for _, insight := range applicationInsights {
		applicationName := insight.Metadata.ApplicationName
		insightsByApplication[applicationName] = append(insightsByApplication[applicationName], insight)
	}

	reports := make([]repositories.ApplicationReport, 0, len(insightsByApplication))

	configByApp := insights.MapApplicationNameToConfiguration(applicationConfiguration)

	for applicationName, insightsForApplication := range insightsByApplication {

		incidentsFromInsights := array.Map(func(insight insights.ApplicationInsightsWithMetadata) repositories.ApplicationIncident {
			return repositories.ApplicationIncident{
				Category:       insight.Insight.Category,
				Summary:        insight.Insight.Summary,
				Recommendation: insight.Insight.Recommendation,
				Timestamp:      insight.Metadata.Timestamp,
				PodName:        insight.Metadata.PodName,
				ContainerName:  insight.Metadata.ContainerName,
				Source:         insight.Metadata.Source,
			}
		})

		incidents := incidentsFromInsights(insightsForApplication)
		report := repositories.ApplicationReport{
			ApplicationName: applicationName,
			Incidents:       incidents,
		}

		config, ok := configByApp[applicationName]
		if ok {
			report.Precision = config.Precision
			report.CustomPrompt = config.CustomPrompt
		}

		reports = append(reports, report)
	}

	return reports, nil

}

func (h *ReportsHandler) generateNodeReports(
	ctx context.Context,
	cluster string,
	fromDate time.Time,
	toDate time.Time,
	nodeConfiguration []*insights.NodeInsightConfiguration,
	maxLength int,
) ([]repositories.NodeReport, error) {

	nodeLogs, err := h.nodeLogsRepository.GetLogs(ctx,
		cluster,
		fromDate,
		toDate)

	if err != nil {
		h.logger.Error("Failed to get node logs", zap.Error(err))
		return nil, err
	}

	filteredNodeLogs := nodeLogs[0:int(math.Min(float64(maxLength), float64(len(nodeLogs))))]

	nodeInsights, err := h.nodeInsightsGenerator.OnDemandNodeInsights(
		filteredNodeLogs,
		nodeConfiguration,
	)
	if err != nil {
		h.logger.Error("Failed to generate node insights", zap.Error(err))
		return nil, err
	}

	insightsByHostname := make(map[string][]insights.NodeInsightsWithMetadata)

	for _, insight := range nodeInsights {
		hostname := insight.Metadata.NodeName
		insightsByHostname[hostname] = append(insightsByHostname[hostname], insight)
	}

	reports := make([]repositories.NodeReport, 0, len(insightsByHostname))

	configByApp := insights.MapNodeNameToConfiguration(nodeConfiguration)

	for hostname, insightsForNode := range insightsByHostname {

		incidentsFromInsights := array.Map(func(insight insights.NodeInsightsWithMetadata) repositories.NodeIncident {
			return repositories.NodeIncident{
				Category:       insight.Insight.Category,
				Summary:        insight.Insight.Summary,
				Recommendation: insight.Insight.Recommendation,
				Timestamp:      insight.Metadata.Timestamp,
				Source:         insight.Metadata.Source,
			}
		})

		incidents := incidentsFromInsights(insightsForNode)
		report := repositories.NodeReport{
			Host:      hostname,
			Incidents: incidents,
		}

		config, ok := configByApp[hostname]
		if ok {
			report.Precision = config.Precision
			report.CustomPrompt = config.CustomPrompt
		}

		reports = append(reports, report)
	}

	return reports, nil

}

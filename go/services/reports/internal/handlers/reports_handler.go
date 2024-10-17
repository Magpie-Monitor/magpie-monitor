package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Magpie-Monitor/magpie-monitor/pkg/routing"
	"github.com/Magpie-Monitor/magpie-monitor/services/reports/internal/services"
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
	router.Methods(http.MethodPost).Path("/scheduled").HandlerFunc(reportsHandler.PostScheduled)
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
	logger         *zap.Logger
	reportsService *services.ReportsService
}

type ReportsHandlerParams struct {
	fx.In
	Logger         *zap.Logger
	ReportsService *services.ReportsService
}

func NewReportsHandler(p ReportsHandlerParams) *ReportsHandler {
	return &ReportsHandler{
		logger:         p.Logger,
		reportsService: p.ReportsService,
	}
}

type reportsPostParams struct {
	ClusterId                *string                                         `json:"clusterId"`
	SinceMs                  *int64                                          `json:"sinceMs"`
	ToMs                     *int64                                          `json:"toMs"`
	ApplicationConfiguration []*repositories.ApplicationInsightConfiguration `json:"applicationConfiguration"`
	NodeConfiguration        []*repositories.NodeInsightConfiguration        `json:"nodeConfiguration"`
	MaxLength                *int                                            `json:"maxLength"`
}

func (h *ReportsHandler) handleResponseHeaderFromRepositoryError(w http.ResponseWriter, err repositories.ReportRepositoryError) {
	switch err.Kind() {
	case repositories.InternalError:
		w.WriteHeader(http.StatusInternalServerError)
	case repositories.ReportNotFound:
		w.WriteHeader(http.StatusNotFound)
	case repositories.InvalidReportId:
		w.WriteHeader(http.StatusBadRequest)
	}
}

func (h *ReportsHandler) GetSingle(w http.ResponseWriter, r *http.Request) {

	ctx := context.Background()
	vars := mux.Vars(r)

	id := vars["id"]
	report, repositoryErr := h.reportsService.GetSingleReport(ctx, id)

	if repositoryErr != nil {

		h.handleResponseHeaderFromRepositoryError(w, *repositoryErr)
		h.logger.Error(
			"Failed to fetch single report by id",
			zap.String("id", id),
			zap.Error(repositoryErr))

		routing.WriteHttpError(w, repositoryErr.Error())
		return
	}

	if report.Status == repositories.ReportState_AwaitingGeneration {
		scheduledReport, err := h.reportsService.RetrieveScheduledReport(report.Id)
		if err != nil {
			h.logger.Error("Failed to retrieve scheduled report", zap.Error(err))
		} else {
			report = scheduledReport
		}
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

	clusterId, isClusterSet := routing.LookupQueryParam(query, "clusterId")
	sinceMs, isSinceMsSet := routing.LookupQueryParam(query, "sinceMs")
	toMs, isToMsSet := routing.LookupQueryParam(query, "toMs")

	filterParams := repositories.FilterParams{}

	if isClusterSet {
		filterParams.ClusterId = &clusterId
	}

	if isSinceMsSet {
		fromDateInt, err := strconv.ParseInt(sinceMs, 10, 64)
		if err != nil {
			h.logger.Warn("Invalid sinceMs query param", zap.Error(err))
			w.WriteHeader(http.StatusBadRequest)
			routing.WriteHttpError(w, "Invalid sinceMs parameter")
			return
		}

		filterParams.SinceMs = &fromDateInt
	}

	if isToMsSet {
		toDateInt, err := strconv.ParseInt(toMs, 10, 64)
		if err != nil {
			h.logger.Warn("Invalid toDate query param", zap.Error(err))
			w.WriteHeader(http.StatusBadRequest)
			routing.WriteHttpError(w, "Invalid toMs parameter")

			return
		}
		filterParams.ToMs = &toDateInt
	}

	reports, repositoryError := h.reportsService.GetAllReports(ctx, filterParams)

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
		routing.WriteHttpError(w, "Failed to parse POST /reports params")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if params.SinceMs == nil {
		w.WriteHeader(http.StatusBadRequest)
		routing.WriteHttpError(w, "Missing sinceMs parameter")
		return
	}

	if params.ToMs == nil {
		w.WriteHeader(http.StatusBadRequest)
		routing.WriteHttpError(w, "Missing toMs parameter")
		return
	}

	if params.ClusterId == nil {
		w.WriteHeader(http.StatusBadRequest)
		routing.WriteHttpError(w, "Missing clusterId parameter")
		return
	}

	if params.MaxLength == nil {
		w.WriteHeader(http.StatusBadRequest)
		routing.WriteHttpError(w, "Missing maxLength parameter")
		return
	}

	report, err := h.reportsService.GenerateAndSaveReport(ctx,
		services.ReportGenerationFilters{
			ClusterId:                *params.ClusterId,
			SinceMs:                  *params.SinceMs,
			ToMs:                     *params.ToMs,
			MaxLength:                *params.MaxLength,
			ApplicationConfiguration: params.ApplicationConfiguration,
			NodeConfiguration:        params.NodeConfiguration,
		})

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

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(reportJson)
}

func (h *ReportsHandler) PostScheduled(w http.ResponseWriter, r *http.Request) {

	defer r.Body.Close()
	var params reportsPostParams
	ctx := context.Background()

	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		h.logger.Error("Failed to parse POST /reports params", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		routing.WriteHttpError(w, fmt.Sprintf("Failed to parse POST /reports params"))
		return
	}

	if params.SinceMs == nil {
		w.WriteHeader(http.StatusBadRequest)
		routing.WriteHttpError(w, "Missing sinceMs parameter")
		return
	}

	if params.ToMs == nil {
		w.WriteHeader(http.StatusBadRequest)
		routing.WriteHttpError(w, "Missing toMs parameter")
		return
	}

	if params.ClusterId == nil {
		w.WriteHeader(http.StatusBadRequest)
		routing.WriteHttpError(w, "Missing clusterId parameter")
		return
	}

	if params.MaxLength == nil {
		w.WriteHeader(http.StatusBadRequest)
		routing.WriteHttpError(w, "Missing maxLength parameter")
		return
	}

	resp, err := h.reportsService.ScheduleReport(ctx,
		services.ReportGenerationFilters{
			ClusterId:                *params.ClusterId,
			SinceMs:                  *params.SinceMs,
			ToMs:                     *params.ToMs,
			MaxLength:                *params.MaxLength,
			ApplicationConfiguration: params.ApplicationConfiguration,
			NodeConfiguration:        params.NodeConfiguration,
		},
	)

	if err != nil {
		h.logger.Error("Failed to generate report", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		routing.WriteHttpError(w, "Internal server error")
		return
	}

	reportJson, err := json.Marshal(resp)
	if err != nil {
		h.logger.Error("Failed encode report into json", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		routing.WriteHttpError(w, "Internal server error")
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(reportJson)
}

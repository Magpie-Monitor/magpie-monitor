package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	messagebroker "github.com/Magpie-Monitor/magpie-monitor/pkg/message-broker"
	"github.com/Magpie-Monitor/magpie-monitor/pkg/routing"
	"github.com/Magpie-Monitor/magpie-monitor/services/reports/internal/brokers"
	"github.com/Magpie-Monitor/magpie-monitor/services/reports/internal/services"
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
	logger                    *zap.Logger
	reportsService            *services.ReportsService
	reportRequestedBroker     messagebroker.MessageBroker[brokers.ReportRequested]
	reportGeneratedBroker     messagebroker.MessageBroker[brokers.ReportGenerated]
	reportRequestFailedBroker messagebroker.MessageBroker[brokers.ReportRequestFailed]
}

type ReportsHandlerParams struct {
	fx.In
	Logger                    *zap.Logger
	ReportsService            *services.ReportsService
	ReportRequestedBroker     messagebroker.MessageBroker[brokers.ReportRequested]
	ReportGeneratedBroker     messagebroker.MessageBroker[brokers.ReportGenerated]
	ReportRequestFailedBroker messagebroker.MessageBroker[brokers.ReportRequestFailed]
}

func NewReportsHandler(p ReportsHandlerParams) *ReportsHandler {
	return &ReportsHandler{
		logger:                    p.Logger,
		reportsService:            p.ReportsService,
		reportRequestedBroker:     p.ReportRequestedBroker,
		reportRequestFailedBroker: p.ReportRequestFailedBroker,
		reportGeneratedBroker:     p.ReportGeneratedBroker,
	}
}

type reportsPostParams struct {
	ClusterId                *string                                     `json:"clusterId"`
	CorrelationId            *string                                     `json:"correlationId"`
	SinceMs                  *int64                                      `json:"sinceMs"`
	ToMs                     *int64                                      `json:"toMs"`
	ApplicationConfiguration []*insights.ApplicationInsightConfiguration `json:"applicationConfiguration"`
	NodeConfiguration        []*insights.NodeInsightConfiguration        `json:"nodeConfiguration"`
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

	if params.CorrelationId == nil {
		w.WriteHeader(http.StatusBadRequest)
		routing.WriteHttpError(w, "Missing correlationId parameter")
		return
	}

	report, err := h.reportsService.GenerateAndSaveReport(ctx,
		services.ReportGenerationFilters{
			ClusterId:                *params.ClusterId,
			CorrelationId:            *params.CorrelationId,
			SinceMs:                  *params.SinceMs,
			ToMs:                     *params.ToMs,
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

// TODO: Remove once the system is migrated to microservices
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

	resp, err := h.reportsService.ScheduleReport(ctx,
		services.ReportGenerationFilters{
			ClusterId:                *params.ClusterId,
			CorrelationId:            *params.CorrelationId,
			SinceMs:                  *params.SinceMs,
			ToMs:                     *params.ToMs,
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

// TODO: Remove everything above, once the system is migrated to microservices
func (h *ReportsHandler) ListenForReportRequests() {

	requests := make(chan brokers.ReportRequested)
	errChan := make(chan error)
	ctx := context.Background()

	go h.reportRequestedBroker.Subscribe(requests, errChan)
	for {
		select {
		case request := <-requests:
			err := h.ScheduleReport(ctx, request.CorrelationId, &request.ReportRequest)
			if err != nil {
				h.logger.Error("Failed to schedule a report", zap.Any("err", err), zap.Any("request", request))
				h.reportRequestFailedBroker.Publish(request.CorrelationId, *err)
			}

		case err := <-errChan:
			h.logger.Error("Recieved a malformed message", zap.Any("err", err))
		}

	}
}

func (h *ReportsHandler) ScheduleReport(ctx context.Context, correlationId string,
	reportRequest *brokers.ReportRequest) *brokers.ReportRequestFailed {

	if reportRequest.SinceMs == nil {
		return brokers.NewReportRequestFailedValidation(
			correlationId,
			"Missing sinceMs parameter",
		)
	}

	if reportRequest.ToMs == nil {
		return brokers.NewReportRequestFailedValidation(
			correlationId,
			"Missing toMs parameter",
		)
	}

	if reportRequest.ClusterId == nil {
		return brokers.NewReportRequestFailedValidation(
			correlationId,
			"Missing clusterId parameter",
		)
	}

	if reportRequest.ClusterId == nil {
		return brokers.NewReportRequestFailedValidation(
			correlationId,
			"Missing maxLength parameter",
		)
	}

	resp, err := h.reportsService.ScheduleReport(ctx,
		services.ReportGenerationFilters{
			ClusterId:                *reportRequest.ClusterId,
			CorrelationId:            correlationId,
			SinceMs:                  *reportRequest.SinceMs,
			ToMs:                     *reportRequest.ToMs,
			ApplicationConfiguration: reportRequest.ApplicationConfiguration,
			NodeConfiguration:        reportRequest.NodeConfiguration,
		},
	)

	if err != nil {
		h.logger.Error("Failed to generate report", zap.Error(err))
		return brokers.NewReportRequestFailedInternalError(
			correlationId,
			fmt.Sprintf("Failed to generate report %s", err.Error()),
		)
	}

	h.logger.Info("Scheduled report", zap.Any("report", resp))

	return nil
}

func (h *ReportsHandler) PollReports() {

	errChannel := make(chan *services.ReportGenerationError)
	reportsChannel := make(chan *repositories.Report)

	go h.reportsService.PollReports(context.Background(), reportsChannel, errChannel)

	for {
		select {
		case report := <-reportsChannel:
			h.reportGeneratedBroker.Publish(report.CorrelationId, brokers.NewReportGenerated(
				report,
			))

		case err := <-errChannel:
			h.reportRequestFailedBroker.Publish(err.Report.CorrelationId, *brokers.NewReportRequestFailedInternalError(
				err.Report.CorrelationId,
				err.Msg,
			))
		}
	}

}

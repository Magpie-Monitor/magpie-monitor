package handlers

import (
	"context"
	"fmt"
	messagebroker "github.com/Magpie-Monitor/magpie-monitor/pkg/message-broker"
	"github.com/Magpie-Monitor/magpie-monitor/services/reports/internal/brokers"
	"github.com/Magpie-Monitor/magpie-monitor/services/reports/internal/services"
	"github.com/Magpie-Monitor/magpie-monitor/services/reports/pkg/insights"
	"github.com/Magpie-Monitor/magpie-monitor/services/reports/pkg/repositories"
	"github.com/gorilla/mux"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"net/http"
)

type ReportsRouter struct {
	mux *mux.Router
}

func NewReportsRouter(reportsHandler *ReportsHandler, rootRouter *mux.Router) *ReportsRouter {

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

func (h *ReportsHandler) ListenForReportRequests() {

	requests := make(chan brokers.ReportRequested)
	errChan := make(chan error)
	ctx := context.Background()

	go h.reportRequestedBroker.Subscribe(ctx, requests, errChan)
	for {
		select {
		case request := <-requests:
			_, err := h.ScheduleReport(ctx, request.CorrelationId, &request.ReportRequest)
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
	reportRequest *brokers.ReportRequest) (*repositories.Report, *brokers.ReportRequestFailed) {

	if reportRequest.SinceMs == nil {
		return nil, brokers.NewReportRequestFailedValidation(
			correlationId,
			"Missing sinceMs parameter",
		)
	}

	if reportRequest.ToMs == nil {
		return nil, brokers.NewReportRequestFailedValidation(
			correlationId,
			"Missing toMs parameter",
		)
	}

	if reportRequest.ClusterId == nil {
		return nil, brokers.NewReportRequestFailedValidation(
			correlationId,
			"Missing clusterId parameter",
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
		return nil, brokers.NewReportRequestFailedInternalError(
			correlationId,
			fmt.Sprintf("Failed to generate report %s", err.Error()),
		)
	}

	h.logger.Info("Scheduled report", zap.Any("report", resp))

	return resp, nil
}

func (h *ReportsHandler) PollReports() {

	errChannel := make(chan *services.ReportGenerationError)
	reportsChannel := make(chan *repositories.Report)

	go h.reportsService.PollReportsPendingGeneration(context.Background(), reportsChannel, errChannel)
	go h.reportsService.PollReportsPendingIncidentMerge(context.Background(), reportsChannel, errChannel)

	for {
		select {
		case report := <-reportsChannel:
			if report.Status == repositories.ReportState_Generated {
				h.reportGeneratedBroker.Publish(report.CorrelationId, brokers.NewReportGenerated(
					report,
				))
			}

		case err := <-errChannel:
			h.reportRequestFailedBroker.Publish(err.Report.CorrelationId, *brokers.NewReportRequestFailedInternalError(
				err.Report.CorrelationId,
				err.Msg,
			))
		}
	}
}

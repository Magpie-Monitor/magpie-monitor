package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	sharedrepositories "github.com/Magpie-Monitor/magpie-monitor/pkg/repositories"
	"github.com/Magpie-Monitor/magpie-monitor/services/reports/pkg/openai"
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
	openAiClient              *openai.Client
}

type ReportsHandlerParams struct {
	fx.In
	Logger                    *zap.Logger
	ReportRepository          repositories.ReportRepository
	ApplicationLogsRepository sharedrepositories.ApplicationLogsRepository
	NodeLogsRepository        sharedrepositories.NodeLogsRepository
	OpenAiClient              *openai.Client
}

func NewReportsHandler(p ReportsHandlerParams) *ReportsHandler {
	return &ReportsHandler{
		logger:                    p.Logger,
		reportRepository:          p.ReportRepository,
		applicationLogsRepository: p.ApplicationLogsRepository,
		nodeLogsRepository:        p.NodeLogsRepository,
		openAiClient:              p.OpenAiClient,
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

	logs, err := h.applicationLogsRepository.GetLogs(ctx,
		params.Cluster,
		time.Unix(params.FromDate, 0),
		time.Unix(params.ToDate, 0))

	if err != nil {
		h.logger.Error("Failed to get logs", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	h.logger.Info("Logs", zap.Any("logs", logs))

	filteredLogs := logs[0:params.MaxLength]
	jsonFilteredLogs, err := json.Marshal(filteredLogs)
	if err != nil {
		h.logger.Error("Failed to encode application logs before sending to model", zap.Error(err))
	}

	h.logger.Sugar().Debugf("Sending logs to the model: %s", string(jsonFilteredLogs))

	openAiResponse, err := h.openAiClient.Complete([]*openai.Message{
		{
			Role: "system",
			Content: `You are a kubernetes cluster system administrator. 
			Given a list of logs from a Kubernetes cluster
			find logs which might suggest any kind of errors or issues. Try to give a possible reason, 
			category of an issue, urgency and possible resolution. Ignore logs which are only 
			informational and are not marked by warnings or errors. Don't provide intruduction. 
			Go straight into describing these logs. The only scenario in which you should include the 
			informationa logs is when they are in a unnatural frequency (based on the timestamp) which might suggest an error.
			Do not mention the information logs if you don't have unnatural amount of logs for a given event.
			As a response for explaination return reports incident, where every inconsitency is an incident`,
		},
		{
			Role: "user",
			Content: fmt.Sprintf(`These are logs from my cluster. 
			Please tell me if they might suggest any kind of issues:
			%s`, jsonFilteredLogs),
		},
	},
	)

	h.logger.Info("Got response from openai", zap.Any("response", openAiResponse))

	if err != nil {
		h.logger.Error("Failed to get reports from openai client", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// TODO: Replace with interface for reports generation based on LLMs response
	// report := generateDummyReportsFromLogs(logs)
	// err = h.reportRepository.InsertReport(ctx, &report)
	// if err != nil {
	// 	h.logger.Error("Failed to generate report", zap.Error(err))
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	return
	// }

	reportJson, err := json.Marshal(openAiResponse)
	if err != nil {
		h.logger.Error("Failed encode report into json", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(reportJson)
}

// func generateDummyReportsFromLogs(logs []*sharedrepositories.NodeLogsDocument) repositories.Report {
//
// 	hostReports := make([]*repositories.HostReport, 0, len(logs))
// 	for _, log := range logs {
// 		hostReports = append(hostReports, &repositories.HostReport{
// 			Host:         log.Name,
// 			CustomPrompt: log.Kind,
// 		})
// 	}
//
// 	report := repositories.Report{
// 		Title:       "title",
// 		StartMs:     21,
// 		EndMs:       43,
// 		HostReports: hostReports}
//
// 	return report
//
// }

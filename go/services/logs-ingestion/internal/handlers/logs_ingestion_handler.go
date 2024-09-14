package handlers

import (
	// "context"
	// "encoding/json"
	"github.com/Magpie-Monitor/magpie-monitor/pkg/repositories"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"net/http"
)

type LogsIngestionRouter struct {
	mux *http.ServeMux
}

func NewLogsIngestionRouter(logsIngestionHanlder *LogsIngestionHandler) *LogsIngestionRouter {
	mux := http.NewServeMux()
	//TODO: Remove once communication with queue is established
	// mux.HandleFunc("POST /", logsIngestionHanlder.Get)

	return &LogsIngestionRouter{
		mux: mux,
	}

}

func (r *LogsIngestionRouter) Pattern() string {
	return "/logs"
}

func (router *LogsIngestionRouter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	router.mux.ServeHTTP(w, r)
}

type LogsIngestionHandler struct {
	logger                    *zap.Logger
	applicationLogsRepository repositories.ApplicationLogsRepository
	nodeLogsRepository        repositories.NodeLogsRepository
}

type LogIngestionHandlerParams struct {
	fx.In
	ApplicationLogsRepository repositories.ApplicationLogsRepository
	NodesLogsRepository       repositories.NodeLogsRepository
	Logger                    *zap.Logger
}

func NewLogsIngestionHandler(p LogIngestionHandlerParams) *LogsIngestionHandler {
	return &LogsIngestionHandler{
		logger:                    p.Logger,
		applicationLogsRepository: p.ApplicationLogsRepository,
		nodeLogsRepository:        p.NodesLogsRepository,
	}
}

//
// func (h *LogsIngestionHandler) Get(w http.ResponseWriter, r *http.Request) {
//
// 	defer r.Body.Close()
//
// 	var logs repositories.NodeLogs
//
// 	ctx := context.Background()
//
// 	err := json.NewDecoder(r.Body).Decode(&logs)
// 	if err != nil {
// 		w.WriteHeader(http.StatusBadRequest)
// 		return
// 	}
//
// 	err = h.nodeLogsRepository.InsertLogs(ctx, logs)
// 	if err != nil {
// 		h.logger.Error("Failed to index logs", zap.Error(err))
// 		w.WriteHeader(http.StatusInternalServerError)
// 		return
// 	}
//
// 	w.WriteHeader(http.StatusOK)
// }

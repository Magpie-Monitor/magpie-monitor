package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Magpie-Monitor/magpie-monitor/pkg/elasticsearch"
	"github.com/Magpie-Monitor/magpie-monitor/pkg/repositories"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"net"
	"net/http"
	"os"
)

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

func (h *LogsIngestionHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	defer r.Body.Close()

	var logs repositories.NodeLogs

	ctx := context.Background()

	err := json.NewDecoder(r.Body).Decode(&logs)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = h.nodeLogsRepository.CreateIndex(ctx, "test_index")

	if err != nil {
		h.logger.Error("Failed to create index", zap.Error(err))
		// w.WriteHeader(http.StatusInternalServerError)
		// return
	}

	err = h.nodeLogsRepository.InsertLogs(ctx, logs)
	if err != nil {
		h.logger.Error("Failed to index document", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func NewServeMux(logIngestion *LogsIngestionHandler) *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle("/logs", logIngestion)
	return mux
}

func NewHTTPServer(lc fx.Lifecycle, mux *http.ServeMux, log *zap.Logger) *http.Server {
	port := os.Getenv("LOGS_INGESTION_SERVICE_HTTP_PORT")

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
			elasticsearch.NewElasticSearchLogsDbClient,
			fx.Annotate(
				repositories.NewElasticSearchNodeLogsRepository,
				fx.As(new(repositories.NodeLogsRepository)),
			),
			fx.Annotate(
				repositories.NewElasticSearchApplicationLogsRepository,
				fx.As(new(repositories.ApplicationLogsRepository)),
			),
			NewLogsIngestionHandler,
			NewHTTPServer,
			NewServeMux,
			zap.NewExample),
		fx.Invoke(func(*http.Server) {}),
	).Run()
}

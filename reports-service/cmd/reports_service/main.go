package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"

	"github.com/Magpie-Monitor/magpie-monitor/reports-service/internal/database"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type LogsIngestionHandler struct {
	log     *zap.Logger
	mongoDb *mongo.Client
}

type LogDto struct {
	Source string `bson:"source"`
	Data   string `bson:"data"`
}

func NewLogsIngestionHandler(log *zap.Logger, mongoDb *mongo.Client) *LogsIngestionHandler {
	return &LogsIngestionHandler{log: log, mongoDb: mongoDb}
}

func (h *LogsIngestionHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	defer r.Body.Close()

	var logBody LogDto

	err := json.NewDecoder(r.Body).Decode(&logBody)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	coll := h.mongoDb.Database("logs").Collection("logs")
	_, err = coll.InsertOne(context.TODO(), logBody)
	if err != nil {
		h.log.Error("Failed to insert log into the db", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func NewServeMux(logIngestion *LogsIngestionHandler) *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle("/log-ingestion", logIngestion)
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
			database.NewMongoDbClient,
			NewHTTPServer,
			NewServeMux,
			NewLogsIngestionHandler,
			zap.NewExample),
		fx.Invoke(func(*http.Server) {}),
	).Run()
}

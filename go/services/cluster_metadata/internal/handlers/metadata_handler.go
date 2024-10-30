package handlers

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/Magpie-Monitor/magpie-monitor/services/cluster_metadata/pkg/repositories"
	"github.com/Magpie-Monitor/magpie-monitor/services/cluster_metadata/pkg/services"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

func NewMetadataRouter(metadataHandler *MetadataHandler, rootRouter *mux.Router) *MetadataRouter {
	router := rootRouter.PathPrefix("/metadata").Subrouter()
	router.Methods(http.MethodPost).Path("/clusters").HandlerFunc(metadataHandler.InsertApplicationMetadata)
	router.Methods(http.MethodPost).Path("/nodes").HandlerFunc(metadataHandler.InsertNodeMetadata)
	router.Methods(http.MethodGet).Path("/healthz").HandlerFunc(metadataHandler.Healthz)

	return &MetadataRouter{
		mux: router,
	}
}

func NewMetadataHandler(log *zap.Logger, service *services.MetadataService) *MetadataHandler {
	clientSecret, present := os.LookupEnv("CLIENT_SECRET")
	if !present {
		panic("No value provided for CLIENT_SECRET env variable")
	}

	return &MetadataHandler{log: log, metadataService: service, clientSecret: clientSecret}
}

type MetadataRouter struct {
	mux *mux.Router
}

func (m *MetadataRouter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m.mux.ServeHTTP(w, r)
}

type ErrorMessage struct {
	TimestampMillis int64  `json:"timestampMillis"`
	Error           string `json:"error"`
}

type MetadataHandler struct {
	log             *zap.Logger
	metadataService *services.MetadataService
	clientSecret    string
}

func (h *MetadataHandler) writeError(w *http.ResponseWriter, msg string, status int) {
	(*w).Header().Set("Content-Type", "application/json")
	(*w).WriteHeader(status)
	json.NewEncoder(*w).Encode(&ErrorMessage{
		TimestampMillis: time.Now().UnixMilli(),
		Error:           msg,
	})
}

func (h *MetadataHandler) InsertApplicationMetadata(w http.ResponseWriter, r *http.Request) {
	m2m := r.Header.Get("X-Client-Secret")
	if m2m != h.clientSecret {
		h.log.Error("Invalid client secret")
		h.writeError(&w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	defer r.Body.Close()

	var metadata repositories.ClusterState

	err := json.NewDecoder(r.Body).Decode(&metadata)
	if err != nil {
		h.log.Error("Error parsing inserted cluster metadata:", zap.Error(err))
		h.writeError(&w, "Error parsing cluster metadata json", http.StatusBadRequest)
		return
	}

	err = h.metadataService.InsertApplicationMetadata(metadata)
	if err != nil {
		h.log.Error("Error inserting cluster metadata:", zap.Error(err))
		h.writeError(&w, "Error inserting cluster metadata", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *MetadataHandler) InsertNodeMetadata(w http.ResponseWriter, r *http.Request) {
	m2m := r.Header.Get("X-Client-Secret")
	if m2m != h.clientSecret {
		h.log.Error("Invalid client secret")
		h.writeError(&w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	defer r.Body.Close()

	var metadata repositories.NodeState

	err := json.NewDecoder(r.Body).Decode(&metadata)
	if err != nil {
		h.log.Error("Error parsing inserted node metadata:", zap.Error(err))
		h.writeError(&w, "Error parsing node metadata json", http.StatusBadRequest)
		return
	}

	err = h.metadataService.InsertNodeMetadata(metadata)
	if err != nil {
		h.log.Error("Error inserting node metadata:", zap.Error(err))
		h.writeError(&w, "Error inserting node metadata", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *MetadataHandler) Healthz(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

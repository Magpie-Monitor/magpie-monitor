package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Magpie-Monitor/magpie-monitor/services/cluster_metadata/internal/entity"
	"github.com/Magpie-Monitor/magpie-monitor/services/cluster_metadata/pkg/services"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

func NewMetadataRouter(metadataHandler *MetadataHandler, rootRouter *mux.Router) *MetadataRouter {
	router := rootRouter.PathPrefix("/api/v1/metadata").Subrouter()
	router.Methods(http.MethodGet).Path("/clusters/{clusterId}/applications").HandlerFunc(metadataHandler.GetClusterMetadataForTimerange)
	router.Methods(http.MethodGet).Path("/clusters/{clusterId}/nodes").HandlerFunc(metadataHandler.GetNodeMetadataForTimerange)
	router.Methods(http.MethodPost).Path("/cluster").HandlerFunc(metadataHandler.InsertClusterMetadata)
	router.Methods(http.MethodPost).Path("/nodes").HandlerFunc(metadataHandler.InsertNodeMetadata)

	return &MetadataRouter{
		mux: router,
	}
}

func NewMetadataHandler(log *zap.Logger, service *services.MetadataService) *MetadataHandler {
	return &MetadataHandler{log: log, metadataService: service}
}

type MetadataRouter struct {
	mux *mux.Router
}

func (m *MetadataRouter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m.mux.ServeHTTP(w, r)
}

type MetadataHandler struct {
	log             *zap.Logger
	metadataService *services.MetadataService
}

func (h *MetadataHandler) GetClusterMetadataForTimerange(w http.ResponseWriter, r *http.Request) {
	sinceMillis, err := strconv.Atoi(r.URL.Query().Get("sinceMillis"))
	if err != nil {
		http.Error(w, "Invalid parameter value for \"sinceMillis\", it has to be an integer", http.StatusBadRequest)
	}

	toMillis, err := strconv.Atoi(r.URL.Query().Get("toMillis"))
	if err != nil {
		http.Error(w, "Invalid parameter value for \"toMillis\", it has to be an integer", http.StatusBadRequest)
	}

	clusterId := mux.Vars(r)["clusterId"]
	metadata, err := h.metadataService.GetClusterMetadataForTimerange(clusterId, sinceMillis, toMillis)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(metadata)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *MetadataHandler) GetNodeMetadataForTimerange(w http.ResponseWriter, r *http.Request) {
	sinceMillis, err := strconv.Atoi(r.URL.Query().Get("sinceMillis"))
	if err != nil {
		http.Error(w, "Invalid parameter value for \"sinceMillis\", it has to be an integer", http.StatusBadRequest)
	}

	toMillis, err := strconv.Atoi(r.URL.Query().Get("toMillis"))
	if err != nil {
		http.Error(w, "Invalid parameter value for \"toMillis\", it has to be an integer", http.StatusBadRequest)
	}

	clusterId := mux.Vars(r)["clusterId"]

	metadata, err := h.metadataService.GetNodeMetadataForTimerange(clusterId, sinceMillis, toMillis)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	err = json.NewEncoder(w).Encode(metadata)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

// TODO - add m2m
func (h *MetadataHandler) InsertClusterMetadata(w http.ResponseWriter, r *http.Request) {
	var metadata entity.ClusterState

	err := json.NewDecoder(r.Body).Decode(&metadata)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.metadataService.InsertClusterMetadata(metadata)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
}

// TODO - add m2m
func (h *MetadataHandler) InsertNodeMetadata(w http.ResponseWriter, r *http.Request) {
	var metadata entity.NodeState

	err := json.NewDecoder(r.Body).Decode(&metadata)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.metadataService.InsertNodeMetadata(metadata)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
}

package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/Magpie-Monitor/magpie-monitor/services/cluster_metadata/pkg/repositories"
	"github.com/Magpie-Monitor/magpie-monitor/services/cluster_metadata/pkg/services"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

func NewMetadataRouter(metadataHandler *MetadataHandler, rootRouter *mux.Router) *MetadataRouter {
	router := rootRouter.PathPrefix("/metadata").Subrouter()
	router.Methods(http.MethodGet).Path("/clusters/{clusterId}/applications").HandlerFunc(metadataHandler.GetClusterMetadataForTimerange)
	router.Methods(http.MethodGet).Path("/clusters/{clusterId}/nodes").HandlerFunc(metadataHandler.GetNodeMetadataForTimerange)
	router.Methods(http.MethodPost).Path("/cluster").HandlerFunc(metadataHandler.InsertClusterMetadata)
	router.Methods(http.MethodPost).Path("/nodes").HandlerFunc(metadataHandler.InsertNodeMetadata)
	router.Methods(http.MethodGet).Path("/healthz").HandlerFunc(metadataHandler.Healthz)

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

type MetadataResponse[T any] struct {
	Response []T `json:"response"`
}

func (h *MetadataHandler) GetClusterMetadataForTimerange(w http.ResponseWriter, r *http.Request) {
	sinceMillis, err := strconv.Atoi(r.URL.Query().Get("sinceMillis"))
	if err != nil {
		h.log.Error("Error parsing sinceMillis:", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	toMillis, err := strconv.Atoi(r.URL.Query().Get("toMillis"))
	if err != nil {
		h.log.Error("Error parsing toMillis:", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	clusterId := mux.Vars(r)["clusterId"]

	metadata, err := h.metadataService.GetClusterMetadataForTimerange(clusterId, sinceMillis, toMillis)
	if err != nil {
		h.log.Error("Error reading cluster metadata:", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
	}

	response := MetadataResponse[repositories.ClusterState]{Response: metadata}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(&response)
	if err != nil {
		h.log.Error("Error parsing cluster metadata", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (h *MetadataHandler) GetNodeMetadataForTimerange(w http.ResponseWriter, r *http.Request) {
	sinceMillis, err := strconv.Atoi(r.URL.Query().Get("sinceMillis"))
	if err != nil {
		h.log.Error("Error parsing sinceMillis:", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	toMillis, err := strconv.Atoi(r.URL.Query().Get("toMillis"))
	if err != nil {
		h.log.Error("Error parsing toMillis:", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	clusterId := mux.Vars(r)["clusterId"]

	metadata, err := h.metadataService.GetNodeMetadataForTimerange(clusterId, sinceMillis, toMillis)
	if err != nil {
		h.log.Error("Error reading node metadata:", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response := MetadataResponse[repositories.NodeState]{Response: metadata}
	log.Println(response)

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(&response)
	if err != nil {
		h.log.Error("Error parsing node metadata:", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// TODO - add m2m
func (h *MetadataHandler) InsertClusterMetadata(w http.ResponseWriter, r *http.Request) {
	var metadata repositories.ClusterState

	err := json.NewDecoder(r.Body).Decode(&metadata)
	if err != nil {
		h.log.Error("Error reading cluster metadata:", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = h.metadataService.InsertClusterMetadata(metadata)
	if err != nil {
		h.log.Error("Error inserting cluster metadata:", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// TODO - add m2m
func (h *MetadataHandler) InsertNodeMetadata(w http.ResponseWriter, r *http.Request) {
	var metadata repositories.NodeState

	err := json.NewDecoder(r.Body).Decode(&metadata)
	if err != nil {
		h.log.Error("Error reading node metadata:", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = h.metadataService.InsertNodeMetadata(metadata)
	if err != nil {
		h.log.Error("Error inserting node metadata:", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *MetadataHandler) Healthz(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

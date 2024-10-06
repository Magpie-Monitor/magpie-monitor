package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Magpie-Monitor/magpie-monitor/services/cluster_metadata/internal/entity"
	"github.com/Magpie-Monitor/magpie-monitor/services/cluster_metadata/pkg/services"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type MetadataRouter struct {
	mux *mux.Router
}

func NewMetadataRouter(metadataHandler MetadataHandler, rootRouter *mux.Router) *MetadataRouter {
	router := rootRouter.PathPrefix("/metadata").Subrouter()
	router.Methods(http.MethodGet).Path("/cluster").HandlerFunc(metadataHandler.GetClusterMetadataForTimerange)
	router.Methods(http.MethodGet).Path("/nodes").HandlerFunc(metadataHandler.GetNodeMetadataForTimerange)
	router.Methods(http.MethodPost).Path("/cluster").HandlerFunc(metadataHandler.InsertClusterMetadata)
	router.Methods(http.MethodGet).Path("/nodes").HandlerFunc(metadataHandler.InsertNodeMetadata)

	return &MetadataRouter{
		mux: router,
	}
}

func (m *MetadataRouter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m.mux.ServeHTTP(w, r)
}

type MetadataHandler struct {
	log             *zap.Logger
	metadataService services.MetadataService
}

func NewMetadataHandler(log *zap.Logger, service services.MetadataService) MetadataHandler {
	return MetadataHandler{log: log, metadataService: service}
}

// TODO - change to URL params
type ClusterMetadataRequest struct {
	ClusterName string `json:"clusterName"`
	SinceMillis int64  `json:"sinceMillis"`
	ToMillis    int64  `json:"toMillis"`
}

func (h *MetadataHandler) GetClusterMetadataForTimerange(w http.ResponseWriter, r *http.Request) {
	var req ClusterMetadataRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	metadata, err := h.metadataService.GetClusterMetadataForTimerange(req.ClusterName, req.SinceMillis, req.ToMillis)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(metadata)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// TODO - change to URL
type NodeMetadataRequest struct {
	NodeName    string `json:"nodeName"`
	SinceMillis int64  `json:"sinceMillis"`
	ToMillis    int64  `json:"toMillis"`
}

func (h *MetadataHandler) GetNodeMetadataForTimerange(w http.ResponseWriter, r *http.Request) {
	var req NodeMetadataRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	metadata, err := h.metadataService.GetNodeMetadataForTimerange(req.NodeName, req.SinceMillis, req.ToMillis)
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

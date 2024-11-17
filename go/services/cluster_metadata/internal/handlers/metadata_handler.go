package handlers

import (
	"net/http"

	"github.com/Magpie-Monitor/magpie-monitor/services/cluster_metadata/pkg/services"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

func NewMetadataRouter(metadataHandler *MetadataHandler, rootRouter *mux.Router) *MetadataRouter {
	router := rootRouter.PathPrefix("/metadata").Subrouter()
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

func (h *MetadataHandler) Healthz(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/Magpie-Monitor/magpie-monitor/services/reports/internal/services"
	"github.com/Magpie-Monitor/magpie-monitor/services/reports/pkg/repositories"
	"github.com/gorilla/mux"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type IncidentsRouter struct {
	mux *mux.Router
}

func NewIncidentsRouter(incidentsHanlder *IncidentsHandler, rootRouter *mux.Router) *IncidentsRouter {
	nodeIncidentsRouter := rootRouter.PathPrefix("/node-incidents").Subrouter()
	nodeIncidentsRouter.Methods(http.MethodGet).Path("/{id}").HandlerFunc(incidentsHanlder.GetSingleNodeIncident)

	applicationIncidentsRouter := rootRouter.PathPrefix("/application-incidents").Subrouter()
	applicationIncidentsRouter.Methods(http.MethodGet).Path("/{id}").HandlerFunc(incidentsHanlder.GetSingleApplicationIncident)

	return &IncidentsRouter{
		mux: rootRouter,
	}
}

func (router *IncidentsRouter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	router.mux.ServeHTTP(w, r)
}

type IncidentsHandler struct {
	logger                      *zap.Logger
	applicationIncidentsService *services.IncidentsService[repositories.ApplicationIncident]
	nodeIncidentsService        *services.IncidentsService[repositories.NodeIncident]
}

type IncidentsHandlerParams struct {
	fx.In
	Logger                      *zap.Logger
	ApplicationIncidentsService *services.IncidentsService[repositories.ApplicationIncident]
	NodeIncidentsService        *services.IncidentsService[repositories.NodeIncident]
}

func NewIncidentHandler(params IncidentsHandlerParams) *IncidentsHandler {
	return &IncidentsHandler{
		logger:                      params.Logger,
		applicationIncidentsService: params.ApplicationIncidentsService,
		nodeIncidentsService:        params.NodeIncidentsService,
	}
}

func (h *IncidentsHandler) handleResponseHeaderFromRepositoryError(w http.ResponseWriter, err repositories.IncidentRepositoryError) {
	switch err.Kind() {
	case repositories.IncidentInternalError:
		w.WriteHeader(http.StatusInternalServerError)
	case repositories.IncidentNotFound:
		w.WriteHeader(http.StatusNotFound)
	case repositories.InvalidIncidentId:
		w.WriteHeader(http.StatusBadRequest)
	}
}

func NewIncidentsHandler(p IncidentsHandlerParams) *IncidentsHandler {
	return &IncidentsHandler{
		logger:                      p.Logger,
		applicationIncidentsService: p.ApplicationIncidentsService,
		nodeIncidentsService:        p.NodeIncidentsService,
	}
}

func (h *IncidentsHandler) GetSingleNodeIncident(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	vars := mux.Vars(r)

	w.Header().Add("Content-Type", "application/json")
	id := vars["id"]
	incident, err := h.nodeIncidentsService.GetSingle(ctx, id)

	if err != nil {
		h.handleResponseHeaderFromRepositoryError(w, *err)
		return
	}

	encodedIncident, encErr := json.Marshal(incident)
	if encErr != nil {
		h.logger.Error("Failed to encode an incident", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, encErr = w.Write(encodedIncident)
	if encErr != nil {
		h.logger.Error("Failed to write an incident to the http response", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *IncidentsHandler) GetSingleApplicationIncident(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	vars := mux.Vars(r)

	id := vars["id"]
	incident, err := h.applicationIncidentsService.GetSingle(ctx, id)
	w.Header().Add("Content-Type", "application/json")

	if err != nil {
		h.handleResponseHeaderFromRepositoryError(w, *err)
		return
	}

	encodedIncident, encErr := json.Marshal(incident)
	if encErr != nil {
		h.logger.Error("Failed to encode an incident", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, encErr = w.Write(encodedIncident)
	if encErr != nil {
		h.logger.Error("Failed to write an incident to the http response", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

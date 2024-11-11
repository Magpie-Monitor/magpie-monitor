package handlers

import (
	"github.com/Magpie-Monitor/magpie-monitor/services/reports/internal/services"
	"github.com/Magpie-Monitor/magpie-monitor/services/reports/pkg/repositories"
	"github.com/gorilla/mux"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"net/http"
)

type IncidentsRouter struct {
	mux *mux.Router
}

func NewIncidentsRouter(incidentsHanlder *IncidentsHandler, rootRouter *mux.Router) *IncidentsRouter {

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

func NewIncidentsHandler(p IncidentsHandlerParams) *IncidentsHandler {
	return &IncidentsHandler{
		logger:                      p.Logger,
		applicationIncidentsService: p.ApplicationIncidentsService,
		nodeIncidentsService:        p.NodeIncidentsService,
	}
}

package services

import (
	"context"
	"github.com/Magpie-Monitor/magpie-monitor/services/reports/pkg/repositories"
	"go.uber.org/zap"
)

type IncidentsService[T any] struct {
	logger             *zap.Logger
	incidentRepository repositories.IncidentRepository[T]
}

func (s *IncidentsService[T]) GetSingle(ctx context.Context, id string) (*T, *repositories.IncidentRepositoryError) {
	incident, err := s.incidentRepository.GetIncident(ctx, id)
	if err != nil {
		s.logger.Error("Failed to get an incident", zap.Error(err))
		return nil, err
	}

	return incident, nil
}

func NewNodeIncidentsService(logger *zap.Logger, incidentRepository repositories.IncidentRepository[repositories.NodeIncident]) *IncidentsService[repositories.NodeIncident] {
	return &IncidentsService[repositories.NodeIncident]{
		logger:             logger,
		incidentRepository: incidentRepository,
	}
}

func NewApplicationIncidentsService(logger *zap.Logger, incidentRepository repositories.IncidentRepository[repositories.ApplicationIncident]) *IncidentsService[repositories.ApplicationIncident] {
	return &IncidentsService[repositories.ApplicationIncident]{
		logger:             logger,
		incidentRepository: incidentRepository,
	}
}

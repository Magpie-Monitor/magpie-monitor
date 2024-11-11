package repositories

import (
	"context"
	"fmt"

	"github.com/IBM/fp-go/array"
	"github.com/Magpie-Monitor/magpie-monitor/pkg/repositories"
	"github.com/Magpie-Monitor/magpie-monitor/services/reports/pkg/insights"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
)

type IncidentRepositoryErrorKind string

type Incident interface {
	GetId() string
	GetTitle() string
	GetRecommendation() string
	GetSummary() string
	GetUrgency() insights.Urgency
	GetCategory() string
}

const (
	IncidentNotFound      IncidentRepositoryErrorKind = "INCIDENT_NOT_FOUND"
	InvalidIncidentId     IncidentRepositoryErrorKind = "INVALID_INCIDENT_ID"
	IncidentInternalError IncidentRepositoryErrorKind = "INTERNAL_ERROR"
)

type IncidentRepositoryError struct {
	msg  string
	kind IncidentRepositoryErrorKind
}

func (e *IncidentRepositoryError) Error() string {
	return e.msg
}

func (e *IncidentRepositoryError) Kind() IncidentRepositoryErrorKind {
	return e.kind
}

func NewIncidentNotFoundError(err error) *IncidentRepositoryError {
	return &IncidentRepositoryError{
		msg:  fmt.Sprintf(" incident does not exists: %s", err),
		kind: IncidentNotFound,
	}
}

func NewInvalidIncidentIdError(err error) *IncidentRepositoryError {
	return &IncidentRepositoryError{
		msg:  fmt.Sprintf("Invalid application incident id: %s", err),
		kind: InvalidIncidentId,
	}
}

func NewIncidentInternalError(err error) *IncidentRepositoryError {
	return &IncidentRepositoryError{
		msg:  err.Error(),
		kind: IncidentInternalError,
	}
}

type IncidentRepository[T any] interface {
	InsertIncidents(ctx context.Context, incidents []*T) ([]string, *IncidentRepositoryError)
	GetIncident(ctx context.Context, id string) (*T, *IncidentRepositoryError)
	GetIncidentsByIds(ctx context.Context, ids []string) ([]*T, *IncidentRepositoryError)
}

type MongoDbIncidentRepository[T any] struct {
	mongoDbCollection *repositories.MongoDbCollection[*T]
	logger            *zap.Logger
}

func (r *MongoDbIncidentRepository[T]) InsertIncidents(ctx context.Context, incidents []*T) ([]string, *IncidentRepositoryError) {

	documents := make([]interface{}, 0, len(incidents))
	for _, incident := range incidents {
		documents = append(documents, incident)
	}

	ids, err := r.mongoDbCollection.InsertDocuments(documents)

	if err != nil {
		r.logger.Error("Failed to insert incidents", zap.Error(err))
		return nil, NewIncidentInternalError(err)
	}

	createdIds := array.Map(func(objectId primitive.ObjectID) string {
		return objectId.Hex()
	})(ids)

	return createdIds, nil
}

func (r *MongoDbIncidentRepository[T]) GetIncident(ctx context.Context, id string) (*T, *IncidentRepositoryError) {

	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		r.logger.Error("Failed to parse application incident id", zap.Error(err))
		return nil, NewInvalidIncidentIdError(err)
	}

	incident, err := r.mongoDbCollection.GetDocument(primitive.D{{Key: "_id", Value: objectId}}, primitive.D{})

	if err != nil {
		r.logger.Error("Failed to decode node incident from mongodb", zap.Error(err))
		return nil, NewIncidentNotFoundError(err)
	}

	return incident, nil
}

func (r *MongoDbIncidentRepository[T]) GetIncidentsByIds(ctx context.Context, ids []string) ([]*T, *IncidentRepositoryError) {

	idObjects := array.Map(func(id string) primitive.ObjectID {
		idObj, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			r.logger.Error("Failed to parse application incident ids", zap.Error(err))
		}
		return idObj
	})(ids)

	incidents, err := r.mongoDbCollection.GetDocuments(bson.D{{Key: "_id", Value: bson.M{"$in": idObjects}}}, primitive.D{})
	if err != nil {
		r.logger.Error("Failed to find application incidents by ids", zap.Error(err))
		return nil, NewIncidentInternalError(err)
	}

	return incidents, nil
}

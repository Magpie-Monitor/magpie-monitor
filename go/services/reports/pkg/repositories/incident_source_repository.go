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

type IncidentSourceRepositoryErrorKind string

type IncidentSource interface {
	GetId() string
	GetTitle() string
	GetRecommendation() string
	GetSummary() string
	GetUrgency() insights.Urgency
	GetCategory() string
}

const (
	IncidentSourceNotFound      IncidentSourceRepositoryErrorKind = "INCIDENT_NOT_FOUND"
	InvalidIncidentSourceId     IncidentSourceRepositoryErrorKind = "INVALID_INCIDENT_ID"
	IncidentSourceInternalError IncidentSourceRepositoryErrorKind = "INTERNAL_ERROR"
)

type IncidentSourceRepositoryError struct {
	msg  string
	kind IncidentSourceRepositoryErrorKind
}

func (e *IncidentSourceRepositoryError) Error() string {
	return e.msg
}

func (e *IncidentSourceRepositoryError) Kind() IncidentSourceRepositoryErrorKind {
	return e.kind
}

func NewIncidentSourceNotFoundError(err error) *IncidentSourceRepositoryError {
	return &IncidentSourceRepositoryError{
		msg:  fmt.Sprintf(" incident does not exists: %s", err),
		kind: IncidentSourceNotFound,
	}
}

func NewInvalidIncidentSourceIdError(err error) *IncidentSourceRepositoryError {
	return &IncidentSourceRepositoryError{
		msg:  fmt.Sprintf("Invalid application incident id: %s", err),
		kind: InvalidIncidentSourceId,
	}
}

func NewIncidentSourceInternalError(err error) *IncidentSourceRepositoryError {
	return &IncidentSourceRepositoryError{
		msg:  err.Error(),
		kind: IncidentSourceInternalError,
	}
}

type IncidentSourceRepository[T any] interface {
	InsertIncidentSources(ctx context.Context, sources []*T) ([]string, *IncidentSourceRepositoryError)
	GetIncidentSource(ctx context.Context, id string) (*T, *IncidentSourceRepositoryError)
	GetIncidentSourcesByIds(ctx context.Context, ids []string) ([]*T, *IncidentSourceRepositoryError)
}

type MongoDbIncidentSourceRepository[T any] struct {
	mongoDbCollection *repositories.MongoDbCollection[*T]
	logger            *zap.Logger
}

func (r *MongoDbIncidentSourceRepository[T]) InsertIncidentSources(ctx context.Context, incidents []*T) ([]string, *IncidentSourceRepositoryError) {

	documents := make([]interface{}, 0, len(incidents))
	for _, incident := range incidents {
		documents = append(documents, incident)
	}

	ids, err := r.mongoDbCollection.InsertDocuments(documents)

	if err != nil {
		r.logger.Error("Failed to insert incidents", zap.Error(err))
		return nil, NewIncidentSourceInternalError(err)
	}

	createdIds := array.Map(func(objectId primitive.ObjectID) string {
		return objectId.Hex()
	})(ids)

	return createdIds, nil
}

func (r *MongoDbIncidentSourceRepository[T]) GetIncidentSource(ctx context.Context, id string) (*T, *IncidentSourceRepositoryError) {

	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		r.logger.Error("Failed to parse application incident id", zap.Error(err))
		return nil, NewInvalidIncidentSourceIdError(err)
	}

	incident, err := r.mongoDbCollection.GetDocument(primitive.D{{Key: "_id", Value: objectId}}, primitive.D{})

	if err != nil {
		r.logger.Error("Failed to decode node incident from mongodb", zap.Error(err))
		return nil, NewIncidentSourceNotFoundError(err)
	}

	return incident, nil
}

func (r *MongoDbIncidentSourceRepository[T]) GetIncidentSourcesByIds(ctx context.Context, ids []string) ([]*T, *IncidentSourceRepositoryError) {

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
		return nil, NewIncidentSourceInternalError(err)
	}

	return incidents, nil
}

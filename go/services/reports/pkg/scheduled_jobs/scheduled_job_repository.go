package scheduledjobs

import (
	"context"
	"fmt"
	"github.com/IBM/fp-go/array"
	"github.com/Magpie-Monitor/magpie-monitor/pkg/repositories"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
)

type ScheduledJobRepositoryErrorKind string

type ScheduledJob interface {
	GetId() string
	SetId(newId string)
	GetStatus() string
}

const (
	ScheduledJobNotFound      ScheduledJobRepositoryErrorKind = "SCHEDULED_JOB_NOT_FOUND"
	InvalidScheduledJobId     ScheduledJobRepositoryErrorKind = "INVALID_SCHEDULED_JOB_ID"
	ScheduledJobInternalError ScheduledJobRepositoryErrorKind = "INTERNAL_ERROR"
)

type ScheduledJobRepositoryError struct {
	msg  string
	kind ScheduledJobRepositoryErrorKind
}

func (e *ScheduledJobRepositoryError) Error() string {
	return e.msg
}

func (e *ScheduledJobRepositoryError) Kind() ScheduledJobRepositoryErrorKind {
	return e.kind
}

func NewScheduledJobNotFoundError(err error) *ScheduledJobRepositoryError {
	return &ScheduledJobRepositoryError{
		msg:  fmt.Sprintf(" incident does not exists: %s", err),
		kind: ScheduledJobNotFound,
	}
}

func NewInvalidScheduledJobIdError(err error) *ScheduledJobRepositoryError {
	return &ScheduledJobRepositoryError{
		msg:  fmt.Sprintf("Invalid application incident id: %s", err),
		kind: InvalidScheduledJobId,
	}
}

func NewScheduledJobInternalError(err error) *ScheduledJobRepositoryError {
	return &ScheduledJobRepositoryError{
		msg:  err.Error(),
		kind: ScheduledJobInternalError,
	}
}

type ScheduledJobRepository[T ScheduledJob] interface {
	InsertScheduledJobs(ctx context.Context, jobs []T) ([]string, *ScheduledJobRepositoryError)
	GetScheduledJob(ctx context.Context, id string) (*T, *ScheduledJobRepositoryError)
	GetScheduledJobsByIds(ctx context.Context, ids []string) ([]T, *ScheduledJobRepositoryError)
	GetScheduledJobsByStatus(ctx context.Context, status string) ([]T, *ScheduledJobRepositoryError)
	UpdateScheduledJob(ctx context.Context, job T) *ScheduledJobRepositoryError
}

type MongoDbScheduledJobRepository[T ScheduledJob] struct {
	mongoDbCollection *repositories.MongoDbCollection[T]
	logger            *zap.Logger
}

func NewMongoDbScheduledJobRepository[T ScheduledJob](mongoDbCollection *repositories.MongoDbCollection[T], logger *zap.Logger) *MongoDbScheduledJobRepository[T] {
	return &MongoDbScheduledJobRepository[T]{
		mongoDbCollection: mongoDbCollection,
		logger:            logger,
	}
}

func (r *MongoDbScheduledJobRepository[T]) InsertScheduledJobs(ctx context.Context, incidents []T) ([]string, *ScheduledJobRepositoryError) {

	documents := make([]interface{}, 0, len(incidents))

	if len(incidents) == 0 {
		return make([]string, 0, 0), nil
	}

	for _, incident := range incidents {
		documents = append(documents, incident)
	}

	ids, err := r.mongoDbCollection.InsertDocuments(documents)

	if err != nil {
		r.logger.Error("Failed to insert jobs", zap.Error(err))
		return nil, NewScheduledJobInternalError(err)
	}

	createdIds := array.Map(func(objectId primitive.ObjectID) string {
		return objectId.Hex()
	})(ids)

	return createdIds, nil
}

func (r *MongoDbScheduledJobRepository[T]) GetScheduledJob(ctx context.Context, id string) (*T, *ScheduledJobRepositoryError) {

	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		r.logger.Error("Failed to parse application incident id", zap.Error(err))
		return nil, NewInvalidScheduledJobIdError(err)
	}

	incident, err := r.mongoDbCollection.GetDocument(primitive.D{{Key: "_id", Value: objectId}}, primitive.D{})

	if err != nil {
		r.logger.Error("Failed to decode node incident from mongodb", zap.Error(err))
		return nil, NewScheduledJobNotFoundError(err)
	}

	return &incident, nil
}

func (r *MongoDbScheduledJobRepository[T]) GetScheduledJobsByIds(ctx context.Context, ids []string) ([]T, *ScheduledJobRepositoryError) {

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
		return nil, NewScheduledJobInternalError(err)
	}

	return incidents, nil
}

func (r *MongoDbScheduledJobRepository[T]) GetScheduledJobsByStatus(ctx context.Context, status string) ([]T, *ScheduledJobRepositoryError) {

	incidents, err := r.mongoDbCollection.GetDocuments(bson.D{{Key: "status", Value: status}}, bson.D{{Key: "scheduledAt", Value: 1}})
	if err != nil {
		r.logger.Error("Failed to find application incidents by status", zap.Error(err), zap.Any("status", status))
		return nil, NewScheduledJobInternalError(err)
	}

	return incidents, nil
}

func (r *MongoDbScheduledJobRepository[T]) UpdateScheduledJob(ctx context.Context, scheduledJob T) *ScheduledJobRepositoryError {
	jobId := scheduledJob.GetId()

	id, err := primitive.ObjectIDFromHex(jobId)
	if err != nil {
		r.logger.Error("Failed to encode job id", zap.Error(err))
		return NewInvalidScheduledJobIdError(err)
	}

	scheduledJob.SetId("")

	err = r.mongoDbCollection.ReplaceDocument(ctx, id, scheduledJob)
	if err != nil {
		r.logger.Error("Failed to update scheduled job", zap.Error(err))
		return NewInvalidScheduledJobIdError(err)
	}

	scheduledJob.SetId(jobId)

	return nil
}

// var _ ScheduledJobRepository[any] = MongoDbScheduledJobRepository[O]{}

package repositories

import (
	"context"
	"fmt"

	"github.com/IBM/fp-go/array"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

const (
	APPLICATION_INCIDENT_COLLECTION = "application_incidents"
)

type ApplicationIncidentRepositoryErrorKind string

const (
	ApplicationIncidentNotFound      ApplicationIncidentRepositoryErrorKind = "INCIDENT_NOT_FOUND"
	InvalidApplicationIncidentId     ApplicationIncidentRepositoryErrorKind = "INVALID_INCIDENT_ID"
	ApplicationIncidentInternalError ApplicationIncidentRepositoryErrorKind = "INTERNAL_ERROR"
)

type ApplicationIncidentRepositoryError struct {
	msg  string
	kind ApplicationIncidentRepositoryErrorKind
}

func (e *ApplicationIncidentRepositoryError) Error() string {
	return e.msg
}

func (e *ApplicationIncidentRepositoryError) Kind() ApplicationIncidentRepositoryErrorKind {
	return e.kind
}

func NewApplicationIncidentNotFoundError(err error) *ApplicationIncidentRepositoryError {
	return &ApplicationIncidentRepositoryError{
		msg:  fmt.Sprintf("Application incident does not exists: %s", err),
		kind: ApplicationIncidentNotFound,
	}
}

func NewInvalidApplicationIncidentIdError(err error) *ApplicationIncidentRepositoryError {
	return &ApplicationIncidentRepositoryError{
		msg:  fmt.Sprintf("Invalid application incident id: %s", err),
		kind: InvalidApplicationIncidentId,
	}
}

func NewApplicationIncidentInternalError(err error) *ApplicationIncidentRepositoryError {
	return &ApplicationIncidentRepositoryError{
		msg:  err.Error(),
		kind: ApplicationIncidentInternalError,
	}
}

type ApplicationIncidentSource struct {
	Timestamp     int64  `bson:"timestamp" json:"timestamp"`
	PodName       string `bson:"podName" json:"podName"`
	ContainerName string `bson:"containerName" json:"containerName"`
	Image         string `bson:"image" json:"image"`
	Content       string `bson:"content" json:"content"`
}

type ApplicationIncident struct {
	Id             string                      `bson:"_id,omitempty" json:"id"`
	Category       string                      `bson:"category" json:"category"`
	Summary        string                      `bson:"summary" json:"summary"`
	Recommendation string                      `bson:"recommendation" json:"recommendation"`
	Urgency        Urgency                     `bson:"urgency" json:"urgency"`
	Sources        []ApplicationIncidentSource `bson:"sources" json:"sources"`
}

type ApplicationIncidentRepository interface {
	InsertIncidents(ctx context.Context, incidents []*ApplicationIncident) ([]string, error)
	GetIncident(ctx context.Context, id string) (*ApplicationIncident, error)
	GetIncidentsByIds(ctx context.Context, ids []string) ([]*ApplicationIncident, error)
}

type MongoDbApplicationIncidentRepository struct {
	mongoDbClient *mongo.Client
	logger        *zap.Logger
}

type ApplicationIncidentParams struct {
	fx.In
	ApplicationIncidentsDbMongoClient *mongo.Client
	Logger                            *zap.Logger
}

func NewMongoDbApplicationIncidentRepository(p ApplicationIncidentParams) *MongoDbApplicationIncidentRepository {

	return &MongoDbApplicationIncidentRepository{
		mongoDbClient: p.ApplicationIncidentsDbMongoClient,
		logger:        p.Logger,
	}
}

func ProvideAsApplicationIncidentRepository(f any) any {
	return fx.Annotate(
		f,
		fx.As(new(ApplicationIncidentRepository)),
	)
}

func (r *MongoDbApplicationIncidentRepository) InsertIncidents(ctx context.Context, incidents []*ApplicationIncident) ([]string, error) {

	coll := r.mongoDbClient.Database(REPORTS_DB_NAME).Collection(APPLICATION_INCIDENT_COLLECTION)

	documents := make([]interface{}, 0, len(incidents))
	for _, incident := range incidents {
		documents = append(documents, incident)
	}

	res, err := coll.InsertMany(ctx, documents)
	if err != nil {
		r.logger.Error("Failed to insert node incidents", zap.Error(err))
		return nil, err
	}

	createdIds := array.Map(func(objectId interface{}) string {
		return objectId.(primitive.ObjectID).Hex()
	})(res.InsertedIDs)

	return createdIds, nil
}

func (r *MongoDbApplicationIncidentRepository) GetIncident(ctx context.Context, id string) (*ApplicationIncident, error) {

	coll := r.mongoDbClient.Database(REPORTS_DB_NAME).Collection(APPLICATION_INCIDENT_COLLECTION)

	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		r.logger.Error("Failed to parse application incident id", zap.Error(err))
		return nil, NewInvalidApplicationIncidentIdError(err)
	}

	documents := coll.FindOne(ctx, bson.M{"_id": objectId})

	var result *ApplicationIncident
	err = documents.Decode(&result)
	if err != nil {
		r.logger.Error("Failed to decode node incident from mongodb", zap.Error(err))
		return nil, NewApplicationIncidentNotFoundError(err)
	}

	return result, nil
}

func (r *MongoDbApplicationIncidentRepository) GetIncidentsByIds(ctx context.Context, ids []string) ([]*ApplicationIncident, error) {
	coll := r.mongoDbClient.Database(REPORTS_DB_NAME).Collection(APPLICATION_INCIDENT_COLLECTION)

	idObjects := array.Map(func(id string) primitive.ObjectID {
		idObj, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			r.logger.Error("Failed to parse application incident ids", zap.Error(err))
		}
		return idObj
	})(ids)

	cur, err := coll.Find(ctx, bson.M{"_id": bson.M{"$in": idObjects}})
	if err != nil {
		r.logger.Error("Failed to find node incidents by ids", zap.Error(err))
		return nil, err
	}

	var results []*ApplicationIncident
	if err = cur.All(ctx, &results); err != nil {
		r.logger.Error("Failed to decode node incidents by ids", zap.Error(err))
		return nil, err
	}

	return results, nil
}

// Compile-time check if MongoDbApplicationIncidentRepository implements
// the ApplicationIncidentRepository interface
var _ ApplicationIncidentRepository = &MongoDbApplicationIncidentRepository{}

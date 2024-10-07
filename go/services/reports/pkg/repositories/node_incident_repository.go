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
	NODE_INCIDENT_COLLECTION = "node_incidents"
)

type NodeIncidentRepositoryErrorKind string

const (
	NodeIncidentNotFound      NodeIncidentRepositoryErrorKind = "REPORT_NOT_FOUND"
	InvalidNodeIncidentId     NodeIncidentRepositoryErrorKind = "INVALID_REPORT_ID"
	NodeIncidentInternalError NodeIncidentRepositoryErrorKind = "INTERNAL_ERROR"
)

type NodeIncidentRepositoryError struct {
	msg  string
	kind NodeIncidentRepositoryErrorKind
}

func (e *NodeIncidentRepositoryError) Error() string {
	return e.msg
}

func (e *NodeIncidentRepositoryError) Kind() NodeIncidentRepositoryErrorKind {
	return e.kind
}

func NewNodeIncidentNotFoundError(err error) *NodeIncidentRepositoryError {
	return &NodeIncidentRepositoryError{
		msg:  fmt.Sprintf("Node incident does not exists: %s", err),
		kind: NodeIncidentNotFound,
	}
}

func NewInvalidNodeIncidentIdError(err error) *NodeIncidentRepositoryError {
	return &NodeIncidentRepositoryError{
		msg:  fmt.Sprintf("Invalid node incident id: %s", err),
		kind: InvalidNodeIncidentId,
	}
}

func NewNodeIncidentInternalError(err error) *NodeIncidentRepositoryError {
	return &NodeIncidentRepositoryError{
		msg:  err.Error(),
		kind: NodeIncidentInternalError,
	}
}

type NodeIncidentSource struct {
	Timestamp int64  `bson:"timestamp" json:"timestamp"`
	NodeName  string `bson:"nodeName" json:"nodeName"`
	Content   string `bson:"content" json:"content"`
}

type NodeIncident struct {
	Id             string               `bson:"_id,omitempty" json:"id"`
	Category       string               `bson:"category" json:"category"`
	Summary        string               `bson:"summary" json:"summary"`
	Recommendation string               `bson:"recommendation" json:"recommendation"`
	Urgency        Urgency              `bson:"urgency" json:"urgency"`
	Sources        []NodeIncidentSource `bson:"sources" json:"sources"`
}

type NodeIncidentRepository interface {
	InsertIncidents(ctx context.Context, incidents []*NodeIncident) ([]string, error)
	GetIncident(ctx context.Context, id string) (*NodeIncident, error)
	GetIncidentsByIds(ctx context.Context, ids []string) ([]*NodeIncident, error)
}

type MongoDbNodeIncidentRepository struct {
	mongoDbClient *mongo.Client
	logger        *zap.Logger
}

type NodeIncidentParams struct {
	fx.In
	NodeIncidentsDbMongoClient *mongo.Client
	Logger                     *zap.Logger
}

func NewMongoDbNodeIncidentRepository(p NodeIncidentParams) *MongoDbNodeIncidentRepository {

	return &MongoDbNodeIncidentRepository{
		mongoDbClient: p.NodeIncidentsDbMongoClient,
		logger:        p.Logger,
	}
}

func ProvideAsNodeIncidentRepository(f any) any {
	return fx.Annotate(
		f,
		fx.As(new(NodeIncidentRepository)),
	)
}

func (r *MongoDbNodeIncidentRepository) InsertIncidents(ctx context.Context, incidents []*NodeIncident) ([]string, error) {

	coll := r.mongoDbClient.Database(REPORTS_DB_NAME).Collection(NODE_INCIDENT_COLLECTION)

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

func (r *MongoDbNodeIncidentRepository) GetIncident(ctx context.Context, id string) (*NodeIncident, error) {

	coll := r.mongoDbClient.Database(REPORTS_DB_NAME).Collection(NODE_INCIDENT_COLLECTION)

	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		r.logger.Error("Failed to parse node incident id", zap.Error(err))
		return nil, NewInvalidNodeIncidentIdError(err)
	}

	documents := coll.FindOne(ctx, bson.M{"_id": objectId})

	var result *NodeIncident
	err = documents.Decode(&result)
	if err != nil {
		r.logger.Error("Failed to decode node incident from mongodb", zap.Error(err))
		return nil, NewNodeIncidentNotFoundError(err)
	}

	return result, nil
}

func (r *MongoDbNodeIncidentRepository) GetIncidentsByIds(ctx context.Context, ids []string) ([]*NodeIncident, error) {
	coll := r.mongoDbClient.Database(REPORTS_DB_NAME).Collection(NODE_INCIDENT_COLLECTION)

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

	var results []*NodeIncident
	if err = cur.All(ctx, &results); err != nil {
		r.logger.Error("Failed to decode node incidents by ids", zap.Error(err))
		return nil, err
	}

	return results, nil
}

// Compile-time check if MongoDbNodeIncidentRepository implements
// the NodeIncidentRepository interface
var _ NodeIncidentRepository = &MongoDbNodeIncidentRepository{}

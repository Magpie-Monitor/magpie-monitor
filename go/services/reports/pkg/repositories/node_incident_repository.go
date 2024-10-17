package repositories

import (
	"github.com/Magpie-Monitor/magpie-monitor/pkg/repositories"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type NodeIncidentSource struct {
	Timestamp int64  `bson:"timestamp" json:"timestamp"`
	Content   string `bson:"content" json:"content"`
	Filename  string `bson:"filename" json:"filename"`
}

type NodeIncident struct {
	Id             string               `bson:"_id,omitempty" json:"id"`
	ClusterId      string               `bson:"clusterId" json:"clusterId"`
	NodeName       string               `bson:"nodeName" json:"nodeName"`
	Category       string               `bson:"category" json:"category"`
	Summary        string               `bson:"summary" json:"summary"`
	Recommendation string               `bson:"recommendation" json:"recommendation"`
	Urgency        Urgency              `bson:"urgency" json:"urgency"`
	Sources        []NodeIncidentSource `bson:"sources" json:"sources"`
}

func NewNodeIncidentsCollection(log *zap.Logger, client *mongo.Client) *repositories.MongoDbCollection[*NodeIncident] {
	return &repositories.MongoDbCollection[*NodeIncident]{Log: log, Db: "reports", Col: "node_incidents", Client: client}
}

type NodeIncidentParams struct {
	fx.In
	IncidentsDbMongoCollection *repositories.MongoDbCollection[*NodeIncident]
	Logger                     *zap.Logger
}

func NewMongoDbNodeIncidentRepository(p NodeIncidentParams) *MongoDbIncidentRepository[NodeIncident] {

	return &MongoDbIncidentRepository[NodeIncident]{
		mongoDbCollection: p.IncidentsDbMongoCollection,
		logger:            p.Logger,
	}
}

func ProvideAsNodeIncidentRepository(f any) any {
	return fx.Annotate(
		f,
		fx.As(new(IncidentRepository[NodeIncident])),
	)
}

// Compile-time check if MongoDbIncidentRepository implements
// the IncidentRepository interface
var _ IncidentRepository[NodeIncident] = &MongoDbIncidentRepository[NodeIncident]{}

package repositories

import (
	"github.com/Magpie-Monitor/magpie-monitor/pkg/repositories"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type NodeIncidentSource struct {
	Id          string `bson:"_id,omitempty" json:"id"`
	ReportId    string `bson:"reportId" json:"reportId"`
	IncidentId  string `bson:"incidentId" json:"incidentId"`
	Timestamp   int64  `bson:"timestamp" json:"timestamp"`
	Filename    string `bson:"filename" json:"filename"`
	SourceLog   string `bson:"sourceLog" json:"sourceLog"`
	SourceLogId string `bson:"sourceLogId" json:"sourceLogId"`
}

func NewNodeIncidentSourcesCollection(
	log *zap.Logger,
	client *mongo.Client) *repositories.MongoDbCollection[*NodeIncidentSource] {

	return &repositories.MongoDbCollection[*NodeIncidentSource]{
		Log: log, Db: "reports", Col: "node_incident_sources", Client: client}
}

type NodeIncidentSourceParams struct {
	fx.In
	IncidentsDbMongoCollection *repositories.MongoDbCollection[*NodeIncidentSource]
	Logger                     *zap.Logger
}

func NewMongoDbNodeIncidentSourceRepository(p NodeIncidentSourceParams) *MongoDbIncidentSourceRepository[NodeIncidentSource] {

	return &MongoDbIncidentSourceRepository[NodeIncidentSource]{
		mongoDbCollection: p.IncidentsDbMongoCollection,
		logger:            p.Logger,
	}
}

func ProvideAsNodeIncidentSourceRepository(f any) any {
	return fx.Annotate(
		f,
		fx.As(new(IncidentSourceRepository[NodeIncidentSource])),
	)
}

// Compile-time check if MongoDbIncidentRepository implements
// the IncidentRepository[NodeIncidentSource] interface
var _ IncidentSourceRepository[NodeIncidentSource] = &MongoDbIncidentSourceRepository[NodeIncidentSource]{}

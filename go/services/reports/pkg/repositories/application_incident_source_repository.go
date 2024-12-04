package repositories

import (
	"github.com/Magpie-Monitor/magpie-monitor/pkg/repositories"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type ApplicationIncidentSource struct {
	Id            string `bson:"_id,omitempty" json:"id"`
	ReportId      string `bson:"reportId" json:"reportId"`
	IncidentId    string `bson:"incidentId" json:"incidentId"`
	CorrelationId string `bson:"correlationId" json:"correlationId"`
	Timestamp     int64  `bson:"timestamp" json:"timestamp"`
	PodName       string `bson:"podName" json:"podName"`
	ContainerName string `bson:"containerName" json:"containerName"`
	Image         string `bson:"image" json:"image"`
	SourceLog     string `bson:"sourceLog" json:"sourceLog"`
	SourceLogId   string `bson:"sourceLogId" json:"sourceLogId"`
}

func NewApplicationIncidentSourcesCollection(
	log *zap.Logger,
	client *mongo.Client) *repositories.MongoDbCollection[*ApplicationIncidentSource] {

	return &repositories.MongoDbCollection[*ApplicationIncidentSource]{
		Log: log, Db: "reports", Col: "application_incident_sources", Client: client}
}

type ApplicationIncidentSourceParams struct {
	fx.In
	IncidentsDbMongoCollection *repositories.MongoDbCollection[*ApplicationIncidentSource]
	Logger                     *zap.Logger
}

func NewMongoDbApplicationIncidentSourceRepository(p ApplicationIncidentSourceParams) *MongoDbIncidentSourceRepository[ApplicationIncidentSource] {

	return &MongoDbIncidentSourceRepository[ApplicationIncidentSource]{
		mongoDbCollection: p.IncidentsDbMongoCollection,
		logger:            p.Logger,
	}
}

func ProvideAsApplicationIncidentSourceRepository(f any) any {
	return fx.Annotate(
		f,
		fx.As(new(IncidentSourceRepository[ApplicationIncidentSource])),
	)
}

// Compile-time check if MongoDbIncidentRepository implements
// the IncidentRepository[ApplicationIncidentSource] interface
var _ IncidentSourceRepository[ApplicationIncidentSource] = &MongoDbIncidentSourceRepository[ApplicationIncidentSource]{}

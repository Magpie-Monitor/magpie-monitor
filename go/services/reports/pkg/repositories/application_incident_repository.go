package repositories

import (
	"github.com/Magpie-Monitor/magpie-monitor/pkg/repositories"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type ApplicationIncidentSource struct {
	Timestamp     int64  `bson:"timestamp" json:"timestamp"`
	PodName       string `bson:"podName" json:"podName"`
	ContainerName string `bson:"containerName" json:"containerName"`
	Image         string `bson:"image" json:"image"`
	Content       string `bson:"content" json:"content"`
}

type ApplicationIncident struct {
	Id              string                      `bson:"_id,omitempty" json:"id"`
	ApplicationName string                      `bson:"applicationName" json:"applicationName"`
	ClusterId       string                      `bson:"clusterId" json:"clusterId"`
	Category        string                      `bson:"category" json:"category"`
	Summary         string                      `bson:"summary" json:"summary"`
	Recommendation  string                      `bson:"recommendation" json:"recommendation"`
	Urgency         Urgency                     `bson:"urgency" json:"urgency"`
	Sources         []ApplicationIncidentSource `bson:"sources" json:"sources"`
}

func NewApplicationIncidentsCollection(log *zap.Logger, client *mongo.Client) *repositories.MongoDbCollection[*ApplicationIncident] {
	return &repositories.MongoDbCollection[*ApplicationIncident]{Log: log, Db: "reports", Col: "application_incidents", Client: client}
}

type ApplicationIncidentParams struct {
	fx.In
	IncidentsDbMongoCollection *repositories.MongoDbCollection[*ApplicationIncident]
	Logger                     *zap.Logger
}

func NewMongoDbApplicationIncidentRepository(p ApplicationIncidentParams) *MongoDbIncidentRepository[ApplicationIncident] {

	return &MongoDbIncidentRepository[ApplicationIncident]{
		mongoDbCollection: p.IncidentsDbMongoCollection,
		logger:            p.Logger,
	}
}

func ProvideAsApplicationIncidentRepository(f any) any {
	return fx.Annotate(
		f,
		fx.As(new(IncidentRepository[ApplicationIncident])),
	)
}

// Compile-time check if MongoDbIncidentRepository implements
// the IncidentRepository[ApplicationIncident] interface
var _ IncidentRepository[ApplicationIncident] = &MongoDbIncidentRepository[ApplicationIncident]{}
package repositories

import (
	"github.com/Magpie-Monitor/magpie-monitor/pkg/repositories"
	"github.com/Magpie-Monitor/magpie-monitor/services/reports/pkg/insights"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// type NodeIncidentSource struct {
// 	Timestamp int64  `bson:"timestamp" json:"timestamp"`
// 	Content   string `bson:"content" json:"content"`
// 	Filename  string `bson:"filename" json:"filename"`
// }

type NodeIncident struct {
	Id             string               `bson:"_id,omitempty" json:"id"`
	Title          string               `bson:"title" json:"title"`
	Accuracy       insights.Accuracy    `bson:"accuracy" json:"accuracy"`
	CustomPrompt   string               `bson:"customPrompt" json:"customPrompt"`
	ClusterId      string               `bson:"clusterId" json:"clusterId"`
	NodeName       string               `bson:"nodeName" json:"nodeName"`
	Category       string               `bson:"category" json:"category"`
	Summary        string               `bson:"summary" json:"summary"`
	Recommendation string               `bson:"recommendation" json:"recommendation"`
	Urgency        insights.Urgency     `bson:"urgency" json:"urgency"`
	Sources        []NodeIncidentSource `bson:"-" json:"-"`
	SourceIds      []string             `bson:"sourceIds" json:"sourceIds"`
}

func (i *NodeIncident) GetRecommendation() string {
	return i.Recommendation
}

func (i *NodeIncident) GetTitle() string {
	return i.Title
}

func (i *NodeIncident) GetSummary() string {
	return i.Summary
}

func (i *NodeIncident) GetId() string {
	return i.Id
}

func (i *NodeIncident) GetCategory() string {
	return i.Category
}

func (i *NodeIncident) SetId(newId string) {
	i.Id = newId
}

func (i *NodeIncident) GetUrgency() insights.Urgency {
	return i.Urgency
}

func NewNodeIncidentsCollection(log *zap.Logger, client *mongo.Client) *repositories.MongoDbCollection[*NodeIncident] {
	return &repositories.MongoDbCollection[*NodeIncident]{Log: log, Db: "reports", Col: "node_incidents", Client: client}
}

type NodeIncidentParams struct {
	fx.In
	IncidentsDbMongoCollection *repositories.MongoDbCollection[*NodeIncident]
	Logger                     *zap.Logger
}

func NewMongoDbNodeIncidentRepository(p NodeIncidentParams) *MongoDbIncidentRepository[*NodeIncident] {

	return &MongoDbIncidentRepository[*NodeIncident]{
		mongoDbCollection: p.IncidentsDbMongoCollection,
		logger:            p.Logger,
	}
}

func ProvideAsNodeIncidentRepository(f any) any {
	return fx.Annotate(
		f,
		fx.As(new(IncidentRepository[*NodeIncident])),
	)
}

// Compile-time check if MongoDbIncidentRepository implements
// the IncidentRepository interface
var _ IncidentRepository[*NodeIncident] = &MongoDbIncidentRepository[*NodeIncident]{}

var _ Incident = &NodeIncident{}

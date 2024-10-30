package repositories

import (
	"github.com/Magpie-Monitor/magpie-monitor/pkg/repositories"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

func NewClusterMetadataCollection(log *zap.Logger, client *mongo.Client) *repositories.MongoDbCollection[ApplicationState] {
	return &repositories.MongoDbCollection[ApplicationState]{Log: log, Db: "METADATA", Col: "APPLICATION_METADATA", Client: client}
}

func NewNodeMetadataCollection(log *zap.Logger, client *mongo.Client) *repositories.MongoDbCollection[NodeState] {
	return &repositories.MongoDbCollection[NodeState]{Log: log, Db: "METADATA", Col: "NODE_METADATA", Client: client}
}

func NewApplicationAggregatedMetadataCollection(log *zap.Logger, client *mongo.Client) *repositories.MongoDbCollection[AggregatedApplicationMetadata] {
	return &repositories.MongoDbCollection[AggregatedApplicationMetadata]{Log: log, Db: "METADATA", Col: "AGGREGATED_APPLICATION_METADATA", Client: client}
}

func NewNodeAggregatedMetadataCollection(log *zap.Logger, client *mongo.Client) *repositories.MongoDbCollection[AggregatedNodeMetadata] {
	return &repositories.MongoDbCollection[AggregatedNodeMetadata]{Log: log, Db: "METADATA", Col: "AGGREGATED_NODE_METADATA", Client: client}
}

func NewClusterAggregatedStateCollection(log *zap.Logger, client *mongo.Client) *repositories.MongoDbCollection[AggregatedClusterState] {
	return &repositories.MongoDbCollection[AggregatedClusterState]{Log: log, Db: "METADATA", Col: "AGGREGATED_CLUSTER_STATE", Client: client}
}

type AggregatedApplicationMetadata struct {
	CollectedAtMs int64                 `json:"collectedAtMs" bson:"collectedAtMs"`
	ClusterId     string                `json:"clusterId" bson:"clusterId"`
	Metadata      []ApplicationMetadata `json:"metadata" bson:"metadata"`
}

type ApplicationMetadata struct {
	Name string `json:"name" bson:"name"`
	Kind string `json:"kind" bson:"kind"`
}

type AggregatedNodeMetadata struct {
	CollectedAtMs int64          `json:"collectedAtMs" bson:"collectedAtMs"`
	ClusterId     string         `json:"clusterId" bson:"clusterId"`
	Metadata      []NodeMetadata `json:"metadata" bson:"metadata"`
}

type NodeMetadata struct {
	Name  string        `json:"name"`
	Files []interface{} `json:"files"`
}

type AggregatedClusterState struct {
	CollectedAtMs int64             `json:"collectedAtMs" bson:"collectedAtMs"`
	Metadata      []ClusterMetadata `json:"metadata" bson:"metadata"`
}

type ClusterMetadata struct {
	ClusterId string `json:"clusterId"`
}

type ApplicationState struct {
	CollectedAtMs int64         `json:"collectedAtMs" bson:"collectedAtMs"`
	ClusterId     string        `json:"clusterId" bson:"clusterId"`
	Applications  []Application `json:"applications" bson:"applications"`
}

type Application struct {
	Kind string `json:"kind" bson:"kind"`
	Name string `json:"name" bson:"name"`
}

type NodeState struct {
	ClusterId     string   `json:"clusterId" bson:"clusterId"`
	NodeName      string   `json:"nodeName" bson:"nodeName"`
	CollectedAtMs int64    `json:"collectedAtMs" bson:"collectedAtMs"`
	WatchedFiles  []string `json:"watchedFiles" bson:"watchedFiles"`
}

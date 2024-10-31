package repositories

import (
	"github.com/Magpie-Monitor/magpie-monitor/pkg/repositories"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

const (
	DATABASE = "METADATA"
)

func NewApplicationAggregatedMetadataCollection(log *zap.Logger, client *mongo.Client) *repositories.MongoDbCollection[AggregatedApplicationMetadata] {
	return &repositories.MongoDbCollection[AggregatedApplicationMetadata]{Log: log, Db: DATABASE, Col: "AGGREGATED_APPLICATION_METADATA", Client: client}
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

func NewNodeAggregatedMetadataCollection(log *zap.Logger, client *mongo.Client) *repositories.MongoDbCollection[AggregatedNodeMetadata] {
	return &repositories.MongoDbCollection[AggregatedNodeMetadata]{Log: log, Db: DATABASE, Col: "AGGREGATED_NODE_METADATA", Client: client}
}

type AggregatedNodeMetadata struct {
	CollectedAtMs int64          `json:"collectedAtMs" bson:"collectedAtMs"`
	ClusterId     string         `json:"clusterId" bson:"clusterId"`
	Metadata      []NodeMetadata `json:"metadata" bson:"metadata"`
}

type NodeMetadata struct {
	Name  string        `json:"name" bson:"name"`
	Files []interface{} `json:"files" bson:"files"`
}

func NewClusterAggregatedStateCollection(log *zap.Logger, client *mongo.Client) *repositories.MongoDbCollection[AggregatedClusterMetadata] {
	return &repositories.MongoDbCollection[AggregatedClusterMetadata]{Log: log, Db: DATABASE, Col: "AGGREGATED_CLUSTER_STATE", Client: client}
}

type AggregatedClusterMetadata struct {
	CollectedAtMs int64             `json:"collectedAtMs" bson:"collectedAtMs"`
	Metadata      []ClusterMetadata `json:"metadata" bson:"metadata"`
}

type ClusterMetadata struct {
	ClusterId string `json:"clusterId" bson:"clusterId"`
}

func NewApplicationMetadataCollection(log *zap.Logger, client *mongo.Client) *repositories.MongoDbCollection[ApplicationState] {
	return &repositories.MongoDbCollection[ApplicationState]{Log: log, Db: DATABASE, Col: "APPLICATION_METADATA", Client: client}
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

func NewNodeMetadataCollection(log *zap.Logger, client *mongo.Client) *repositories.MongoDbCollection[NodeState] {
	return &repositories.MongoDbCollection[NodeState]{Log: log, Db: DATABASE, Col: "NODE_METADATA", Client: client}
}

type NodeState struct {
	ClusterId     string   `json:"clusterId" bson:"clusterId"`
	NodeName      string   `json:"nodeName" bson:"nodeName"`
	CollectedAtMs int64    `json:"collectedAtMs" bson:"collectedAtMs"`
	WatchedFiles  []string `json:"watchedFiles" bson:"watchedFiles"`
}

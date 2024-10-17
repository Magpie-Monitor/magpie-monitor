package repositories

import (
	"github.com/Magpie-Monitor/magpie-monitor/pkg/repositories"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

func NewClusterMetadataCollection(log *zap.Logger, client *mongo.Client) *repositories.MongoDbCollection[ClusterState] {
	return &repositories.MongoDbCollection[ClusterState]{Log: log, Db: "METADATA", Col: "CLUSTER_STATE", Client: client}
}

func NewNodeMetadataCollection(log *zap.Logger, client *mongo.Client) *repositories.MongoDbCollection[NodeState] {
	return &repositories.MongoDbCollection[NodeState]{Log: log, Db: "METADATA", Col: "NODE_STATE", Client: client}
}

type ClusterState struct {
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

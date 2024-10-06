package repositories

import (
	"context"

	"github.com/Magpie-Monitor/magpie-monitor/services/cluster_metadata/internal/entity"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

func NewClusterMetadataCollection(log *zap.Logger, client *mongo.Client) *MongoDbCollection[entity.ClusterState] {
	return &MongoDbCollection[entity.ClusterState]{log: log, db: "METADATA", col: "CLUSTER_STATE", client: client}
}

func NewNodeMetadataCollection(log *zap.Logger, client *mongo.Client) *MongoDbCollection[entity.NodeState] {
	return &MongoDbCollection[entity.NodeState]{log: log, db: "METADATA", col: "NODE_STATE", client: client}
}

type ClusterState struct {
	CollectedAtMs int64         `json:"collectedAtMs"`
	ClusterName   string        `json:"clusterName"`
	Applications  []Application `json:"applications"`
}

type Application struct {
	Kind string `json:"kind"`
	Name string `json:"name"`
}

type NodeState struct {
	NodeName      string   `json:"nodeName"`
	CollectedAtMs int64    `json:"collectedAtMs"`
	WatchedFiles  []string `json:"watchedFiles"`
}

type MongoDbCollection[T any] struct {
	log    *zap.Logger
	db     string
	col    string
	client *mongo.Client
}

func (m *MongoDbCollection[T]) GetFilteredDocuments(filter primitive.D) ([]T, error) {
	col := m.client.Database(m.db).Collection(m.col)

	cursor, err := col.Find(context.TODO(), filter)
	if err != nil {
		m.log.Error("Error fetching documents:", zap.String("database", m.db), zap.String("collection", m.col), zap.Error(err))
		return nil, err
	}

	var results []T
	if err = cursor.All(context.TODO(), &results); err != nil {
		m.log.Error("Error parsing filtered documents:", zap.String("database", m.db), zap.String("collection", m.col), zap.Error(err))
		return nil, err
	}

	return results, nil
}

func (m *MongoDbCollection[T]) InsertDocuments(docs []interface{}) error {
	col := m.client.Database(m.db).Collection(m.col)

	_, err := col.InsertMany(context.TODO(), docs)
	if err != nil {
		m.log.Error("Error inserting documents:", zap.String("database", m.db), zap.String("collection", m.col), zap.Error(err))
		return err
	}

	return nil
}

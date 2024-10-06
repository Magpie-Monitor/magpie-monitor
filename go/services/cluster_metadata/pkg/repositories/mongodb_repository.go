package repositories

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

func NewClusterMetadataCollection(log *zap.Logger, client *mongo.Client) *MongoDbCollection[ClusterState] {
	return &MongoDbCollection[ClusterState]{log: log, db: "METADATA", col: "CLUSTER_STATE", client: client}
}

func NewNodeMetadataCollection(log *zap.Logger, client *mongo.Client) *MongoDbCollection[NodeState] {
	return &MongoDbCollection[NodeState]{log: log, db: "METADATA", col: "NODE_STATE", client: client}
}

type ClusterState struct {
	CollectedAtMs int64         `json:"collectedAtMs" bson:"collectedAtMs"`
	ClusterName   string        `json:"clusterName" bson:"clusterName"`
	Applications  []Application `json:"applications" bson:"applications"`
}

type Application struct {
	Kind string `json:"kind" bson:"kind"`
	Name string `json:"name" bson:"name"`
}

type NodeState struct {
	ClusterName   string   `json:"clusterName" bson:"clusterName"`
	NodeName      string   `json:"nodeName" bson:"nodeName"`
	CollectedAtMs int64    `json:"collectedAtMs" bson:"collectedAtMs"`
	WatchedFiles  []string `json:"watchedFiles" bson:"watchedFiles"`
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
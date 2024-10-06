package repositories

import (
	"context"

	"github.com/Magpie-Monitor/magpie-monitor/services/cluster_metadata/internal/entity"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MetadataRepository interface {
}

type MongoDbRepository[T any] struct {
	db     string
	col    string
	client *mongo.Client
	// logger *zap.Logger
}

func NewClusterMetadataRepository(client *mongo.Client) MongoDbRepository[entity.ClusterState] {
	return MongoDbRepository[entity.ClusterState]{db: "METADATA", col: "CLUSTER_STATE", client: client}
}

func NewNodeMetadataRepository(client *mongo.Client) MongoDbRepository[entity.NodeState] {
	return MongoDbRepository[entity.NodeState]{db: "METADATA", col: "NODE_STATE", client: client}
}

func NewMongoRepo[T any](db, col string, client *mongo.Client) MongoDbRepository[T] {
	return MongoDbRepository[T]{db: db, col: col, client: client}
}

func (m *MongoDbRepository[T]) GetFilteredDocuments(filter primitive.D) ([]T, error) {
	col := m.client.Database(m.db).Collection(m.col)

	cursor, err := col.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}

	var results []T
	if err = cursor.All(context.TODO(), &results); err != nil {
		return nil, err
	}

	return results, nil
}

func (m *MongoDbRepository[T]) InsertDocuments(docs []interface{}) error {
	col := m.client.Database(m.db).Collection(m.col)

	_, err := col.InsertMany(context.TODO(), docs)
	if err != nil {
		return err
	}

	return nil
}

// func (m *MongoDbMetadataRepository[T]) GetClusterMetadata(clusterName string, siceMs int64, toMs int64) []T {
// 	col := m.client.Database("METADATA").Collection("CLUSTER_METADATA")

// 	filter := bson.D{
// 		{"$and", bson.A{
// 			bson.D{{"collectedatms", bson.D{{"$gt", siceMs}}}},
// 			bson.D{{"collectedatms", bson.D{{"$lt", toMs}}}},
// 			bson.D{{"clustername", bson.D{{"$eq", clusterName}}}},
// 		}},
// 	}

// 	cursor, err := col.Find(context.TODO(), filter)
// 	if err != nil {
// 		panic(err)
// 	}

// 	var results []T
// 	if err = cursor.All(context.TODO(), &results); err != nil {
// 		panic(err)
// 	}

// 	return results
// }

// func (m *MongoDbMetadataRepository) GetNodeMetadata(nodeName string, siceMs int64, toMs int64) []entity.NodeState {
// 	col := m.client.Database("METADATA").Collection("NODE_METADATA")

// 	filter := bson.D{
// 		{"$and", bson.A{
// 			bson.D{{"collectedatms", bson.D{{"$gt", siceMs}}}},
// 			bson.D{{"collectedatms", bson.D{{"$lt", toMs}}}},
// 			bson.D{{"nodename", bson.D{{"$eq", nodeName}}}},
// 		}},
// 	}

// 	cursor, err := col.Find(context.TODO(), filter)
// 	if err != nil {
// 		panic(err)
// 	}

// 	var results []entity.NodeState
// 	if err = cursor.All(context.TODO(), &results); err != nil {
// 		panic(err)
// 	}

// 	return results
// }

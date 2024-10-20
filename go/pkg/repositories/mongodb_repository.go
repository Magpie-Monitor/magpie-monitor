package repositories

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

type MongoDbCollection[T any] struct {
	Log    *zap.Logger
	Db     string
	Col    string
	Client *mongo.Client
}

func (m *MongoDbCollection[T]) GetDocuments(filter primitive.D, sort primitive.D) ([]T, error) {
	opts := options.Find().SetSort(sort)

	col := m.Client.Database(m.Db).Collection(m.Col)

	cursor, err := col.Find(context.TODO(), filter, opts)
	if err != nil {
		m.Log.Error("Error fetching documents:", zap.String("database", m.Db), zap.String("collection", m.Col), zap.Error(err))
		return nil, err
	}

	results := []T{}
	if err = cursor.All(context.TODO(), &results); err != nil {
		m.Log.Error("Error parsing filtered documents:", zap.String("database", m.Db), zap.String("collection", m.Col), zap.Error(err))
		return nil, err
	}

	return results, nil
}

func (m *MongoDbCollection[T]) GetDocument(filter primitive.D, sort primitive.D) (T, error) {
	opts := options.FindOne().SetSort(sort)

	col := m.Client.Database(m.Db).Collection(m.Col)

	var result T
	err := col.FindOne(context.TODO(), filter, opts).Decode(&result)

	if err != nil {
		m.Log.Error("Error parsing filtered document:", zap.String("database", m.Db), zap.String("collection", m.Col), zap.Error(err))
		return result, err
	}

	return result, nil
}

func (m *MongoDbCollection[T]) GetDistinctDocumentFieldValues(fieldName string, filter bson.D) ([]interface{}, error) {
	col := m.Client.Database(m.Db).Collection(m.Col)
	return col.Distinct(context.TODO(), fieldName, filter)
}

func (m *MongoDbCollection[T]) InsertDocuments(docs []interface{}) ([]primitive.ObjectID, error) {
	col := m.Client.Database(m.Db).Collection(m.Col)

	ids, err := col.InsertMany(context.TODO(), docs)
	if err != nil {
		m.Log.Error("Error inserting documents:", zap.String("database", m.Db), zap.String("collection", m.Col), zap.Error(err))
		return nil, err
	}

	createdIds := make([]primitive.ObjectID, 0, len(ids.InsertedIDs))
	for _, id := range ids.InsertedIDs {
		createdIds = append(createdIds, id.(primitive.ObjectID))
	}

	return createdIds, nil
}

func (m *MongoDbCollection[T]) InsertDocument(doc interface{}) (*primitive.ObjectID, error) {
	col := m.Client.Database(m.Db).Collection(m.Col)

	id, err := col.InsertOne(context.TODO(), doc)
	if err != nil {
		m.Log.Error("Error inserting document:", zap.String("database", m.Db), zap.String("collection", m.Col), zap.Error(err))
		return nil, err
	}
	objectId := id.InsertedID.(primitive.ObjectID)

	return &objectId, nil
}

func (m *MongoDbCollection[T]) ReplaceDocument(ctx context.Context, id primitive.ObjectID, document T) error {
	coll := m.Client.Database(m.Db).Collection(m.Col)
	opts := options.Update().SetUpsert(true)

	_, err := coll.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": document}, opts)

	if err != nil {
		m.Log.Error("Error inserting documents:", zap.String("database", m.Db), zap.String("collection", m.Col), zap.Error(err))
		return err
	}

	return nil
}

func (m *MongoDbCollection[T]) Count(filter bson.D) (int64, error) {
	return m.Client.Database(m.Db).Collection(m.Col).CountDocuments(context.TODO(), filter)
}

package database

import (
	"context"

	"github.com/Magpie-Monitor/magpie-monitor/pkg/mongodb"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

func NewMetadataDbMongoClient(log *zap.Logger) *mongo.Client {

	// envs.ValidateEnvs("Failed to connect to reportsdb", []string{
	// 	"REPORTSDB_USER",
	// 	"REPORTSDB_PASSWORD",
	// 	"REPORTSDB_HOST",
	// 	"REPORTSDB_PORT",
	// })

	// mongoDbUri := mongodb.GetMongoDbUri(
	// 	os.Getenv("REPORTSDB_USER"),
	// 	os.Getenv("REPORTSDB_PASSWORD"),
	// 	os.Getenv("REPORTSDB_HOST"),
	// 	os.Getenv("REPORTSDB_PORT"),
	// )

	mongoDbUri := mongodb.GetMongoDbUri(
		"mongo",
		"mongo",
		"localhost",
		"2222",
	)

	client, err := mongo.Connect(context.TODO(),
		options.Client().ApplyURI(mongoDbUri))

	// log.Info("Connected to reportsdb", zap.String("uri", mongoDbUri))

	if err != nil {
		panic(err)
		// log.Error("Failed to connect to reportsdb", zap.String("uri", mongoDbUri))
		// sh.Shutdown()
		// return nil
	}

	// lc.Append(
	// 	fx.Hook{
	// 		OnStop: func(ctx context.Context) error {
	// 			log.Info("Disconnecting from reportsdb", zap.String("uri", mongoDbUri))
	// 			return client.Disconnect(ctx)
	// 		},
	// 	},
	// )

	return client

}

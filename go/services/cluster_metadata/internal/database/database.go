package database

import (
	"context"

	"github.com/Magpie-Monitor/magpie-monitor/pkg/mongodb"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func NewMetadataDbMongoClient(lc fx.Lifecycle, sh fx.Shutdowner, log *zap.Logger) *mongo.Client {
	// envs.ValidateEnvs("Failed to connect to reportsdb", []string{
	// 	"METADATADB_USER",
	// 	"METADATADB_PASSWORD",
	// 	"METADATADB_HOST",
	// 	"METADATADB_PORT",
	// })

	// mongoDbUri := mongodb.GetMongoDbUri(
	// 	os.Getenv("METADATADB_USER"),
	// 	os.Getenv("METADATADB_PASSWORD"),
	// 	os.Getenv("METADATADB_HOST"),
	// 	os.Getenv("METADATADB_PORT"),
	// )

	mongoDbUri := mongodb.GetMongoDbUri(
		"mongo",
		"mongo",
		"localhost",
		"2222",
	)

	client, err := mongo.Connect(context.TODO(),
		options.Client().ApplyURI(mongoDbUri))

	log.Info("Connected to metadatadb", zap.String("uri", mongoDbUri))
	if err != nil {
		log.Error("Failed to connect to metadatadb", zap.String("uri", mongoDbUri))
		sh.Shutdown()
		return nil
	}

	lc.Append(
		fx.Hook{
			OnStop: func(ctx context.Context) error {
				log.Info("Disconnecting from metadatadb", zap.String("uri", mongoDbUri))
				return client.Disconnect(ctx)
			},
		},
	)

	return client
}

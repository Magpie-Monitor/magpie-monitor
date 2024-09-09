package database

import (
	"context"
	"fmt"
	"github.com/Magpie-Monitor/magpie-monitor/pkg/mongodb"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"os"
)

func NewReportsDbMongoClient(lc fx.Lifecycle, sh fx.Shutdowner, log *zap.Logger) *mongo.Client {

	mongoDbUri := mongodb.GetMongoDbUri(
		os.Getenv("REPORTSDB_USER"),
		os.Getenv("REPORTSDB_PASSWORD"),
		os.Getenv("REPORTSDB_HOST"),
		os.Getenv("REPORTSDB_PORT"),
	)

	client, err := mongo.Connect(context.TODO(),
		options.Client().ApplyURI(mongoDbUri))

	log.Info("Connecting to Mongodb", zap.String("uri", mongoDbUri))

	if err != nil {
		log.Error("Failed to connect to reportsdb", zap.String("uri", mongoDbUri))
		sh.Shutdown()
		fmt.Print()
		return nil
	}

	lc.Append(
		fx.Hook{
			OnStop: func(ctx context.Context) error {
				log.Info("Disconnecting from reportsdb", zap.String("uri", mongoDbUri))
				return client.Disconnect(ctx)
			},
		},
	)

	return client

}

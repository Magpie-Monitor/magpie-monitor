package database

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"os"
)

func getMongoDbUri() string {
	mongoUser := os.Getenv("MONGO_USER")
	mongoPassword := os.Getenv("MONGO_PASSWORD")
	mongoHost := os.Getenv("MONGO_HOST")
	mongoPort := os.Getenv("MONGO_PORT")

	return fmt.Sprintf("mongodb://%s:%s@%s:%s", mongoUser, mongoPassword, mongoHost, mongoPort)
}

func NewMongoDbClient(lc fx.Lifecycle, sh fx.Shutdowner, log *zap.Logger) *mongo.Client {

	mongoDbUri := getMongoDbUri()

	client, err := mongo.Connect(context.TODO(),
		options.Client().ApplyURI(mongoDbUri))

	log.Info("Connecting to Mongodb", zap.String("uri", mongoDbUri))

	if err != nil {
		log.Error("Failed to connect to mongodb", zap.String("uri", mongoDbUri))
		sh.Shutdown()
		fmt.Print()
		return nil
	}

	lc.Append(
		fx.Hook{
			OnStop: func(ctx context.Context) error {
				log.Info("Disconnecting from Mongodb", zap.String("uri", mongoDbUri))
				return client.Disconnect(ctx)
			},
		},
	)

	return client

}

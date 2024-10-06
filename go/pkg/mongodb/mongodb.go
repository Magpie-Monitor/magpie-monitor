package mongodb

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type MongoDbConnectionDetails struct {
	User     string
	Password string
	Host     string
	Port     string
}

func NewMongoDbClient(lc fx.Lifecycle, sh fx.Shutdowner, log *zap.Logger, conn *MongoDbConnectionDetails) *mongo.Client {
	mongoDbUri := GetMongoDbUri(
		conn.User,
		conn.Password,
		conn.Host,
		conn.Port,
	)

	client, err := mongo.Connect(context.TODO(),
		options.Client().ApplyURI(mongoDbUri))

	log.Info("Connected to mongodb", zap.String("uri", mongoDbUri))
	if err != nil {
		log.Error("Failed to connect to mongodb", zap.String("uri", mongoDbUri))
		sh.Shutdown()
		return nil
	}

	lc.Append(
		fx.Hook{
			OnStop: func(ctx context.Context) error {
				log.Info("Disconnecting from mongodb", zap.String("uri", mongoDbUri))
				return client.Disconnect(ctx)
			},
		},
	)

	return client
}

func GetMongoDbUri(user string, password string, host string, port string) string {

	return fmt.Sprintf("mongodb://%s:%s@%s:%s", user, password, host, port)
}

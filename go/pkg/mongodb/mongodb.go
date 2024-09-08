// package mongodb
package mongodb

import (
	// "context"
	"fmt"
	// "go.mongodb.org/mongo-driver/mongo"
	// "go.mongodb.org/mongo-driver/mongo/options"
	// "go.uber.org/fx"
	// "go.uber.org/zap"
)

func GetMongoDbUri(user string, password string, host string, port string) string {
	// mongoUser := os.Getenv("MONGO_USER")
	// mongoPassword := os.Getenv("MONGO_PASSWORD")
	// mongoHost := os.Getenv("MONGO_HOST")
	// mongoPort := os.Getenv("MONGO_PORT")

	return fmt.Sprintf("mongodb://%s:%s@%s:%s", user, password, host, port)
}

//
// type MongoDbClientResult struct {
// 	fx.Out
// 	Client mongo.Client
// }
//
// type Params struct {
// 	fx.In
// 	URI string `name:"mongo_uri"`
// 	log *zap.Logger
// 	lc  fx.Lifecycle
// }

// func NewMongoDbClient(p Params) (MongoDbClientResult, error) {
//
// 	client, err := mongo.Connect(context.TODO(),
// 		options.Client().ApplyURI(p.URI))
//
// 	p.log.Info("Connecting to Mongodb at", zap.String("uri", p.URI))
//
// 	if err != nil {
// 		p.log.Error("Failed to connect to Mongodb at", zap.String("uri", p.URI))
// 		return MongoDbClientResult{}, err
// 	}
//
// 	p.lc.Append(
// 		fx.Hook{
// 			OnStop: func(ctx context.Context) error {
// 				p.log.Info("Disconnecting from Mongodb", zap.String("uri", p.URI))
// 				return client.Disconnect(ctx)
// 			},
// 		},
// 	)
//
// 	return MongoDbClientResult{}, nil
//
// }

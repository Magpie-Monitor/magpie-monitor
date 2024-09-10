package mongodb

import (
	"fmt"
)

func GetMongoDbUri(user string, password string, host string, port string) string {

	return fmt.Sprintf("mongodb://%s:%s@%s:%s", user, password, host, port)
}

package elasticsearch

import (
	"github.com/elastic/go-elasticsearch/v8"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"os"
)

func NewElasticSearchLogsDbClient(sh fx.Shutdowner, log *zap.Logger) *elasticsearch.TypedClient {

	path := "/usr/local/share/ca/ca.crt"
	caCert, err := os.ReadFile(path)
	if err != nil {
		log.Error("Failed to read ca-certificates", zap.Error(err))
		sh.Shutdown()
		return nil
	}

	cfg := elasticsearch.Config{
		Addresses: []string{
			"https://elasticsearch:9200",
		},
		Username: "elastic",
		Password: "password",
		CACert:   caCert,
	}

	es, err := elasticsearch.NewTypedClient(cfg)
	if err != nil {
		log.Error("Failed to connect to elasic search", zap.Error(err))
		sh.Shutdown()
		return nil
	}

	log.Info("Connected to ElasticSearch!")

	return es
}

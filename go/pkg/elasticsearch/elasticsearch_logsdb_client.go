package elasticsearch

import (
	"fmt"
	"github.com/elastic/go-elasticsearch/v8"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"os"
)

func NewElasticSearchLogsDbClient(sh fx.Shutdowner, log *zap.Logger) *elasticsearch.TypedClient {

	certPath := os.Getenv("LOGSDB_CERT_PATH")
	caCert, err := os.ReadFile(certPath)
	if err != nil {
		log.Error("Failed to read logsdb certificate", zap.Error(err))
		sh.Shutdown()
		return nil
	}

	esPort := os.Getenv("LOGSDB_PORT")
	esHost := os.Getenv("LOGSDB_HOST")
	esUser := os.Getenv("LOGSDB_USER")
	esPassword := os.Getenv("LOGSDB_PASSWORD")

	cfg := elasticsearch.Config{
		Addresses: []string{
			fmt.Sprintf("https://%s:%s", esHost, esPort),
		},
		Username: esUser,
		Password: esPassword,
		CACert:   caCert,
	}

	es, err := elasticsearch.NewTypedClient(cfg)
	if err != nil {
		log.Error("Failed to connect logsdb", zap.Error(err))
		sh.Shutdown()
		return nil
	}

	log.Info("Connected to logsdb!")

	return es
}

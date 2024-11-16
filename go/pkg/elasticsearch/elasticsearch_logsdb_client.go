package elasticsearch

import (
	"fmt"
	"github.com/Magpie-Monitor/magpie-monitor/pkg/envs"
	"github.com/elastic/go-elasticsearch/v8"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"os"
)

func NewElasticSearchLogsDbClient(sh fx.Shutdowner, log *zap.Logger) *elasticsearch.TypedClient {

	envs.ValidateEnvs(
		"Failed to connect to logsdb.",
		[]string{
			"LOGSDB_CERT_PATH",
			"LOGSDB_PORT",
			"LOGSDB_HOST",
			"LOGSDB_USER",
			"LOGSDB_PASSWORD",
		},
	)

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

	log.Info(fmt.Sprintf("Connected to logsdb at %s", cfg.Addresses[0]))

	return es
}

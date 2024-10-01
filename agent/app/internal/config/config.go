package config

import (
	"flag"
	"log"
	"os"
)

type arrayFlags []string

func (i *arrayFlags) String() string {
	return "Array flags"
}

func (i *arrayFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}

type GlobalConfig struct {
	Mode                  string
	NodeName              string
	ClusterName           string
	ScrapeIntervalSeconds int
}

type RedisConfig struct {
	Url      string
	Password string
	Database int
}

type BrokerConfig struct {
	Url       string
	Username  string
	Password  string
	PodTopic  string
	NodeTopic string
	BatchSize int
}

type Config struct {
	Global             GlobalConfig
	Redis              RedisConfig
	Broker             BrokerConfig
	WatchedFiles       []string
	ExcludedNamespaces []string
}

func NewConfig() Config {
	nodeName := os.Getenv("NODE_NAME")

	mode := flag.String("scrape", "pods", "Mode in which log collector runs, either \"nodes\" to scrape nodes or \"pods\" to scrape pods.")
	clusterName := flag.String("clusterFriendlyName", "unknown", "Friendly name of your cluster, visible in Magpie Cloud.")

	scrapeIntervalSeconds := flag.Int("scrapeIntervalSeconds", 10, "Interval between scraping logs from files in \"nodes\" mode or pods in \"pods\" mode.")

	redisUrl := flag.String("redisUrl", "", "Redis URL in cluster DNS format, that is: service.namespace.svc.cluster.local:port")
	redisPassword := flag.String("redisPassword", "", "Password to Redis instance pointed by --redisUrl flag.")
	redisDatabase := flag.Int("redisDatabase", 0, "Database number for Redis instance.")

	remoteWriteBrokerUrl := flag.String("remoteWriteBrokerUrl", "localhost:9094", "URL of remote write broker.")
	remoteWritePodTopic := flag.String("remoteWritePodTopic", "pods", "Broker topic to which pod logs will be sent.")
	remoteWriteNodeTopic := flag.String("remoteWriteNodeTopic", "nodes", "Broker topic to which node logs will be sent.")
	remoteWriteBatchSize := flag.Int("remoteWriteBatchSize", 0, "Number of messages that are buffered and sent to broker in a single batch.")

	remoteWriteBrokerUsername := flag.String("remoteWriteBrokerUsername", "username", "SASL authentication broker username.")
	remoteWriteBrokerPassword := flag.String("remoteWriteBrokerPassword", "password", "SASL authentication broker password.")

	var watchedFiles arrayFlags
	flag.Var(&watchedFiles, "file", "Log files that are watched for log collector running in \"nodes\" mode.")

	var remoteWriteUrls arrayFlags
	flag.Var(&remoteWriteUrls, "remoteWriteUrl", "URL to which logs are pushed using remote write protocol.")

	var excludedNamespaces arrayFlags
	flag.Var(&excludedNamespaces, "excludedNamespace", "Namespace which is excluded from log scraping when agent runs in \"pods\" mode.")

	flag.Parse()

	log.Println("Agent configured to run in mode: ", *mode)
	if *mode == "nodes" {
		log.Println("Node agent running on node: ", nodeName)
	}
	log.Println("Redis url: ", *redisUrl)

	return Config{
		Global: GlobalConfig{
			Mode:                  *mode,
			NodeName:              nodeName,
			ClusterName:           *clusterName,
			ScrapeIntervalSeconds: *scrapeIntervalSeconds,
		},
		Redis: RedisConfig{
			Url:      *redisUrl,
			Password: *redisPassword,
			Database: *redisDatabase,
		},
		Broker: BrokerConfig{
			Url:       *remoteWriteBrokerUrl,
			Username:  *remoteWriteBrokerUsername,
			Password:  *remoteWriteBrokerPassword,
			PodTopic:  *remoteWritePodTopic,
			NodeTopic: *remoteWriteNodeTopic,
			BatchSize: *remoteWriteBatchSize,
		},
		WatchedFiles:       watchedFiles,
		ExcludedNamespaces: excludedNamespaces,
	}
}

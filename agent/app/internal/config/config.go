package config

import (
	"flag"
	"log"
	"os"

	nodeData "github.com/Magpie-Monitor/magpie-monitor/agent/internal/agent/node/data"
	"github.com/Magpie-Monitor/magpie-monitor/agent/internal/agent/pods/data"
)

type arrayFlags []string

func (i *arrayFlags) String() string {
	return "Array flags"
}

func (i *arrayFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}

type Channels struct {
	ApplicationLogsChannel     chan data.Chunk
	ApplicationMetadataChannel chan data.ClusterState
	NodeLogsChannel            chan nodeData.Chunk
	NodeMetadataChannel        chan nodeData.NodeState
}

func NewChannels() *Channels {
	return &Channels{
		ApplicationLogsChannel:     make(chan data.Chunk, 10),
		ApplicationMetadataChannel: make(chan data.ClusterState, 10),
		NodeLogsChannel:            make(chan nodeData.Chunk, 10),
		NodeMetadataChannel:        make(chan nodeData.NodeState, 10),
	}
}

func (c *Channels) Close() {
	close(c.ApplicationLogsChannel)
	close(c.ApplicationMetadataChannel)
	close(c.NodeLogsChannel)
	close(c.NodeMetadataChannel)
}

type GlobalConfig struct {
	Mode                               string
	NodeName                           string
	ClusterName                        string
	LogScrapeIntervalSeconds           int
	MetadataScrapeIntervalSeconds      int
	PodMetadataRemoteWriteUrl          string
	NodeMetadataRemoteWriteUrl         string
	ClusterMetadataServiceClientSecret string
	RunningLocally                     bool
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

	runningLocally := flag.Bool("runningLocally", false, "Determines whether an agent is running locally in a dev environment.")

	mode := flag.String("scrape", "pods", "Mode in which log collector runs, either \"nodes\" to scrape nodes or \"pods\" to scrape pods.")
	clusterName := flag.String("clusterFriendlyName", "unknown", "Friendly name of your cluster, visible in Magpie Cloud.")

	logScrapeIntervalSeconds := flag.Int("logScrapeIntervalSeconds", 10, "Interval between scraping logs from files in \"nodes\" mode or pods in \"pods\" mode.")
	metadataScrapeIntervalSeconds := flag.Int("metadataScrapeIntervalSeconds", 10, "Interval between scraping nodes metadata in \"nodes\" mode or cluster metadata in \"pods\".")

	redisUrl := flag.String("redisUrl", "", "Redis URL in cluster DNS format, that is: service.namespace.svc.cluster.local:port")
	redisPassword := flag.String("redisPassword", "", "Password to Redis instance pointed by --redisUrl flag.")
	redisDatabase := flag.Int("redisDatabase", 0, "Database number for Redis instance.")

	remoteWriteBrokerUrl := flag.String("remoteWriteBrokerUrl", "localhost:9094", "URL of remote write broker.")
	remoteWriteApplicationTopic := flag.String("remoteWriteApplicationTopic", "pods", "Broker topic to which pod logs will be sent.")
	remoteWriteNodeTopic := flag.String("remoteWriteNodeTopic", "nodes", "Broker topic to which node logs will be sent.")
	remoteWriteBatchSize := flag.Int("remoteWriteBatchSize", 0, "Number of messages that are buffered and sent to broker in a single batch.")

	remoteWriteBrokerUsername := flag.String("remoteWriteBrokerUsername", "username", "SASL authentication broker username.")
	remoteWriteBrokerPassword := flag.String("remoteWriteBrokerPassword", "password", "SASL authentication broker password.")

	podRemoteWriteMetadataUrl := flag.String("podRemoteWriteMetadataUrl", "", "URL for cluster metadata remote write.")
	nodeRemoteWriteMetadataUrl := flag.String("nodeRemoteWriteMetadataUrl", "", "URL for node metadata remote write.")

	clusterMetadataServiceClientSecret := flag.String("clusterMetadataServiceClientSecret", "", "Client secret for metadata service remote write.")

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
			Mode:                               *mode,
			NodeName:                           nodeName,
			ClusterName:                        *clusterName,
			LogScrapeIntervalSeconds:           *logScrapeIntervalSeconds,
			MetadataScrapeIntervalSeconds:      *metadataScrapeIntervalSeconds,
			PodMetadataRemoteWriteUrl:          *podRemoteWriteMetadataUrl,
			NodeMetadataRemoteWriteUrl:         *nodeRemoteWriteMetadataUrl,
			ClusterMetadataServiceClientSecret: *clusterMetadataServiceClientSecret,
			RunningLocally:                     *runningLocally,
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
			PodTopic:  *remoteWriteApplicationTopic,
			NodeTopic: *remoteWriteNodeTopic,
			BatchSize: *remoteWriteBatchSize,
		},
		WatchedFiles:       watchedFiles,
		ExcludedNamespaces: excludedNamespaces,
	}
}

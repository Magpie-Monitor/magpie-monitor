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
	ApplicationMetadataChannel chan data.ApplicationState
	NodeLogsChannel            chan nodeData.Chunk
	NodeMetadataChannel        chan nodeData.NodeState
}

func NewChannels() *Channels {
	return &Channels{
		ApplicationLogsChannel:     make(chan data.Chunk, 10),
		ApplicationMetadataChannel: make(chan data.ApplicationState, 10),
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
	ClusterId                          string
	LogScrapeIntervalSeconds           int
	MetadataScrapeIntervalSeconds      int
	PodMetadataRemoteWriteUrl          string
	NodeMetadataRemoteWriteUrl         string
	ClusterMetadataServiceClientSecret string
	RunningMode                        string
	MaxPodPacketSizeBytes              int
	MaxContainerPacketSizeBytes        int
	NodePacketSizeBytes                int
}

type RedisConfig struct {
	Url      string
	Password string
	Database int
}

type BrokerConfig struct {
	Url                      string
	Username                 string
	Password                 string
	ApplicationTopic         string
	NodeTopic                string
	ApplicationMetadataTopic string
	NodeMetadataTopic        string
	BatchSize                int
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

	runningMode := flag.String("runningMode", "remote", "Determines whether an agent is running locally in a dev environment. Set to \"local\" when running locally and \"remote\" when not.")

	mode := flag.String("scrape", "pods", "Mode in which log collector runs, either \"nodes\" to scrape nodes or \"pods\" to scrape pods.")
	clusterId := flag.String("clusterFriendlyName", "unknown", "Friendly name of your cluster, visible in Magpie Cloud.")

	logScrapeIntervalSeconds := flag.Int("logScrapeIntervalSeconds", 10, "Interval between scraping logs from files in \"nodes\" mode or pods in \"pods\" mode.")
	metadataScrapeIntervalSeconds := flag.Int("metadataScrapeIntervalSeconds", 10, "Interval between scraping nodes metadata in \"nodes\" mode or cluster metadata in \"pods\".")

	redisUrl := flag.String("redisUrl", "", "Redis URL in cluster DNS format, that is: service.namespace.svc.cluster.local:port")
	redisPassword := flag.String("redisPassword", "", "Password to Redis instance pointed by --redisUrl flag.")
	redisDatabase := flag.Int("redisDatabase", 0, "Database number for Redis instance.")

	remoteWriteBrokerUrl := flag.String("remoteWriteBrokerUrl", "localhost:9094", "URL of remote write broker.")

	remoteWriteApplicationTopic := flag.String("remoteWriteApplicationTopic", "pods", "Broker topic to which pod logs will be sent.")
	remoteWriteNodeTopic := flag.String("remoteWriteNodeTopic", "nodes", "Broker topic to which node logs will be sent.")

	remoteWriteApplicationMetadataTopic := flag.String("remoteWriteApplicationMetadataTopic", "application_metadata", "Broker topic to which pod logs will be sent.")
	remoteWriteNodeMetadataTopic := flag.String("remoteWriteNodeMetadataTopic", "node_metadata", "Broker topic to which node logs will be sent.")

	remoteWriteBatchSize := flag.Int("remoteWriteBatchSize", 0, "Number of messages that are buffered and sent to broker in a single batch.")

	remoteWriteBrokerUsername := flag.String("remoteWriteBrokerUsername", "username", "SASL authentication broker username.")
	remoteWriteBrokerPassword := flag.String("remoteWriteBrokerPassword", "password", "SASL authentication broker password.")

	podRemoteWriteMetadataUrl := flag.String("podRemoteWriteMetadataUrl", "", "URL for cluster metadata remote write.")
	nodeRemoteWriteMetadataUrl := flag.String("nodeRemoteWriteMetadataUrl", "", "URL for node metadata remote write.")

	clusterMetadataServiceClientSecret := flag.String("clusterMetadataServiceClientSecret", "", "Client secret for metadata service remote write.")

	maxPodPacketSizeBytes := flag.Int("maxPodPacketSizeBytes", 5_000, "Maximal size of a single pod packet sent to Kafka in bytes.")
	maxContainerPacketSizeBytes := flag.Int("maxContainerPacketSizeBytes", 1_000, "Maximal size of a single container packet in bytes.")

	nodePacketSizeBytes := flag.Int("nodePacketSizeBytes", 1_000, "Size of the log packet for node logs in bytes.")

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
			ClusterId:                          *clusterId,
			LogScrapeIntervalSeconds:           *logScrapeIntervalSeconds,
			MetadataScrapeIntervalSeconds:      *metadataScrapeIntervalSeconds,
			PodMetadataRemoteWriteUrl:          *podRemoteWriteMetadataUrl,
			NodeMetadataRemoteWriteUrl:         *nodeRemoteWriteMetadataUrl,
			ClusterMetadataServiceClientSecret: *clusterMetadataServiceClientSecret,
			RunningMode:                        *runningMode,
			MaxPodPacketSizeBytes:              *maxPodPacketSizeBytes,
			MaxContainerPacketSizeBytes:        *maxContainerPacketSizeBytes,
			NodePacketSizeBytes:                *nodePacketSizeBytes,
		},
		Redis: RedisConfig{
			Url:      *redisUrl,
			Password: *redisPassword,
			Database: *redisDatabase,
		},
		Broker: BrokerConfig{
			Url:                      *remoteWriteBrokerUrl,
			Username:                 *remoteWriteBrokerUsername,
			Password:                 *remoteWriteBrokerPassword,
			ApplicationTopic:         *remoteWriteApplicationTopic,
			NodeTopic:                *remoteWriteNodeTopic,
			ApplicationMetadataTopic: *remoteWriteApplicationMetadataTopic,
			NodeMetadataTopic:        *remoteWriteNodeMetadataTopic,
			BatchSize:                *remoteWriteBatchSize,
		},
		WatchedFiles:       watchedFiles,
		ExcludedNamespaces: excludedNamespaces,
	}
}

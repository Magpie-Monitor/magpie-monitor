package wrapper

import (
	"fmt"
	agent2 "github.com/Magpie-Monitor/magpie-monitor/agent/internal/agent/node/agent"
	data2 "github.com/Magpie-Monitor/magpie-monitor/agent/internal/agent/node/data"
	"github.com/Magpie-Monitor/magpie-monitor/agent/internal/agent/pods/agent"
	"github.com/Magpie-Monitor/magpie-monitor/agent/internal/agent/pods/data"
	"github.com/Magpie-Monitor/magpie-monitor/agent/internal/config"
	"github.com/Magpie-Monitor/magpie-monitor/agent/internal/remote_write"
	"k8s.io/apimachinery/pkg/util/json"
	"log"
)

type AgentWrapper struct {
	config             config.Config
	podWriter          remote_write.RemoteWriter
	nodeWriter         remote_write.RemoteWriter
	podMetadataWriter  remote_write.RemoteWriter
	nodeMetadataWriter remote_write.RemoteWriter
}

func NewAgentWrapper(config config.Config) AgentWrapper {
	return AgentWrapper{
		config: config,
		podWriter: remote_write.NewStreamWriter(config.Broker.Url, config.Broker.PodTopic, config.Broker.Username,
			config.Broker.Password, config.Broker.BatchSize),
		nodeWriter: remote_write.NewStreamWriter(config.Broker.Url, config.Broker.NodeTopic, config.Broker.Username,
			config.Broker.Password, config.Broker.BatchSize),
		podMetadataWriter:  remote_write.NewMetadataWriter(config.Global.PodMetadataRemoteWriteUrl),
		nodeMetadataWriter: remote_write.NewMetadataWriter(config.Global.NodeMetadataRemoteWriteUrl),
	}
}

func (a *AgentWrapper) Start() {
	mode := a.config.Global.Mode
	if mode == "nodes" {
		log.Println("Watched files: ", a.config.WatchedFiles)
		if len(a.config.WatchedFiles) == 0 {
			panic("Node agent doesn't have any files configured, please point watched files in the cfg.")
		}
		a.startNodeAgent()
	} else if mode == "pods" {
		a.startPodAgent()
	} else {
		panic(fmt.Sprintf("Mode: %s not supported", mode))
	}
}

func (a *AgentWrapper) startNodeAgent() {
	logChannel := make(chan data2.Chunk)
	metadataChannel := make(chan data2.NodeState)

	nodeAgent := agent2.NewReader(a.config.Global.NodeName, a.config.WatchedFiles, a.config.Global.LogScrapeIntervalSeconds,
		a.config.Global.MetadataScrapeIntervalSeconds, logChannel, metadataChannel, a.config.Redis.Url, a.config.Redis.Password, a.config.Redis.Database)
	go nodeAgent.Start()

	go a.watchNodeLogsChannel(logChannel)
	a.watchNodeMetadataChannel(metadataChannel)
}

func (a *AgentWrapper) watchNodeLogsChannel(logChannel chan data2.Chunk) {
	for metadata := range logChannel {
		log.Println("Collected node logs: ", metadata)
		a.writeChunk(metadata, a.nodeWriter)
	}
}

func (a *AgentWrapper) watchNodeMetadataChannel(metadataChannel chan data2.NodeState) {
	for metadata := range metadataChannel {
		log.Println("Collected node metadata: ", metadata)
		a.writeChunk(metadata, a.nodeMetadataWriter)
	}
}

func (a *AgentWrapper) startPodAgent() {
	logChannel := make(chan data.Chunk)
	metadataChannel := make(chan data.ClusterState)

	podAgent := agent.NewAgent(a.config.ExcludedNamespaces, a.config.Global.LogScrapeIntervalSeconds, a.config.Global.MetadataScrapeIntervalSeconds, logChannel, metadataChannel)

	go podAgent.Start()
	go a.watchClusterLogsChannel(logChannel)
	a.watchClusterMetadataChannel(metadataChannel)
}

func (a *AgentWrapper) watchClusterLogsChannel(logChannel chan data.Chunk) {
	for chunk := range logChannel {
		log.Println("Collected pod chunk: ", chunk)
		a.writeChunk(chunk, a.podWriter)
	}
}

func (a *AgentWrapper) watchClusterMetadataChannel(metadataChannel chan data.ClusterState) {
	for metadata := range metadataChannel {
		log.Println("Collected cluster metadata: ", metadata)
		a.writeChunk(metadata, a.podMetadataWriter)
	}
}

func (a *AgentWrapper) writeChunk(chunk interface{}, writer remote_write.RemoteWriter) {
	jsonChunk, err := json.Marshal(chunk)
	if err != nil {
		log.Println("Error converting chunk to JSON: ", err)
	} else {
		writer.Write(string(jsonChunk))
	}
}

package wrapper

import (
	"fmt"
	"github.com/Magpie-Monitor/magpie-monitor/agent/internal/agent/node"
	"github.com/Magpie-Monitor/magpie-monitor/agent/internal/agent/pods"
	"github.com/Magpie-Monitor/magpie-monitor/agent/internal/config"
	"github.com/Magpie-Monitor/magpie-monitor/agent/internal/remoteWrite"
	"k8s.io/apimachinery/pkg/util/json"
	"log"
)

type AgentWrapper struct {
	config     config.Config
	podWriter  remoteWrite.RemoteWriter
	nodeWriter remoteWrite.RemoteWriter
}

func NewAgentWrapper(config config.Config) AgentWrapper {
	return AgentWrapper{
		config: config,
		podWriter: remoteWrite.NewStreamWriter(config.Broker.Url, config.Broker.PodTopic, config.Broker.Username,
			config.Broker.Password, config.Broker.BatchSize),
		nodeWriter: remoteWrite.NewStreamWriter(config.Broker.Url, config.Broker.NodeTopic, config.Broker.Username,
			config.Broker.Password, config.Broker.BatchSize),
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
	logChannel := make(chan node.Chunk)

	agent := node.NewReader(a.config.Global.NodeName, a.config.WatchedFiles, a.config.Global.ScrapeIntervalSeconds,
		logChannel, a.config.Redis.Url, a.config.Redis.Password, a.config.Redis.Database)
	go agent.WatchFiles()

	for chunk := range logChannel {
		log.Println("Collected node chunk: ", chunk)
		a.writeChunk(chunk, a.nodeWriter)
	}
}

func (a *AgentWrapper) startPodAgent() {
	logChannel := make(chan pods.Chunk)
	agent := pods.NewAgent(a.config.ExcludedNamespaces, a.config.Global.ScrapeIntervalSeconds, logChannel)
	go agent.Start()

	for chunk := range logChannel {
		log.Println("Collected pod chunk: ", chunk)
		a.writeChunk(chunk, a.podWriter)
	}
}

func (a *AgentWrapper) writeChunk(chunk any, writer remoteWrite.RemoteWriter) {
	jsonChunk, err := json.Marshal(chunk)
	if err != nil {
		log.Println("Error converting chunk to JSON: ", err)
	} else {
		writer.Write(string(jsonChunk))
	}
}

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
	config       config.Config
	remoteWriter remoteWrite.RemoteWriter
}

func NewAgentWrapper(config config.Config) AgentWrapper {
	return AgentWrapper{
		config:       config,
		remoteWriter: remoteWrite.NewRemoteWriter(config.RemoteWriteUrls, config.RemoteWriteRetryInterval, config.RemoteWriteMaxRetries),
	}
}

func (a *AgentWrapper) Start() {
	if a.config.Mode == "nodes" {
		log.Println("Watched files: ", a.config.WatchedFiles)
		if len(a.config.WatchedFiles) == 0 {
			panic("Node agent doesn't have any files configured, please point watched files in the cfg.")
		}
		a.startNodeAgent()
	} else if a.config.Mode == "pods" {
		a.startPodAgent()
	} else {
		panic(fmt.Sprintf("Mode: %s not supported", a.config.Mode))
	}
}

func (a *AgentWrapper) startNodeAgent() {
	logChannel := make(chan node.Chunk)

	agent := node.NewReader(a.config.WatchedFiles, a.config.ScrapeInterval, nil, logChannel, a.config.RedisUrl, a.config.RedisPassword)
	go agent.WatchFiles()

	for chunk := range logChannel {
		log.Println("Collected node chunk: ", chunk)
		a.writeChunk(chunk)
	}
}

func (a *AgentWrapper) startPodAgent() {
	logChannel := make(chan pods.Chunk)
	agent := pods.NewAgent(a.config.ExcludedNamespaces, a.config.ScrapeInterval, logChannel)
	go agent.Start()

	for chunk := range logChannel {
		log.Println("Collected pod chunk: ", chunk)
		a.writeChunk(chunk)
	}
}

func (a *AgentWrapper) writeChunk(chunk any) {
	jsonChunk, err := json.Marshal(chunk)
	if err != nil {
		log.Println("Error converting chunk to JSON: ", err)
	} else {
		a.remoteWriter.Write(string(jsonChunk))
	}
}

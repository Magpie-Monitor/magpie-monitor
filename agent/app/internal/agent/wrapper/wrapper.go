package wrapper

import (
	"fmt"
	"k8s.io/apimachinery/pkg/util/json"
	"log"
	"logather/internal/agent/entity"
	"logather/internal/agent/node"
	"logather/internal/agent/pods"
	"logather/internal/config"
	"logather/internal/remoteWrite"
)

type AgentWrapper struct {
	config       config.Config
	remoteWriter remoteWrite.RemoteWriter
}

func NewAgentWrapper(config config.Config) AgentWrapper {
	return AgentWrapper{config: config, remoteWriter: remoteWrite.NewRemoteWriter(config.RemoteWriteUrls)}
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
	logChannel := make(chan entity.Chunk)

	agent := node.NewReader(a.config.WatchedFiles, nil, logChannel, a.config.RedisUrl)
	go agent.WatchFiles()

	for chunk := range logChannel {
		fmt.Println(chunk)
		a.writeChunk(chunk)
	}
}

func (a *AgentWrapper) startPodAgent() {
	logChannel := make(chan entity.Chunk)
	agent := pods.NewAgent(nil, 30, logChannel)
	go agent.Start()

	for chunk := range logChannel {
		fmt.Println(chunk)
		a.writeChunk(chunk)
	}
}

func (a *AgentWrapper) writeChunk(chunk entity.Chunk) {
	jsonChunk, err := json.Marshal(chunk)
	if err != nil {
		log.Println("Error converting chunk to JSON: ", err)
	} else {
		a.remoteWriter.Write(string(jsonChunk))
	}
}

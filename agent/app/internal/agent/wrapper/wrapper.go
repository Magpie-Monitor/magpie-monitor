package wrapper

import (
	"fmt"
	"log"

	"github.com/Magpie-Monitor/magpie-monitor/agent/internal/agent/collector"
	nodeAgent "github.com/Magpie-Monitor/magpie-monitor/agent/internal/agent/node/agent"
	"github.com/Magpie-Monitor/magpie-monitor/agent/internal/agent/pods/agent"
	"github.com/Magpie-Monitor/magpie-monitor/agent/internal/config"
)

type AgentWrapper struct {
	config    *config.Config
	channels  *config.Channels
	collector collector.DataCollector
}

func NewAgentWrapper(cfg *config.Config) AgentWrapper {
	channels := config.NewChannels()
	return AgentWrapper{
		config:    cfg,
		channels:  channels,
		collector: collector.NewDataCollector(cfg, channels),
	}
}

func (a *AgentWrapper) Start() {
	defer a.channels.Close()

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
	nodesAgent := nodeAgent.NewReader(a.config, a.channels.NodeLogsChannel, a.channels.NodeMetadataChannel)
	go nodesAgent.Start()
	a.collector.CollectNodes()
}

func (a *AgentWrapper) startPodAgent() {
	podAgent := agent.NewAgent(a.config, a.channels.ApplicationLogsChannel, a.channels.ApplicationMetadataChannel)
	go podAgent.Start()
	a.collector.CollectCluster()
}

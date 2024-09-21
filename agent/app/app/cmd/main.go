package main

import (
	"github.com/Magpie-Monitor/magpie-monitor/agent/internal/agent/wrapper"
	"github.com/Magpie-Monitor/magpie-monitor/agent/internal/config"
)

func main() {
	cfg := config.Config{Mode: "pods", RemoteWriteBrokerUrl: "127.0.0.1:9094", RemoteWritePodTopic: "pods", RemoteWriteNodeTopic: "nodes"}
	//cfg := config.NewConfig()
	agentWrapper := wrapper.NewAgentWrapper(cfg)
	agentWrapper.Start()
}

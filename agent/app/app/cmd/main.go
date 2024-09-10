package main

import (
	"github.com/Magpie-Monitor/magpie-monitor/agent/internal/agent/wrapper"
	"github.com/Magpie-Monitor/magpie-monitor/agent/internal/config"
)

func main() {
	//cfg := config.Config{Mode: "pods", RemoteWriteUrls: []string{"http://localhost:8080/api/v1/ingest"}}
	cfg := config.GetConfig()
	agentWrapper := wrapper.NewAgentWrapper(cfg)
	agentWrapper.Start()
}

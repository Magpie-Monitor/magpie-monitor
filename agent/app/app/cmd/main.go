package main

import (
	"github.com/Magpie-Monitor/magpie-monitor/agent/internal/agent/wrapper"
	"github.com/Magpie-Monitor/magpie-monitor/agent/internal/config"
)

func main() {
	//cfg := config.Config{Global: config.GlobalConfig{Mode: "nodes", MetadataScrapeIntervalSeconds: 10}, WatchedFiles: []string{"test.log"}}
	cfg := config.NewConfig()
	agentWrapper := wrapper.NewAgentWrapper(cfg)
	agentWrapper.Start()
}

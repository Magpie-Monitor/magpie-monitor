package main

import (
	"github.com/Magpie-Monitor/magpie-monitor/agent/internal/agent/wrapper"
	"github.com/Magpie-Monitor/magpie-monitor/agent/internal/config"
)

func main() {
	cfg := config.NewConfig()
	agentWrapper := wrapper.NewAgentWrapper(cfg)
	agentWrapper.Start()
}

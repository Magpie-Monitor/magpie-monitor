package main

import (
	"logather/internal/agent/wrapper"
	"logather/internal/config"
)

func main() {
	cfg := config.GetConfig()
	agentWrapper := wrapper.NewAgentWrapper(cfg)
	agentWrapper.Start()
}

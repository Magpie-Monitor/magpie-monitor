package main

import (
	"logather/internal/agent/wrapper"
	"logather/internal/config"
)

func main() {
	//cfg := config.GetConfig()
	cfg := config.Config{Mode: "pods"}
	agentWrapper := wrapper.NewAgentWrapper(cfg)
	agentWrapper.Start()
}

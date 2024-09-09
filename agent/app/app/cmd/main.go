package main

import (
	"logather/internal/agent/wrapper"
	"logather/internal/config"
)

// cfg := config.Config{Mode: "pods", RemoteWriteUrls: []string{"http://localhost:8080/api/v1/ingest"}}
func main() {
	cfg := config.GetConfig()
	agentWrapper := wrapper.NewAgentWrapper(cfg)
	agentWrapper.Start()
}

package main

import (
	"flag"
	"fmt"
	"k8s.io/client-go/util/homedir"
	"log"
	"logather/internal/agent/node"
	"logather/internal/agent/pods"
	"logather/internal/transformer"
	"os/exec"
	"path/filepath"
)

func main() {
	//RunNodeAgentDemo()
	mode := *flag.String("scrape", "nodes", "Mode in which log collector runs, either \"nodes\" to scrape nodes or \"pods\" to scrape pods.")
	flag.Parse()

	if mode == "nodes" {
		RunNodeAgentDemo()
	} else if mode == "pods" {
		RunPodAgentDemo()
	} else {
		panic(fmt.Sprintf("Mode: %s not supported", mode))
	}
}

func RunNodeAgentDemo() {
	go func() {
		cmd := exec.Command("bash", "-c", "./generate-logs.sh")
		if err := cmd.Run(); err != nil {
			panic(err)
		}
	}()

	//watchedFiles := []string{"test.log", "test2.log"}
	watchedFiles := []string{"/logs/btmp"}
	RunNodeAgent(watchedFiles)
}

func RunPodAgentDemo() {
	//dir := "/var/log"
	dir := ""
	RunPodAgent(dir)
}

func RunNodeAgent(watchedFiles []string) {
	c := make(chan node.IncrementalFetch)

	t1 := transformer.DummyTransformer{}

	transformers := []transformer.Transformer{t1}

	agent := node.NewReader(watchedFiles, transformers, nil, c)
	agent.WatchFiles()

	// prints gathered log lines
	for elem := range c {
		log.Println(elem.Content)
	}
}

func RunPodAgent(logsDir string) {
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}

	agent := pods.NewAgent(*kubeconfig, nil, 2, logsDir)
	agent.Start()
}

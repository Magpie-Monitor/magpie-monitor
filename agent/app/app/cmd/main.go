package main

import (
	"flag"
	"fmt"
	"k8s.io/client-go/util/homedir"
	"log"
	"logather/internal/agent/node"
	"logather/internal/agent/pods"
	"logather/internal/transformer"
	"path/filepath"
)

type arrayFlags []string

func (i *arrayFlags) String() string {
	return i.String()
}

func (i *arrayFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}

func main() {
	mode := *flag.String("scrape", "nodes", "Mode in which log collector runs, either \"nodes\" to scrape nodes or \"pods\" to scrape pods.")

	var watchedFiles arrayFlags
	flag.Var(&watchedFiles, "file", "Log files that are watched for log collector running in \"nodes\" mode.")

	flag.Parse()

	log.Println("Starting agent in mode: ", mode)

	if mode == "nodes" {
		log.Println("Watched files: ", watchedFiles)
		if len(watchedFiles) == 0 {
			panic("Node agent doesn't have any files configured, please point watched files in the config.")
		}
		RunNodeAgent(watchedFiles)
	} else if mode == "pods" {
		RunPodAgent("/logs")
	} else {
		panic(fmt.Sprintf("Mode: %s not supported", mode))
	}
}

//func RunNodeAgentDemo(watchedFiles []string) {
//	go func() {
//		cmd := exec.Command("bash", "-c", "./generate-logs.sh")
//		if err := cmd.Run(); err != nil {
//			panic(err)
//		}
//	}()
//
//	//watchedFiles := []string{"/logs/btmp"}
//	RunNodeAgent(watchedFiles)
//}

//func RunPodAgentDemo() {
//	//dir := "/var/log"
//	dir := ""
//	RunPodAgent(dir)
//}

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

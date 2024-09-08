package main

import (
	"flag"
	"fmt"
	"k8s.io/client-go/util/homedir"
	"log"
	"logather/internal/agent/node"
	"logather/internal/agent/pods"
	"logather/internal/remoteWrite"
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
	redisUrl := *flag.String("redisUrl", "", "Redis URL in cluster DNS format, that is: service.namespace.svc.cluster.local:port")

	var watchedFiles arrayFlags
	flag.Var(&watchedFiles, "file", "Log files that are watched for log collector running in \"nodes\" mode.")

	flag.Parse()

	log.Println("Starting agent in mode: ", mode)
	log.Println("Redis url: ", redisUrl)

	if mode == "nodes" {
		log.Println("Watched files: ", watchedFiles)
		if len(watchedFiles) == 0 {
			panic("Node agent doesn't have any files configured, please point watched files in the config.")
		}
		RunNodeAgent(watchedFiles, redisUrl)
	} else if mode == "pods" {
		RunPodAgent("./logs")
	} else {
		panic(fmt.Sprintf("Mode: %s not supported", mode))
	}
}

func RunNodeAgent(watchedFiles []string, redisUrl string) {
	logChannel := make(chan node.IncrementalFetch)

	agent := node.NewReader(watchedFiles, nil, logChannel, redisUrl)
	agent.WatchFiles()

	var buffer = make(map[string]node.IncrementalFetch)
	writer := remoteWrite.NewRemoteWriter([]string{"http://localhost:8080/api/v1/ingest"}) // TODO - revisit

	for incrementalFetch := range logChannel {
		bufferedIncrementalFetch, ok := buffer[incrementalFetch.Dir]
		if ok {
			bufferedContent := bufferedIncrementalFetch.Content
			incrementalFetch.Content = bufferedContent + incrementalFetch.Content
			buffer[incrementalFetch.Dir] = incrementalFetch
		} else {
			buffer[incrementalFetch.Dir] = incrementalFetch
		}

		if len(incrementalFetch.Content) > 100 {
			writer.Write(incrementalFetch)
			incrementalFetch.Content = ""
			buffer[incrementalFetch.Dir] = incrementalFetch
		}
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

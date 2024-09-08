package wrapper

import (
	"flag"
	"fmt"
	"k8s.io/client-go/util/homedir"
	"log"
	"logather/internal/agent/node"
	"logather/internal/agent/pods"
	"logather/internal/config"
	"logather/internal/remoteWrite"
	"path/filepath"
)

type AgentWrapper struct {
	config config.Config
}

func NewAgentWrapper(config config.Config) AgentWrapper {
	return AgentWrapper{config: config}
}

func (a *AgentWrapper) Start() {
	if a.config.Mode == "nodes" {
		log.Println("Watched files: ", a.config.WatchedFiles)
		if len(a.config.WatchedFiles) == 0 {
			panic("Node agent doesn't have any files configured, please point watched files in the cfg.")
		}
		a.startNodeAgent()
	} else if a.config.Mode == "pods" {
		a.startPodAgent()
	} else {
		panic(fmt.Sprintf("Mode: %s not supported", a.config.Mode))
	}
}

func (a *AgentWrapper) startNodeAgent() {
	logChannel := make(chan node.IncrementalFetch)

	agent := node.NewReader(a.config.WatchedFiles, nil, logChannel, a.config.RedisUrl)
	agent.WatchFiles()

	var buffer = make(map[string]node.IncrementalFetch)
	writer := remoteWrite.NewRemoteWriter(a.config.RemoteWriteUrls)

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

func (a *AgentWrapper) startPodAgent() {
	//config, err := rest.InClusterConfig() // TODO - https://github.com/kubernetes/client-go/tree/master/examples/in-cluster-client-configuration
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}

	agent := pods.NewAgent(*kubeconfig, nil, 2)
	agent.Start()
}

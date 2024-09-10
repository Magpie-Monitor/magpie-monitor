package pods

import (
	"context"
	"flag"
	"fmt"
	"github.com/Magpie-Monitor/magpie-monitor/agent/internal/agent/entity"
	v2 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"log"
	"path/filepath"
	"slices"
	"time"
)

type Agent struct {
	excludedNamespaces        []string
	includedNamespaces        []string
	collectionIntervalSeconds int
	collectionDirectory       string
	client                    *kubernetes.Clientset
	readTimestamps            map[string]int64
	results                   chan entity.Chunk
}

func NewAgent(excludedNamespaces []string, collectionIntervalSeconds int, results chan entity.Chunk) *Agent {
	return &Agent{
		excludedNamespaces:        excludedNamespaces,
		collectionIntervalSeconds: collectionIntervalSeconds,
		readTimestamps:            make(map[string]int64),
		results:                   results,
	}
}

func (a *Agent) Start() {
	a.authenticate()
	a.fetchNamespaces()
	a.gatherLogs()
}

func (a *Agent) authenticate() {
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}

	var config *rest.Config
	if len(*kubeconfig) > 0 {
		c, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
		if err != nil {
			log.Println("Failed to create kubernetes API client from kubeconfig")
			panic(err.Error())
		}
		config = c
	} else {
		c, err := rest.InClusterConfig()
		if err != nil {
			log.Println("Failed to create kubernetes API client from ServiceAccount token")
			panic(err.Error())
		}
		config = c
	}

	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	a.client = client
}

func (a *Agent) fetchNamespaces() {
	a.includedNamespaces = make([]string, 0)
	namespaces, err := a.client.CoreV1().
		Namespaces().
		List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(fmt.Sprintf("Error fetching namespaces: %s", err.Error()))
	}

	for _, namespace := range namespaces.Items {
		if !slices.Contains(a.excludedNamespaces, namespace.Name) {
			a.includedNamespaces = append(a.includedNamespaces, namespace.Name)
		}
	}
}

func (a *Agent) gatherLogs() {
	for {
		for _, namespace := range a.includedNamespaces {
			log.Println("Fetching logs for namespace: ", namespace)

			deployments, err := a.client.AppsV1().
				Deployments(namespace).
				List(context.TODO(), metav1.ListOptions{})
			if err != nil {
				log.Println("Error fetching Deployments: ", err)
				log.Println("Skipping iteration")
			} else {
				a.fetchDeploymentLogsSinceSeconds(namespace, deployments.Items)
			}

			statefulSets, err := a.client.AppsV1().StatefulSets(namespace).List(context.TODO(), metav1.ListOptions{})
			if err != nil {
				log.Println("Error fetching StatefulSets: ", err)
				log.Println("Skipping iteration")
			} else {
				a.fetchStatefulSetLogsSinceSeconds(namespace, statefulSets.Items)
			}

			daemonSets, err := a.client.AppsV1().DaemonSets(namespace).List(context.TODO(), metav1.ListOptions{})
			if err != nil {
				log.Println("Error fetching DaemonSets: ", err)
				log.Println("Skipping iteration")
			} else {
				a.fetchDaemonSetLogsSinceSeconds(namespace, daemonSets.Items)
			}
		}

		log.Println("Sleeping for: ", a.collectionIntervalSeconds, " seconds")
		time.Sleep(time.Duration(a.collectionIntervalSeconds) * time.Second)
	}
}

func (a *Agent) calculateLastReadTimeDiff(name string) int64 {
	var readDiff = int64(a.collectionIntervalSeconds)
	now := time.Now()
	unix := now.Unix()
	previousRead, ok := a.readTimestamps[name]
	if ok {
		readDiff = unix - previousRead
	}
	return readDiff
}

func (a *Agent) setReadTimestamp(name string) {
	now := time.Now()
	a.readTimestamps[name] = now.Unix()
}

func (a *Agent) fetchDeploymentLogsSinceSeconds(namespace string, deployments []v2.Deployment) {
	for _, deployment := range deployments {
		name := deployment.Name
		selectors := deployment.Spec.Selector

		log.Println("Deployment: ", name)

		readDiff := a.calculateLastReadTimeDiff(name)
		logs := a.fetchLogsSinceSeconds(selectors, namespace, &readDiff)
		a.setReadTimestamp(name)

		if len(logs) > 0 {
			a.results <- entity.Chunk{Kind: "Deployment", Name: name, Namespace: namespace, Content: logs}
		}
	}
}

func (a *Agent) fetchStatefulSetLogsSinceSeconds(namespace string, statefulSets []v2.StatefulSet) {
	log.Println("Fetching logs from StatefulSets")
	for _, statefulSet := range statefulSets {
		name := statefulSet.Name
		selectors := statefulSet.Spec.Selector

		log.Println("StatefulSet: ", name)
		readDiff := a.calculateLastReadTimeDiff(name)
		logs := a.fetchLogsSinceSeconds(selectors, namespace, &readDiff)
		a.setReadTimestamp(name)

		if len(logs) > 0 {
			a.results <- entity.Chunk{Kind: "StatefulSet", Name: name, Namespace: namespace, Content: logs}
		}
	}
}

func (a *Agent) fetchDaemonSetLogsSinceSeconds(namespace string, daemonSets []v2.DaemonSet) {
	log.Println("Fetching logs from DaemonSets")
	for _, daemonSet := range daemonSets {
		name := daemonSet.Name
		selectors := daemonSet.Spec.Selector

		log.Println("DaemonSet: ", name)
		readDiff := a.calculateLastReadTimeDiff(name)
		logs := a.fetchLogsSinceSeconds(selectors, namespace, &readDiff)
		a.setReadTimestamp(name)

		if len(logs) > 0 {
			a.results <- entity.Chunk{Kind: "StatefulSet", Name: name, Namespace: namespace, Content: logs}
		}
	}
}

func (a *Agent) fetchLogsSinceSeconds(selector *metav1.LabelSelector, namespace string, sinceSeconds *int64) string {
	var result string

	pods, _ := a.client.CoreV1().
		Pods(namespace).
		List(context.TODO(),
			metav1.ListOptions{
				LabelSelector: labels.Set(selector.MatchLabels).String(),
			})
	for _, pod := range pods.Items {
		log.Println("Fetching logs for pod: ", pod.Name)

		for _, container := range pod.Spec.Containers {
			log.Println("Fetching logs for container: ", container.Name)
			// TODO - explore since time and log streaming
			// TODO - explore stdout + stderr
			logs := a.client.CoreV1().Pods(namespace).GetLogs(pod.Name, &v1.PodLogOptions{Container: container.Name, SinceSeconds: sinceSeconds}).Do(context.TODO())
			rawLogs, _ := logs.Raw()
			result += string(rawLogs)
		}
	}

	return result
}

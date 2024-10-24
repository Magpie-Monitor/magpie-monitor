package agent

import (
	"context"
	"flag"
	"fmt"
	"log"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/Magpie-Monitor/magpie-monitor/agent/internal/agent/pods/data"
	"github.com/Magpie-Monitor/magpie-monitor/agent/internal/config"
	v2 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

type Agent struct {
	clusterId                         string
	excludedNamespaces                []string
	includedNamespaces                []string
	logCollectionIntervalSeconds      int
	metadataCollectionIntervalSeconds int
	client                            *kubernetes.Clientset
	readTimestamps                    map[string]int64
	readTimes                         map[string]time.Time
	results                           chan<- data.Chunk
	metadata                          chan<- data.ClusterState
	runningMode                       string
}

func NewAgent(cfg *config.Config, logsChan chan<- data.Chunk, metadataChan chan<- data.ClusterState) *Agent {
	return &Agent{
		clusterId:                         cfg.Global.ClusterId,
		excludedNamespaces:                cfg.ExcludedNamespaces,
		logCollectionIntervalSeconds:      cfg.Global.LogScrapeIntervalSeconds,
		metadataCollectionIntervalSeconds: cfg.Global.MetadataScrapeIntervalSeconds,
		readTimestamps:                    make(map[string]int64),
		readTimes:                         make(map[string]time.Time),
		results:                           logsChan,
		metadata:                          metadataChan,
		runningMode:                       cfg.Global.RunningMode,
	}
}

func (a *Agent) Start() {
	a.authenticate()
	a.fetchNamespaces()
	go a.gatherLogs()
	go a.gatherClusterMetadata()
}

func (a *Agent) authenticate() {
	var config *rest.Config

	if a.runningMode == "local" {
		var kubeconfig *string

		if home := homedir.HomeDir(); home != "" {
			kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
		} else {
			kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
		}

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
			a.fetchLogsForNamespace(namespace)
		}

		log.Println("Sleeping for: ", a.logCollectionIntervalSeconds, " seconds")
		time.Sleep(time.Duration(a.logCollectionIntervalSeconds) * time.Second)
	}
}

func (a *Agent) fetchLogsForNamespace(namespace string) {
	log.Println("Fetching logs for namespace: ", namespace)

	deployments, err := a.client.AppsV1().
		Deployments(namespace).
		List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Println("Error fetching Deployments: ", err)
		log.Println("Skipping iteration")
	} else {
		a.fetchDeploymentLogsSinceTime(namespace, deployments.Items)
	}

	statefulSets, err := a.client.AppsV1().StatefulSets(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Println("Error fetching StatefulSets: ", err)
		log.Println("Skipping iteration")
	} else {
		a.fetchStatefulSetLogsSinceTime(namespace, statefulSets.Items)
	}

	daemonSets, err := a.client.AppsV1().DaemonSets(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Println("Error fetching DaemonSets: ", err)
		log.Println("Skipping iteration")
	} else {
		a.fetchDaemonSetLogsSinceTime(namespace, daemonSets.Items)
	}
}

func (a *Agent) fetchDeploymentLogsSinceTime(namespace string, deployments []v2.Deployment) {
	for _, deployment := range deployments {
		selectors := deployment.Spec.Selector
		logs := a.fetchPodLogsSinceTime(selectors, namespace)
		a.sendResult(data.Deployment, deployment.Name, namespace, logs)
	}
}

func (a *Agent) fetchStatefulSetLogsSinceTime(namespace string, statefulSets []v2.StatefulSet) {
	log.Println("Fetching logs from StatefulSets")
	for _, statefulSet := range statefulSets {
		selectors := statefulSet.Spec.Selector
		logs := a.fetchPodLogsSinceTime(selectors, namespace)
		a.sendResult(data.StatefulSet, statefulSet.Name, namespace, logs)
	}
}

func (a *Agent) fetchDaemonSetLogsSinceTime(namespace string, daemonSets []v2.DaemonSet) {
	log.Println("Fetching logs from DaemonSets")
	for _, daemonSet := range daemonSets {
		selectors := daemonSet.Spec.Selector
		logs := a.fetchPodLogsSinceTime(selectors, namespace)
		a.sendResult(data.DaemonSet, daemonSet.Name, namespace, logs)
	}
}

func (a *Agent) fetchPodLogsSinceTime(selector *metav1.LabelSelector, namespace string) []data.Pod {
	res := make([]data.Pod, 0)

	// TODO - error handling, abstraction over K8S API
	pods, _ := a.client.CoreV1().
		Pods(namespace).
		List(
			context.TODO(),
			metav1.ListOptions{LabelSelector: labels.Set(selector.MatchLabels).String()},
		)
	for _, pod := range pods.Items {
		log.Println("Fetching logs for pod: ", pod.Name)

		containers := make([]data.Container, 0, len(pod.Spec.Containers))
		for _, container := range pod.Spec.Containers {
			log.Println("Fetching logs for container: ", container.Name)
			c := a.fetchContainerLogsSinceTime(&container, pod.Name, namespace)
			containers = append(containers, c)
		}

		res = append(res, data.Pod{Name: pod.Name, Containers: containers})
	}

	return res
}

func (a *Agent) fetchContainerLogsSinceTime(container *v1.Container, podName, namespace string) data.Container {
	sinceTime := a.getReadTimestamp(podName, container.Name)

	// Sleep till all the logs from current second arrive.
	// Precision of the logs API is within seconds,
	// so to not fetch logs twice, we have to gather all logs
	// from the ongoing second. Then, we cut off the logs from the
	// following second, so they are fetched in next iteration.

	time.Sleep(time.Duration(999999999 - sinceTime.Nanosecond()))

	beforeTs := time.Now().UnixNano()
	// TODO - abstraction over this part
	logs := a.client.CoreV1().
		Pods(namespace).
		GetLogs(
			podName,
			&v1.PodLogOptions{
				Container:  container.Name,
				SinceTime:  &metav1.Time{Time: sinceTime},
				Timestamps: true},
		).Do(context.TODO())
	afterTs := time.Now().UnixNano()

	if logs.Error() != nil {
		log.Println("Error fetching logs for Pod: ", podName, " container: ", container.Name)
		return data.Container{}
	}

	// Subtract request time from the next fetch time,
	// incorporating logs that were emitted at the time of making GetLogs() call.
	now := afterTs - (afterTs - beforeTs)
	a.setReadTimestamp(podName, container.Name, now)

	rawLogs, err := logs.Raw()
	if err != nil {
		log.Println("Failed to fetch raw logs for container: ", container.Name)
		return data.Container{}
	}

	deduplicatedLogs, err := a.deduplicate(string(rawLogs))
	if err != nil {
		log.Println("Failed to fetch container: ", container.Name, " logs")
		return data.Container{}
	}

	return data.Container{Name: container.Name, Image: container.Image, Content: deduplicatedLogs}
}

func (a *Agent) setReadTimestamp(podName, containerName string, timestampUnixMicro int64) {
	a.readTimestamps[a.getTimestampKey(podName, containerName)] = timestampUnixMicro
}

func (a *Agent) getReadTimestamp(podName, containerName string) time.Time {
	val, ok := a.readTimestamps[a.getTimestampKey(podName, containerName)]
	if ok {
		return time.Unix(0, val)
	}
	return time.Now()
}

func (a *Agent) getTimestampKey(podName, containerName string) string {
	return fmt.Sprintf("%s-%s", podName, containerName)
}

func (a *Agent) sendResult(kind data.ApplicationKind, name, namespace string, pods []data.Pod) {
	a.results <- data.Chunk{
		ClusterId:     a.clusterId,
		Kind:          kind,
		CollectedAtMs: time.Now().UnixMilli(),
		Name:          name,
		Namespace:     namespace,
		Pods:          pods,
	}
}

// API Server returns logs with second precision, while in our case nanosecond precision is preferred.
// Therefore, on every fetch we wait till the end of the ongoing second and then fetch the logs.
// As a result, if the throughput of logs is high, we get log lines from the next second, which we don't need.
// Deduplication removes those excessive lines, as they will be fetched in the next iteration.
func (a *Agent) deduplicate(logs string) (string, error) {
	split := strings.Split(strings.Trim(logs, "\n"), "\n")

	// Too little log lines for duplication to occur, no need to deduplicate.
	if len(split) < 2 {
		return logs, nil
	}

	lastLogSeconds, err := a.getSecondFromLogTimestamp(split[len(split)-1])
	if err != nil {
		log.Println("Deduplication failed on last log timestamp extraction")
		return "", err
	}

	for i := len(split) - 2; i >= 0; i-- {
		logLine := split[i]
		currentLogSeconds, err := a.getSecondFromLogTimestamp(logLine)
		if err != nil {
			log.Println("Deduplication failed on current log timestamp extraction")
			return "", err
		}

		// Logs from the next second were fetched, they have to be removed.
		// They will be fetched in the next iteration.
		if lastLogSeconds > currentLogSeconds {
			return strings.Join(split[:i+1], "\n"), nil
		}
	}

	return logs, nil
}

// Returns second value for a log line beginning with RFC3339 timestamp,
// ex. "2006-01-02T15:04:05.123123Z $S0ME_LOG" should return 5.
func (a *Agent) getSecondFromLogTimestamp(logLine string) (int, error) {
	timestamp := strings.Split(logLine, " ")[0]
	parsedTime, err := time.Parse(time.RFC3339, timestamp)
	if err != nil {
		log.Println("Error parsing log timestamp: ", timestamp)
		return 0, err
	}
	return parsedTime.Second(), nil
}

func (a *Agent) gatherClusterMetadata() {
	// TODO - create kubernetes API client
	for {
		state := data.NewClusterState(a.clusterId)
		for _, namespace := range a.includedNamespaces {
			deployments, err := a.client.AppsV1().
				Deployments(namespace).
				List(context.TODO(), metav1.ListOptions{})
			if err != nil {
				log.Println("Error fetching Deployments: ", err)
			} else {
				state.AppendDeployments(&deployments.Items)
			}

			statefulSets, err := a.client.AppsV1().StatefulSets(namespace).List(context.TODO(), metav1.ListOptions{})
			if err != nil {
				log.Println("Error fetching StatefulSets: ", err)
			} else {
				state.AppendStatefulSets(&statefulSets.Items)
			}

			daemonSets, err := a.client.AppsV1().DaemonSets(namespace).List(context.TODO(), metav1.ListOptions{})
			if err != nil {
				log.Println("Error fetching DaemonSets: ", err)
			} else {
				state.AppendDaemonSets(&daemonSets.Items)
			}

			state.SetTimestamp()

		}

		a.metadata <- state
		time.Sleep(time.Duration(a.metadataCollectionIntervalSeconds) * time.Second)
	}
}

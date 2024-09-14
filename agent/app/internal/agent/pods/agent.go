package pods

import (
	"context"
	"flag"
	"fmt"
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
	"strings"
	"time"
)

type ObjectKind string

const (
	Deployment  ObjectKind = "Deployment"
	StatefulSet ObjectKind = "StatefulSet"
	DaemonSet   ObjectKind = "DaemonSet"
)

func (o ObjectKind) String() string {
	switch o {
	case Deployment:
		return "Deployment"
	case StatefulSet:
		return "StatefulSet"
	case DaemonSet:
		return "DaemonSet"
	default:
		return "unknown"
	}
}

type Agent struct {
	clusterName               string
	excludedNamespaces        []string
	includedNamespaces        []string
	collectionIntervalSeconds int
	collectionDirectory       string
	client                    *kubernetes.Clientset
	readTimestamps            map[string]int64
	readTimes                 map[string]time.Time
	results                   chan PodChunk
}

func NewAgent(excludedNamespaces []string, collectionIntervalSeconds int, results chan PodChunk) *Agent {
	return &Agent{
		excludedNamespaces:        excludedNamespaces,
		collectionIntervalSeconds: collectionIntervalSeconds,
		readTimestamps:            make(map[string]int64),
		readTimes:                 make(map[string]time.Time),
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
	//a.includedNamespaces = append(a.includedNamespaces, "mock-ns")

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

		log.Println("Sleeping for: ", a.collectionIntervalSeconds, " seconds")
		time.Sleep(time.Duration(a.collectionIntervalSeconds) * time.Second)
	}
}

func (a *Agent) fetchDeploymentLogsSinceTime(namespace string, deployments []v2.Deployment) {
	for _, deployment := range deployments {
		selectors := deployment.Spec.Selector
		logs := a.fetchPodLogsSinceTime(selectors, namespace)
		a.sendResult(Deployment, deployment.Name, namespace, logs)
	}
}

func (a *Agent) fetchStatefulSetLogsSinceTime(namespace string, statefulSets []v2.StatefulSet) {
	log.Println("Fetching logs from StatefulSets")
	for _, statefulSet := range statefulSets {
		selectors := statefulSet.Spec.Selector
		logs := a.fetchPodLogsSinceTime(selectors, namespace)
		a.sendResult(StatefulSet, statefulSet.Name, namespace, logs)
	}
}

func (a *Agent) fetchDaemonSetLogsSinceTime(namespace string, daemonSets []v2.DaemonSet) {
	log.Println("Fetching logs from DaemonSets")
	for _, daemonSet := range daemonSets {
		selectors := daemonSet.Spec.Selector
		logs := a.fetchPodLogsSinceTime(selectors, namespace)
		a.sendResult(DaemonSet, daemonSet.Name, namespace, logs)
	}
}

func (a *Agent) fetchPodLogsSinceTime(selector *metav1.LabelSelector, namespace string) []Pod {
	res := make([]Pod, 0)

	pods, _ := a.client.CoreV1().
		Pods(namespace).
		List(context.TODO(),
			metav1.ListOptions{
				LabelSelector: labels.Set(selector.MatchLabels).String(),
			})
	for _, pod := range pods.Items {
		log.Println("Fetching logs for pod: ", pod.Name)

		containers := make([]Container, 0)
		for _, container := range pod.Spec.Containers {
			log.Println("Fetching logs for container: ", container.Name)
			c := a.fetchContainerLogsSinceTime(container, pod.Name, namespace)
			containers = append(containers, c)
		}

		res = append(res, Pod{Name: pod.Name, Containers: containers})
	}

	return res
}

func (a *Agent) fetchContainerLogsSinceTime(container v1.Container, podName, namespace string) Container {
	log.Println("Fetching logs for container: ", container.Name)

	sinceTime := a.getReadTimestamp(podName, container.Name)

	// Sleep till all the logs from current second arrive.
	// Precision of the logs API is within seconds,
	// so to not fetch logs twice, we have to gather all logs
	// from the ongoing second. Then, we cut off the logs from the
	// following second, so they are fetched in next iteration.
	time.Sleep(time.Duration(999999999 - sinceTime.Nanosecond()))

	before := time.Now().UnixNano()
	logs := a.client.CoreV1().
		Pods(namespace).
		GetLogs(podName,
			&v1.PodLogOptions{Container: container.Name, SinceTime: &metav1.Time{Time: sinceTime}, Timestamps: true}).
		Do(context.TODO())
	after := time.Now().UnixNano()

	if logs.Error() != nil {
		log.Println("Error fetching logs for Pod: ", podName, " container: ", container.Name)
		return Container{}
	} else {
		// Subtract request time from the next fetch time,
		// incorporating logs that were emitted at the time of making GetLogs() call.
		now := after - (after - before)
		a.setReadTimestamp(podName, container.Name, now)

		rawLogs, err := logs.Raw()
		if err != nil {
			log.Println("Failed to fetch raw logs for container: ", container.Name)
			return Container{}
		}

		deduplicatedLogs, err := a.deduplicate(string(rawLogs))
		if err != nil {
			log.Println("Failed to fetch container: ", container.Name, " logs")
			return Container{}
		}

		return Container{Name: container.Name, Image: container.Image, Content: deduplicatedLogs}
	}
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

func (a *Agent) sendResult(kind ObjectKind, name, namespace string, pods []Pod) {
	a.results <- PodChunk{
		Cluster:   a.clusterName,
		Kind:      kind.String(),
		Timestamp: time.Now().UnixNano(),
		Name:      name,
		Namespace: namespace,
		Pods:      pods,
	}
}

func (a *Agent) deduplicate(logs string) (string, error) {
	split := strings.Split(logs, "\n")
	// Last line is always empty.
	split = split[:len(split)-1]

	if len(split) >= 2 {
		lastLogSecond, err := a.getSecondFromLogTimestamp(split[len(split)-2])
		if err != nil {
			log.Println("Deduplication failed on last log timestamp extraction")
			return "", err
		}

		for i := len(split) - 3; i >= 0; i-- {
			logLine := split[i]
			currentLogSecond, err := a.getSecondFromLogTimestamp(logLine)
			if err != nil {
				log.Println("Deduplication failed on current log timestamp extraction")
				return "", err
			}

			// Logs from the next second were fetched, they have to be removed.
			// They will be fetched in the next iteration.
			if lastLogSecond > currentLogSecond {
				return strings.Join(split[:i+1], "\n"), nil
			}
		}
	}

	return logs, nil
}

func (a *Agent) getSecondFromLogTimestamp(log string) (int, error) {
	timestamp := strings.Split(log, " ")[0]
	parsedTime, err := time.Parse(time.RFC3339, timestamp)
	if err != nil {
		return 0, err
	}
	return parsedTime.Second(), nil
}

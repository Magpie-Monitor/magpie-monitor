package agent

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Magpie-Monitor/magpie-monitor/agent/internal/agent/pods/data"
	"github.com/Magpie-Monitor/magpie-monitor/agent/internal/config"
	v2 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	kube "github.com/Magpie-Monitor/magpie-monitor/agent/pkg/kubernetes"
)

type Agent struct {
	clusterId                         string
	excludedNamespaces                []string
	includedNamespaces                []string
	logCollectionIntervalSeconds      int
	metadataCollectionIntervalSeconds int
	readTimestamps                    map[string]int64
	readTimes                         map[string]time.Time
	results                           chan<- data.Chunk
	metadata                          chan<- data.ApplicationState
	runningMode                       string
	maxPodPacketSizeBytes             int
	maxContainerPacketSizeBytes       int
	kubernetesClient                  kube.KubernetesApiClient
}

func NewAgent(cfg *config.Config, logsChan chan<- data.Chunk, metadataChan chan<- data.ApplicationState) *Agent {
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
		maxPodPacketSizeBytes:             cfg.Global.MaxPodPacketSizeBytes,
		maxContainerPacketSizeBytes:       cfg.Global.MaxContainerPacketSizeBytes,
		kubernetesClient:                  kube.NewKubernetesApiClient(cfg.Global.RunningMode),
	}
}

func (a *Agent) Start() {
	a.fetchNamespaces()
	go a.gatherLogs()
	go a.gatherClusterMetadata()
}

func (a *Agent) fetchNamespaces() {
	namespaces, err := a.kubernetesClient.GetNamespaces(a.excludedNamespaces)

	if err != nil {
		panic(fmt.Sprintf("Error fetching namespaces: %s", err.Error()))
	}

	a.includedNamespaces = namespaces
}

func (a *Agent) gatherLogs() {
	for {
		for _, namespace := range a.includedNamespaces {
			log.Printf("Fetching logs for namespace %s", namespace)
			a.fetchLogsForNamespace(namespace)
		}

		log.Printf("Logs gathered, sleeping for %d seconds", a.logCollectionIntervalSeconds)
		time.Sleep(time.Duration(a.logCollectionIntervalSeconds) * time.Second)
	}
}

func (a *Agent) fetchLogsForNamespace(namespace string) {
	err := a.fetchLogsForDeployments(namespace)
	if err != nil {
		log.Printf("Error fetching deployments, err=%s", err.Error())
	}

	err = a.fetchLogsForStatefulSets(namespace)
	if err != nil {
		log.Printf("Error fetching statefulSets, err=%s", err.Error())
	}

	err = a.fetchLogsForDaemonSets(namespace)
	if err != nil {
		log.Printf("Error fetching daemonSets, err=%s", err.Error())
	}
}

func (a *Agent) fetchLogsForDeployments(namespace string) error {
	deployments, err := a.kubernetesClient.GetDeployments(namespace)
	if err != nil {
		return err
	}

	a.fetchDeploymentLogsSinceTime(namespace, deployments)

	return nil
}

func (a *Agent) fetchLogsForDaemonSets(namespace string) error {
	statefulSets, err := a.kubernetesClient.GetStatefulSets(namespace)
	if err != nil {
		return err
	}

	a.fetchStatefulSetLogsSinceTime(namespace, statefulSets)

	return nil
}

func (a *Agent) fetchLogsForStatefulSets(namespace string) error {
	daemonSets, err := a.kubernetesClient.GetDaemonSets(namespace)
	if err != nil {
		return err
	}

	a.fetchDaemonSetLogsSinceTime(namespace, daemonSets)

	return nil
}

func (a *Agent) fetchDeploymentLogsSinceTime(namespace string, deployments []v2.Deployment) {
	for _, deployment := range deployments {
		logPackets, err := a.fetchPodLogsSinceTime(deployment.Spec.Selector, namespace)
		if err != nil {
			log.Printf("Error fetching deployment logs for deployment=%s, namespace=%s err=%s", deployment.Name, namespace, err.Error())
			continue
		}

		for _, packet := range logPackets {
			a.sendResult(data.Deployment, deployment.Name, namespace, packet)
		}
	}
}

func (a *Agent) fetchStatefulSetLogsSinceTime(namespace string, statefulSets []v2.StatefulSet) {
	for _, statefulSet := range statefulSets {
		logPackets, err := a.fetchPodLogsSinceTime(statefulSet.Spec.Selector, namespace)
		if err != nil {
			log.Printf("Error fetching statefulSet logs for statefulSet=%s, namespace=%s err=%s", statefulSet.Name, namespace, err.Error())
			continue
		}

		for _, packet := range logPackets {
			a.sendResult(data.StatefulSet, statefulSet.Name, namespace, packet)
		}
	}
}

func (a *Agent) fetchDaemonSetLogsSinceTime(namespace string, daemonSets []v2.DaemonSet) {
	for _, daemonSet := range daemonSets {
		logPackets, err := a.fetchPodLogsSinceTime(daemonSet.Spec.Selector, namespace)
		if err != nil {
			log.Printf("Error fetching daemonSet logs for daemonSet=%s, namespace=%s err=%s", daemonSet.Name, namespace, err.Error())
			continue
		}

		for _, packet := range logPackets {
			a.sendResult(data.DaemonSet, daemonSet.Name, namespace, packet)
		}
	}
}

func (a *Agent) fetchPodLogsSinceTime(selector *metav1.LabelSelector, namespace string) ([][]data.Pod, error) {
	pods, err := a.kubernetesClient.GetPods(selector, namespace)
	if err != nil {
		return nil, err
	}

	return a.getPodLogsPackets(pods), nil
}

func (a *Agent) getPodLogsPackets(pods []v1.Pod) [][]data.Pod {
	podPackets := make([][]data.Pod, 0)

	for _, pod := range pods {
		containers := a.fetchContainerLogsForPod(pod)
		podPacket := a.splitPodContainerLogsIntoPackets(pod.Name, containers)
		podPackets = append(podPackets, podPacket)
	}

	return podPackets
}

func (a *Agent) splitPodContainerLogsIntoPackets(podName string, containers []data.Container) []data.Pod {
	var (
		podPacket                []data.Pod
		currentPacketLen         = 0
		containerPacketsTotalLen = len(containers) * a.maxContainerPacketSizeBytes
		currentPacketFreeBytes   = a.maxPodPacketSizeBytes
	)

	// Pod fits into the packet.
	if currentPacketFreeBytes >= containerPacketsTotalLen {
		currentPacketLen += containerPacketsTotalLen
		podPacket = append(podPacket, data.Pod{Name: podName, Containers: containers})
		return podPacket
	}

	containerPackets := a.splitContainerIntoPackets(containers)

	for _, packet := range containerPackets {
		podPacket = append(podPacket, data.Pod{Name: podName, Containers: packet})
	}

	return podPacket
}

func (a *Agent) splitContainerIntoPackets(containers []data.Container) [][]data.Container {
	var (
		containerPackets       [][]data.Container
		containerPacket        []data.Container
		currentPacketFreeBytes = a.maxPodPacketSizeBytes
		currentPacketLen       = 0
		containerPacketLen     = a.maxContainerPacketSizeBytes
	)

	for _, container := range containers {
		if currentPacketLen+containerPacketLen > currentPacketFreeBytes {
			containerPackets = append(containerPackets, containerPacket)
			containerPacket = make([]data.Container, 0)
			containerPacket = append(containerPacket, container)
			currentPacketLen = containerPacketLen
			continue
		}

		containerPacket = append(containerPacket, container)
		currentPacketLen += containerPacketLen
	}

	return append(containerPackets, containerPacket)
}

func (a *Agent) fetchContainerLogsForPod(pod v1.Pod) []data.Container {
	containers := make([]data.Container, 0, len(pod.Spec.Containers))

	for _, container := range pod.Spec.Containers {
		c := a.fetchContainerLogsSinceTime(&container, pod.Name, pod.Namespace)
		containers = append(containers, c...)
	}

	return containers
}

func (a *Agent) fetchContainerLogsSinceTime(container *v1.Container, podName, namespace string) []data.Container {
	sinceTime := a.getReadTimestamp(podName, container.Name)

	// Sleep till all the logs from current second arrive.
	// Precision of the logs API is within seconds,
	// so to not fetch logs twice, we have to gather all logs
	// from the ongoing second. Then, we cut off the logs from the
	// following second, so they are fetched in next iteration.

	time.Sleep(time.Duration(999999999 - sinceTime.Nanosecond()))

	beforeTs := time.Now().UnixNano()
	logs, err := a.kubernetesClient.GetContainerLogsSinceTime(podName, container.Name, namespace, sinceTime, true)
	afterTs := time.Now().UnixNano()

	if err != nil {
		log.Println("Error fetching logs for Pod: ", podName, " container: ", container.Name)
		return nil
	}

	// Subtract request time from the next fetch time,
	// incorporating logs that were emitted at the time of making GetLogs() call.
	now := afterTs - (afterTs - beforeTs)
	a.setReadTimestamp(podName, container.Name, now)

	deduplicatedLogs, err := a.deduplicate(logs)
	if err != nil {
		log.Println("Failed to fetch container: ", container.Name, " logs")
		return nil
	}

	return a.splitLogsIntoContainerPackets(container.Name, container.Image, deduplicatedLogs)
}

func (a *Agent) splitLogsIntoContainerPackets(containerName, containerImage, logs string) []data.Container {
	if len(logs) == 0 {
		return []data.Container{{Name: containerName, Image: containerImage, Content: logs}}
	}

	var (
		logPackets         []string
		currentPacket      string
		maxPacketSizeBytes = a.maxContainerPacketSizeBytes
		currentPacketLen   = 0
	)

	logLines := strings.Split(logs, "\n")
	for _, line := range logLines {
		if currentPacketLen < maxPacketSizeBytes {
			currentPacket += "\n" + line
			currentPacketLen += len(line)
			continue
		}

		logPackets = append(logPackets, currentPacket)
		currentPacket = line
		currentPacketLen = len(line)
	}

	logPackets = append(logPackets, currentPacket)

	var containers []data.Container
	for _, packet := range logPackets {
		containers = append(containers, data.Container{Name: containerName, Image: containerImage, Content: packet})
	}

	return containers
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
	for {
		applicationState := data.NewApplicationState(a.clusterId)

		for _, namespace := range a.includedNamespaces {
			a.getApplicationStateForNamespace(namespace, &applicationState)
		}

		a.metadata <- applicationState
		time.Sleep(time.Duration(a.metadataCollectionIntervalSeconds) * time.Second)
	}
}

func (a *Agent) getApplicationStateForNamespace(namespace string, applicationState *data.ApplicationState) {
	err := a.appendDeploymentState(namespace, applicationState)
	if err != nil {
		log.Printf("Error fetching deployment state, err=%s", err.Error())
	}

	err = a.appendStatefulSetState(namespace, applicationState)
	if err != nil {
		log.Printf("Error fetching statefulSet state, err=%s", err.Error())
	}

	err = a.appendDaemonSetState(namespace, applicationState)
	if err != nil {
		log.Printf("Error fetching daemonSet state, err=%s", err.Error())
	}

	applicationState.SetTimestamp()
}

func (a *Agent) appendDeploymentState(namespace string, applicationState *data.ApplicationState) error {
	deployments, err := a.kubernetesClient.GetDeployments(namespace)

	if err != nil {
		return err
	}

	applicationState.AppendDeployments(&deployments)

	return nil
}

func (a *Agent) appendStatefulSetState(namespace string, applicationState *data.ApplicationState) error {
	statefulSets, err := a.kubernetesClient.GetStatefulSets(namespace)

	if err != nil {
		return err
	}

	applicationState.AppendStatefulSets(&statefulSets)

	return nil

}

func (a *Agent) appendDaemonSetState(namespace string, applicationState *data.ApplicationState) error {
	daemonSets, err := a.kubernetesClient.GetDaemonSets(namespace)

	if err != nil {
		return err
	}

	applicationState.AppendDaemonSets(&daemonSets)

	return nil
}

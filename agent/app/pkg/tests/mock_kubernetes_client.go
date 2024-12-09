package tests

import (
	"fmt"
	"time"

	"github.com/Magpie-Monitor/magpie-monitor/agent/pkg/kubernetes"
	v2 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	LOG_KEY_FMT = "%s-%s-%s"
)

type Logs struct {
	Namespace string
	Pod       string
	Container string
	Logs      string
}

func NewMockKubernetesApiClient(
	namespaces []string,
	deployments []v2.Deployment,
	statefulSets []v2.StatefulSet,
	daemonSets []v2.DaemonSet,
	pods []v1.Pod,
	logs []Logs,
) kubernetes.KubernetesApiClient {
	return &MockClient{
		namespaces:   namespaces,
		deployments:  deployments,
		statefulSets: statefulSets,
		daemonSets:   daemonSets,
		pods:         pods,
		logs:         getLogs(logs),
	}
}

type MockClient struct {
	namespaces   []string
	deployments  []v2.Deployment
	statefulSets []v2.StatefulSet
	daemonSets   []v2.DaemonSet
	pods         []v1.Pod
	logs         map[string]Logs
}

func getLogs(logs []Logs) map[string]Logs {
	res := make(map[string]Logs, len(logs))

	for _, log := range logs {
		res[getLogKey(log.Namespace, log.Pod, log.Container)] = log
	}

	return res
}

func getLogKey(namespace, pod, container string) string {
	return fmt.Sprintf(LOG_KEY_FMT, namespace, pod, container)
}

func (c *MockClient) GetDeployments(namespace string) ([]v2.Deployment, error) {
	return c.deployments, nil
}

func (c *MockClient) GetStatefulSets(namespace string) ([]v2.StatefulSet, error) {
	return c.statefulSets, nil
}

func (c *MockClient) GetDaemonSets(namespace string) ([]v2.DaemonSet, error) {
	return c.daemonSets, nil
}

func (c *MockClient) GetNamespaces(excludedNamespaces []string) ([]string, error) {
	return c.namespaces, nil
}

func (c *MockClient) GetPods(selector *metav1.LabelSelector, namespace string) ([]v1.Pod, error) {
	return c.pods, nil
}

func (c *MockClient) GetContainerLogsSinceTime(podName, containerName, namespace string, sinceTime time.Time, timestamps bool) (string, error) {
	return c.logs[getLogKey(namespace, podName, containerName)].Logs, nil
}

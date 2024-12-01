package agent

import (
	"testing"

	"github.com/Magpie-Monitor/magpie-monitor/agent/internal/agent/pods/data"
	"github.com/Magpie-Monitor/magpie-monitor/agent/pkg/tests"
	"github.com/stretchr/testify/assert"
	v2 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestLogSecondExtraction(t *testing.T) {

	testCase := []struct {
		name            string
		log             string
		expectedSeconds int
	}{
		{
			name:            "Test log second extraction",
			log:             "2006-01-02T15:04:05.123Z MY LOG\n",
			expectedSeconds: 5,
		},
		{
			name:            "Test log second extraction",
			log:             "2006-01-02T15:04:59.123Z MY LOG\n",
			expectedSeconds: 59,
		},
		{
			name:            "Test log second extraction",
			log:             "2006-12-02T15:04:00.123Z MY LOG\n",
			expectedSeconds: 0,
		},
	}

	for _, test := range testCase {
		testFunc := func(t *testing.T) {

			agent := Agent{}
			seconds, err := agent.getSecondFromLogTimestamp(test.log)
			if err != nil {
				t.Error("Error getting second from log timestamp")
			}

			assert.Equal(t, test.expectedSeconds, seconds)
		}

		t.Run(test.name, testFunc)
	}
}

func TestLogDeduplication(t *testing.T) {

	testCase := []struct {
		name         string
		logs         string
		expectedLogs string
	}{
		{
			name:         "Test deduplication",
			logs:         "2006-01-02T15:04:05.123Z MY LOG\n2006-01-02T15:04:05.456Z MY LOG\n2006-01-02T15:04:06.089Z MY LOG\n",
			expectedLogs: "2006-01-02T15:04:05.123Z MY LOG\n2006-01-02T15:04:05.456Z MY LOG",
		},
		{
			name:         "Test deduplication with lowest possible line number",
			logs:         "2006-01-02T15:04:05.123Z MY LOG\n2006-01-02T15:04:05.456Z MY LOG\n2006-01-02T15:04:06.089Z MY LOG\n",
			expectedLogs: "2006-01-02T15:04:05.123Z MY LOG\n2006-01-02T15:04:05.456Z MY LOG",
		},
		{
			name:         "Test deduplication with non duplicated logs",
			logs:         "2006-01-02T15:04:05.123Z MY LOG\n",
			expectedLogs: "2006-01-02T15:04:05.123Z MY LOG\n",
		},

		{
			name:         "Test deduplication with empty log",
			logs:         "",
			expectedLogs: "",
		},
	}

	for _, test := range testCase {
		testFunc := func(t *testing.T) {

			agent := Agent{}
			result, err := agent.deduplicate(test.logs)
			if err != nil {
				t.Fatalf("Error deduplicating logs %v", err)
			}

			assert.Equal(t, result, test.expectedLogs)
		}

		t.Run(test.name, testFunc)
	}
}

func TestPodLogsPacketSplit(t *testing.T) {

	var logs = `11111111111111111 \n
	222222222222222222222 \n
	333333333333333333333 \n
	444444444444444444444 \n`

	testCases := []struct {
		name                     string
		logs                     string
		maxPodPacketSizeBytes    int
		containerPacketSizeBytes int
		expectedPackets          int
	}{
		{
			name:                     "Test pod logs packet split single",
			logs:                     logs,
			maxPodPacketSizeBytes:    100,
			containerPacketSizeBytes: 20,
			expectedPackets:          1,
		},
		{
			name:                     "Test pod logs packet split even",
			logs:                     logs,
			maxPodPacketSizeBytes:    20,
			containerPacketSizeBytes: 20,
			expectedPackets:          4,
		},
		{
			name:                     "Test pod logs packet split uneven",
			logs:                     logs,
			maxPodPacketSizeBytes:    60,
			containerPacketSizeBytes: 20,
			expectedPackets:          2,
		},
	}

	for _, test := range testCases {
		testFunc := func(t *testing.T) {

			agent := Agent{maxPodPacketSizeBytes: test.maxPodPacketSizeBytes, maxContainerPacketSizeBytes: test.maxPodPacketSizeBytes}

			containers := agent.splitLogsIntoContainerPackets("test", "test", test.logs)

			packets := agent.splitPodContainerLogsIntoPackets(
				"test",
				containers,
			)

			assert.Equal(t, len(packets), test.expectedPackets)
		}

		t.Run(test.name, testFunc)
	}

}

func TestContainerLogsPacketSplit(t *testing.T) {

	var logs = `11111111111111111 \n
	222222222222222222222 \n
	333333333333333333333 \n
	444444444444444444444 \n`

	testCases := []struct {
		name                     string
		logs                     string
		containerPacketSizeBytes int
		expectedContainers       int
	}{
		{
			name:                     "Test_Container_Logs_Packet_Split_Even",
			logs:                     logs,
			containerPacketSizeBytes: 20,
			expectedContainers:       4,
		},
		{
			name:                     "Test container logs packet split single",
			logs:                     logs,
			containerPacketSizeBytes: 200,
			expectedContainers:       1,
		},
		{
			name:                     "Test container logs packet split none",
			logs:                     "",
			containerPacketSizeBytes: 1,
			expectedContainers:       1,
		},
		{
			name:                     "Test container logs packet split none zero",
			logs:                     "",
			containerPacketSizeBytes: 0,
			expectedContainers:       1,
		},
	}

	for _, test := range testCases {
		testFunc := func(t *testing.T) {
			agent := Agent{maxContainerPacketSizeBytes: test.containerPacketSizeBytes}

			containers := len(agent.splitLogsIntoContainerPackets("test", "test", test.logs))

			assert.Equal(t, containers, test.expectedContainers)
		}

		t.Run(test.name, testFunc)
	}
}

func TestDeploymentLogsSinceTime(t *testing.T) {

	var (
		namespace  = "namespace"
		deployment = "deployment"
		pod        = "pod"
		container  = "container"
		logs       = `2006-01-02T15:04:05.123Z MY LOG\n
2006-01-02T15:04:05.123Z MY LOG 1\n
2006-01-02T15:04:05.123Z MY LOG 2\n
2006-01-02T15:04:05.123Z MY LOG 3\n
`
		clusterId  = "cluster"
		objectMeta = metav1.ObjectMeta{
			Name:      pod,
			Namespace: namespace,
		}
		podSpec = v1.PodSpec{
			Containers: []v1.Container{
				{
					Name: container,
				},
			},
		}
	)

	testCase := []struct {
		name       string
		namespace  string
		deployment v2.Deployment
		pod        v1.Pod
		logs       tests.Logs
	}{
		{
			name:      "test",
			namespace: namespace,
			deployment: v2.Deployment{
				metav1.TypeMeta{},
				metav1.ObjectMeta{
					Name: deployment,
				},
				v2.DeploymentSpec{
					Template: v1.PodTemplateSpec{
						objectMeta,
						podSpec,
					},
				},
				v2.DeploymentStatus{},
			},
			pod: v1.Pod{
				metav1.TypeMeta{},
				objectMeta,
				podSpec,
				v1.PodStatus{},
			},
			logs: tests.Logs{
				Namespace: namespace,
				Pod:       pod,
				Container: container,
				Logs:      logs,
			},
		},
	}

	for _, test := range testCase {

		testFunc := func(t *testing.T) {
			res := make(chan data.Chunk)

			agent := Agent{
				kubernetesClient: tests.NewMockKubernetesApiClient(
					[]string{test.namespace},
					[]v2.Deployment{test.deployment},
					[]v2.StatefulSet{},
					[]v2.DaemonSet{},
					[]v1.Pod{test.pod},
					[]tests.Logs{test.logs},
				),
				results:                     res,
				readTimestamps:              make(map[string]int64, 0),
				maxPodPacketSizeBytes:       1000,
				maxContainerPacketSizeBytes: 1000,
				clusterId:                   clusterId,
			}

			go agent.fetchDeploymentLogsSinceTime(test.namespace, []v2.Deployment{test.deployment})

			msg := <-res

			assert.Equal(t, msg.ClusterId, clusterId)
			assert.Equal(t, "\n"+logs, msg.Pods[0].Containers[0].Content)
		}

		t.Run(test.name, testFunc)
	}
}

func TestStatefulSetLogsSinceTime(t *testing.T) {
	var (
		namespace   = "namespace"
		statefulSet = "statefulSet"
		pod         = "pod"
		container   = "container"
		logs        = `2006-01-02T15:04:05.123Z MY LOG\n
2006-01-02T15:04:05.123Z MY LOG 1\n
2006-01-02T15:04:05.123Z MY LOG 2\n
2006-01-02T15:04:05.123Z MY LOG 3\n
`
		clusterId  = "cluster"
		objectMeta = metav1.ObjectMeta{
			Name:      pod,
			Namespace: namespace,
		}
		podSpec = v1.PodSpec{
			Containers: []v1.Container{
				{
					Name: container,
				},
			},
		}
	)

	testCase := []struct {
		name        string
		namespace   string
		statefulSet v2.StatefulSet
		pod         v1.Pod
		logs        tests.Logs
	}{
		{
			name:      "test",
			namespace: namespace,
			statefulSet: v2.StatefulSet{
				metav1.TypeMeta{},
				metav1.ObjectMeta{
					Name: statefulSet,
				},
				v2.StatefulSetSpec{
					Template: v1.PodTemplateSpec{
						objectMeta,
						podSpec,
					},
				},
				v2.StatefulSetStatus{},
			},
			pod: v1.Pod{
				metav1.TypeMeta{},
				objectMeta,
				podSpec,
				v1.PodStatus{},
			},
			logs: tests.Logs{
				Namespace: namespace,
				Pod:       pod,
				Container: container,
				Logs:      logs,
			},
		},
	}

	for _, test := range testCase {

		testFunc := func(t *testing.T) {
			res := make(chan data.Chunk)

			agent := Agent{
				kubernetesClient: tests.NewMockKubernetesApiClient(
					[]string{test.namespace},
					[]v2.Deployment{},
					[]v2.StatefulSet{test.statefulSet},
					[]v2.DaemonSet{},
					[]v1.Pod{test.pod},
					[]tests.Logs{test.logs},
				),
				results:                     res,
				readTimestamps:              make(map[string]int64, 0),
				maxPodPacketSizeBytes:       1000,
				maxContainerPacketSizeBytes: 1000,
				clusterId:                   clusterId,
			}

			go agent.fetchStatefulSetLogsSinceTime(test.namespace, []v2.StatefulSet{test.statefulSet})

			msg := <-res

			assert.Equal(t, msg.ClusterId, clusterId)
			assert.Equal(t, "\n"+logs, msg.Pods[0].Containers[0].Content)
		}

		t.Run(test.name, testFunc)
	}
}

func TestDaemonSetLogsSinceTime(t *testing.T) {
	var (
		namespace = "namespace"
		daemonSet = "daemonSet"
		pod       = "pod"
		container = "container"
		logs      = `2006-01-02T15:04:05.123Z MY LOG\n
2006-01-02T15:04:05.123Z MY LOG 1\n
2006-01-02T15:04:05.123Z MY LOG 2\n
2006-01-02T15:04:05.123Z MY LOG 3\n
`
		clusterId  = "cluster"
		objectMeta = metav1.ObjectMeta{
			Name:      pod,
			Namespace: namespace,
		}
		podSpec = v1.PodSpec{
			Containers: []v1.Container{
				{
					Name: container,
				},
			},
		}
	)

	testCase := []struct {
		name      string
		namespace string
		daemonSet v2.DaemonSet
		pod       v1.Pod
		logs      tests.Logs
	}{
		{
			name:      "test",
			namespace: namespace,
			daemonSet: v2.DaemonSet{
				metav1.TypeMeta{},
				metav1.ObjectMeta{
					Name: daemonSet,
				},
				v2.DaemonSetSpec{
					Template: v1.PodTemplateSpec{
						objectMeta,
						podSpec,
					},
				},
				v2.DaemonSetStatus{},
			},
			pod: v1.Pod{
				metav1.TypeMeta{},
				objectMeta,
				podSpec,
				v1.PodStatus{},
			},
			logs: tests.Logs{
				Namespace: namespace,
				Pod:       pod,
				Container: container,
				Logs:      logs,
			},
		},
	}

	for _, test := range testCase {

		testFunc := func(t *testing.T) {
			res := make(chan data.Chunk)

			agent := Agent{
				kubernetesClient: tests.NewMockKubernetesApiClient(
					[]string{test.namespace},
					[]v2.Deployment{},
					[]v2.StatefulSet{},
					[]v2.DaemonSet{test.daemonSet},
					[]v1.Pod{test.pod},
					[]tests.Logs{test.logs},
				),
				results:                     res,
				readTimestamps:              make(map[string]int64, 0),
				maxPodPacketSizeBytes:       1000,
				maxContainerPacketSizeBytes: 1000,
				clusterId:                   clusterId,
			}

			go agent.fetchDaemonSetLogsSinceTime(test.namespace, []v2.DaemonSet{test.daemonSet})

			msg := <-res

			assert.Equal(t, msg.ClusterId, clusterId)
			assert.Equal(t, "\n"+logs, msg.Pods[0].Containers[0].Content)
		}

		t.Run(test.name, testFunc)
	}
}

func TestFetchNamespaces(t *testing.T) {

	testCase := []struct {
		name               string
		namespaces         []string
		excludedNamespaces []string
		result             []string
	}{
		{
			name:               "Fetch namespaces with excluded namespaces",
			namespaces:         []string{"ns-1", "ns-2", "ns-3"},
			excludedNamespaces: []string{"ns-3", "ns-4"},
			result:             []string{"ns-1", "ns-2"},
		},
		{
			name:               "Fetch namespaces without excluded namespaces",
			namespaces:         []string{"ns-1", "ns-2", "ns-3"},
			excludedNamespaces: []string{},
			result:             []string{"ns-1", "ns-2", "ns-3"},
		},
		{
			name:               "Fetch empty namespaces",
			namespaces:         []string{},
			excludedNamespaces: []string{},
			result:             []string{},
		},
		{
			name:               "Fetch namespaces without overlapping excluded namespaces",
			namespaces:         []string{"ns-1", "ns-2"},
			excludedNamespaces: []string{"ns-3", "ns-4"},
			result:             []string{"ns-1", "ns-2"},
		},
	}

	for _, test := range testCase {

		testFunc := func(t *testing.T) {
			agent := Agent{
				excludedNamespaces: test.excludedNamespaces,
				kubernetesClient: tests.NewMockKubernetesApiClient(
					test.namespaces,
					[]v2.Deployment{},
					[]v2.StatefulSet{},
					[]v2.DaemonSet{},
					[]v1.Pod{},
					[]tests.Logs{},
				),
			}

			agent.fetchNamespaces()

			res := agent.includedNamespaces
			assert.Equal(t, res, test.namespaces)
		}

		t.Run(test.name, testFunc)
	}

}

func TestGatherClusterMetadata(t *testing.T) {

	testCase := []struct {
		name         string
		namespaces   []string
		deployments  []v2.Deployment
		statefulSets []v2.StatefulSet
		daemonSets   []v2.DaemonSet
		clusterId    string
		result       data.ApplicationState
	}{
		{
			name:         "Gather cluster metadata with multiple applications",
			namespaces:   []string{"cluster"},
			deployments:  tests.InitializeDeployments([]string{"dp"}),
			statefulSets: tests.InitializeStatefulSets([]string{"sts", "sts-1"}),
			daemonSets:   tests.InitializeDaemonSets([]string{"ds"}),
			clusterId:    "cluster",
			result: data.ApplicationState{
				ClusterId: "cluster",
				Applications: []data.Application{
					{
						Kind: data.Deployment,
						Name: "dp",
					},
					{
						Kind: data.StatefulSet,
						Name: "sts",
					},
					{
						Kind: data.StatefulSet,
						Name: "sts-1",
					},
					{
						Kind: data.DaemonSet,
						Name: "ds",
					},
				},
			},
		},
		{
			name:         "Gather cluster metadata without applications",
			namespaces:   []string{"cluster"},
			deployments:  make([]v2.Deployment, 0),
			statefulSets: make([]v2.StatefulSet, 0),
			daemonSets:   make([]v2.DaemonSet, 0),
			clusterId:    "cluster",
			result: data.ApplicationState{
				ClusterId:    "cluster",
				Applications: []data.Application{},
			},
		},
	}

	for _, test := range testCase {

		testFunc := func(t *testing.T) {
			metadata := make(chan data.ApplicationState)

			agent := Agent{
				includedNamespaces: []string{"test"},
				kubernetesClient: tests.NewMockKubernetesApiClient(
					test.namespaces,
					test.deployments,
					test.statefulSets,
					test.daemonSets,
					[]v1.Pod{},
					[]tests.Logs{},
				),
				metadata:                          metadata,
				metadataCollectionIntervalSeconds: 5,
				clusterId:                         test.clusterId,
			}

			go agent.gatherClusterMetadata()

			msg := <-metadata

			assert.Equal(t, msg.ClusterId, test.result.ClusterId)
			assert.Equal(t, msg.Applications, test.result.Applications)
		}

		t.Run(test.name, testFunc)
	}
}

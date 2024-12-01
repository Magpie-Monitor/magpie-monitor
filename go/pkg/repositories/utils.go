package repositories

import (
	"context"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"testing"
)

var (
	DEFAULT_TEST_APPLICATION_LOGS_COLLECTED_AT_DATE = 10
	DEFAULT_TEST_APPLICATION_LOGS_KIND              = "test-kind"
	DEFAULT_TEST_CLUSTER_ID                         = "default-test-cluster"
	DEFAULT_TEST_APPLICATION_POD_NAME               = "default-test-pod-name"
	DEFAULT_TEST_APPLICATION_IMAGE_NAME             = "default-test-image"
	DEFAULT_TEST_APPLICATION_LOGS_CONTENT           = "Application logs test content"
	DEFAULT_TEST_APPLICATION_CONTAINER_NAME         = "test-container"
	DEFAULT_TEST_APPLICATION_NAMESPACE              = "test-namespace"
	DEFAULT_TEST_APPLICATION_NAME                   = "test-application"

	DEFAULT_APPLICATION_LOG = ApplicationLogs{
		Name:          DEFAULT_TEST_APPLICATION_NAME,
		ClusterId:     DEFAULT_TEST_CLUSTER_ID,
		Kind:          DEFAULT_TEST_APPLICATION_LOGS_KIND,
		CollectedAtMs: int64(DEFAULT_TEST_APPLICATION_LOGS_COLLECTED_AT_DATE),
		Namespace:     DEFAULT_TEST_APPLICATION_NAMESPACE,
		Pods: []*PodLogs{
			{
				Name: DEFAULT_TEST_APPLICATION_POD_NAME,
				Containers: []*ContainerLogs{
					{
						Name:    DEFAULT_TEST_APPLICATION_CONTAINER_NAME,
						Image:   DEFAULT_TEST_APPLICATION_IMAGE_NAME,
						Content: DEFAULT_TEST_APPLICATION_LOGS_CONTENT,
					},
				},
			},
		},
	}

	DEFAULT_TEST_NODE_LOGS_COLLECTED_AT_DATE = 10
	DEFAULT_TEST_NODE_LOGS_CLUSTER_ID        = 10
	DEFAULT_TEST_NODE_LOGS_KIND              = "test-kind"
	DEFAULT_TEST_NODE_LOGS_CONTENT           = "Node logs test content"
	DEFAULT_TEST_NODE_LOGS_FILENAME          = "/var/logs/test-journal"
	DEFAULT_TEST_NODE_NAMESPACE              = "test-namespace"
	DEFAULT_TEST_NODE_NAME                   = "test-node"

	DEFAULT_NODE_LOG = NodeLogs{
		ClusterId:     DEFAULT_TEST_CLUSTER_ID,
		Kind:          DEFAULT_TEST_NODE_LOGS_KIND,
		CollectedAtMs: int64(DEFAULT_TEST_NODE_LOGS_COLLECTED_AT_DATE),
		Name:          DEFAULT_TEST_NODE_NAME,
		Filename:      DEFAULT_TEST_NODE_LOGS_FILENAME,
		Content:       DEFAULT_TEST_NODE_LOGS_CONTENT,
	}
)

func GetDefaultNodeTestLogsFromCluster(clusterId string) []*NodeLogs {

	defaultLog := DEFAULT_NODE_LOG

	defaultLog.ClusterId = clusterId

	return []*NodeLogs{&defaultLog}
}

func GetDefaultApplicationTestLogsFromCluster(clusterId string) []*ApplicationLogs {

	defaultLog := DEFAULT_APPLICATION_LOG

	defaultLog.ClusterId = clusterId

	return []*ApplicationLogs{&defaultLog}
}

func PrefillApplicationLogs(
	t *testing.T,
	logger *zap.Logger,
	applicationLogsRepository ApplicationLogsRepository,
	logs []*ApplicationLogs) []string {

	ctx := context.Background()

	indexes := make(map[string]bool, 0)

	for _, log := range logs {
		index := getApplicationLogsIndexName(log)
		indexes[index] = true
	}

	for index := range indexes {
		applicationLogsRepository.RemoveIndex(ctx, index)
	}

	insertedLogsIds := make([]string, 0, 0)

	for _, log := range logs {
		ids, err := applicationLogsRepository.InsertLogs(ctx, log)
		assert.NoError(t, err, "Failed to prefill application logs")
		insertedLogsIds = append(insertedLogsIds, ids...)
	}

	return insertedLogsIds
}

// Returns ids of created logs
func PrefillNodeLogs(
	t *testing.T,
	logger *zap.Logger,
	nodeLogsRepository NodeLogsRepository,
	logs []*NodeLogs) []string {

	ctx := context.Background()

	indexes := make(map[string]bool, 0)

	for _, log := range logs {
		index := getNodeLogsIndexName(log)
		indexes[index] = true
	}

	for index := range indexes {
		nodeLogsRepository.RemoveIndex(ctx, index)
		// assert.NoError(t, err, "Failed to remove index")
	}
	ids := make([]string, 0, 0)

	for _, log := range logs {
		insertedLogId, err := nodeLogsRepository.InsertLogs(ctx, log)
		ids = append(ids, insertedLogId)
		assert.NoError(t, err, "Failed to prefill application logs")
	}

	return ids
}

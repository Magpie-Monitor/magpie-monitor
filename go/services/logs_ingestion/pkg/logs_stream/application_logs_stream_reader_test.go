package logsstream_test

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/Magpie-Monitor/magpie-monitor/pkg/repositories"
	"github.com/Magpie-Monitor/magpie-monitor/pkg/tests"
	"github.com/Magpie-Monitor/magpie-monitor/services/logs_ingestion/pkg/config"
	_ "github.com/Magpie-Monitor/magpie-monitor/services/logs_ingestion/pkg/config"
	logsstream "github.com/Magpie-Monitor/magpie-monitor/services/logs_ingestion/pkg/logs_stream"
	"github.com/stretchr/testify/assert"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

const LOGS_INGESTION_TEST_INDEX = "ingestion-test"

func TestApplicationLogsStreamReader(t *testing.T) {

	type TestDependencies struct {
		fx.In
		Reader                    *logsstream.KafkaApplicationLogsStreamReader
		Writer                    *tests.KafkaLogsStreamWriter
		Logger                    *zap.Logger
		ApplicationLogsRepository repositories.ApplicationLogsRepository
	}

	test := func(dependencies TestDependencies) {
		if dependencies.Reader == nil {
			t.Fatal("Failed to load reader")
		}

		dependencies.Logger.Info("Success", zap.Any("reader", dependencies.Reader))

		ctx := context.Background()

		dependencies.ApplicationLogsRepository.RemoveIndex(ctx, fmt.Sprintf("%s-applications-1970-1", LOGS_INGESTION_TEST_INDEX))

		// err = dependencies.ApplicationLogsRepository.RemoveIndex(ctx, "test-cluster-nodes-1970-1")
		// assert.NoError(t, err, "Failed to remove application logs index")

		testNodeLog := repositories.NodeLogs{
			Id:            "log-1",
			ClusterId:     LOGS_INGESTION_TEST_INDEX,
			Kind:          "Node",
			CollectedAtMs: 10,
			Name:          "test-node",
			Filename:      "/etc/var/journalctl",
			Content:       "Test content",
		}

		encodedNodeLog, err := json.Marshal(testNodeLog)
		if err != nil {
			t.Fatal("Failed to encode node logs")
		}

		dependencies.Writer.WriteNodeLogs(ctx, "testkey-1", string(encodedNodeLog))

		testApplicationLog := repositories.ApplicationLogs{
			ClusterId:     LOGS_INGESTION_TEST_INDEX,
			Kind:          "Node",
			CollectedAtMs: 10,
			Namespace:     "test-node",
			Pods: []*repositories.PodLogs{
				&repositories.PodLogs{
					Containers: []*repositories.ContainerLogs{
						&repositories.ContainerLogs{
							Name:    "test-container",
							Image:   "test-image",
							Content: "Test application content",
						},
					},
				},
			},
		}

		encodedApplicationLog, err := json.Marshal(testApplicationLog)
		if err != nil {
			t.Fatal("Failed to encode node logs")
		}

		dependencies.Writer.WriteApplicationLogs(ctx, "testkey-1", string(encodedApplicationLog))

		go dependencies.Reader.Listen()

		// Wait 5 seconds for entries to be fetched from Kafka and inserted into Elastic
		time.Sleep(time.Second * 5)

		applicationLogs, err := dependencies.ApplicationLogsRepository.GetLogs(
			ctx,
			LOGS_INGESTION_TEST_INDEX,
			time.UnixMilli(0),
			time.UnixMilli(20),
		)
		if err != nil {
			t.Fatal("Failed to fetch application logs")
		}

		// Check if an element is returned
		assert.Len(t, applicationLogs, 1)

		// Check if the name matches
		assert.Equal(t, "test-node", applicationLogs[0].Namespace)
	}

	tests.RunTest(test, t, config.AppModule)
}

package logsstream_test

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Magpie-Monitor/magpie-monitor/pkg/envs"
	"github.com/Magpie-Monitor/magpie-monitor/pkg/repositories"
	"github.com/Magpie-Monitor/magpie-monitor/pkg/tests"
	"github.com/Magpie-Monitor/magpie-monitor/services/logs_ingestion/pkg/config"
	_ "github.com/Magpie-Monitor/magpie-monitor/services/logs_ingestion/pkg/config"
	logsstream "github.com/Magpie-Monitor/magpie-monitor/services/logs_ingestion/pkg/logs_stream"
	"github.com/stretchr/testify/assert"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"testing"
	"time"
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

	integrationTestWaitModifier := envs.ConvertToInt(tests.INTEGRATION_TEST_WAIT_MODIFIER_KEY)

	test := func(dependencies TestDependencies) {
		if dependencies.Reader == nil {
			t.Fatal("Failed to load reader")
		}

		dependencies.Logger.Info("Success", zap.Any("reader", dependencies.Reader))

		ctx := context.Background()

		dependencies.ApplicationLogsRepository.RemoveIndex(ctx, fmt.Sprintf("%s-applications-1970-1", LOGS_INGESTION_TEST_INDEX))

		testApplicationLog := repositories.ApplicationLogs{
			ClusterId:     LOGS_INGESTION_TEST_INDEX,
			Kind:          "Node",
			CollectedAtMs: 10,
			Namespace:     "test-node",
			Pods: []*repositories.PodLogs{
				&repositories.PodLogs{
					Name: "test-pod",
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

		expectedLogDocument := repositories.ApplicationLogsDocument{
			ClusterId:     LOGS_INGESTION_TEST_INDEX,
			Kind:          "Node",
			CollectedAtMs: 10,
			Namespace:     "test-node",
			PodName:       "test-pod",
			ContainerName: "test-container",
			Content:       "Test application content",
			Image:         "test-image",
		}

		encodedApplicationLog, err := json.Marshal(testApplicationLog)
		if err != nil {
			t.Fatal("Failed to encode node logs")
		}

		dependencies.Writer.WriteApplicationLogs(ctx, "testkey-1", string(encodedApplicationLog))

		go dependencies.Reader.Listen()

		// Wait 10 seconds for entries to be fetched from Kafka and inserted into Elastic
		time.Sleep(time.Second * 20 * time.Duration(integrationTestWaitModifier))

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

		// Unset id to check only predefined values
		applicationLogs[0].Id = ""

		// Check if the log is correctly fetched and transformed
		assert.Equal(t, expectedLogDocument, *applicationLogs[0], "Expected log does not match the actual application log")
	}

	tests.RunTest(test, t, config.AppModule)
}

package logsstream_test

import (
	"context"
	"encoding/json"
	"fmt"
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

func TestNodeLogsStreamReader(t *testing.T) {

	type TestDependencies struct {
		fx.In
		Reader             *logsstream.KafkaNodeLogsStreamReader
		Writer             *tests.KafkaLogsStreamWriter
		Logger             *zap.Logger
		NodeLogsRepository repositories.NodeLogsRepository
	}

	test := func(dependencies TestDependencies) {
		if dependencies.Reader == nil {
			t.Fatal("Failed to load reader")
		}

		ctx := context.Background()

		dependencies.NodeLogsRepository.RemoveIndex(ctx, fmt.Sprintf("%s-nodes-1970-1", LOGS_INGESTION_TEST_INDEX))

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

		// encodedApplicationLog, err := json.Marshal(testApplicationLog)
		// if err != nil {
		// 	t.Fatal("Failed to encode node logs")
		// }

		// dependencies.Writer.WriteApplicationLogs(ctx, "testkey-1", string(encodedApplicationLog))

		go dependencies.Reader.Listen()

		// Wait 5 seconds for entries to be fetched from Kafka and inserted into Elastic
		time.Sleep(time.Second * 5)

		nodeLogs, err := dependencies.NodeLogsRepository.GetLogs(
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
		assert.Equal(t, "test-node", applicationLogs[0].Name)
	}

	tests.RunTest(test, t, config.AppModule)
}

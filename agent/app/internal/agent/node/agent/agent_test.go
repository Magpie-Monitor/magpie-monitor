package agent

import (
	"os"
	"testing"
	"time"

	"github.com/Magpie-Monitor/magpie-monitor/agent/internal/agent/node/data"
	"github.com/Magpie-Monitor/magpie-monitor/agent/pkg/tests"
	"github.com/stretchr/testify/assert"
)

func TestLogsSplit(t *testing.T) {

	testCases := []struct {
		name            string
		logs            string
		packetSizeBytes int
		expectedPackets int
	}{
		{
			name:            "Test logs packet split",
			logs:            "111111111111111111 \n 222222222222222222 \n 3333333333333333333 \n 4444444444444444444",
			packetSizeBytes: 80,
			expectedPackets: 2,
		},
		{
			name:            "Test logs packet split single",
			logs:            "111111111111111111 \n 222222222222222222 \n 3333333333333333333 \n 4444444444444444444",
			packetSizeBytes: 200,
			expectedPackets: 1,
		},
		{
			name:            "Test logs packet split none",
			logs:            "",
			packetSizeBytes: 200,
			expectedPackets: 1,
		},
	}

	for _, test := range testCases {
		testFunc := func(t *testing.T) {
			agent := IncrementalReader{packetSizeBytes: test.packetSizeBytes}

			packets := agent.splitLogsIntoPackets("test", "test", "test", test.logs)

			assert.Equal(t, test.expectedPackets, len(packets))
		}

		t.Run(test.name, testFunc)
	}
}

func TestGatherNodeMatadata(t *testing.T) {

	testCases := []struct {
		name         string
		clusterId    string
		nodeName     string
		watchedFiles []string
	}{
		{
			name:         "Test gather metadata with multiple files",
			clusterId:    "cluster",
			nodeName:     "node",
			watchedFiles: []string{"file1", "file2", "file3"},
		},
		{
			name:         "Test gather metadata without files",
			clusterId:    "cluster",
			nodeName:     "node",
			watchedFiles: []string{},
		},
	}

	for _, test := range testCases {
		testFunc := func(t *testing.T) {
			metadata := make(chan data.NodeState)

			agent := IncrementalReader{metadata: metadata, nodeName: test.nodeName, clusterId: test.clusterId, files: test.watchedFiles}

			go agent.gatherNodeMetadata()

			msg := <-metadata

			assert.Equal(t, test.clusterId, msg.ClusterId)
			assert.Equal(t, test.nodeName, msg.NodeName)
			assert.Equal(t, test.watchedFiles, msg.WatchedFiles)
		}

		t.Run(test.name, testFunc)
	}
}

func TestWatchFile(t *testing.T) {

	testCases := []struct {
		name       string
		fileName   string
		clusterId  string
		nodeName   string
		logContent string
	}{
		{
			name: "Test gather metadata with multiple files",
			// /files/file is created as part of Dockerfile "test"
			fileName:   "/files/file",
			clusterId:  "cluster",
			nodeName:   "node",
			logContent: "Test line",
		},
		{
			name:      "Test gather metadata without files",
			clusterId: "cluster",
			nodeName:  "node",
			// /files/file is created as part of Dockerfile "test"
			fileName:   "/files/file",
			logContent: "Test line\nTest line\nTest line\nTest line\nTest line\nTest line\n",
		},
	}

	for _, test := range testCases {
		testFunc := func(t *testing.T) {
			results := make(chan data.Chunk)

			agent := IncrementalReader{
				results:               results,
				nodeName:              test.nodeName,
				clusterId:             test.clusterId,
				redis:                 tests.NewMockRedis(),
				scrapeIntervalSeconds: 3,
				packetSizeBytes:       200,
				files:                 []string{test.fileName},
			}

			go agent.watchFiles()
			time.Sleep(3 * time.Second)

			file, err := os.OpenFile(test.fileName, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
			if err != nil {
				t.Error("Error opening a file")
			}

			_, err = file.WriteString(test.logContent)
			if err != nil {
				t.Error("Error writing to file")
			}

			msg := <-results
			// Add string escape character to the end of the content
			assert.EqualValues(t, test.logContent+"\x00", msg.Content)
			assert.Equal(t, test.clusterId, msg.ClusterId)
			assert.Equal(t, test.nodeName, msg.Name)
			assert.Equal(t, test.fileName, msg.Filename)

			err = os.Truncate(test.fileName, 0)
			if err != nil {
				t.Error("Error wiping content of the file")
			}
		}

		t.Run(test.name, testFunc)
	}
}

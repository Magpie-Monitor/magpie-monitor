package openai_test

import (
	"github.com/Magpie-Monitor/magpie-monitor/pkg/tests"
	"github.com/Magpie-Monitor/magpie-monitor/services/reports/pkg/config"
	"github.com/Magpie-Monitor/magpie-monitor/services/reports/pkg/openai"
	"github.com/stretchr/testify/assert"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"testing"
)

func TestClientSplitCompletionReqestsByBatchSize(t *testing.T) {

	type TestDependencies struct {
		fx.In
		Logger *zap.Logger
		Client *openai.Client
	}

	testCases := []struct {
		completionRequests                 map[string]*openai.CompletionRequest
		batchSizeBytes                     int
		expectedSplittedCompletionRequests []map[string]*openai.CompletionRequest
	}{
		{
			batchSizeBytes: 180,
			completionRequests: map[string]*openai.CompletionRequest{
				"app-1": {
					// Total length of 171 bytes
					Messages: []*openai.Message{
						{
							Content: "This is the first message content",
						},
						{
							Content: "This is a second message",
						},
					},
				},
				"app-2": {
					// Total length of 171 bytes
					Messages: []*openai.Message{
						{
							Content: "This is the first message content",
						},
						{
							Content: "This is a second message",
						},
					},
				},
				"app-3": {
					// Total length of 171 bytes
					Messages: []*openai.Message{
						{
							Content: "This is the first message content",
						},
						{
							Content: "This is a second message",
						},
					},
				},
			},
			expectedSplittedCompletionRequests: []map[string]*openai.CompletionRequest{
				{
					"app-1": {
						// Total length of 171 bytes
						Messages: []*openai.Message{
							{
								Content: "This is the first message content",
							},
							{
								Content: "This is a second message",
							},
						},
					},
				},
				{
					"app-2": {
						// Total length of 171 bytes
						Messages: []*openai.Message{
							{
								Content: "This is the first message content",
							},
							{
								Content: "This is a second message",
							},
						},
					},
				},
				{

					"app-3": {
						// Total length of 171 bytes
						Messages: []*openai.Message{
							{
								Content: "This is the first message content",
							},
							{
								Content: "This is a second message",
							},
						},
					},
				},
			},
		},
	}

	test := func(dependencies TestDependencies) {
		for _, tc := range testCases {
			dependencies.Client.BatchSizeBytes = tc.batchSizeBytes
			splitted, err := dependencies.Client.SplitCompletionReqestsByBatchSize(tc.completionRequests)
			assert.NoError(t, err, "Failed to split completion requests")
			assert.ElementsMatch(t, tc.expectedSplittedCompletionRequests, splitted)

		}
	}

	tests.RunTest(test, t, config.AppModule)
}

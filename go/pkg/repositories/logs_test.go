package repositories

import (
	"reflect"
	"testing"
)

type TestLog struct {
	Content string
}

func (l *TestLog) GetContent() string {
	return l.Content
}

func TestDecode(t *testing.T) {
	testsCases := []struct {
		description string
		rawString   string
		allLogs     []*TestLog
		expect      [][]*TestLog
	}{
		{
			description: "Given list of logs should separate them by max packet size",
			allLogs: []*TestLog{
				{
					Content: "0123456789",
				},
				{
					Content: "0123456789",
				},
				{
					Content: "0123456789",
				},
				{
					Content: "0123456789",
				},
			},
			expect: [][]*TestLog{
				{
					&TestLog{
						Content: "0123456789",
					},
					&TestLog{
						Content: "0123456789",
					},
				},
				{
					&TestLog{
						Content: "0123456789",
					},
					&TestLog{
						Content: "0123456789",
					},
				},
			},
		},
	}

	for _, tc := range testsCases {
		t.Run(tc.description, func(t *testing.T) {
			logPackets := SplitLogsIntoPackets(tc.allLogs, 20)

			// testCaseReader := strings.NewReader(tc.rawString)
			// results := make([]decodingTestCase, 0)
			//
			// err := jsonl.NewJsonLinesDecoder(testCaseReader).Decode(&results)
			// if err != nil {
			// 	t.Fatalf("Failed to decode an rawString into jsonl")
			// }
			if !reflect.DeepEqual(tc.expect, logPackets) {
				t.Fatalf("Wanted %+v, got %+v", tc.expect, logPackets)
			}

		})
	}

}

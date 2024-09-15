package pods

import "testing"

func Test_Second_Extraction(t *testing.T) {
	a := Agent{}
	log := "2006-01-02T15:04:05.123Z MY LOG\n"
	expectedSeconds := 5
	seconds, err := a.getSecondFromLogTimestamp(log)

	if err != nil {
		t.Fatalf("Error %v", err)
	}

	if seconds != expectedSeconds {
		t.Fatalf("Result: \n%d \nnot equal to expected result: \n%d", seconds, expectedSeconds)
	}
}

func Test_Deduplication(t *testing.T) {
	// Kubernetes API returns logs with tailing newline.
	logs := "2006-01-02T15:04:05.123Z MY LOG\n2006-01-02T15:04:05.456Z MY LOG\n2006-01-02T15:04:06.089Z MY LOG\n"
	expectedLogs := "2006-01-02T15:04:05.123Z MY LOG\n2006-01-02T15:04:05.456Z MY LOG"

	testDeduplication(logs, expectedLogs, t)
}

func Test_Deduplication_With_Lowest_Possible_Line_Number(t *testing.T) {
	// Kubernetes API returns logs with tailing newline.
	logs := "2006-01-02T15:04:05.123Z MY LOG\n2006-01-02T15:04:06.0056Z MY LOG\n"
	expectedLogs := "2006-01-02T15:04:05.123Z MY LOG"

	testDeduplication(logs, expectedLogs, t)
}

func Test_Deduplication_With_Non_Duplicated_Logs(t *testing.T) {
	// Kubernetes API returns logs with tailing newline.
	logs := "2006-01-02T15:04:05.123Z MY LOG\n"
	expectedLogs := "2006-01-02T15:04:05.123Z MY LOG\n"

	testDeduplication(logs, expectedLogs, t)
}

func Test_Deduplication_With_Empty_Log(t *testing.T) {
	logs := ""
	expectedLogs := ""

	testDeduplication(logs, expectedLogs, t)
}

func testDeduplication(logs, expectedLogs string, t *testing.T) {
	agent := Agent{}
	result, err := agent.deduplicate(logs)
	if err != nil {
		t.Fatalf("Error %v", err)
	}

	if result != expectedLogs {
		t.Fatalf("Result: \n%s \nnot equal to expected result: \n%s", result, expectedLogs)
	}
}

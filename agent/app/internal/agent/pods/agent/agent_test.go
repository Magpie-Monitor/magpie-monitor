package agent

import (
	"testing"
)

func Test_Log_Second_Extraction(t *testing.T) {
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

var logs = `11111111111111111 \n
		222222222222222222222 \n
		333333333333333333333 \n
		444444444444444444444 \n`

func Test_Pod_Logs_Packet_Split_Single(t *testing.T) {
	testPodContainerPacketSplit(
		podLogsSplitParams{
			Logs:                     logs,
			MaxPodPacketSizeBytes:    100,
			ContainerPacketSizeBytes: 20,
			ExpectedPackets:          1,
		},
		t,
	)
}

func Test_Pod_Logs_Packet_Split_Even(t *testing.T) {
	testPodContainerPacketSplit(
		podLogsSplitParams{
			Logs:                     logs,
			MaxPodPacketSizeBytes:    20,
			ContainerPacketSizeBytes: 20,
			ExpectedPackets:          4,
		},
		t,
	)
}

func Test_Pod_Logs_Packet_Split_Uneven(t *testing.T) {
	testPodContainerPacketSplit(
		podLogsSplitParams{
			Logs:                     logs,
			MaxPodPacketSizeBytes:    60,
			ContainerPacketSizeBytes: 20,
			ExpectedPackets:          2,
		},
		t,
	)
}

type podLogsSplitParams struct {
	Logs                     string
	MaxPodPacketSizeBytes    int
	ContainerPacketSizeBytes int
	ExpectedPackets          int
}

func testPodContainerPacketSplit(params podLogsSplitParams, t *testing.T) {
	agent := Agent{maxPodPacketSizeBytes: params.MaxPodPacketSizeBytes, maxContainerPacketSizeBytes: params.ContainerPacketSizeBytes}

	containers := agent.splitLogsIntoContainerPackets("test", "test", params.Logs)

	packets := agent.splitPodContainerLogsIntoPackets(
		"test",
		containers,
	)

	if len(packets) != params.ExpectedPackets {
		t.Fatalf("Number of packets: %d not equal to expected result: %d", len(packets), params.ExpectedPackets)
	}
}

func Test_Container_Logs_Packet_Split_Even(t *testing.T) {
	testContainerLogsPacketSplit(
		containerLogsSplitParams{
			Logs:                     "11111111111111111 \n 222222222222222222222 \n 333333333333333333333 \n 444444444444444444444 \n",
			ContainerPacketSizeBytes: 20,
			ExpectedContainers:       4,
		},
		t,
	)
}

func Test_Container_Logs_Packet_Split_Single(t *testing.T) {
	testContainerLogsPacketSplit(
		containerLogsSplitParams{
			Logs:                     "11111111111111111 \n 222222222222222222222 \n 333333333333333333333 \n 444444444444444444444 \n",
			ContainerPacketSizeBytes: 200,
			ExpectedContainers:       1,
		},
		t,
	)
}

func Test_Container_Logs_Packet_Split_None(t *testing.T) {
	testContainerLogsPacketSplit(
		containerLogsSplitParams{
			Logs:                     "",
			ContainerPacketSizeBytes: 1,
			ExpectedContainers:       1,
		},
		t,
	)
}

func Test_Container_Logs_Packet_Split_None_Zero(t *testing.T) {
	testContainerLogsPacketSplit(
		containerLogsSplitParams{
			Logs:                     "",
			ContainerPacketSizeBytes: 0,
			ExpectedContainers:       1,
		},
		t,
	)
}

type containerLogsSplitParams struct {
	Logs                     string
	ContainerPacketSizeBytes int
	ExpectedContainers       int
}

func testContainerLogsPacketSplit(params containerLogsSplitParams, t *testing.T) {
	agent := Agent{maxContainerPacketSizeBytes: params.ContainerPacketSizeBytes}

	containers := len(agent.splitLogsIntoContainerPackets("test", "test", params.Logs))
	expectedContainers := params.ExpectedContainers

	if containers != expectedContainers {
		t.Fatalf("Number of containers: %d not equal to expected result: %d", containers, expectedContainers)
	}
}

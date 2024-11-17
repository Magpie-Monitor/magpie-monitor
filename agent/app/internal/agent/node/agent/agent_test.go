package agent

import (
	"testing"
)

func Test_Logs_Packet_Split(t *testing.T) {
	testLogsSplit(
		logsPacketSplitParams{
			PacketSizeBytes: 80,
			ExpectedPackets: 2,
			Logs:            "111111111111111111 \n 222222222222222222 \n 3333333333333333333 \n 4444444444444444444",
		},
		t,
	)
}

func Test_Logs_Packet_Split_Single(t *testing.T) {
	testLogsSplit(
		logsPacketSplitParams{
			PacketSizeBytes: 200,
			ExpectedPackets: 1,
			Logs:            "111111111111111111 \n 222222222222222222 \n 3333333333333333333 \n 4444444444444444444",
		},
		t,
	)
}

func Test_Logs_Packet_Split_None(t *testing.T) {
	testLogsSplit(
		logsPacketSplitParams{
			PacketSizeBytes: 200,
			ExpectedPackets: 1,
			Logs:            "",
		},
		t,
	)
}

type logsPacketSplitParams struct {
	PacketSizeBytes int
	ExpectedPackets int
	Logs            string
}

func testLogsSplit(params logsPacketSplitParams, t *testing.T) {
	agent := IncrementalReader{packetSizeBytes: params.PacketSizeBytes}

	packets := agent.splitLogsIntoPackets("test", "test", "test", params.Logs)

	if len(packets) != params.ExpectedPackets {
		t.Fatalf("Number of packets: %d doesn't match expected packets: %d", len(packets), params.ExpectedPackets)
	}
}

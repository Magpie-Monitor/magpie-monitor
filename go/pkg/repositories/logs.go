package repositories

type Log interface {
	GetContent() string
}

// Split logs into subpackets of size not greater than maxPacketSize
func SplitLogsIntoPackets[T Log](logs []T, maxPacketSize int) [][]T {

	var logsPackets [][]T
	var lastPacketSize = 0
	var currentPacket []T

	for _, log := range logs {

		if lastPacketSize+len(log.GetContent()) > maxPacketSize {

			logsPackets = append(logsPackets, currentPacket)
			currentPacket = []T{log}
			lastPacketSize = len(log.GetContent())

		} else {
			currentPacket = append(currentPacket, log)
			lastPacketSize += len(log.GetContent())
		}

	}

	return logsPackets
}

package splitting

func SplitStringsIntoPackets(entries []string, maxPacketSize int) [][]string {

	var packets [][]string
	var lastPacketSize = 0
	var currentPacket []string

	for _, entry := range entries {

		if lastPacketSize+len(entry) > maxPacketSize {

			packets = append(packets, currentPacket)
			currentPacket = []string{entry}
			lastPacketSize = len(entry)

		} else {
			currentPacket = append(currentPacket, entry)
			lastPacketSize += len(entry)
		}

	}

	packets = append(packets, currentPacket)

	return packets
}

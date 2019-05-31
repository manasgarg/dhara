package stream

type MessageId struct {
	partOne uint64
	partTwo uint32
}

type Message struct {
	id   MessageId
	body []byte
}

func compareMessageIds(m1 *MessageId, m2 *MessageId) int {
	if m1.partOne < m2.partOne || (m1.partOne == m2.partOne && m1.partTwo < m2.partTwo) {
		return -1
	} else if m1.partOne == m2.partOne && m1.partTwo == m2.partTwo {
		return 0
	} else {
		return 1
	}
}

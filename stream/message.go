package stream

type MessageId struct {
	PartOne int64
	PartTwo int32
}

type Message struct {
	Id   MessageId
	Body []byte
}

func CompareMessageIds(m1 *MessageId, m2 *MessageId) int {
	if m1.PartOne < m2.PartOne || (m1.PartOne == m2.PartOne && m1.PartTwo < m2.PartTwo) {
		return -1
	} else if m1.PartOne == m2.PartOne && m1.PartTwo == m2.PartTwo {
		return 0
	} else {
		return 1
	}
}

package stream

import (
	"sort"
	"time"
)

type MemStream struct {
	id             string
	createTime     time.Time
	updateTime     time.Time
	messageCount   time.Time
	startMessageId *MessageId
	endMessageId   *MessageId

	messages []*Message
}

var (
	memStreams = make(map[string]*MemStream)
)

func GetMemStream(streamId string) (*MemStream, error) {
	if memStreams[streamId] != nil {
		return memStreams[streamId], nil
	}

	now := time.Now()
	s := &MemStream{
		id:         streamId,
		createTime: now,
		updateTime: now,
		messages:   make([]*Message, 0),
	}

	memStreams[streamId] = s

	return s, nil
}

func (s *MemStream) AddMessage(messageBody []byte) *MessageId {
	now := time.Now()
	timestamp := uint64(now.UnixNano() / 1000)
	var counter uint32 = 0

	if s.endMessageId != nil && s.endMessageId.partOne >= timestamp {
		s.endMessageId.partTwo++

		timestamp = s.endMessageId.partOne
		counter = s.endMessageId.partTwo
	}

	messageId := MessageId{timestamp, counter}

	s.messages = append(s.messages, &Message{id: messageId, body: messageBody})
	s.updateTime = now
	s.endMessageId = &messageId

	return &messageId
}

func (s *MemStream) GetMessages(startMessageId *MessageId, count int) []*Message {
	if len(s.messages) == 0 || compareMessageIds(startMessageId, s.endMessageId) == 1 {
		return []*Message{}
	}

	begin := sort.Search(len(s.messages), func(i int) bool {
		return compareMessageIds(&s.messages[i].id, startMessageId) >= 0
	})

	end := begin + count
	if end > len(s.messages)-begin {
		end = len(s.messages) - begin
	}

	messages := make([]*Message, end-begin)
	copy(messages, s.messages[begin:end])

	return messages
}

package stream

import (
	"sort"
	"time"
)

const (
	DEFAULT_MAX_MESSAGES = 10000
)

type MemStream struct {
	Id             string
	CreateTime     time.Time
	UpdateTime     time.Time
	MessageCount   time.Time
	StartMessageId *MessageId
	EndMessageId   *MessageId

	Messages []*Message
}

var (
	streams = make(map[string]*MemStream)
)

func GetMemStream(streamId string) (*MemStream, error) {
	if streams[streamId] != nil {
		return streams[streamId], nil
	}

	now := time.Now()
	s := &MemStream{
		Id:         streamId,
		CreateTime: now,
		UpdateTime: now,
		Messages:   make([]*Message, 0),
	}

	streams[streamId] = s

	return s, nil
}

func (s *MemStream) AddMessage(messageBody []byte) *MessageId {
	now := time.Now()
	timestamp := now.UnixNano() / 1000
	var counter int32 = 0

	if s.EndMessageId != nil && s.EndMessageId.PartOne >= timestamp {
		s.EndMessageId.PartTwo++

		timestamp = s.EndMessageId.PartOne
		counter = s.EndMessageId.PartTwo
	}

	messageId := MessageId{timestamp, counter}

	s.Messages = append(s.Messages, &Message{Id: messageId, Body: messageBody})
	s.UpdateTime = now
	s.EndMessageId = &messageId

	return &messageId
}

func (s *MemStream) GetMessages(startMessageId *MessageId, count int) []*Message {
	if len(s.Messages) == 0 || CompareMessageIds(startMessageId, s.EndMessageId) == 1 {
		return []*Message{}
	}

	begin := sort.Search(len(s.Messages), func(i int) bool {
		return CompareMessageIds(&s.Messages[i].Id, startMessageId) >= 0
	})

	end := begin + count
	if end > len(s.Messages)-begin {
		end = len(s.Messages) - begin
	}

	messages := make([]*Message, end-begin)
	copy(messages, s.Messages[begin:end])

	return messages
}

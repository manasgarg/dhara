package stream

import (
	"time"
)

type DiskStream struct {
	Id             string
	CreateTime     time.Time
	UpdateTime     time.Time
	MessageCount   time.Time
	StartMessageId MessageId
	EndMessageId   MessageId

	Dir                       string
	MaxMessageSize            int32
	MaxSegmentSize            int32
	MaxMessageCountPerSegment int32
	MaxDurationPerSegment     int32
	MessageDelimiter          string

	Segments []*Segment
}

package stream

type Segment struct {
	Id             string
	LogPath        string
	IndexPath      string
	StartMessageId MessageId
	EndMessageId   MessageId
	Size           int32
	MessageCount   int32
	Duration       int32
	IsOpen         bool
}

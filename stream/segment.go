package stream

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/manasgarg/dhara/utils"
)

const (
	indexVersionMajor uint16 = 0
	indexVersionMinor uint16 = 1
)

type DiskSegment struct {
	id             uint64
	logPath        string
	indexPath      string
	baseMessageId  MessageId
	startMessageId MessageId
	endMessageId   MessageId
	size           uint32
	messageCount   uint32
	duration       uint32
	isOpen         bool
}

type IndexRecord struct {
	idPartOne     uint32
	idPartTwo     uint32
	messageLength uint32
	logOffset     uint32
}

func LoadSegmentInfo(segmentId uint64, s *DiskStream) (*DiskSegment, error) {
	segment := &DiskSegment{
		id:        segmentId,
		logPath:   filepath.Join(s.dir, fmt.Sprintf("%s-%020d.log", s.id, segmentId)),
		indexPath: filepath.Join(s.dir, fmt.Sprintf("%s-%020d.index", s.id, segmentId)),
	}

	if _, err := os.Stat(segment.logPath); os.IsNotExist(err) {
		utils.SLogger.Errorw("Error to find log file for segment.", "log_path", segment.logPath)
		return nil, err
	}

	return segment, nil
}

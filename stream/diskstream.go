package stream

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/manasgarg/dhara/utils"
	"github.com/spf13/viper"
)

type DiskStream struct {
	id             string
	createTime     time.Time
	updateTime     time.Time
	messageCount   time.Time
	startMessageId MessageId
	endMessageId   MessageId

	dir                       string
	maxMessageSize            int32
	maxSegmentSize            int32
	maxMessageCountPerSegment int32
	maxDurationPerSegment     int32
	messageDelimiter          []byte

	segments       []*DiskSegment
	currentSegment *DiskSegment
}

var (
	diskStreams = make(map[string]*DiskStream)
)

func (s *DiskStream) initOnDisk() error {
	root := viper.GetString("stream_root_dir")

	s.dir = filepath.Join(root, s.id)
	if _, err := os.Stat(s.dir); os.IsNotExist(err) {
		err = os.MkdirAll(s.dir, 0755)
		if err != nil {
			utils.SLogger.Fatalw("Failed to create stream directory", "stream", s.id, "directory", s.dir)
			return err
		}
	}

	return nil
}

func (s *DiskStream) discoverSegments() error {
	files, _ := filepath.Glob(filepath.Join(s.dir, "*.index"))
	if len(files) == 0 {
	}

	sort.SliceStable(files, func(i, j int) bool { return files[i] < files[j] })
	for _, file := range files {
		var segmentNumber uint64 = 0
		_, err := fmt.Sscanf(file, "%d.index", &segmentNumber)
		if err != nil {
			utils.SLogger.Fatalw("Segment file seems corrupted. Failed to extract segment number.",
				"stream", s.id, "directory", s.dir, "segment_file", file)

			continue
		}

		if segment, err := LoadSegmentInfo(segmentNumber, s); err == nil {
			s.segments = append(s.segments, segment)
		}
	}

	return nil
}

func GetDiskStream(streamId string) (*DiskStream, error) {
	if diskStreams[streamId] != nil {
		return diskStreams[streamId], nil
	}

	return NewDiskStream(streamId)
}

func NewDiskStream(streamId string) (*DiskStream, error) {
	now := time.Now()
	s := &DiskStream{
		id:               streamId,
		createTime:       now,
		updateTime:       now,
		messageDelimiter: []byte{'\n'},
		segments:         make([]*DiskSegment, 0),
	}

	s.initOnDisk()
	s.discoverSegments()

	diskStreams[streamId] = s

	return s, nil
}

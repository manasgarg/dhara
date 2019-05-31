package stream

import (
	"bytes"
	"encoding/binary"
	"os"
	"sync"

	"github.com/manasgarg/dhara/utils"
)

var (
	dharaIndexMagic = []byte{'d', 'h', 'a', 'r', 'a', '-', '-', '-'}
)

type Index struct {
	path           string
	fp             *os.File
	startMessageId MessageId

	mux sync.Mutex
}

func Open(startMessageId MessageId, path string) (idx *Index, err error) {
	idx = &Index{
		startMessageId: startMessageId,
		path:           path,
	}

	idx.fp, err = os.Open(path)
	if err != nil {
		utils.SLogger.Errorw("Error in reading index file.", "index_path", path, "error", err)
		return nil, err
	}

	buf := make([]byte, 1024)
	n, err := idx.fp.Read(buf)
	if n != 1024 {
		utils.SLogger.Errorw("Corrupted index file.", "index_path", path, "error", err)
		return nil, err
	}

	if !bytes.Equal(dharaIndexMagic, buf[:10]) {
		utils.SLogger.Errorw("Corrupted index file. Magic code missing.", "index_path", path)
		return nil, err
	}

	major := binary.BigEndian.Uint16(buf[10:12])
	minor := binary.BigEndian.Uint16(buf[12:14])
	if major > indexVersionMajor || (major == indexVersionMajor && minor > indexVersionMinor) {
		utils.SLogger.Errorw("Index version too high", "major", major, "minor", minor,
			"supported_major", indexVersionMajor, "supported_minor", indexVersionMinor)
		return nil, err
	}

	return idx, nil
}

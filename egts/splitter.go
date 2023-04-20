package egts

import (
	"bufio"
	"encoding/binary"

	"github.com/gotrackery/protocol/common"
)

var (
	_ common.FrameSplitter = (*Splitter)(nil)
)

// Splitter implements common.FrameSplitter contract to extract EGTS data packet from incoming bytes.
type Splitter struct {
	badData []byte
	err     error
}

// NewSplitter creates a new Splitter instance for EGTS protocol.
func NewSplitter() *Splitter {
	return &Splitter{}
}

// Splitter implements bufio.SplitFunc contract to extract EGTS data packet from incoming bytes stream.
func (s *Splitter) Splitter() bufio.SplitFunc {
	return func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		const headerLen = 10
		if atEOF && len(data) == 0 {
			return 0, nil, nil
		}

		if len(data) > 1 && data[0] != allowedFirstByte1 {
			s.badData = data
			s.err = common.ErrBadData
			return 0, nil, s.err
		}

		if len(data) < headerLen {
			// Request more data.
			return 0, nil, nil
		}

		bodyLen := binary.LittleEndian.Uint16(data[5:7])
		pkgLen := uint16(data[3])
		if bodyLen > 0 {
			pkgLen += bodyLen + 2
		}
		if len(data) < int(pkgLen) {
			// Request more data.
			return 0, nil, nil
		}

		if atEOF {
			s.badData = data
			s.err = common.ErrBadData
			return 0, nil, s.err
		}

		// Finally got all data, return it.
		return int(pkgLen), data[0:pkgLen], nil
	}
}

// Error returns error if any registered.
// Use it to check that data corresponds to EGTS protocol.
func (s *Splitter) Error() error {
	return s.err
}

// BadData returns bad data if any registered.
// Use it to log which bytes couldn't be parsed as EGTS protocol.
func (s *Splitter) BadData() []byte {
	return s.badData
}

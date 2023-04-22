package wialonips

import (
	"bufio"
	"bytes"

	"github.com/gotrackery/protocol/common"
)

var (
	_ common.FrameSplitter = (*Splitter)(nil)
)

var (
	allowedFirstByte1 byte = 0x31 // 1 - version 1.1 UDP
	allowedFirstByte2 byte = 0x32 // 2 - version 2.0 UDP
)

// Splitter implements common.FrameSplitter contract to extract WialonIPS data packet from incoming bytes.
type Splitter struct {
	badData []byte
	err     error
}

// NewSplitter creates a new Splitter instance for WialonIPS protocol.
func NewSplitter() *Splitter {
	return &Splitter{}
}

// Splitter implements bufio.SplitFunc contract to extract WialonIPS data packet from incoming bytes stream.
func (s *Splitter) Splitter() bufio.SplitFunc {
	return func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		if atEOF && len(data) == 0 {
			return 0, nil, nil
		}
		if len(data) > 0 && allowedFirstByte(data[0]) {
			s.badData = data
			s.err = common.ErrBadData
			return 0, nil, s.err
		}
		if i := bytes.IndexByte(data, '\n'); i >= 0 {
			if len(data) > 1 && data[i-1] == '\r' {
				// We have a full newline-terminated line.
				return i + 1, data[0 : i+1], nil
			}
		}
		// If we're at EOF, we have a final, non-terminated line. Return it.
		if atEOF {
			s.badData = data
			s.err = common.ErrBadData
			return 0, nil, s.err
		}
		// Request more data.
		return 0, nil, nil
	}
}

func allowedFirstByte(first byte) bool {
	return first != packetTypeDelimiter[0] &&
		first != allowedFirstByte1 &&
		first != allowedFirstByte2 &&
		first != compressionMark
}

// Error returns error if any registered.
// Use it to check that data corresponds to WialonIPS protocol.
func (s *Splitter) Error() error {
	return s.err
}

// BadData returns bad data if any registered.
// Use it to log which bytes couldn't be parsed as WialonIPS protocol.
func (s *Splitter) BadData() []byte {
	return s.badData
}

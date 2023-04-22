package wialonretr

import (
	"bufio"
	"encoding/binary"

	"github.com/gotrackery/protocol/common"
)

var (
	_ common.FrameSplitter = (*Splitter)(nil)
)

// Splitter implements common.FrameSplitter contract to extract WialonRetranslator data packet from incoming bytes.
type Splitter struct {
	badData []byte
	err     error
}

// NewSplitter creates a new Splitter instance for WialonRetranslator protocol.
func NewSplitter() *Splitter {
	return &Splitter{}
}

// Splitter implements bufio.SplitFunc contract to extract WialonRetranslator data packet from incoming bytes stream.
// Data packets start with 4 little endian bytes of packet length.
func (s *Splitter) Splitter() bufio.SplitFunc {
	return func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		if atEOF && len(data) == 0 {
			return 0, nil, nil
		}

		if len(data) > 0 && data[0] == 0 {
			s.badData = data
			s.err = common.ErrBadData
			return 0, nil, s.err
		}

		// calc length and wait when whole packet is read
		var bodyLen uint32
		if len(data) > packageHeaderLen {
			bodyLen = binary.LittleEndian.Uint32(data[:packageHeaderLen])
		}
		totalLen := int(packageHeaderLen + bodyLen)
		if len(data) >= totalLen {
			return totalLen, data[:totalLen], nil
		}

		// If we're at EOF, we have a final.
		if atEOF {
			s.badData = data
			s.err = common.ErrBadData
			return 0, nil, s.err
		}
		// Request more data.
		return 0, nil, nil
	}
}

// Error returns error if any registered.
// Use it to check that data corresponds to WialonRetranslator protocol.
func (s *Splitter) Error() error {
	return s.err
}

// BadData returns bad data if error was registered.
// Use it to log which bytes couldn't be parsed as WialonRetranslator protocol.
func (s *Splitter) BadData() []byte {
	return s.badData
}

package wialonretr

import (
	"encoding/binary"
	"io"

	"github.com/gotrackery/protocol"
)

// ScanPackage implements bufio.SplitFunc contract to extract WialonRetranslator data packet from incoming bytes stream.
func ScanPackage(data []byte, atEOF bool) (advance int, token []byte, err error) {
	const (
		headerLen = 4
	)
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}

	if len(data) > 0 && data[0] == 0 {
		return 0, nil, protocol.ErrInconsistentData // Not possible to pass data to Scanner.Bytes().
		// It is possible to implement own scanner package to be able to pass data to Scanner.Bytes() and log it then.
	}

	// calc length and wait when whole packet is read
	var bodyLen uint32
	if len(data) > headerLen {
		bodyLen = binary.LittleEndian.Uint32(data[:headerLen])
	}
	totalLen := int(headerLen + bodyLen)
	if len(data) >= totalLen {
		return totalLen, data[:totalLen], nil
	}

	// If we're at EOF, we have a final.
	if atEOF {
		return len(data), data, io.EOF
	}
	// Request more data.
	return 0, nil, nil
}

package wialonretr

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"
	"time"

	"github.com/gotrackery/protocol"
)

const (
	packageHeaderLen = 4
)

const (
	InfoLocation = 1 << iota
	InfoDigitalInputs
	InfoDigitalOutputs
	InfoAlertBits
	InfoDriverID
)

var (
	eol       = []byte{0x00}       // end-of-line
	blockMark = []byte{0x0B, 0xBB} // block separator
)

// ScanPackage implements bufio.SplitFunc contract to extract WialonRetranslator data packet from incoming bytes stream.
// Data packets start with 4 little endian bytes of packet length.
func ScanPackage(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}

	if len(data) > 0 && data[0] == 0 {
		return 0, nil, protocol.ErrInconsistentData // Not possible to pass data to Scanner.Bytes().
		// It is possible to implement own scanner package to be able to pass data to Scanner.Bytes() and log it then.
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
		return len(data), data, io.EOF
	}
	// Request more data.
	return 0, nil, nil
}

type Package struct {
	IMEI             string
	RegisteredAt     time.Time
	MessagesBitFlags uint32 // bit mask
	err              error
}

func (p *Package) Decode(data []byte) error {
	body := data[packageHeaderLen:] // skip header length
	b, a, f := bytes.Cut(body, eol) // get imei
	if !f {
		p.err = ErrWialonRetranslatorBadDeviceID
		return p.err
	}
	p.IMEI = string(b)

	buf := bytes.NewReader(a)

	var utc int32 // get utc time
	if err := binary.Read(buf, binary.BigEndian, &utc); err != nil {
		p.err = fmt.Errorf("read utc time: %w", err)
		return p.err
	}
	p.RegisteredAt = time.Unix(int64(utc), 0)

	if err := binary.Read(buf, binary.BigEndian, &p.MessagesBitFlags); err != nil {
		p.err = fmt.Errorf("read bit flags of messages: %w", err)
		return p.err
	}

	scanner := bufio.NewScanner(buf)
	scanner.Split(scanBlock)
	for scanner.Scan() {
		block := scanner.Bytes()
		fmt.Println(hex.EncodeToString(block))
	}

	return nil
}

func scanBlock(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}

	if i := bytes.Index(data, blockMark); i != -1 {
		if i == 0 {
			return i + 2, nil, nil // skip first block separator
		}
		return i + 2, data[:i], nil
	}
	if atEOF { // return rest of data
		return len(data), data, nil
	}
	return 0, nil, nil
}

// Alternative implementation. More accuracy and more slow.
// func scanBlock(data []byte, atEOF bool) (advance int, token []byte, err error) {
// 	const blockHeaderLen = 6
// 	if atEOF && len(data) == 0 {
// 		return 0, nil, nil
// 	}
//
// 	// if there is no blockMark then it's not a Wialon Retranslator packet.
// 	if data[0] != blockMark[0] || data[1] != blockMark[1] {
// 		fmt.Println(data[0:2])
// 		return 0, nil, protocol.ErrInconsistentData
// 	}
//
// 	if len(data) < blockHeaderLen && !atEOF { // wait to get whole block header
// 		return 0, nil, nil
// 	}
//
// 	blockLen := binary.BigEndian.Uint32(data[2:6])
// 	ln := int(blockHeaderLen + blockLen)
// 	if len(data) < ln && !atEOF { // wait to get whole block
// 		return 0, nil, nil
// 	}
//
// 	if len(data) >= ln {
// 		return ln, data[:ln], nil
// 	}
//
// 	if atEOF { // return rest of data
// 		return len(data), data, nil
// 	}
// 	return 0, nil, nil
// }

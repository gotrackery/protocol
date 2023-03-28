package wialonips

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/gotrackery/protocol"
	"github.com/sigurn/crc16"
)

// PacketType is the type of package of Wialon IPS.
type PacketType string

const (
	UnknownPacket       PacketType = ""   // Unknown packet type. Used for initializing PacketType.
	LoginPacket         PacketType = "L"  // Login packet.
	ShortenedDataPacket PacketType = "SD" // Shortened data packet.
	DataPacket          PacketType = "D"  // Data packet.
	BlackBoxPacket      PacketType = "B"  // Black box packet.
	PingPacket          PacketType = "P"  // Ping packet.
)

var (
	packageDelimiter        = []byte{0x23} // #
	allowedFirstByte1  byte = 0x31         // 1 - version 1.1 UDP
	allowedFirstByte2  byte = 0x32         // 2 - version 2.0 UDP
	allowedFirstByteFF byte = 0xFF         // 0xFF - compressed
	crc16Table              = crc16.MakeTable(crc16.CRC16_ARC)
)

const (
	ddmmyyhhmiss = "020106150405"
)

// Version is the version of Wialon IPS protocol.
type Version int

const (
	UnknownVersion Version = iota // Unknown version. Used for initializing Version.
	V1_1                          // Version 1.1
	V2_0                          // Version 2.0
)

var (
	ErrWialonIPSUnsupportedPacketType = errors.New("unsupported packet type")
)

// String returns the string representation of the version.
func (s Version) String() string {
	if s < V1_1 || s > V2_0 {
		return fmt.Sprintf("Version(%d)", int(s))
	}
	var statuses = [...]string{"unknown", "1.1", "2.0"}
	return statuses[s]
}

// IsValid returns true if the version is valid.
func (s Version) IsValid() bool {
	switch s {
	case V1_1, V2_0:
		return true
	case UnknownVersion:
		return false
	}
	return false
}

// Message is the message of Wialon IPS protocol.
type Message interface {
	Decode(data []byte) error
	// Version returns version of Wialon IPS protocol. To avoid reflection use.
	Version() Version
	IMEI() string
	Response() []byte
}

// ScanPackage implements bufio.SplitFunc contract to extract Wialon IPS data packet from incoming bytes stream.
func ScanPackage(data []byte, atEOF bool) (advance int, token []byte, err error) { //nolint:cyclop
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}
	if len(data) > 0 &&
		data[0] != packageDelimiter[0] &&
		data[0] != allowedFirstByte1 &&
		data[0] != allowedFirstByte2 &&
		data[0] != allowedFirstByteFF {
		return len(data), data, protocol.ErrInconsistentData // let handle it, log for example
	}
	if i := bytes.IndexByte(data, '\n'); i >= 0 {
		if len(data) > 1 && data[i-1] == '\r' {
			// We have a full newline-terminated line.
			return i + 1, data[0 : i+1], nil
		}
	}
	// If we're at EOF, we have a final, non-terminated line. Return it.
	if atEOF {
		return len(data), data, nil
	}
	// Request more data.
	return 0, nil, nil
}

// Package is the package of Wialon IPS protocol.
// All data is received in text format as a packet which looks as follows:
// #PT#msgCRC\r\n.
type Package struct {
	Type    PacketType
	Version Version
	IMEI    string // IMEI is the unique identifier of the device.
	Message Message
}

// NewPackage creates a new package of Wialon IPS protocol.
func NewPackage(v Version, imei string) Package {
	return Package{Version: v, IMEI: imei}
}

// Decode decodes bytes to the package of Wialon IPS protocol.
func (p *Package) Decode(data []byte) error { //nolint:cyclop
	// ToDo add deflate
	bytesSet := bytes.SplitN(data, packageDelimiter, 3) //nolint:gomnd
	if len(bytesSet) != 3 {                             //nolint:gomnd
		return fmt.Errorf("invalid package structure: %w", protocol.ErrInconsistentData)
	}
	// ToDo UDP bytesSet[0] contains Protocol_version;imei - v2.0 and imei - v1.0
	p.parsePackageType(bytesSet[1])
	if p.Type == UnknownPacket {
		return ErrWialonIPSUnsupportedPacketType
	}

	msg := bytes.TrimRight(bytesSet[2], "\r\n")

	switch p.Type {
	case UnknownPacket:
		return ErrWialonIPSUnsupportedPacketType
	case LoginPacket:
		p.Message = &LoginMessage{message: message{imei: p.IMEI, ver: p.Version}}
	case ShortenedDataPacket:
		p.Message = &ShortenedDataMessage{message: message{imei: p.IMEI, ver: p.Version}}
	case DataPacket:
		p.Message = &DataMessage{ShortenedDataMessage: ShortenedDataMessage{message: message{imei: p.IMEI, ver: p.Version}}}
	case BlackBoxPacket:
		p.Message = &BlackBoxMessage{message: message{imei: p.IMEI, ver: p.Version}}
	case PingPacket:
		return ErrWialonIPSUnsupportedPacketType // TODO implement this.
	}

	err := p.Message.Decode(msg)
	if err != nil {
		p.Message = nil
		return fmt.Errorf("failed to decode message: %w", err)
	}

	if p.Type == LoginPacket {
		p.Version = p.Message.Version()
		p.IMEI = p.Message.IMEI()
	}

	return nil
}

func (p *Package) parsePackageType(data []byte) {
	p.Type = PacketType(data)
	switch p.Type { //nolint:exhaustive
	case LoginPacket, ShortenedDataPacket, DataPacket, BlackBoxPacket, PingPacket:
		return
	default:
		p.Type = UnknownPacket
	}
}

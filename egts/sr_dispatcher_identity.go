package egts

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

// SrDispatcherIdentity structure of subrecord of EGTS_SR_DISPATCHER_IDENTITY type, which is used
// only by the authorized TS when requesting authorization on the authorizing TS and contains credentials
// by the authorized ACH.
type SrDispatcherIdentity struct {
	DispatcherType uint8  `json:"DT"`
	DispatcherID   uint32 `json:"DID"`
	Description    string `json:"DSCR"`
}

// Decode parses the set of bytes into EGTS_SR_DISPATCHER_IDENTITY structure.
func (d *SrDispatcherIdentity) Decode(content []byte) error {
	var err error

	buf := bytes.NewBuffer(content)

	if d.DispatcherType, err = buf.ReadByte(); err != nil {
		return fmt.Errorf("failed to get dispatcher type: %w", err)
	}

	tmpIntBuf := make([]byte, 4)
	if _, err = buf.Read(tmpIntBuf); err != nil {
		return fmt.Errorf("failed to get a unique dispatcher ID: %w", err)
	}
	d.DispatcherID = binary.LittleEndian.Uint32(tmpIntBuf)

	d.Description = buf.String()

	return nil
}

// Encode returns the set of bytes of the EGTS_SR_DISPATCHER_IDENTITY structure.
func (d *SrDispatcherIdentity) Encode() ([]byte, error) {
	var (
		result []byte
		err    error
	)

	buf := new(bytes.Buffer)

	if err = buf.WriteByte(d.DispatcherType); err != nil {
		return result, fmt.Errorf("failed to record dispatcher type: %w", err)
	}

	if err = binary.Write(buf, binary.LittleEndian, d.DispatcherID); err != nil {
		return result, fmt.Errorf("failed to write unique dispatcher identifier: %w", err)
	}

	if _, err = buf.WriteString(d.Description); err != nil {
		return result, fmt.Errorf("failed to record a unique short description: %w", err)
	}

	return buf.Bytes(), nil
}

// Length returns the length of the EGTS_SR_DISPATCHER_IDENTITY structure.
func (d *SrDispatcherIdentity) Length() uint16 {
	var result uint16

	if recBytes, err := d.Encode(); err != nil {
		result = uint16(0)
	} else {
		result = uint16(len(recBytes))
	}

	return result
}

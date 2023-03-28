package egts

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

const (
	// Delimiter is the delimiter used in the EGTS packet.
	Delimiter       = uint8(0)
	textSectionSize = 10
)

// SrModuleData structure of EGTS_AUTH_SERVICE subrecord of EGTS_SR_MODULE_DATA type.
type SrModuleData struct {
	ModuleType      int8   `json:"MT"`
	VendorID        uint32 `json:"VID"`
	FirmwareVersion uint16 `json:"FWV"`
	SoftwareVersion uint16 `json:"SWV"`
	Modification    byte   `json:"MD"`
	State           byte   `json:"ST"`
	SerialNumber    string `json:"SRN"`
	_               byte   `json:"-"`
	Description     string `json:"DSCR"`
	_               byte   `json:"-"`
}

// Decode decodes EGTS_SR_MODULE_DATA subrecord from incoming bytes stream.
func (e *SrModuleData) Decode(content []byte) error {
	var err error
	buf := bytes.NewReader(content)

	moduleType, err := buf.ReadByte()
	if err != nil {
		return fmt.Errorf("failed to get the module type: %w", err)
	}
	e.ModuleType = int8(moduleType)

	tmpBuf := make([]byte, 4)
	if _, err = buf.Read(tmpBuf); err != nil {
		return fmt.Errorf("failed to get the vendor id: %w", err)
	}
	e.VendorID = binary.LittleEndian.Uint32(tmpBuf)

	tmpBuf = make([]byte, 2)
	if _, err = buf.Read(tmpBuf); err != nil {
		return fmt.Errorf("failed to get the firmware version: %w", err)
	}
	e.FirmwareVersion = binary.LittleEndian.Uint16(tmpBuf)

	tmpBuf = make([]byte, 2)
	if _, err = buf.Read(tmpBuf); err != nil {
		return fmt.Errorf("failed to get the software version: %w", err)
	}
	e.SoftwareVersion = binary.LittleEndian.Uint16(tmpBuf)

	e.Modification, err = buf.ReadByte()
	if err != nil {
		return fmt.Errorf("failed to get the modification: %w", err)
	}

	e.State, err = buf.ReadByte()
	if err != nil {
		return fmt.Errorf("failed to get the state: %w", err)
	}

	serialNumber := make([]byte, 0, textSectionSize)
	for {
		b, err := buf.ReadByte()
		if err != nil {
			return fmt.Errorf("failed to get the serial number: %w", err)
		}
		if b == Delimiter {
			break
		}
		serialNumber = append(serialNumber, b)
	}
	e.SerialNumber = string(serialNumber)

	description := make([]byte, 0, textSectionSize)
	for {
		b, err := buf.ReadByte()
		if err != nil {
			return fmt.Errorf("failed to get a description: %w", err)
		}
		if b == Delimiter {
			break
		}
		description = append(description, b)
	}
	e.Description = string(description)

	return nil
}

// Encode encodes EGTS_SR_MODULE_DATA subrecord to outgoing bytes stream.
func (e *SrModuleData) Encode() ([]byte, error) {
	var (
		result []byte
		err    error
	)
	buf := new(bytes.Buffer)

	if err = binary.Write(buf, binary.LittleEndian, e.ModuleType); err != nil {
		return result, fmt.Errorf("failed to write the module type: %w", err)
	}
	if err = binary.Write(buf, binary.LittleEndian, e.VendorID); err != nil {
		return result, fmt.Errorf("failed to write the vendor id: %w", err)
	}
	if err = binary.Write(buf, binary.LittleEndian, e.FirmwareVersion); err != nil {
		return result, fmt.Errorf("failed to write the firmware version: %w", err)
	}
	if err = binary.Write(buf, binary.LittleEndian, e.SoftwareVersion); err != nil {
		return result, fmt.Errorf("failed to write the software version: %w", err)
	}
	if err = binary.Write(buf, binary.LittleEndian, e.Modification); err != nil {
		return result, fmt.Errorf("failed to write the modification: %w", err)
	}
	if err = binary.Write(buf, binary.LittleEndian, e.State); err != nil {
		return result, fmt.Errorf("failed to write the state: %w", err)
	}
	if err = binary.Write(buf, binary.LittleEndian, []byte(e.SerialNumber)); err != nil {
		return result, fmt.Errorf("failed to write the serial number: %w", err)
	}
	if err = binary.Write(buf, binary.LittleEndian, Delimiter); err != nil {
		return result, fmt.Errorf("failed to write the delimiter: %w", err)
	}
	if err = binary.Write(buf, binary.LittleEndian, []byte(e.Description)); err != nil {
		return result, fmt.Errorf("failed to write the description: %w", err)
	}
	if err = binary.Write(buf, binary.LittleEndian, Delimiter); err != nil {
		return result, fmt.Errorf("failed to write the delimiter: %w", err)
	}

	return buf.Bytes(), nil
}

// Length returns the length of EGTS_SR_MODULE_DATA subrecord.
func (e *SrModuleData) Length() uint16 {
	var result uint16

	if recBytes, err := e.Encode(); err == nil {
		result = uint16(len(recBytes))
	}

	return result
}

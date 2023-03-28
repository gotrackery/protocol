package egts

import (
	"encoding/binary"
	"errors"
)

// SrAbsAnSensData is a subrecord structure of EGTS_SR_ABS_AN_SENS_DATA type, which is used by the subscriber's
// terminal to transmit data about the state of one analog input.
type SrAbsAnSensData struct {
	SensorNumber uint8  `json:"SensorNumber"`
	Value        uint32 `json:"Value"`
}

// Decode parses the set of bytes into EGTS_SR_ABS_AN_SENS_DATA structure.
func (e *SrAbsAnSensData) Decode(content []byte) error {
	if len(content) < int(e.Length()) {
		return errors.New("incorrect data size")
	}
	e.SensorNumber = content[0]
	e.Value = binary.LittleEndian.Uint32(content) >> 8
	return nil
}

// Encode returns the set of bytes of the EGTS_SR_ABS_AN_SENS_DATA structure.
func (e *SrAbsAnSensData) Encode() ([]byte, error) {
	return []byte{
		e.SensorNumber,
		byte(e.Value),
		byte(e.Value >> 8),
		byte(e.Value >> 16),
	}, nil
}

// Length returns the length of the EGTS_SR_ABS_AN_SENS_DATA structure.
func (e *SrAbsAnSensData) Length() uint16 {
	return 4
}

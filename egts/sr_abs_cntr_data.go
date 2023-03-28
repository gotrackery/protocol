package egts

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

// SrAbsCntrData structure of EGTS_SR_ABS_CNTR_DATA type subrecord, which is used
// subscriber terminal to transmit data to the hardware and software complex about
// State of one count input.
type SrAbsCntrData struct {
	CounterNumber uint8  `json:"CN"`
	CounterValue  uint32 `json:"CNV"`
}

// Decode parses the set of bytes into EGTS_SR_ABS_CNTR_DATA structure.
func (e *SrAbsCntrData) Decode(content []byte) error {
	var (
		err error
	)
	buf := bytes.NewReader(content)

	if e.CounterNumber, err = buf.ReadByte(); err != nil {
		return fmt.Errorf("failed to get a count input number: %w", err)
	}

	counterVal := make([]byte, 3)
	if _, err = buf.Read(counterVal); err != nil {
		return fmt.Errorf("failed to get the value of the counting input reading: %w", err)
	}

	counterVal = append(counterVal, 0x00)
	e.CounterValue = binary.LittleEndian.Uint32(counterVal)

	return nil
}

// Encode returns the set of bytes of EGTS_SR_ABS_CNTR_DATA structure.
func (e *SrAbsCntrData) Encode() ([]byte, error) {
	var (
		err    error
		result []byte
	)
	buf := new(bytes.Buffer)

	if err = buf.WriteByte(e.CounterNumber); err != nil {
		return result, fmt.Errorf("failed to write the number of the counting input: %w", err)
	}

	counterVal := make([]byte, 4)
	binary.LittleEndian.PutUint32(counterVal, e.CounterValue)
	if _, err = buf.Write(counterVal[:3]); err != nil {
		return result, fmt.Errorf("failed to write the number of the counting input: %w", err)
	}

	result = buf.Bytes()
	return result, nil
}

// Length returns the length of EGTS_SR_ABS_CNTR_DATA structure.
func (e *SrAbsCntrData) Length() uint16 {
	var result uint16

	if recBytes, err := e.Encode(); err != nil {
		result = uint16(0)
	} else {
		result = uint16(len(recBytes))
	}

	return result
}

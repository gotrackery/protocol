package egts

import (
	"bytes"
	"fmt"
)

// SrResultCode is the structure of EGTS_SR_RESULT_CODE subrecord, which is used by the telematics
// platform to inform the AC about the results of the AC authentication procedure.
type SrResultCode struct {
	ResultCode uint8 `json:"RCD"`
}

// Decode parses the set of bytes into EGTS_SR_RESULT_CODE structure.
func (s *SrResultCode) Decode(content []byte) error {
	var (
		err error
	)
	buf := bytes.NewBuffer(content)

	if s.ResultCode, err = buf.ReadByte(); err != nil {
		return fmt.Errorf("failed to get result code: %w", err)
	}

	return nil
}

// Encode returns the set of bytes of the EGTS_SR_RESULT_CODE structure.
func (s *SrResultCode) Encode() ([]byte, error) {
	var (
		result []byte
		err    error
	)
	buf := new(bytes.Buffer)

	if err = buf.WriteByte(s.ResultCode); err != nil {
		return result, fmt.Errorf("failed to write result code: %w", err)
	}

	result = buf.Bytes()
	return result, nil
}

// Length returns the length of the EGTS_SR_RESULT_CODE structure.
func (s *SrResultCode) Length() uint16 {
	var result uint16

	if recBytes, err := s.Encode(); err != nil {
		result = uint16(0)
	} else {
		result = uint16(len(recBytes))
	}

	return result
}

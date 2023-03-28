package egts

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

// SrResponse subrecord structure of EGTS_SR_RESPONSE type, which is used to confirm
// reception of the results of service support processing.
type SrResponse struct {
	ConfirmedRecordNumber uint16 `json:"CRN"`
	RecordStatus          uint8  `json:"RST"`
}

// Decode parses the set of bytes into EGTS_SR_RESPONSE structure.
func (s *SrResponse) Decode(content []byte) error {
	var (
		err error
	)
	buf := bytes.NewBuffer(content)

	tmpIntBuf := make([]byte, 2)
	if _, err = buf.Read(tmpIntBuf); err != nil {
		return fmt.Errorf("failed to get to be confirmed record number: %w", err)
	}
	s.ConfirmedRecordNumber = binary.LittleEndian.Uint16(tmpIntBuf)

	if s.RecordStatus, err = buf.ReadByte(); err != nil {
		return fmt.Errorf("failed to get record processing status: %w", err)
	}

	sfd := ServiceDataSet{}
	if err = sfd.Decode(buf.Bytes()); err != nil {
		return fmt.Errorf("failed to decode service data set: %w", err)
	}
	return nil
}

// Encode returns the set of bytes of EGTS_SR_RESPONSE structure.
func (s *SrResponse) Encode() ([]byte, error) {
	var (
		result []byte
		err    error
	)
	buf := new(bytes.Buffer)

	if err = binary.Write(buf, binary.LittleEndian, s.ConfirmedRecordNumber); err != nil {
		return result, fmt.Errorf("failed to write the number of the record to be confirmed: %w", err)
	}

	if err = buf.WriteByte(s.RecordStatus); err != nil {
		return result, fmt.Errorf("failed to write processing status: %w", err)
	}

	result = buf.Bytes()
	return result, nil
}

// Length returns the length of EGTS_SR_RESPONSE structure.
func (s *SrResponse) Length() uint16 {
	var result uint16

	if recBytes, err := s.Encode(); err != nil {
		result = uint16(0)
	} else {
		result = uint16(len(recBytes))
	}

	return result
}

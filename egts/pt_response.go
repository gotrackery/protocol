package egts

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

// PtResponse substructure of EGTS_PT_RESPONSE type.
type PtResponse struct {
	ResponsePacketID uint16     `json:"RPID"`
	ProcessingResult uint8      `json:"PR"`
	SDR              BinaryData `json:"SDR"`
}

// Decode decodes the bytes into EGTS_PT_RESPONSE type struct.
func (s *PtResponse) Decode(content []byte) error {
	var (
		err error
	)
	buf := bytes.NewBuffer(content)

	tmpIntBuf := make([]byte, 2)
	if _, err = buf.Read(tmpIntBuf); err != nil {
		return fmt.Errorf("failed to get the packet identifier from the response: %w", err)
	}
	s.ResponsePacketID = binary.LittleEndian.Uint16(tmpIntBuf)

	if s.ProcessingResult, err = buf.ReadByte(); err != nil {
		return fmt.Errorf("failed to get processing result code: %w", err)
	}

	// if there is a service level, because it is optional
	if buf.Len() > 0 {
		s.SDR = &ServiceDataSet{}
		if err = s.SDR.Decode(buf.Bytes()); err != nil {
			return fmt.Errorf("failed to decode service data set: %w", err)
		}
	}

	return nil
}

// Encode encodes the EGTS_PT_RESPONSE type struct into bytes.
func (s *PtResponse) Encode() ([]byte, error) {
	var (
		result   []byte
		sdrBytes []byte
		err      error
	)
	buf := new(bytes.Buffer)

	if err = binary.Write(buf, binary.LittleEndian, s.ResponsePacketID); err != nil {
		return result, fmt.Errorf("failed to write packet identifier in response: %w", err)
	}

	if err = buf.WriteByte(s.ProcessingResult); err != nil {
		return result, fmt.Errorf("failed to write the result of processing to the package: %w", err)
	}

	if s.SDR != nil {
		if sdrBytes, err = s.SDR.Encode(); err != nil {
			return result, fmt.Errorf("failed to encode service data set: %w", err)
		}
		buf.Write(sdrBytes)
	}

	result = buf.Bytes()
	return result, nil
}

// Length returns the length of the EGTS_PT_RESPONSE type struct.
func (s *PtResponse) Length() uint16 {
	var result uint16

	if recBytes, err := s.Encode(); err != nil {
		result = uint16(0)
	} else {
		result = uint16(len(recBytes))
	}

	return result
}

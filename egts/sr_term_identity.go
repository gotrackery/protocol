package egts

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strconv"
)

// SrTermIdentity structure of EGTS_SR_TERM_IDENTITY type subrecord, which is used by the AC when requesting
// authorization to the telematics platform and contains the AC credentials.
type SrTermIdentity struct {
	TerminalIdentifier       uint32 `json:"TID"`
	MNE                      string `json:"MNE"`
	BSE                      string `json:"BSE"`
	NIDE                     string `json:"NIDE"`
	SSRA                     string `json:"SSRA"`
	LNGCE                    string `json:"LNGCE"`
	IMSIE                    string `json:"IMSIE"`
	IMEIE                    string `json:"IMEIE"`
	HDIDE                    string `json:"HDIDE"`
	HomeDispatcherIdentifier uint16 `json:"HDID"`
	IMEI                     string `json:"IMEI"`
	IMSI                     string `json:"IMSI"`
	LanguageCode             string `json:"LNGC"`
	NetworkIdentifier        []byte `json:"NID"`
	BufferSize               uint16 `json:"BS"`
	MobileNumber             string `json:"MSISDN"`
}

// Decode parses the bytes into EGTS_SR_TERM_IDENTITY structure.
func (e *SrTermIdentity) Decode(content []byte) error {
	var (
		err   error
		flags byte
	)
	buf := bytes.NewReader(content)

	tmpBuf := make([]byte, 4)
	if _, err = buf.Read(tmpBuf); err != nil {
		return fmt.Errorf("failed to get terminal ID on authorization: %w", err)
	}
	e.TerminalIdentifier = binary.LittleEndian.Uint32(tmpBuf)

	if flags, err = buf.ReadByte(); err != nil {
		return fmt.Errorf("failed to read the flags byte term identify: %w", err)
	}
	flagBits := fmt.Sprintf("%08b", flags)
	e.MNE = flagBits[:1]
	e.BSE = flagBits[1:2]
	e.NIDE = flagBits[2:3]
	e.SSRA = flagBits[3:4]
	e.LNGCE = flagBits[4:5]
	e.IMSIE = flagBits[5:6]
	e.IMEIE = flagBits[6:7]
	e.HDIDE = flagBits[7:]

	if e.HDIDE == "1" {
		tmpBuf = make([]byte, 2)
		if _, err = buf.Read(tmpBuf); err != nil {
			return fmt.Errorf("failed to get the home telematics platform ID during authorization: %w", err)
		}
		e.HomeDispatcherIdentifier = binary.LittleEndian.Uint16(tmpBuf)
	}

	if e.IMEIE == "1" {
		tmpBuf = make([]byte, 15)
		if _, err = buf.Read(tmpBuf); err != nil {
			return fmt.Errorf("failed to get IMEI when authorizing: %w", err)
		}
		e.IMEI = string(tmpBuf)
	}

	if e.IMSIE == "1" {
		tmpBuf = make([]byte, 16)
		if _, err = buf.Read(tmpBuf); err != nil {
			return fmt.Errorf("failed to get IMSI during authorization: %w", err)
		}
		e.IMSI = string(tmpBuf)
	}

	if e.LNGCE == "1" {
		tmpBuf = make([]byte, 3)
		if _, err = buf.Read(tmpBuf); err != nil {
			return fmt.Errorf("failed to get language code during authorization: %w", err)
		}
		e.LanguageCode = string(tmpBuf)
	}

	if e.NIDE == "1" {
		e.NetworkIdentifier = make([]byte, 3)
		if _, err = buf.Read(e.NetworkIdentifier); err != nil {
			return fmt.Errorf("failed to get network ID code when authorizing: %w", err)
		}
	}

	if e.BSE == "1" {
		tmpBuf = make([]byte, 2)
		if _, err = buf.Read(tmpBuf); err != nil {
			return fmt.Errorf("failed to get maximum buffer size during authorization: %w", err)
		}
		e.BufferSize = binary.LittleEndian.Uint16(tmpBuf)
	}

	if e.MNE == "1" {
		tmpBuf = make([]byte, 15)
		if _, err = buf.Read(tmpBuf); err != nil {
			return fmt.Errorf("failed to get the phone number of the mobile subscriber: %w", err)
		}
		e.MobileNumber = string(tmpBuf)
	}

	return nil
}

// Encode returns the bytes of the EGTS_SR_TERM_IDENTITY structure.
func (e *SrTermIdentity) Encode() ([]byte, error) {
	var (
		result []byte
		flags  uint64
		err    error
	)
	buf := new(bytes.Buffer)

	if err = binary.Write(buf, binary.LittleEndian, e.TerminalIdentifier); err != nil {
		return result, fmt.Errorf("failed to write terminal ID during authorization: %w", err)
	}

	flags, err = strconv.ParseUint(e.MNE+e.BSE+e.NIDE+e.SSRA+e.LNGCE+e.IMSIE+e.IMEIE+e.HDIDE, 2, 8)
	if err != nil {
		return result, fmt.Errorf("failed to convert subrecord flags to a number: %w", err)
	}
	if err = buf.WriteByte(uint8(flags)); err != nil {
		return result, fmt.Errorf("failed to write the flags byte term identify: %w", err)
	}

	if e.HDIDE == "1" {
		if err = binary.Write(buf, binary.LittleEndian, e.HomeDispatcherIdentifier); err != nil {
			return result, fmt.Errorf(
				"failed to write the ID of the home telematics platform during authorization: %w", err)
		}
	}

	if e.IMEIE == "1" {
		if _, err = buf.Write([]byte(e.IMEI)); err != nil {
			return result, fmt.Errorf("failed to write IMEI when authorizing: %w", err)
		}
	}

	if e.IMSIE == "1" {
		if _, err = buf.Write([]byte(e.IMSI)); err != nil {
			return result, fmt.Errorf("failed to write IMSI during authorization: %w", err)
		}
	}

	if e.LNGCE == "1" {
		if _, err = buf.Write([]byte(e.LanguageCode)); err != nil {
			return result, fmt.Errorf("failed to write IMSI during authorization: %w", err)
		}
	}

	if e.NIDE == "1" {
		if _, err = buf.Write(e.NetworkIdentifier); err != nil {
			return result, fmt.Errorf("failed to write operator network ID code during authorization: %w", err)
		}
	}

	if e.BSE == "1" {
		if err = binary.Write(buf, binary.LittleEndian, e.BufferSize); err != nil {
			return result, fmt.Errorf("failed to write maximum buffer size during authorization: %w", err)
		}
	}

	if e.MNE == "1" {
		if _, err = buf.Write([]byte(e.MobileNumber)); err != nil {
			return result, fmt.Errorf("failed to record the phone number of the mobile subscriber: %w", err)
		}
	}

	result = buf.Bytes()
	return result, nil
}

// Length gets the length of the encoded subrecord.
func (e *SrTermIdentity) Length() uint16 {
	var result uint16

	if recBytes, err := e.Encode(); err != nil {
		result = uint16(0)
	} else {
		result = uint16(len(recBytes))
	}

	return result
}

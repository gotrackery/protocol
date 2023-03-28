package egts

import (
	"bytes"
	"fmt"
	"strings"
)

// SrAuthInfo subrecord structure of EGTS_SR_AUTH_INFO type, which is intended for transmitting to
// the telematics platform of the CA authentication data using the previously transmitted
// from the platform's side to implement data encryption.
type SrAuthInfo struct {
	UserName       string `json:"UNM"`
	UserPassword   string `json:"UPSW"`
	ServerSequence string `json:"SS"`
}

// Decode parses the set of bytes into EGTS_SR_AUTH_INFO structure.
func (e *SrAuthInfo) Decode(content []byte) error {
	var (
		err    error
		tmpStr string
	)
	// string field separator from GOST 54619 - 2011 section EGTS_SR_AUTH_INFO
	sep := byte(0x00)

	buf := bytes.NewBuffer(content)
	tmpStr, err = buf.ReadString(sep)
	if err != nil {
		return fmt.Errorf("failed to read user name sr_auth_info: %w", err)
	}
	e.UserName = strings.TrimSuffix(tmpStr, string(sep))

	tmpStr, err = buf.ReadString(sep)
	if err != nil {
		return fmt.Errorf("failed to read password sr_auth_info: %w", err)
	}
	e.UserPassword = strings.TrimSuffix(tmpStr, string(sep))

	if buf.Len() > 0 {
		tmpStr, err = buf.ReadString(sep)
		if err != nil {
			return fmt.Errorf("failed to read SS from sr_auth_info: %w", err)
		}
		e.ServerSequence = strings.TrimSuffix(tmpStr, string(sep))
	}

	return nil
}

// Encode encodes the EGTS_SR_AUTH_INFO structure into the set of bytes.
func (e *SrAuthInfo) Encode() ([]byte, error) {
	var (
		result []byte
	)
	// string field separator from GOST 54619 - 2011 section EGTS_SR_AUTH_INFO
	sep := byte(0x00)

	result = append(result, []byte(e.UserName)...)
	result = append(result, sep)

	result = append(result, []byte(e.UserPassword)...)
	result = append(result, sep)

	// optional field, availability depends on encryption algorithm used
	// special server sequence of bytes, passed in EGTS_SR_AUTH_PARAMS subrecord
	if e.ServerSequence != "" {
		result = append(result, []byte(e.ServerSequence)...)
		result = append(result, sep)
	}

	return result, nil
}

// Length returns the length of the EGTS_SR_AUTH_INFO structure.
func (e *SrAuthInfo) Length() uint16 {
	var result uint16

	if recBytes, err := e.Encode(); err != nil {
		result = uint16(0)
	} else {
		result = uint16(len(recBytes))
	}

	return result
}

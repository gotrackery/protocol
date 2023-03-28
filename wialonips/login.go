package wialonips

import (
	"bytes"
	"fmt"
)

var _ Message = (*LoginMessage)(nil)

// LoginMessage is a WialonIPS login message.
// The packet is used for the device authorization on the server.
// Every TCP connection starts with sending this packet from the device to the server.
// Other data should be transferred only after the server confirms the successful authorization of the device.
// The login package looks as follows:
// #L#Protocol_version;IMEI;Password;CRC16\r\n.
type LoginMessage struct {
	Password string
	message
}

// Decode decodes a WialonIPS message.
func (l *LoginMessage) Decode(data []byte) error {
	l.ver = UnknownVersion
	bytesSet := bytes.Split(data, fieldsDelimiter)
	length := len(bytesSet)
	if length < 2 { //nolint:gomnd
		l.err = ErrWialonIPSInvalidLoginMessage // "0"
		return l.err
	}

	if length == 2 { //nolint:gomnd
		l.ver = V1_1
	} else {
		l.ver = V2_0 // bytesSet[0] contains version number
		bytesSet = bytesSet[1:length]
	}

	if l.ver == V2_0 {
		l.err = validateCRC(data, fieldsDelimiter)
		if l.err != nil {
			return l.err
		}
	}

	l.imei = string(bytesSet[0])
	l.Password = string(bytesSet[1])

	return nil
}

// Response returns a WialonIPS response message.
func (l *LoginMessage) Response() []byte {
	return []byte(fmt.Sprintf(responseTemplate, LoginPacket, MapErrToRespCode(LoginPacket, l.err)))
}

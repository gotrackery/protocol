package wialonips

import (
	"bytes"
	"fmt"

	"github.com/sigurn/crc16"
)

const (
	na               = "NA"
	intType          = 0x31 // 1
	floatType        = 0x32 // 2
	stringType       = 0x33 // 3
	responseTemplate = "#A%s#%s\r\n"
)

var (
	fieldsDelimiter   = []byte{0x3b} // ;
	analogDelimiter   = []byte{0x2c} // ,
	valuesDelimiter   = []byte{0x2c} // ,
	paramsDelimiter   = []byte{0x3a} // :
	blackBoxDelimiter = []byte{0x7c} // |
)

type message struct {
	imei string
	ver  Version
	err  error
}

func (m *message) Version() Version {
	return m.ver
}

func (m *message) IMEI() string {
	return m.imei
}

func (m *message) Error() error {
	return m.err
}

func validateCRC(data []byte, delim []byte) error {
	i := bytes.LastIndex(data, delim)
	crc, err := parseCRC(data[i+1:])
	if err != nil {
		return fmt.Errorf("failed parsing crc %s: %w", string(data[i+1:]), ErrWialonIPSCRC16Validation)
	}
	if crc != crc16.Checksum(data[:i+1], crc16Table) {
		return fmt.Errorf("failed crc validation: %w", ErrWialonIPSCRC16Validation)
	}
	return nil
}

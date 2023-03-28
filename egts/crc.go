package egts

import (
	"github.com/sigurn/crc16"
	"github.com/sigurn/crc8"
)

var (
	crc8Table  = crc8.MakeTable(crc8.Params{Poly: 0x31, Init: 0xFF, Check: 0xF7, Name: "CRC-8"})
	crc16Table = crc16.MakeTable(crc16.CRC16_CCITT_FALSE)
)

// CRC8 calculates the CRC-8 checksum of the data.
func CRC8(data []byte) byte {
	return crc8.Checksum(data, crc8Table)
}

// CRC16 calculates the CRC-16 checksum of the data.
func CRC16(data []byte) uint16 {
	return crc16.Checksum(data, crc16Table)
}

package egts

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_crc8(t *testing.T) {
	crc := CRC8([]byte("123456789"))
	checkVal := byte(0xF7)

	assert.Equal(t, crc, checkVal)
}

func Test_crc16(t *testing.T) {
	crc := CRC16([]byte("123456789"))
	checkVal := uint16(0x29b1)

	assert.Equal(t, crc, checkVal)
}

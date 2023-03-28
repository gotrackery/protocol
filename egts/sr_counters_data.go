package egts

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strconv"
)

// SrCountersData is a subrecord structure of the EGTS_SR_COUNTERS_DATA type, which is used by the hardware and
// System for transmitting the count inputs values to the subscriber's terminal.
type SrCountersData struct {
	CounterFieldExists1 string `json:"CFE1"`
	CounterFieldExists2 string `json:"CFE2"`
	CounterFieldExists3 string `json:"CFE3"`
	CounterFieldExists4 string `json:"CFE4"`
	CounterFieldExists5 string `json:"CFE5"`
	CounterFieldExists6 string `json:"CFE6"`
	CounterFieldExists7 string `json:"CFE7"`
	CounterFieldExists8 string `json:"CFE8"`
	Counter1            uint32 `json:"CN1"`
	Counter2            uint32 `json:"CN2"`
	Counter3            uint32 `json:"CN3"`
	Counter4            uint32 `json:"CN4"`
	Counter5            uint32 `json:"CN5"`
	Counter6            uint32 `json:"CN6"`
	Counter7            uint32 `json:"CN7"`
	Counter8            uint32 `json:"CN8"`
}

// Decode decodes the EGTS_SR_COUNTERS_DATA subrecord into SrCountersData struct.
func (c *SrCountersData) Decode(content []byte) error {
	var (
		err        error
		flags      byte
		counterVal []byte
	)
	buf := bytes.NewReader(content)

	if flags, err = buf.ReadByte(); err != nil {
		return fmt.Errorf("failed to get a byte of digital outputs sr_counters_data: %w", err)
	}
	flagBits := fmt.Sprintf("%08b", flags)

	c.CounterFieldExists8 = flagBits[:1]
	c.CounterFieldExists7 = flagBits[1:2]
	c.CounterFieldExists6 = flagBits[2:3]
	c.CounterFieldExists5 = flagBits[3:4]
	c.CounterFieldExists4 = flagBits[4:5]
	c.CounterFieldExists3 = flagBits[5:6]
	c.CounterFieldExists2 = flagBits[6:7]
	c.CounterFieldExists1 = flagBits[7:]

	tmpBuf := make([]byte, 3)
	if c.CounterFieldExists1 == "1" {
		if _, err = buf.Read(tmpBuf); err != nil {
			return fmt.Errorf("failed to get CN1 reading: %w", err)
		}
		counterVal = append(tmpBuf, 0x00)
		c.Counter1 = binary.LittleEndian.Uint32(counterVal)
	}

	if c.CounterFieldExists2 == "1" {
		if _, err = buf.Read(tmpBuf); err != nil {
			return fmt.Errorf("failed to get CN2 reading: %w", err)
		}
		counterVal = append(tmpBuf, 0x00)
		c.Counter1 = binary.LittleEndian.Uint32(counterVal)
	}

	if c.CounterFieldExists3 == "1" {
		if _, err = buf.Read(tmpBuf); err != nil {
			return fmt.Errorf("failed to get CN3 reading: %w", err)
		}
		counterVal = append(tmpBuf, 0x00)
		c.Counter3 = binary.LittleEndian.Uint32(counterVal)
	}

	if c.CounterFieldExists4 == "1" {
		if _, err = buf.Read(tmpBuf); err != nil {
			return fmt.Errorf("failed to get CN4 reading: %w", err)
		}
		counterVal = append(tmpBuf, 0x00)
		c.Counter4 = binary.LittleEndian.Uint32(counterVal)
	}

	if c.CounterFieldExists5 == "1" {
		if _, err = buf.Read(tmpBuf); err != nil {
			return fmt.Errorf("failed to get CN5 reading: %w", err)
		}
		counterVal = append(tmpBuf, 0x00)
		c.Counter5 = binary.LittleEndian.Uint32(counterVal)
	}

	if c.CounterFieldExists6 == "1" {
		if _, err = buf.Read(tmpBuf); err != nil {
			return fmt.Errorf("failed to get CN6 reading: %w", err)
		}
		counterVal = append(tmpBuf, 0x00)
		c.Counter6 = binary.LittleEndian.Uint32(counterVal)
	}

	if c.CounterFieldExists7 == "1" {
		if _, err = buf.Read(tmpBuf); err != nil {
			return fmt.Errorf("failed to get CN7 reading: %w", err)
		}
		counterVal = append(tmpBuf, 0x00)
		c.Counter7 = binary.LittleEndian.Uint32(counterVal)
	}

	if c.CounterFieldExists8 == "1" {
		if _, err = buf.Read(tmpBuf); err != nil {
			return fmt.Errorf("failed to get CN8 reading: %w", err)
		}
		counterVal = append(tmpBuf, 0x00)
		c.Counter8 = binary.LittleEndian.Uint32(counterVal)
	}
	return nil
}

// Encode encodes the SrCountersData struct into EGTS_SR_COUNTERS_DATA subrecord.
func (c *SrCountersData) Encode() ([]byte, error) {
	var (
		err    error
		flags  uint64
		result []byte
	)
	buf := new(bytes.Buffer)
	flagsBits := c.CounterFieldExists8 +
		c.CounterFieldExists7 +
		c.CounterFieldExists6 +
		c.CounterFieldExists5 +
		c.CounterFieldExists4 +
		c.CounterFieldExists3 +
		c.CounterFieldExists2 +
		c.CounterFieldExists1

	if flags, err = strconv.ParseUint(flagsBits, 2, 8); err != nil {
		return result, fmt.Errorf("failed to generate a byte of analog outputs counters_data: %w", err)
	}
	if err = buf.WriteByte(uint8(flags)); err != nil {
		return result, fmt.Errorf("failed to write bytes of analog outputs counters_data: %w", err)
	}

	sensVal := make([]byte, 4)
	if c.CounterFieldExists1 == "1" {
		binary.LittleEndian.PutUint32(sensVal, c.Counter1)
		if _, err = buf.Write(sensVal[:3]); err != nil {
			return result, fmt.Errorf("failed to write CN1 reading: %w", err)
		}
	}

	if c.CounterFieldExists2 == "1" {
		binary.LittleEndian.PutUint32(sensVal, c.Counter2)
		if _, err = buf.Write(sensVal[:3]); err != nil {
			return result, fmt.Errorf("failed to write CN2 reading: %w", err)
		}
	}

	if c.CounterFieldExists3 == "1" {
		binary.LittleEndian.PutUint32(sensVal, c.Counter3)
		if _, err = buf.Write(sensVal[:3]); err != nil {
			return result, fmt.Errorf("failed to write CN3 reading: %w", err)
		}
	}

	if c.CounterFieldExists4 == "1" {
		binary.LittleEndian.PutUint32(sensVal, c.Counter4)
		if _, err = buf.Write(sensVal[:3]); err != nil {
			return result, fmt.Errorf("failed to write CN4 reading: %w", err)
		}
	}

	if c.CounterFieldExists5 == "1" {
		binary.LittleEndian.PutUint32(sensVal, c.Counter5)
		if _, err = buf.Write(sensVal[:3]); err != nil {
			return result, fmt.Errorf("failed to write CN5 reading: %w", err)
		}
	}

	if c.CounterFieldExists6 == "1" {
		binary.LittleEndian.PutUint32(sensVal, c.Counter6)
		if _, err = buf.Write(sensVal[:3]); err != nil {
			return result, fmt.Errorf("failed to write CN6 reading: %w", err)
		}
	}

	if c.CounterFieldExists7 == "1" {
		binary.LittleEndian.PutUint32(sensVal, c.Counter7)
		if _, err = buf.Write(sensVal[:3]); err != nil {
			return result, fmt.Errorf("failed to write CN7 reading: %w", err)
		}
	}

	if c.CounterFieldExists8 == "1" {
		binary.LittleEndian.PutUint32(sensVal, c.Counter8)
		if _, err = buf.Write(sensVal[:3]); err != nil {
			return result, fmt.Errorf("failed to write CN8 reading: %w", err)
		}
	}

	result = buf.Bytes()

	return result, nil
}

// Length returns the length of the EGTS_SR_COUNTERS_DATA.
func (c *SrCountersData) Length() uint16 {
	var result uint16

	if recBytes, err := c.Encode(); err != nil {
		result = uint16(0)
	} else {
		result = uint16(len(recBytes))
	}

	return result
}

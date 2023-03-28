package egts

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strconv"
)

// SrLiquidLevelSensor subrecord structure of EGTS_SR_LIQUID_LEVEL_SENSOR type, which is used
// subscriber terminal to transmit the data on DUH readings to the hardware-software complex.
type SrLiquidLevelSensor struct {
	LiquidLevelSensorErrorFlag string `json:"LLSEF"`
	LiquidLevelSensorValueUnit string `json:"LLSVU"`
	RawDataFlag                string `json:"RDF"`
	LiquidLevelSensorNumber    uint8  `json:"LLSN"`
	ModuleAddress              uint16 `json:"MADDR"`
	LiquidLevelSensorData      uint32 `json:"LLSD"`
}

// Decode parses the set of bytes into EGTS_SR_LIQUID_LEVEL_SENSOR structure.
func (e *SrLiquidLevelSensor) Decode(content []byte) error {
	var (
		err     error
		flags   byte
		sensNum uint64
	)
	buf := bytes.NewReader(content)

	if flags, err = buf.ReadByte(); err != nil {
		return fmt.Errorf("failed to get liquid_level flags byte: %w", err)
	}
	flagBits := fmt.Sprintf("%08b", flags)

	e.LiquidLevelSensorErrorFlag = flagBits[1:2]
	e.LiquidLevelSensorValueUnit = flagBits[2:4]
	e.RawDataFlag = flagBits[4:5]

	if sensNum, err = strconv.ParseUint(flagBits[5:], 2, 8); err != nil {
		return fmt.Errorf("failed to read the LLS number: %w", err)
	}
	e.LiquidLevelSensorNumber = uint8(sensNum)

	bytesTmpBuf := make([]byte, 2)
	if _, err = buf.Read(bytesTmpBuf); err != nil {
		return fmt.Errorf("failed to get the address of LLS module: %w", err)
	}
	e.ModuleAddress = binary.LittleEndian.Uint16(bytesTmpBuf)

	bytesTmpBuf = make([]byte, 4)
	if _, err = buf.Read(bytesTmpBuf); err != nil {
		return fmt.Errorf("failed to get a LLS reading: %w", err)
	}
	e.LiquidLevelSensorData = binary.LittleEndian.Uint32(bytesTmpBuf)

	return nil
}

// Encode encodes the EGTS_SR_LIQUID_LEVEL_SENSOR structure into the set of bytes.
func (e *SrLiquidLevelSensor) Encode() ([]byte, error) {
	var (
		err    error
		flags  uint64
		result []byte
	)
	buf := new(bytes.Buffer)

	flagsBits := "0" + e.LiquidLevelSensorErrorFlag + e.LiquidLevelSensorValueUnit +
		e.RawDataFlag + fmt.Sprintf("%03b", e.LiquidLevelSensorNumber)
	if flags, err = strconv.ParseUint(flagsBits, 2, 8); err != nil {
		return result, fmt.Errorf("failed to generate the ext_pos_data flags byte: %w", err)
	}

	if err = buf.WriteByte(uint8(flags)); err != nil {
		return result, fmt.Errorf("failed to write the ext_pos_data flags byte: %w", err)
	}

	if err = binary.Write(buf, binary.LittleEndian, e.ModuleAddress); err != nil {
		return result, fmt.Errorf("failed to write the address of LLS module: %w", err)
	}

	if err = binary.Write(buf, binary.LittleEndian, e.LiquidLevelSensorData); err != nil {
		return result, fmt.Errorf("failed to write LLS readings: %w", err)
	}

	result = buf.Bytes()

	return result, nil
}

// Length returns the length of the EGTS_SR_LIQUID_LEVEL_SENSOR structure.
func (e *SrLiquidLevelSensor) Length() uint16 {
	var result uint16

	if recBytes, err := e.Encode(); err != nil {
		result = uint16(0)
	} else {
		result = uint16(len(recBytes))
	}

	return result
}

package egts

import (
	"bytes"
	"fmt"
	"strconv"
)

// SrStateData is the structure of subrecord of EGTS_SR_STATE_DATA type, used to transmit to the
// information about the subscriber terminal state (current operation mode,
// voltage of the main and backup power supplies, etc.).
type SrStateData struct {
	State                  uint8  `json:"ST"`
	MainPowerSourceVoltage uint8  `json:"MPSV"`
	BackUpBatteryVoltage   uint8  `json:"BBV"`
	InternalBatteryVoltage uint8  `json:"IBV"`
	NMS                    string `json:"NMS"`
	IBU                    string `json:"IBU"`
	BBU                    string `json:"BBU"`
}

// Decode parses the set of bytes into EGTS_SR_STATE_DATA structure.
func (e *SrStateData) Decode(content []byte) error {
	var (
		err   error
		flags byte
	)

	buf := bytes.NewReader(content)
	if e.State, err = buf.ReadByte(); err != nil {
		return fmt.Errorf("failed to read the current mode of operation: %w", err)
	}

	if e.MainPowerSourceVoltage, err = buf.ReadByte(); err != nil {
		return fmt.Errorf("the value of the main power supply voltage could not be read: %w", err)
	}

	if e.BackUpBatteryVoltage, err = buf.ReadByte(); err != nil {
		return fmt.Errorf("failed to get the backup battery voltage value: %w", err)
	}

	if e.InternalBatteryVoltage, err = buf.ReadByte(); err != nil {
		return fmt.Errorf("failed to get the voltage value of the internal battery: %w", err)
	}

	if flags, err = buf.ReadByte(); err != nil {
		return fmt.Errorf("failed to get state_data flags byte: %w", err)
	}
	flagBits := fmt.Sprintf("%08b", flags)
	e.NMS = flagBits[5:6]
	e.IBU = flagBits[6:7]
	e.BBU = flagBits[7:]

	return nil
}

// Encode encodes the EGTS_SR_STATE_DATA structure into the set of bytes.
func (e *SrStateData) Encode() ([]byte, error) {
	var (
		err    error
		flags  uint64
		result []byte
	)
	buf := new(bytes.Buffer)

	if err = buf.WriteByte(e.State); err != nil {
		return result, fmt.Errorf("the current operating mode could not be written: %w", err)
	}

	if err = buf.WriteByte(e.MainPowerSourceVoltage); err != nil {
		return result, fmt.Errorf("the value of the main power supply voltage could not be written: %w", err)
	}

	if err = buf.WriteByte(e.BackUpBatteryVoltage); err != nil {
		return result, fmt.Errorf("failed to write backup battery voltage value: %w", err)
	}

	if err = buf.WriteByte(e.InternalBatteryVoltage); err != nil {
		return result, fmt.Errorf("failed to record the voltage value of the internal battery: %w", err)
	}

	if flags, err = strconv.ParseUint("00000"+e.NMS+e.IBU+e.BBU, 2, 8); err != nil {
		return result, fmt.Errorf("failed to generate state_data flags byte: %w", err)
	}

	if err = buf.WriteByte(uint8(flags)); err != nil {
		return result, fmt.Errorf("failed to write state_data flags byte: %w", err)
	}

	result = buf.Bytes()
	return result, nil
}

// Length returns the length of the EGTS_SR_STATE_DATA structure.
func (e *SrStateData) Length() uint16 {
	var result uint16

	if recBytes, err := e.Encode(); err != nil {
		result = uint16(0)
	} else {
		result = uint16(len(recBytes))
	}

	return result
}

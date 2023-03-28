package egts

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strconv"
)

// SrExtPosData structure of EGTS_SR_EXT_POS_DATA type subrecord, which is used by the subscriber's
// terminal when transmitting additional location data.
type SrExtPosData struct {
	NavigationSystemFieldExists   string `json:"NSFE"`
	SatellitesFieldExists         string `json:"SFE"`
	PdopFieldExists               string `json:"PFE"`
	HdopFieldExists               string `json:"HFE"`
	VdopFieldExists               string `json:"VFE"`
	VerticalDilutionOfPrecision   uint16 `json:"VDOP"`
	HorizontalDilutionOfPrecision uint16 `json:"HDOP"`
	PositionDilutionOfPrecision   uint16 `json:"PDOP"`
	Satellites                    uint8  `json:"SAT"`
	NavigationSystem              uint16 `json:"NS"`
}

// Decode decodes EGTS_SR_EXT_POS_DATA subrecord from incoming bytes stream.
func (e *SrExtPosData) Decode(content []byte) error {
	var (
		err   error
		flags byte
	)
	tmpBuf := make([]byte, 2)
	buf := bytes.NewReader(content)

	// flags byte
	if flags, err = buf.ReadByte(); err != nil {
		return fmt.Errorf("failed to get the ext_pos_data flags byte: %w", err)
	}
	flagBits := fmt.Sprintf("%08b", flags)

	e.NavigationSystemFieldExists = flagBits[3:4]
	e.SatellitesFieldExists = flagBits[4:5]
	e.PdopFieldExists = flagBits[5:6]
	e.HdopFieldExists = flagBits[6:7]
	e.VdopFieldExists = flagBits[7:]

	if e.VdopFieldExists == "1" {
		if _, err = buf.Read(tmpBuf); err != nil {
			return fmt.Errorf("it was not possible to get vdop value: %w", err)
		}
		e.VerticalDilutionOfPrecision = binary.LittleEndian.Uint16(tmpBuf)
	}

	if e.HdopFieldExists == "1" {
		if _, err = buf.Read(tmpBuf); err != nil {
			return fmt.Errorf("it was not possible to get hdop value: %w", err)
		}
		e.HorizontalDilutionOfPrecision = binary.LittleEndian.Uint16(tmpBuf)
	}

	if e.PdopFieldExists == "1" {
		if _, err = buf.Read(tmpBuf); err != nil {
			return fmt.Errorf("failed to get pdop value: %w", err)
		}
		e.PositionDilutionOfPrecision = binary.LittleEndian.Uint16(tmpBuf)
	}

	if e.SatellitesFieldExists == "1" {
		if e.Satellites, err = buf.ReadByte(); err != nil {
			return fmt.Errorf("failed to get the number of visible satellites: %w", err)
		}
	}

	if e.NavigationSystemFieldExists == "1" {
		if _, err = buf.Read(tmpBuf); err != nil {
			return fmt.Errorf("failed to get satellite bit flags: %w", err)
		}
		e.NavigationSystem = binary.LittleEndian.Uint16(tmpBuf)
	}

	return nil
}

// Encode encodes EGTS_SR_EXT_POS_DATA subrecord to bytes stream.
func (e *SrExtPosData) Encode() ([]byte, error) {
	var (
		err    error
		flags  uint64
		result []byte
	)

	buf := new(bytes.Buffer)

	// flags byte
	flagsBits := "000" + e.NavigationSystemFieldExists + e.SatellitesFieldExists +
		e.PdopFieldExists + e.HdopFieldExists + e.VdopFieldExists
	if flags, err = strconv.ParseUint(flagsBits, 2, 8); err != nil {
		return result, fmt.Errorf("failed to generate the ext_pos_data flags byte: %w", err)
	}

	if err = buf.WriteByte(uint8(flags)); err != nil {
		return result, fmt.Errorf("failed to write the ext_pos_data flags byte: %w", err)
	}

	if e.VdopFieldExists == "1" {
		if err = binary.Write(buf, binary.LittleEndian, e.VerticalDilutionOfPrecision); err != nil {
			return result, fmt.Errorf("it was not possible to write vdop value: %w", err)
		}
	}

	if e.HdopFieldExists == "1" {
		if err = binary.Write(buf, binary.LittleEndian, e.HorizontalDilutionOfPrecision); err != nil {
			return result, fmt.Errorf("it was not possible to write hdop value: %w", err)
		}
	}

	if e.PdopFieldExists == "1" {
		if err = binary.Write(buf, binary.LittleEndian, e.PositionDilutionOfPrecision); err != nil {
			return result, fmt.Errorf("it was not possible to write pdop value: %w", err)
		}
	}

	if e.SatellitesFieldExists == "1" {
		if err = buf.WriteByte(e.Satellites); err != nil {
			return result, fmt.Errorf("failed to write the number of visible satellites: %w", err)
		}
	}

	if e.NavigationSystemFieldExists == "1" {
		if err = binary.Write(buf, binary.LittleEndian, e.NavigationSystem); err != nil {
			return result, fmt.Errorf("failed to write satellite bit flagsм: %w", err)
		}
	}

	result = buf.Bytes()
	return result, nil
}

// Length returns the length of the EGTS_SR_EXT_POS_DATA subrecord.
func (e *SrExtPosData) Length() uint16 {
	var result uint16

	if recBytes, err := e.Encode(); err != nil {
		result = uint16(0)
	} else {
		result = uint16(len(recBytes))
	}

	return result
}

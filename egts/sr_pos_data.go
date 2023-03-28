package egts

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strconv"
	"time"
)

const (
	// LOHSEast is the east longitude of the WGS84 reference point.
	LOHSEast = "0"
	// LOHSWest is the west longitude of the WGS84 reference point.
	LOHSWest = "1"
	// LAHSNorth is the north latitude of the WGS84 reference point.
	LAHSNorth = "0"
	// LAHSSouth is the south latitude of the WGS84 reference point.
	LAHSSouth = "1"
	// MVParking is the parking state of the vehicle.
	MVParking = "0"
	// MVMoving is the moving state of the vehicle.
	MVMoving = "1"
	// BBActual is the actual coordinates of the vehicle.
	BBActual = "0"
	// BBMemory is the coordinates of the vehicle from black box memory.
	BBMemory = "1"
	// FIX2D is the 2D coordinates.
	FIX2D = "0"
	// FIX3D is the 3D coordinates.
	FIX3D = "1"
	// CSWGS84 is the WGS84 coordinates.
	CSWGS84 = "0"
	// CSPZ90 is the PZ90 coordinates.
	CSPZ90 = "1"
	// VLDInvalid marks that coordinates are not valid.
	VLDInvalid = "0"
	// VLDValid marks that coordinates are valid.
	VLDValid = "1"
	// ALTSAboveSea is the altitude above sea level.
	ALTSAboveSea = 0
	// ALTSBelowSea is the altitude below sea level.
	ALTSBelowSea = 1
)

// SrPosData is structure of EGTS_SR_POS_DATA subrecord type,
// which is used by the subscriber's terminal when transmitting basic location data.
type SrPosData struct {
	// NavigationTime (NTM) - navigation time (number of seconds since 00:00:00 01.01.2010 UTC).
	NavigationTime time.Time `json:"NTM"`
	// Latitude (LAT) modulo, degrees/90 * 0xFFFFFFFFFFFF and the integer part is taken.
	Latitude float64 `json:"LAT"`
	// Longitude (LONG) modulo, degrees/180 - 0xFFFFFFFFFFFF and the integer part is taken.
	Longitude float64 `json:"LONG"`
	/* FLG - Flags section */

	// ALTE - bit flag determines the presence of the ALT field in the subrecord:
	// 1 - the ALT field is transmitted;
	// 0 - is not transmitted.
	ALTE string `json:"ALTE"`
	// LOHS - A bit flag defines hemispheric longitude:
	// 0 - eastern longitude:
	// 1 - west longitude.
	LOHS string `json:"LOHS"`
	// LAHS - the bit flag defines the hemisphere latitude:
	// 0 - north latitude;
	// 1 - south latitude.
	LAHS string `json:"LAHS"`
	// MV - bit flag, sign of movement:
	// 1 - movement;
	// 0 - vehicle is in parking mode.
	MV string `json:"MV"`
	// BB - bit flag, sign of sending data from memory ("black box"):
	// 0 - actual data;
	// 1 - data from memory ("black box").
	BB string `json:"BB"`
	// FIX - bit field, type of coordinate determination:
	// 0 - 2D fix;
	// 1 - 3D fix.
	FIX string `json:"FIX"`
	// CS - bit field, the type of system used:
	// 0 - WGS-84 coordinate system;
	// 1 - state geocentric coordinate system (ПЗ-90.02).
	CS string `json:"CS"`
	// VLD - bit flag, a sign of "validity" of coordinate data:
	// 1 - data are "valid";
	// 0 - "invalid" data.
	VLD string `json:"VLD"`
	// DirectionHighestBit - (DIRH) the highest bit (8) of the DIR parameter.
	DirectionHighestBit uint8 `json:"DIRH"`
	// AltitudeSign - (ALTS) bit flag, defines the altitude relative to sea level and makes sense
	// only when the ALTE flag is set:
	// 0 - point above sea level;
	// 1 - below sea level.
	AltitudeSign uint8 `json:"ALTS"`
	// Speed (SPD) - speed in km/h in increments of 0.1 km/h (14 low bits are used).
	Speed uint16 `json:"SPD"`
	// Direction (DIR) - direction of movement. Defined as the angle in degrees, which is counted clockwise
	// between the north direction of the geographic meridian and the direction of motion at the measurement point (
	// additionally, the most significant bit is in the DIRH field).
	Direction byte `json:"DIR"`
	// Odometer (ODM) - пройденное расстояние (пробег) в км, с дискретностью 0,1 км.
	Odometer uint32 `json:"ODM"`
	// DigitalInputs (DIN) - bit flags, define the state of the main digital inputs 1 ... 8 (if the bit is 1,
	// then the corresponding input is active, if 0,
	// then it is inactive). This field is included for convenience of use and traffic saving when working in the
	// transport monitoring systems of the basic level.
	DigitalInputs byte `json:"DIN"`
	// Source (SRC) - defines the source (event) that initiated the sending of this navigation information.
	Source byte `json:"SRC"`
	// Altitude (ALT) - altitude above sea level, m (optional parameter,
	// the presence of which is determined by the ALTE bit flag).
	Altitude uint32 `json:"ALT"`
	// SourceData (SRCD) - data characterizing the source (event) from the SRC field.
	// The presence and interpretation of the value of this field is determined by the SRC field.
	SourceData int16 `json:"SRCD"`
}

// Decode parses bytes into a subrecord structure.
func (e *SrPosData) Decode(content []byte) (err error) {
	var (
		flags byte
		speed uint64
	)
	buf := bytes.NewReader(content)

	startDate := time.Date(2010, time.January, 1, 0, 0, 0, 0, time.UTC)
	tmpUint32Buf := make([]byte, 4)
	if _, err = buf.Read(tmpUint32Buf); err != nil {
		return fmt.Errorf("failed to get the navigation time: %w", err)
	}
	preFieldVal := binary.LittleEndian.Uint32(tmpUint32Buf)
	e.NavigationTime = startDate.Add(time.Duration(preFieldVal) * time.Second)

	if _, err = buf.Read(tmpUint32Buf); err != nil {
		return fmt.Errorf("failed to get latitude: %w", err)
	}

	preFieldVal = binary.LittleEndian.Uint32(tmpUint32Buf)
	e.Latitude = float64(preFieldVal) * 90 / 0xFFFFFFFF

	if _, err = buf.Read(tmpUint32Buf); err != nil {
		return fmt.Errorf("failed to get longitude: %w", err)
	}
	preFieldVal = binary.LittleEndian.Uint32(tmpUint32Buf)
	e.Longitude = float64(preFieldVal) * 180 / 0xFFFFFFFF

	// байт флагов
	if flags, err = buf.ReadByte(); err != nil {
		return fmt.Errorf("failed to get the pos_data flags byte: %w", err)
	}
	flagBits := fmt.Sprintf("%08b", flags)
	e.ALTE = flagBits[:1]
	e.LOHS = flagBits[1:2]
	e.LAHS = flagBits[2:3]
	e.MV = flagBits[3:4]
	e.BB = flagBits[4:5]
	e.CS = flagBits[5:6]
	e.FIX = flagBits[6:7]
	e.VLD = flagBits[7:]

	// скорость
	tmpUint16Buf := make([]byte, 2)
	if _, err = buf.Read(tmpUint16Buf); err != nil {
		return fmt.Errorf("failed to get speed: %w", err)
	}
	spd := binary.LittleEndian.Uint16(tmpUint16Buf)
	e.DirectionHighestBit = uint8(spd >> 15 & 0x1)
	e.AltitudeSign = uint8(spd >> 14 & 0x1)

	speedBits := fmt.Sprintf("%016b", spd)
	if speed, err = strconv.ParseUint(speedBits[2:], 2, 16); err != nil {
		return fmt.Errorf("failed to decrypt bit rate: %w", err)
	}

	// т.к. скорость с дискретностью 0,1 км
	e.Speed = uint16(speed) / 10

	if e.Direction, err = buf.ReadByte(); err != nil {
		return fmt.Errorf("failed to get the direction of travel: %w", err)
	}
	e.Direction |= e.DirectionHighestBit << 7

	bytesTmpBuf := make([]byte, 3)
	if _, err = buf.Read(bytesTmpBuf); err != nil {
		return fmt.Errorf("failed to get the traveled distance (mileage) in km: %w", err)
	}
	bytesTmpBuf = append(bytesTmpBuf, 0x00)
	e.Odometer = binary.LittleEndian.Uint32(bytesTmpBuf)

	if e.DigitalInputs, err = buf.ReadByte(); err != nil {
		return fmt.Errorf("failed to receive bit flags, determine the state of the main discrete inputs: %w", err)
	}

	if e.Source, err = buf.ReadByte(); err != nil {
		return fmt.Errorf("failed to get the source (event) that initiated the parcel: %w", err)
	}

	if e.ALTE == "1" {
		bytesTmpBuf = []byte{0, 0, 0, 0}
		if _, err = buf.Read(bytesTmpBuf); err != nil {
			return fmt.Errorf("failed to get the altitude above sea level: %w", err)
		}
		e.Altitude = binary.LittleEndian.Uint32(bytesTmpBuf)
	}

	// something wrong with it from real data
	// err = binary.Read(buf, binary.LittleEndian, &e.Source)
	// if err != nil {
	// 	return fmt.Errorf("failed to get data characterizing the source: %w", err)
	// }

	return nil
}

// Encode converts a subwrite to a set of bytes.
func (e *SrPosData) Encode() (result []byte, err error) {
	var flags uint64

	buf := new(bytes.Buffer)

	startDate := time.Date(2010, time.January, 1, 0, 0, 0, 0, time.UTC)
	if err = binary.Write(buf, binary.LittleEndian, uint32(e.NavigationTime.Sub(startDate).Seconds())); err != nil {
		return nil, fmt.Errorf("failed to record navigation time: %w", err)
	}

	if err = binary.Write(buf, binary.LittleEndian, uint32(e.Latitude/90*0xFFFFFFFF)); err != nil {
		return nil, fmt.Errorf("failed to record latitude: %w", err)
	}

	if err = binary.Write(buf, binary.LittleEndian, uint32(e.Longitude/180*0xFFFFFFFF)); err != nil {
		return nil, fmt.Errorf("failed to record longitude: %w", err)
	}

	// байт флагов
	flags, err = strconv.ParseUint(e.ALTE+e.LOHS+e.LAHS+e.MV+e.BB+e.CS+e.FIX+e.VLD, 2, 8)
	if err != nil {
		return nil, fmt.Errorf("failed to generate the pos_data flag byte: %w", err)
	}

	if err = buf.WriteByte(uint8(flags)); err != nil {
		return nil, fmt.Errorf("failed to record flags: %w", err)
	}

	// скорость
	speed := e.Speed*10 | uint16(e.DirectionHighestBit)<<15 // 15 bit
	speed |= uint16(e.AltitudeSign) << 14                   // 14 bit
	spd := make([]byte, 2)
	binary.LittleEndian.PutUint16(spd, speed)
	if _, err = buf.Write(spd); err != nil {
		return nil, fmt.Errorf("failed to record speed: %w", err)
	}

	dir := e.Direction &^ (e.DirectionHighestBit << 7)
	if err = binary.Write(buf, binary.LittleEndian, dir); err != nil {
		return nil, fmt.Errorf("failed to record the direction of travel: %w", err)
	}

	bytesTmpBuf := make([]byte, 4)
	binary.LittleEndian.PutUint32(bytesTmpBuf, e.Odometer)
	if _, err = buf.Write(bytesTmpBuf[:3]); err != nil {
		return nil, fmt.Errorf("failed to record the traveled distance (mileage) in km: %w", err)
	}

	if err = binary.Write(buf, binary.LittleEndian, e.DigitalInputs); err != nil {
		return nil, fmt.Errorf("failed to write bit flags, determine the state of the main discrete inputs: %w", err)
	}

	if err = binary.Write(buf, binary.LittleEndian, e.Source); err != nil {
		return nil, fmt.Errorf("failed to record the source (event) that initiated the parcel: %w", err)
	}

	if e.ALTE == "1" {
		bytesTmpBuf = []byte{0, 0, 0, 0}
		binary.LittleEndian.PutUint32(bytesTmpBuf, e.Altitude)
		if _, err = buf.Write(bytesTmpBuf[:3]); err != nil {
			return nil, fmt.Errorf("failed to record altitude: %w", err)
		}
	}

	// something wrong with it from real data
	// err = binary.Write(buf, binary.LittleEndian, e.SourceData)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to record source data: %w", err)
	// }

	return buf.Bytes(), nil
}

// Length gets the length of the encoded subrecord.
func (e *SrPosData) Length() uint16 {
	var result uint16

	if recBytes, err := e.Encode(); err != nil {
		result = uint16(0)
	} else {
		result = uint16(len(recBytes))
	}

	return result
}

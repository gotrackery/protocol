package wialonretr

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

var _ Data = (*PositionInfo)(nil)

// PositionInfo represents WialonRetranslator posinfo data block.
type PositionInfo struct {
	Lon    float64
	Lat    float64
	Alt    float64
	Speed  int16
	Course int16
	Sats   int8
}

// Decode decodes WialonRetranslator posinfo data block from bytes.
func (p *PositionInfo) Decode(data []byte) error {
	var (
		le struct {
			Lon float64
			Lat float64
			Alt float64
		}
		be struct {
			Speed  int16
			Course int16
			Sats   int8
		}
	)
	reader := bytes.NewReader(data)
	err := binary.Read(reader, binary.LittleEndian, &le)
	if err != nil {
		return fmt.Errorf("read lon, lat, alt: %w", err)
	}
	p.Lon = le.Lon
	p.Lat = le.Lat
	p.Alt = le.Alt

	err = binary.Read(reader, binary.BigEndian, &be)
	if err != nil {
		return fmt.Errorf("read speed, course, sats: %w", err)
	}
	p.Speed = be.Speed
	p.Course = be.Course
	p.Sats = be.Sats

	return nil
}

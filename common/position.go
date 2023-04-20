package common

import (
	"time"

	"github.com/peterstace/simplefeatures/geom"
	"github.com/rs/zerolog"
	"gopkg.in/guregu/null.v4"
)

const (
	Proto      = "proto"
	Odometer   = "odometer"
	Satellites = "sats"
	HDOP       = "hdop"
	VDOP       = "vdop"
	PDOP       = "pdop"
	NavSystem  = "navsys"
	Move       = "move"
	DigInput   = "dinput"
	DigOutput  = "doutput"
	AnInput    = "ainput"
)

var _ zerolog.LogObjectMarshaler = (*Position)(nil)

// Location is a generic location struct.
type Location struct {
	geom.Coordinates
	Valid bool
}

// Cellular is a generic cellular struct.
type Cellular struct {
	// CellID - CID, CI - A GSM Cell ID
	CellID int64
	// LAC - Location Area Code
	LAC int64
	// MCC - Mobile Country Code
	MCC int64
	// MNC - Mobile Network Code
	MNC int64
}

// Position is a generic position struct.
type Position struct {
	Location
	Cellular   *Cellular
	Protocol   string
	DeviceID   string
	Attributes Attributes
	DeviceTime time.Time
	Speed      null.Float
	Course     null.Float
}

// MarshalZerologObject implements zerolog.LogObjectMarshaler.
func (p Position) MarshalZerologObject(e *zerolog.Event) {
	e.
		Str("proto", p.Protocol).
		Str("device", p.DeviceID).
		Time("at", p.DeviceTime).
		Float64("lon", p.X).
		Float64("lat", p.Y).
		Float64("alt", p.Z).
		Bool("valid", p.Valid)
}

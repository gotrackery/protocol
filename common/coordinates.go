package common

import (
	"fmt"
	"math"
	"strconv"

	"github.com/peterstace/simplefeatures/geom"
)

// CardinalAxis is the cardinal axis of a coordinate system.
type CardinalAxis string

const (
	North CardinalAxis = "N"
	South CardinalAxis = "S"
	East  CardinalAxis = "E"
	West  CardinalAxis = "W"
)

// ParseCardinalAxis parses a string into a CardinalAxis.
func ParseCardinalAxis(raw string) (cp CardinalAxis, err error) {
	cp = CardinalAxis(raw)
	switch cp {
	case North, South, East, West:
	default:
		cp = ""
		err = fmt.Errorf("not allowed cardinal value %s", raw) //nolint:goerr113
	}
	return cp, err
}

// Sign returns the floating point sign of the coordinate.
func (c CardinalAxis) Sign() float64 {
	switch c {
	case North, East:
		return 1.0
	case South, West:
		return -1.0
	default:
		return 0.0
	}
}

// AxisWGS84 is the WGS84 axis of a coordinate system.
type AxisWGS84 struct {
	Coordinate float64
	Cardinal   CardinalAxis
}

// Float64 returns the floating point value of the coordinate.
// Example: 5544.6025;N
// 55 is a degree value.
// 44.6025 / 60 = 0,743375 is a minute value. N is north latitude (positive sign).
// 55 + 0,743375 = +55,743375
// .
func (a AxisWGS84) Float64() float64 {
	intDeg := math.Trunc(a.Coordinate / 100)     //nolint:gomnd
	deg := intDeg + (a.Coordinate-100*intDeg)/60 //nolint:gomnd
	return math.Copysign(deg, a.Cardinal.Sign())
}

// ParseAxisWGS84 parses a string into a AxisWGS84.
func ParseAxisWGS84(coord string, c CardinalAxis) (AxisWGS84, error) {
	f, err := strconv.ParseFloat(coord, 64)
	if err != nil {
		return AxisWGS84{}, fmt.Errorf("failed to parse coordinate value %s: %w", coord, err)
	}
	if f < 0.0 || f > 18000.0 {
		return AxisWGS84{}, fmt.Errorf("not allowed coordinate value %s", coord) //nolint:goerr113
	}
	return AxisWGS84{Coordinate: f, Cardinal: c}, nil
}

// PointWGS84 is the WGS84 point of a coordinate system.
type PointWGS84 struct {
	Lon, Lat AxisWGS84
	Valid    bool
}

// ParsePointWGS84 parses a pair string lon/lat into a PointWGS84.
func ParsePointWGS84(lon string, cardLon CardinalAxis, lat string, cardLat CardinalAxis) (PointWGS84, error) {
	axisLon, err := ParseAxisWGS84(lon, cardLon)
	if err != nil {
		return PointWGS84{}, fmt.Errorf("failed to parse longitude value %s: %w", lon, err)
	}
	axisLat, err := ParseAxisWGS84(lat, cardLat)
	if err != nil {
		return PointWGS84{}, fmt.Errorf("failed to parse latitude value %s: %w", lat, err)
	}
	return PointWGS84{
		Lon:   axisLon,
		Lat:   axisLat,
		Valid: true,
	}, nil
}

// LocationXY is the XY location of a coordinate system.
func (w PointWGS84) LocationXY() Location {
	return Location{
		Coordinates: geom.Coordinates{
			XY: geom.XY{
				X: w.Lon.Float64(),
				Y: w.Lat.Float64(),
			},
			Type: geom.DimXY,
		},
		Valid: w.Valid,
	}
}

// LocationXYZ is the XYZ location of a coordinate system.
func (w PointWGS84) LocationXYZ(z float64) Location {
	l := w.LocationXY()
	l.Type = geom.DimXYZ
	l.Z = z
	return l
}

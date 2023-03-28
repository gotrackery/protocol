package wialonips

import (
	"bytes"
	"fmt"
	"time"

	"github.com/gotrackery/protocol/generic"
	"gopkg.in/guregu/null.v4"
)

var _ Message = (*ShortenedDataMessage)(nil)

// ShortenedDataMessage is a WialonIPS shortened data message.
// The packet contains only navigation data and looks as follows:
// #SD#Date;Time;Lat1;Lat2;Lon1;Lon2;Speed;Course;Alt;Sats;CRC16\r\n.
type ShortenedDataMessage struct {
	RegisteredAt time.Time
	Point        generic.PointWGS84
	// Speed is type of int according specs, but there is found it comes in float xxx.yyy thru retranslation.
	Speed  null.Float
	Course null.Int
	// Altitude is type of int according specs, but there is found it comes in float xxx.yyy thru retranslation.
	Altitude null.Float
	Sat      null.Int
	message
}

// Decode decodes a WialonIPS message.
func (s *ShortenedDataMessage) Decode(data []byte) error {
	bytesSet := bytes.Split(data, fieldsDelimiter)
	length := len(bytesSet)
	if (length != 10 && s.ver == V1_1) || (length != 11 && s.ver == V2_0) {
		s.err = ErrWialonIPSInvalidDataMessage // -1
		return s.err
	}

	if s.ver == V2_0 {
		s.err = validateCRC(data, fieldsDelimiter)
		if s.err != nil {
			return s.err
		}
	}

	s.RegisteredAt, s.Point, s.Speed, s.Course, s.Altitude, s.Sat, s.err = parseBaseFields(bytesSet[0:10])
	return s.err
}

// Response returns a WialonIPS response message.
func (s *ShortenedDataMessage) Response() []byte {
	return []byte(fmt.Sprintf(responseTemplate, ShortenedDataPacket, MapErrToRespCode(ShortenedDataPacket, s.err)))
}

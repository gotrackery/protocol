package wialonips

import (
	"bytes"
	"fmt"

	"github.com/gotrackery/protocol/common"
	"gopkg.in/guregu/null.v4"
)

var _ Message = (*DataMessage)(nil)

// DataMessage is a WialonIPS data message.
// The packet contains additional data structures and looks as follows:
// #D#Date;Time;Lat1;Lat2;Lon1;Lon2;Speed;Course;Alt;Sats;HDOP;Inputs; Outputs;ADC;Ibutton;Params;CRC16\r\n
// Each parameter has the following structure:
// Name:Type:Value
// Examples of additional parameters: count1:1:564, fuel:2:45.8, hw:3:V4.5, SOS:1:1.
type DataMessage struct {
	ShortenedDataMessage
	HDOP       null.Float
	Inputs     null.Int // bit map.
	Outputs    null.Int // bit map.
	ADC        []null.Float
	IButton    null.String
	Attributes common.Attributes
}

// Decode decodes a WialonIPS message.
func (d *DataMessage) Decode(data []byte) error {
	bytesSet := bytes.Split(data, fieldsDelimiter)
	length := len(bytesSet)
	if (length != 16 && d.ver == V1_1) || (length != 17 && d.ver == V2_0) {
		d.err = ErrWialonIPSInvalidDataMessage // -1
		return d.err
	}

	if d.ver == V2_0 {
		d.err = validateCRC(data, fieldsDelimiter)
		if d.err != nil {
			return d.err
		}
	}

	d.RegisteredAt, d.Point, d.Speed, d.Course, d.Altitude, d.Sat, d.err = parseBaseFields(bytesSet[0:10])
	if d.err != nil {
		return d.err
	}

	d.HDOP, d.Inputs, d.Outputs, d.ADC, d.IButton, d.err = parseAdditionalFields(bytesSet[10:15])
	if d.err != nil {
		return d.err
	}

	d.Attributes, d.err = parseAttrs(bytesSet[15])
	return d.err
}

// Response returns a WialonIPS response message.
func (d *DataMessage) Response() []byte {
	return []byte(fmt.Sprintf(responseTemplate, DataPacket, MapErrToRespCode(DataPacket, d.err)))
}

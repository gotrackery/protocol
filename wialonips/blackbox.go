package wialonips

import (
	"bytes"
	"fmt"
	"strconv"
)

var _ Message = (*BlackBoxMessage)(nil)

// BlackBoxMessage is a WialonIPS black box message.
type BlackBoxMessage struct {
	ShortenedMessages []ShortenedDataMessage
	DataMessages      []DataMessage
	message
}

// Decode decodes a WialonIPS message.
// The black box packet is used to transmit messages for the past period.
// The maximum number of messages that can be transmitted in one packet is 5000. The packet looks as follows:
// #B#Date;Time;Lat1;Lat2;Lon1;Lon2;Speed;Course;Alt;Sats|Date;Time;Lat1;Lat 2;Lon1;Lon2;Speed;Course;Alt;Sats|Date
// ;Time;Lat1;Lat2;Lon1;Lon2;Speed; Course;Alt;Sats|CRC16\r\n.
func (bb *BlackBoxMessage) Decode(data []byte) error {
	const shortDataLen = 10
	bytesSet := bytes.Split(data, blackBoxDelimiter)
	if bb.ver == V2_0 {
		bb.err = validateCRC(data, blackBoxDelimiter)
		if bb.err != nil {
			return bb.err
		}
		bytesSet = bytesSet[0 : len(bytesSet)-1]
	}

	probe := bytes.Count(bytesSet[0], fieldsDelimiter)
	isShortened := probe == shortDataLen
	if isShortened {
		bb.ShortenedMessages = make([]ShortenedDataMessage, 0, len(bytesSet))
		for _, d := range bytesSet {
			msg := ShortenedDataMessage{}
			bb.err = msg.Decode(d)
			if bb.err != nil {
				return bb.err
			}
			bb.ShortenedMessages = append(bb.ShortenedMessages, msg)
		}
		return nil
	}

	bb.DataMessages = make([]DataMessage, 0, len(bytesSet))
	for _, d := range bytesSet {
		msg := DataMessage{}
		bb.err = msg.Decode(d)
		if bb.err != nil {
			return bb.err
		}
		bb.DataMessages = append(bb.DataMessages, msg)
	}
	return nil
}

// Response returns a WialonIPS response message.
func (bb *BlackBoxMessage) Response() []byte {
	return []byte(fmt.Sprintf(responseTemplate, BlackBoxPacket,
		strconv.Itoa(len(bb.DataMessages)+len(bb.ShortenedMessages))))
}

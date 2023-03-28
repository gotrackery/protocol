package egts

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strconv"
	"time"
)

// Convert the navigation time to the format required by the standard: number of seconds since 00:00:00 01.01.2010 UTC.
var timeOffset = time.Date(2010, time.January, 1, 0, 0, 0, 0, time.UTC)

// ServiceDataRecord record containing monitoring information.
type ServiceDataRecord struct {
	RecordLength             uint16    `json:"RL"`
	RecordNumber             uint16    `json:"RN"`
	SourceServiceOnDevice    string    `json:"SSOD"`
	RecipientServiceOnDevice string    `json:"RSOD"`
	Group                    string    `json:"GRP"`
	RecordProcessingPriority string    `json:"RPP"`
	TimeFieldExists          string    `json:"TMFE"`
	EventIDFieldExists       string    `json:"EVFE"`
	ObjectIDFieldExists      string    `json:"OBFE"`
	ObjectIdentifier         uint32    `json:"OID"`
	EventIdentifier          uint32    `json:"EVID"`
	Time                     time.Time `json:"TM"`
	SourceServiceType        byte      `json:"SST"`
	RecipientServiceType     byte      `json:"RST"`
	RecordDataSet            `json:"RD"`
}

// ServiceDataSet set of consecutive records with information.
type ServiceDataSet []ServiceDataRecord

// Decode decodes the given byte slice into a ServiceDataRecord.
func (s *ServiceDataSet) Decode(serviceDS []byte) error {
	var (
		err   error
		flags byte
	)
	buf := bytes.NewReader(serviceDS)

	for buf.Len() > 0 {
		sdr := ServiceDataRecord{}
		tmpIntBuf := make([]byte, 2)
		if _, err = buf.Read(tmpIntBuf); err != nil {
			return fmt.Errorf("failed to get the SDR record length: %w", err)
		}
		sdr.RecordLength = binary.LittleEndian.Uint16(tmpIntBuf)

		if _, err = buf.Read(tmpIntBuf); err != nil {
			return fmt.Errorf("failed to get SDR record number: %w", err)
		}
		sdr.RecordNumber = binary.LittleEndian.Uint16(tmpIntBuf)

		if flags, err = buf.ReadByte(); err != nil {
			return fmt.Errorf("failed to read the SDR flags byte: %w", err)
		}
		flagBits := fmt.Sprintf("%08b", flags)
		sdr.SourceServiceOnDevice = flagBits[:1]
		sdr.RecipientServiceOnDevice = flagBits[1:2]
		sdr.Group = flagBits[2:3]
		sdr.RecordProcessingPriority = flagBits[3:5]
		sdr.TimeFieldExists = flagBits[5:6]
		sdr.EventIDFieldExists = flagBits[6:7]
		sdr.ObjectIDFieldExists = flagBits[7:]

		if sdr.ObjectIDFieldExists == "1" {
			oid := make([]byte, 4)
			if _, err := buf.Read(oid); err != nil {
				return fmt.Errorf("failed to get SDR object identifier: %w", err)
			}
			sdr.ObjectIdentifier = binary.LittleEndian.Uint32(oid)
		}

		if sdr.EventIDFieldExists == "1" {
			event := make([]byte, 4)
			if _, err := buf.Read(event); err != nil {
				return fmt.Errorf("failed to get SDR event identifier: %w", err)
			}
			sdr.EventIdentifier = binary.LittleEndian.Uint32(event)
		}

		// Convert the navigation time to the format required by the standard:
		// number of seconds since 00:00:00 01.01.2010 UTC.
		if sdr.TimeFieldExists == "1" {
			tm := make([]byte, 4)
			if _, err := buf.Read(tm); err != nil {
				return fmt.Errorf("failed to get record generation time on the sender side of the SDR: %w", err)
			}
			preFieldVal := binary.LittleEndian.Uint32(tm)
			sdr.Time = timeOffset.Add(time.Duration(preFieldVal) * time.Second)
		}

		if sdr.SourceServiceType, err = buf.ReadByte(); err != nil {
			return fmt.Errorf("failed to read the identifier of the SDR sending service type: %w", err)
		}

		if sdr.RecipientServiceType, err = buf.ReadByte(); err != nil {
			return fmt.Errorf("failed to read the identifier of the SDR recipient service type: %w", err)
		}

		if buf.Len() != 0 {
			rds := RecordDataSet{}
			rdsBytes := make([]byte, sdr.RecordLength)
			if _, err = buf.Read(rdsBytes); err != nil {
				return fmt.Errorf("failed to read the SDR record data: %w", err)
			}

			if err = rds.Decode(rdsBytes); err != nil {
				return fmt.Errorf("failed to decode the SDR record data: %w", err)
			}
			sdr.RecordDataSet = rds
		}

		*s = append(*s, sdr)
	}
	return nil
}

// Encode encodes the given ServiceDataRecord into a byte slice.
func (s *ServiceDataSet) Encode() ([]byte, error) {
	var (
		result []byte
		flags  uint64
	)

	buf := new(bytes.Buffer)

	for _, sdr := range *s {
		rd, err := sdr.RecordDataSet.Encode()
		if err != nil {
			return result, fmt.Errorf("failed to encode the SDR record data: %w", err)
		}

		if sdr.RecordLength == 0 {
			sdr.RecordLength = uint16(len(rd))
		}
		if err = binary.Write(buf, binary.LittleEndian, sdr.RecordLength); err != nil {
			return result, fmt.Errorf("failed to write SDR record length: %w", err)
		}

		if err = binary.Write(buf, binary.LittleEndian, sdr.RecordNumber); err != nil {
			return result, fmt.Errorf("failed to write SDR record number: %w", err)
		}

		// составной байт
		flagsBits := sdr.SourceServiceOnDevice + sdr.RecipientServiceOnDevice + sdr.Group + sdr.RecordProcessingPriority +
			sdr.TimeFieldExists + sdr.EventIDFieldExists + sdr.ObjectIDFieldExists
		if flags, err = strconv.ParseUint(flagsBits, 2, 8); err != nil {
			return result, fmt.Errorf("failed to generate SDR flags byte: %w", err)
		}
		if err = buf.WriteByte(uint8(flags)); err != nil {
			return result, fmt.Errorf("failed to write SDR flags: %w", err)
		}

		if sdr.ObjectIDFieldExists == "1" {
			if err = binary.Write(buf, binary.LittleEndian, sdr.ObjectIdentifier); err != nil {
				return result, fmt.Errorf("failed to write SDR object identifier: %w", err)
			}
		}

		if sdr.EventIDFieldExists == "1" {
			if err = binary.Write(buf, binary.LittleEndian, sdr.EventIdentifier); err != nil {
				return result, fmt.Errorf("failed to write SDR event identifier: %w", err)
			}
		}

		if sdr.TimeFieldExists == "1" {
			tm := uint32(sdr.Time.Unix() - timeOffset.Unix())
			if err = binary.Write(buf, binary.LittleEndian, tm); err != nil {
				return result, fmt.Errorf(
					"failed to write the time of recording on the sender side of the SDR: %w", err)
			}
		}

		if err = buf.WriteByte(sdr.SourceServiceType); err != nil {
			return result, fmt.Errorf("failed to write the identifier of the SDR sending service type: %w", err)
		}

		if err = buf.WriteByte(sdr.RecipientServiceType); err != nil {
			return result, fmt.Errorf("failed to write the identifier of the SDR recipient service type: %w", err)
		}

		buf.Write(rd)
	}

	result = buf.Bytes()

	return result, nil
}

// Length returns the length of the encoded byte slice.
func (s *ServiceDataSet) Length() uint16 {
	var result uint16

	if recBytes, err := s.Encode(); err != nil {
		result = uint16(0)
	} else {
		result = uint16(len(recBytes))
	}

	return result
}

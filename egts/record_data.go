package egts

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

// RecordData structure of the sub-section of the ServiceDataRecord record.
type RecordData struct {
	SubrecordType   byte       `json:"SRT"`
	SubrecordLength uint16     `json:"SRL"`
	SubrecordData   BinaryData `json:"SRD"`
}

// RecordDataSet describes an array with subrecords of the EGTS protocol.
type RecordDataSet []RecordData

// Decode parses the set of bytes into RecordDataSet structure.
func (rds *RecordDataSet) Decode(recDS []byte) error {
	var (
		err error
	)
	buf := bytes.NewBuffer(recDS)
	for buf.Len() > 0 {
		rd := RecordData{}
		if rd.SubrecordType, err = buf.ReadByte(); err != nil {
			return fmt.Errorf("failed to get subrecord data record type: %w", err)
		}

		tmpIntBuf := make([]byte, 2)
		if _, err = buf.Read(tmpIntBuf); err != nil {
			return fmt.Errorf("failed to get subrecord data record length: %w", err)
		}
		rd.SubrecordLength = binary.LittleEndian.Uint16(tmpIntBuf)

		subRecordBytes := buf.Next(int(rd.SubrecordLength))

		switch rd.SubrecordType {
		case SrPosDataType:
			rd.SubrecordData = &SrPosData{}
		case SrTermIdentityType:
			rd.SubrecordData = &SrTermIdentity{}
		case SrModuleDataType:
			rd.SubrecordData = &SrModuleData{}
		case SrRecordResponseType:
			rd.SubrecordData = &SrResponse{}
		case SrResultCodeType:
			rd.SubrecordData = &SrResultCode{}
		case SrExtPosDataType:
			rd.SubrecordData = &SrExtPosData{}
		case SrAdSensorsDataType:
			rd.SubrecordData = &SrAdSensorsData{}
		case SrType20:
			// the indication is indirect in the specifications it is not
			if rd.SubrecordLength == uint16(5) {
				rd.SubrecordData = &SrStateData{}
			} else {
				// TODO: add EGTS_SR_ACCEL_DATA
				return fmt.Errorf("unimplemented section EGTS_SR_ACCEL_DATA: %d. Length: %d. Contents: %X",
					rd.SubrecordType, rd.SubrecordLength, subRecordBytes)
			}
		case SrStateDataType:
			rd.SubrecordData = &SrStateData{}
		case SrLiquidLevelSensorType:
			rd.SubrecordData = &SrLiquidLevelSensor{}
		case SrAbsCntrDataType:
			rd.SubrecordData = &SrAbsCntrData{}
		case SrAuthInfoType:
			rd.SubrecordData = &SrAuthInfo{}
		case SrCountersDataType:
			rd.SubrecordData = &SrCountersData{}
		case SrEgtsPlusDataType:
			rd.SubrecordData = &StorageRecord{}
		case SrAbsAnSensDataType:
			rd.SubrecordData = &SrAbsAnSensData{}
		case SrDispatcherIdentityType:
			rd.SubrecordData = &SrDispatcherIdentity{}
		default:
			return fmt.Errorf("unknown subrecord type: %d. Length: %d. Contents: %X", rd.SubrecordType,
				rd.SubrecordLength, subRecordBytes)
		}

		if err = rd.SubrecordData.Decode(subRecordBytes); err != nil {
			return fmt.Errorf("failed to decode subrecord data: %w", err)
		}
		*rds = append(*rds, rd)
	}

	return nil
}

// Encode returns the set of bytes of the RecordDataSet structure.
func (rds *RecordDataSet) Encode() ([]byte, error) {
	var (
		result []byte
		err    error
	)
	buf := new(bytes.Buffer)

	for _, rd := range *rds {
		if rd.SubrecordType == 0 {
			switch rd.SubrecordData.(type) {
			case *SrPosData:
				rd.SubrecordType = SrPosDataType
			case *SrTermIdentity:
				rd.SubrecordType = SrTermIdentityType
			case *SrResponse:
				rd.SubrecordType = SrRecordResponseType
			case *SrResultCode:
				rd.SubrecordType = SrResultCodeType
			case *SrExtPosData:
				rd.SubrecordType = SrExtPosDataType
			case *SrAdSensorsData:
				rd.SubrecordType = SrAdSensorsDataType
			case *SrStateData:
				rd.SubrecordType = SrStateDataType
			case *SrLiquidLevelSensor:
				rd.SubrecordType = SrLiquidLevelSensorType
			case *SrAbsCntrData:
				rd.SubrecordType = SrAbsCntrDataType
			case *SrAuthInfo:
				rd.SubrecordType = SrAuthInfoType
			case *SrCountersData:
				rd.SubrecordType = SrCountersDataType
			case *StorageRecord:
				rd.SubrecordType = SrEgtsPlusDataType
			case *SrAbsAnSensData:
				rd.SubrecordType = SrAbsAnSensDataType
			default:
				return result, fmt.Errorf("there is no known code for this type of subrecord: %T", rd.SubrecordData)
			}
		}

		if err = binary.Write(buf, binary.LittleEndian, rd.SubrecordType); err != nil {
			return result, fmt.Errorf("failed to write subrecord type: %w", err)
		}

		if rd.SubrecordLength == 0 {
			rd.SubrecordLength = rd.SubrecordData.Length()
		}
		if err = binary.Write(buf, binary.LittleEndian, rd.SubrecordLength); err != nil {
			return result, fmt.Errorf("failed to write subrecord length: %w", err)
		}

		var srd []byte
		srd, err = rd.SubrecordData.Encode()
		if err != nil {
			return result, fmt.Errorf("failed to encode subrecord data: %w", err)
		}
		_, err = buf.Write(srd)
		if err != nil {
			return result, fmt.Errorf("failed to write subrecord data: %w", err)
		}
	}

	result = buf.Bytes()
	return result, nil
}

// Length returns the length of the RecordDataSet structure.
func (rds *RecordDataSet) Length() uint16 {
	var result uint16

	if recBytes, err := rds.Encode(); err != nil {
		result = uint16(0)
	} else {
		result = uint16(len(recBytes))
	}

	return result
}

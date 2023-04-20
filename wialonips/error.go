package wialonips

import "errors"

var (
	ErrWialonIPSUnsupportedPacketType = errors.New("unsupported packet type")
	ErrWialonIPSInvalidLoginMessage   = errors.New("invalid login #L# message")
	ErrWialonIPSInvalidSDMessage      = errors.New("invalid shortened data #SD# message")
	ErrWialonIPSInvalidDataMessage    = errors.New("invalid data #D# message")
	ErrWialonIPSParseDateTime         = errors.New("invalid date or time data")
	ErrWialonIPSParsePoint            = errors.New("invalid coordinates data")
	ErrWialonIPSParseSCA              = errors.New("invalid speed, course or altitude data")
	ErrWialonIPSParseSats             = errors.New("invalid sats data")
	ErrWialonIPSParseHDOP             = errors.New("invalid hdop data")
	ErrWialonIPSParseInOutput         = errors.New("invalid input or output data")
	ErrWialonIPSParseADC              = errors.New("invalid adc data")
	ErrWialonIPSParseAttribute        = errors.New("invalid parameter data")
	ErrWialonIPSCRC16Validation       = errors.New("CRC16 validation not passed")
)

// MapErrToRespCode maps errors to WialonIPS respond codes.
//
//nolint:funlen,gocognit,gocyclo,gocognit,cyclop
func MapErrToRespCode(p PacketType, err error) string {
	switch {
	case p == LoginPacket && err == nil:
		return "1"
	case p == LoginPacket && errors.Is(err, ErrWialonIPSInvalidLoginMessage):
		return "0"
	// case p == LoginPacket && errors.Is(err, ToDo password verification):
	// 	return "01"
	case p == LoginPacket && errors.Is(err, ErrWialonIPSCRC16Validation):
		return "10"

	case p == ShortenedDataPacket && err == nil:
		return "1"
	case p == ShortenedDataPacket && errors.Is(err, ErrWialonIPSInvalidSDMessage):
		return "-1"
	case p == ShortenedDataPacket && errors.Is(err, ErrWialonIPSParseDateTime):
		return "0"
	case p == ShortenedDataPacket && errors.Is(err, ErrWialonIPSParsePoint):
		return "10"
	case p == ShortenedDataPacket && errors.Is(err, ErrWialonIPSParseSCA):
		return "11"
	case p == ShortenedDataPacket && errors.Is(err, ErrWialonIPSParseSats):
		return "12"
	case p == ShortenedDataPacket && errors.Is(err, ErrWialonIPSCRC16Validation):
		return "13"

	case p == DataPacket && err == nil:
		return "1"
	case p == DataPacket && errors.Is(err, ErrWialonIPSInvalidDataMessage):
		return "-1"
	case p == DataPacket && errors.Is(err, ErrWialonIPSParseDateTime):
		return "0"
	case p == DataPacket && errors.Is(err, ErrWialonIPSParsePoint):
		return "10"
	case p == DataPacket && errors.Is(err, ErrWialonIPSParseSCA):
		return "11"
	case p == DataPacket && (errors.Is(err, ErrWialonIPSParseSats) || errors.Is(err, ErrWialonIPSParseHDOP)):
		return "12"
	case p == DataPacket && errors.Is(err, ErrWialonIPSParseInOutput):
		return "13"
	case p == DataPacket && errors.Is(err, ErrWialonIPSParseADC):
		return "14"
	case p == DataPacket && errors.Is(err, ErrWialonIPSParseAttribute):
		return "15"
	case p == DataPacket && errors.Is(err, ErrWialonIPSCRC16Validation):
		return "16"
	}

	return "-1"
}

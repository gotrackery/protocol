package wialonips

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gotrackery/protocol/common"
	"gopkg.in/guregu/null.v4"
)

//nolint:nakedret
func parseBaseFields(data [][]byte) (
	t time.Time, p common.PointWGS84, speed null.Float, course null.Int, alt null.Float, sat null.Int, err error) {
	t, err = parseTime(data[0:2])
	if err != nil {
		err = fmt.Errorf("parse time: %w", errors.Join(err, ErrWialonIPSParseDateTime)) // "0"
		return
	}

	p, err = parsePoint(data[2:6])
	if err != nil {
		err = fmt.Errorf("parse point: %w", errors.Join(err, ErrWialonIPSParsePoint)) // "10"
		return
	}

	speed, err = parseFloat(data[6])
	if err != nil {
		err = fmt.Errorf("parse speed: %w", errors.Join(err, ErrWialonIPSParseSCA)) // "11"
		return
	}

	course, err = parseInt(data[7])
	if err != nil {
		// ToDo if i>359 return error
		err = fmt.Errorf("parse course: %w", errors.Join(err, ErrWialonIPSParseSCA)) // "11"
		return
	}

	alt, err = parseFloat(data[8])
	if err != nil {
		err = fmt.Errorf("parse altitude: %w", errors.Join(err, ErrWialonIPSParseSCA)) // "11"
		return
	}

	sat, err = parseInt(data[9])
	if err != nil {
		err = fmt.Errorf("parse sats: %w", errors.Join(err, ErrWialonIPSParseSats)) // "12"
		return
	}

	return
}

func parseAdditionalFields(data [][]byte) (
	hdop null.Float, inputs, outputs null.Int, adc []null.Float, ibutton null.String, err error) {
	hdop, err = parseFloat(data[0])
	if err != nil {
		err = fmt.Errorf("parse hdop: %w", errors.Join(err, ErrWialonIPSParseHDOP)) // 12
		return
	}

	inputs, err = parseInt(data[1])
	if err != nil {
		err = fmt.Errorf("parse input: %w", errors.Join(err, ErrWialonIPSParseInOutput)) // 13
		return
	}

	outputs, err = parseInt(data[2])
	if err != nil {
		err = fmt.Errorf("parse output: %w", errors.Join(err, ErrWialonIPSParseInOutput)) // 13
		return
	}

	adc, err = parseADC(data[3])
	if err != nil {
		err = fmt.Errorf("parse adc: %w", errors.Join(err, ErrWialonIPSParseADC)) // 14
		return
	}

	ibutton = parseString(data[4])
	return
}

func parseADC(data []byte) (adc []null.Float, err error) {
	if len(data) == 0 {
		return
	}
	bytesSet := bytes.Split(data, analogDelimiter)
	adc = make([]null.Float, 0, len(bytesSet))
	for _, v := range bytesSet {
		f, err := parseFloat(v)
		if err != nil {
			return nil, err
		}
		adc = append(adc, f)
	}
	return
}

func parseAttrs(data []byte) (a common.Attributes, err error) {
	if len(data) == 0 || (len(data) == 2 && string(data) == na) {
		return
	}

	bytesSet := bytes.Split(data, valuesDelimiter)
	a = make(common.Attributes)
	for _, attr := range bytesSet {
		k, v, err := parseAttr(attr)
		if err != nil {
			return nil, fmt.Errorf("parse attribute: %w", errors.Join(err, ErrWialonIPSParseAttribute))
		}
		a[k] = v
	}
	return
}

func parseAttr(attr []byte) (k string, v interface{}, err error) {
	var bytesSet = bytes.Split(attr, paramsDelimiter)
	if len(bytesSet) != 3 { //nolint:gomnd
		return "", nil, errors.New("invalid attribute format") //nolint:goerr113
	}
	value, err := parseValue(bytesSet[1][0], bytesSet[2])
	if err != nil {
		return "", nil, err
	}

	return string(bytesSet[0]), value, nil
}

func parseValue(b byte, data []byte) (val interface{}, err error) {
	switch b {
	case intType:
		return strconv.ParseInt(string(data), 10, 64) //nolint:wrapcheck
	case floatType:
		return strconv.ParseFloat(string(data), 64) //nolint:wrapcheck
	case stringType:
		return string(data), nil
	}
	return nil, fmt.Errorf("invalid parameter type %v", b) //nolint:goerr113
}

func parseTime(data [][]byte) (t time.Time, err error) {
	t, err = time.Parse(ddmmyyhhmiss, string(bytes.Join(data, []byte{})))
	if err != nil && string(data[0]) == na {
		// according spec
		return time.Now(), nil
	}
	return
}

func parsePoint(data [][]byte) (common.PointWGS84, error) {
	latC, err := common.ParseCardinalAxis(string(data[1]))
	if err != nil {
		if string(data[1]) == na {
			latC = common.North
		} else {
			return common.PointWGS84{}, fmt.Errorf("parse latitude: %w", err)
		}
	}
	lonC, err := common.ParseCardinalAxis(string(data[3]))
	if err != nil {
		if string(data[1]) == na {
			lonC = common.East
		} else {
			return common.PointWGS84{}, fmt.Errorf("parse logitude: %w", err)
		}
	}
	wgs84, err := common.ParsePointWGS84(string(data[2]), lonC, string(data[0]), latC)
	if err != nil {
		if string(data[0]) == na && string(data[2]) == na {
			wgs84, _ = common.ParsePointWGS84("0.0", lonC, "0.0", latC)
			wgs84.Valid = false
		} else {
			return common.PointWGS84{}, fmt.Errorf("parse point: %w", err)
		}
	}
	return wgs84, nil
}

func parseInt(data []byte) (null.Int, error) {
	s := string(data)
	if s == na || s == "" {
		return null.NewInt(0, false), nil
	}
	i, err := strconv.ParseInt(s, 10, 0)
	if err != nil {
		return null.NewInt(0, false), err
	}
	return null.NewInt(i, true), nil
}

func parseFloat(data []byte) (null.Float, error) {
	s := string(data)
	if s == na || s == "" {
		return null.NewFloat(0.0, false), nil
	}
	s = strings.Replace(s, ",", ".", 1)
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return null.NewFloat(0.0, false), err
	}
	return null.NewFloat(f, true), nil
}

func parseString(data []byte) null.String {
	s := string(data)
	if s == na {
		return null.NewString("", false)
	}
	return null.NewString(s, true)
}

func parseCRC(data []byte) (crc uint16, err error) {
	s := string(data)
	if len(s)%2 != 0 {
		s = fmt.Sprintf("0%s", s)
	}
	data, err = hex.DecodeString(s)
	if err != nil {
		return 0, fmt.Errorf("failed parse crc: %w", err)
	}
	if len(data) != 2 { //nolint:gomnd
		return 0, fmt.Errorf("invalid crc length %v", len(data)) //nolint:goerr113
	}
	return binary.BigEndian.Uint16(data), nil
}

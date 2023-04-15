package wialonretr

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math"
)

type dataType byte

const (
	dataTypeText dataType = iota + 1
	dataTypeBinary
	dataTypeInt32
	dataTypeFloat64
	dataTypeInt64
)

const (
	PosInfoName    = "posinfo"
	ImageName      = "imag"
	AvlInputsName  = "avl_inputs"
	AvlOutputsName = "avl_outputs"
	AvlDriverName  = "avl_driver"
)

const (
	shownParam byte = iota
	hiddenParam
)

// Data provides contract for decoding data blocks of WialonRetranslator protocol.
type Data interface {
	Decode([]byte) error
}

var _ Data = (*DataBlock)(nil)

// DataBlock represents WialonRetranslator data block.
type DataBlock struct {
	securityParam byte
	name          string
	Value         interface{}
}

type DataBlocks map[string]DataBlock

func (db *DataBlocks) AddPosInfo(pi PositionInfo) {
	db.initIfNil()
	(*db)[PosInfoName] = DataBlock{securityParam: hiddenParam, name: PosInfoName, Value: pi}
}

func (db *DataBlocks) AddImage(i Image) {
	db.initIfNil()
	(*db)[ImageName] = DataBlock{securityParam: hiddenParam, name: ImageName, Value: i}
}

func (db *DataBlocks) AddDataBlock(name string, b DataBlock) {
	db.initIfNil()
	(*db)[b.name] = DataBlock{securityParam: determineSecParam(name), name: name, Value: b.Value}
}

func (db *DataBlocks) initIfNil() {
	if *db == nil {
		*db = make(map[string]DataBlock)
	}
}

func determineSecParam(name string) byte {
	switch name {
	case PosInfoName, AvlInputsName, AvlOutputsName, AvlDriverName:
		return hiddenParam
	case ImageName:
		return shownParam
	}
	return shownParam
}

// Decode decodes WialonRetranslator data block from bytes.
func (d *DataBlock) Decode(data []byte) error { //nolint:cyclop
	const dataBlockHeaderLen = 6
	if len(data) < dataBlockHeaderLen {
		return fmt.Errorf("invalid len %d of data block", len(data)) //nolint:goerr113
	}

	d.securityParam = data[4]

	dt := dataType(data[5])

	b, a, f := bytes.Cut(data[6:], eol)
	if !f {
		return ErrWialonRetranslatorCutBlockName
	}

	d.name = string(b)

	switch dt {
	case dataTypeText:
		d.Value = string(a)
	case dataTypeBinary:
		switch d.name {
		case PosInfoName:
			var pi PositionInfo
			if err := pi.Decode(a); err != nil {
				return fmt.Errorf("decode posinfo: %w", err)
			}
			d.Value = pi
			return nil
		case ImageName:
			var i Image
			if err := i.Decode(a); err != nil {
				return fmt.Errorf("decode image: %w", err)
			}
			d.Value = i
			return nil
		}
		d.Value = a
	case dataTypeInt32:
		d.Value = int32(binary.BigEndian.Uint32(a))
	case dataTypeFloat64:
		d.Value = math.Float64frombits(binary.LittleEndian.Uint64(a))
	case dataTypeInt64:
		d.Value = int64(binary.BigEndian.Uint64(a))
	default:
		return fmt.Errorf("unknown data type %d", dt) //nolint:goerr113
	}

	return nil
}

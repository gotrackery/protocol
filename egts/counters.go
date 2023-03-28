package egts

import (
	"math"
	"sync/atomic"
)

var (
	cntPacketIdentifier uint32
	cntRecordNumber     uint32
)

func nextPacketIdentifier() uint16 {
	if cntPacketIdentifier < math.MaxUint16 {
		atomic.AddUint32(&cntPacketIdentifier, 1)
	} else {
		cntPacketIdentifier = 0
	}
	return uint16(atomic.LoadUint32(&cntPacketIdentifier))
}

func nextRecordNumber() uint16 {
	if cntRecordNumber < math.MaxUint16 {
		atomic.AddUint32(&cntRecordNumber, 1)
	} else {
		cntRecordNumber = 0
	}
	return uint16(atomic.LoadUint32(&cntRecordNumber))
}

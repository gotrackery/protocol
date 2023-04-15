package wialonretr

import (
	"errors"
	"fmt"

	"github.com/gotrackery/protocol"
)

var (
	ErrWialonRetranslatorBadDeviceID  = fmt.Errorf("failed to get device id: %w", protocol.ErrInconsistentData)
	ErrWialonRetranslatorCutBlockName = errors.New("cut block data name")
)

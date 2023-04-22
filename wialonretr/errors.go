package wialonretr

import (
	"errors"
	"fmt"

	"github.com/gotrackery/protocol/common"
)

var (
	ErrWialonRetranslatorBadDeviceID  = fmt.Errorf("failed to get device id: %w", common.ErrBadData)
	ErrWialonRetranslatorCutBlockName = errors.New("cut block data name")
)

package wialonretr

import (
	"fmt"

	"github.com/gotrackery/protocol"
)

var ErrWialonRetranslatorBadDeviceID = fmt.Errorf("failed to get device id: %w", protocol.ErrInconsistentData)

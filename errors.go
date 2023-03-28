package protocol

import (
	"errors"
)

var (
	// ErrInconsistentData is returned when the data is inconsistent with protocol specification.
	// It signals that communication session must be terminated.
	ErrInconsistentData = errors.New("data inconsistent to protocol specification")
)

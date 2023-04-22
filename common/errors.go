package common

import (
	"errors"
)

var (
	// ErrBadData is returned when the data is inconsistent with protocol specification.
	// It signals that communication session can be terminated.
	ErrBadData = errors.New("data is bad to protocol specification")
)

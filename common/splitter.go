package common

import (
	"bufio"
)

// FrameSplitter is a contract to split incoming bytes stream into byte data frames that can be decoded later.
type FrameSplitter interface {
	// Splitter returns bufio.SplitFunc to split incoming bytes stream into byte data frames that can be decoded later.
	Splitter() bufio.SplitFunc
	// Error returns error if any registered.
	// Use it to signal that communication session can be terminated.
	Error() error
	// BadData returns bad data if any registered.
	// Use it to log which bytes couldn't be parsed.
	BadData() []byte
}

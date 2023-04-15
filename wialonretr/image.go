package wialonretr

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

var _ Data = (*Image)(nil)

const title = 0

// Image represents WialonRetranslator image data block.
type Image []byte

// Decode decodes WialonRetranslator image data block from bytes.
func (i *Image) Decode(data []byte) error {
	var img struct {
		title uint64
		size  uint32
		jpeg  []byte
	}
	reader := bytes.NewReader(data)
	err := binary.Read(reader, binary.BigEndian, &img)
	if err != nil {
		return fmt.Errorf("read title, size, jpeg: %w", err)
	}
	if img.title != title {
		return fmt.Errorf("invalid title: %d", img.title) //nolint:goerr113
	}
	*i = img.jpeg
	return nil
}

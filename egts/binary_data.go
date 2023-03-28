package egts

// BinaryData interface for working with binary sections.
type BinaryData interface {
	// Decode parses the set of bytes into the packet structure.
	Decode([]byte) error
	// Encode encodes the packet structure into the set of bytes.
	Encode() ([]byte, error)
	// Length returns the length of the binary data.
	Length() uint16
}

package types

func NewHash() Hash {
	return Hash{}
}

// SetBytes set bytes
// Parameters:
//   - b: []byte, the value of Hash
//
// Returns:
func (h *Hash) SetBytes(b []byte) {
	if len(b) > len(h) {
		b = b[len(b)-HashLength:]
	}
	copy(h[HashLength-len(b):], b)
}

func (h *Hash) SetString(hexString string) {

}

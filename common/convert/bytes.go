package convert

import "lattice-go/common/types"

// BytesToHash convert bytes to hash
// Parameters
//   - b []byte
//
// Returns
//   - types.Hash
func BytesToHash(b []byte) types.Hash {
	var h types.Hash
	copy(h[:], b)
	return h
}

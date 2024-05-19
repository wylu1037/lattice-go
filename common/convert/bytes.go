package convert

import "lattice-go/common/types"

func BytesToHash(b []byte) types.Hash {
	var h types.Hash
	copy(h[:], b)
	return h
}

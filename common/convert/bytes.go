package convert

import (
	"fmt"
	"github.com/btcsuite/btcutil/base58"
	"github.com/ethereum/go-ethereum/common"
	"lattice-go/common/types"
	"strings"
)

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

func AddressToZltc(address common.Address) string {
	return fmt.Sprintf("zltc_%s", base58.CheckEncode(address.Bytes(), types.AddressVersion))
}

func ZltcToAddress(zltc string) (common.Address, error) {
	elem := strings.SplitN(zltc, "_", 2)
	if len(elem) != 2 {
		return common.Address{}, fmt.Errorf("invalid address: %s", zltc)
	}
	if elem[0] != types.AddressTitle {
		return common.Address{}, fmt.Errorf("invalid address: %s", zltc)
	}
	dec, version, err := base58.CheckDecode(elem[1])
	if version != types.AddressVersion || err != nil {
		return common.Address{}, fmt.Errorf("invalid address: %s", zltc)
	}
	return common.BytesToAddress(dec), nil
}

package convert

import (
	"fmt"
	"github.com/btcsuite/btcutil/base58"
	"github.com/ethereum/go-ethereum/common"
	"lattice-go/common/types"
	"strings"
)

func AddressToZltc(address common.Address) string {
	return fmt.Sprintf("zltc_%s", base58.CheckEncode(address.Bytes(), types.AddressVersion))
}

// ZltcToAddress 将ZLTC地址转为ETH地址
// Parameters
//   - s string: zltc地址，Example：zltc_dhdfbm9JEoyDvYoCDVsABiZj52TAo9Ei6
//
// Returns
//   - Address: 0x9293c604c644bfac34f498998cc3402f203d4d6b
//   - error:
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

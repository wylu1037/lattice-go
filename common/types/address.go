package types

import (
	"encoding/hex"
	"fmt"
	"github.com/btcsuite/btcutil/base58"
	"golang.org/x/crypto/sha3"
	"lattice-go/common/constant"
	"strings"
)

func NewAddress(bytes []byte) Address {
	var addr Address
	addr.SetBytes(bytes)
	return addr
}

func (addr *Address) SetBytes(bytes []byte) {
	if len(bytes) > len(addr) {
		bytes = bytes[len(bytes)-AddressLength:]
	}
	copy(addr[AddressLength-len(bytes):], bytes)
}

func (addr *Address) Base58CheckSum() string {
	return base58.CheckEncode(addr[:], AddressVersion)
}

func (addr *Address) String() string {
	return fmt.Sprintf("%s_%s", AddressTitle, addr.Base58CheckSum())
}

func (addr *Address) Hex() string {
	unCheckSummed := hex.EncodeToString(addr[:])
	sha := sha3.NewLegacyKeccak256()
	sha.Write([]byte(unCheckSummed))
	hash := sha.Sum(nil)

	result := []byte(unCheckSummed)
	for i, c := range result {
		if c > '9' {
			hashByte := hash[i/2]
			if i%2 == 0 {
				hashByte >>= 4
			} else {
				hashByte &= 0xf
			}
			if hashByte > 7 {
				result[i] -= 32 // convert to uppercase
			}
		}
	}
	return constant.HexPrefix + string(result)
}

// ZltcToAddress 将ZLTC地址转为ETH地址
// Parameters
//   - s string: zltc地址，Example：zltc_dhdfbm9JEoyDvYoCDVsABiZj52TAo9Ei6
//
// Returns
//   - Address: 0x9293c604c644bfac34f498998cc3402f203d4d6b
//   - error:
func ZltcToAddress(s string) (Address, error) {
	elem := strings.SplitN(s, "_", 2)
	if len(elem) != 2 {
		return Address{}, fmt.Errorf("invalid address: %s", s)
	}
	if elem[0] != AddressTitle {
		return Address{}, fmt.Errorf("invalid address: %s", s)
	}
	dec, version, err := base58.CheckDecode(elem[1])
	if version != AddressVersion || err != nil {
		return Address{}, fmt.Errorf("invalid address: %s", s)
	}
	return NewAddress(dec), nil
}

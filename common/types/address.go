package types

import (
	"fmt"
	"github.com/btcsuite/btcutil/base58"
	"strings"
)

func NewAddress(bytes []byte) Address {
	var addr Address
	addr.SetBytes(bytes)
	return addr
}

func (a *Address) SetBytes(bytes []byte) {
	if len(bytes) > len(a) {
		bytes = bytes[len(bytes)-AddressLength:]
	}
	copy(a[AddressLength-len(bytes):], bytes)
}

func (a *Address) Base58CheckSum() string {
	return base58.CheckEncode(a[:], AddressVersion)
}

func (a *Address) String() string {
	return fmt.Sprintf("%s_%s", AddressTitle, a.Base58CheckSum())
}

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

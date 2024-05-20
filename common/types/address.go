package types

import (
	"fmt"
	"github.com/btcsuite/btcutil/base58"
)

func NewAddress() {

}

func (a Address) Base58CheckSum() string {
	return base58.CheckEncode(a[:], AddressVersion)
}

func (a Address) String() string {
	return fmt.Sprintf("%s_%s", AddressTitle, a.Base58CheckSum())
}

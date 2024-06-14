package constant

import (
	"github.com/ethereum/go-ethereum/common"
)

const (
	HexPrefix      = "0x"
	AddressVersion = 1
	AddressLength  = common.AddressLength // 20 byte
	AddressTitle   = "zltc"
	HashLength     = 32 // 32 byte

	Sm2p256v1SignatureLength      = 97 // uint is char, Example:
	Secp256k1SignatureLength      = 65
	Sm2p256v1SignatureRemark byte = 1
)

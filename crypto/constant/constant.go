package constant

import "lattice-go/common/types"

const (
	Secp256k1 types.Curve = "secp256k1"
	Sm2p256v1 types.Curve = "sm2p256v1"
)

const (
	Sm2p256v1SignatureLength      = 97 // uint is char, Example:
	Secp256k1SignatureLength      = 65
	Sm2p256v1SignatureRemark byte = 1
)

package constant

// Curve Elliptic curve
type Curve string

const (
	Secp256k1 Curve = "secp256k1"
	Sm2p256v1 Curve = "sm2p256v1"
)

const (
	Sm2p256v1SignatureLength = 97 // uint is char, Example:
	Secp256k1SignatureLength = 65
)

package types

const (
	AddressVersion = 1
	AddressLength  = 20 // 20 byte
	AddressTitle   = "zltc"
	HashLength     = 32 // 32 byte
)

// Address define `Address` type
type Address [AddressLength]byte

// Hash define `Hash` type
type Hash [HashLength]byte

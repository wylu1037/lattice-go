package convert

import (
	"crypto/ecdsa"
	"github.com/tjfoc/gmsm/sm2"
)

func EcdsaSKToSm2SK(sk *ecdsa.PrivateKey) *sm2.PrivateKey {
	return &sm2.PrivateKey{
		PublicKey: sm2.PublicKey{
			Curve: sk.Curve,
			X:     sk.X,
			Y:     sk.Y,
		},
		D: sk.D,
	}
}

func EcdsaPKToSm2PK(pk *ecdsa.PublicKey) *sm2.PublicKey {
	return &sm2.PublicKey{
		Curve: pk.Curve,
		X:     pk.X,
		Y:     pk.Y,
	}
}

func Sm2PKToEcdsaPK(pk *sm2.PublicKey) *ecdsa.PublicKey {
	return &ecdsa.PublicKey{
		Curve: pk.Curve,
		X:     pk.X,
		Y:     pk.Y,
	}
}

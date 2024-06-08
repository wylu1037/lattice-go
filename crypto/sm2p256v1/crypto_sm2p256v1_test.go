package sm2p256v1

import (
	"fmt"
	"testing"
)

func TestSm2p256v1Api_GenerateKeyPair(t *testing.T) {
	crypto := New()
	sk, err := crypto.GenerateKeyPair()
	if err != nil {
		t.Error(err)
	}
	t.Log(sk)
	skHexString, err := crypto.SKToHexString(sk)
	if err != nil {
		t.Error(err)
	}
	pkHexString, err := crypto.PKToHexString(&sk.PublicKey)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(skHexString)
	fmt.Println(pkHexString)
}

func TestSm2p256v1Api_HexToSK(t *testing.T) {
	sk := "0xb3e4575b72bffe9e27d7bb75f56cfbefcb0da2bfc2b457369b674c13662b0b9b"
	crypto := New()
	priv, err := crypto.HexToSK(sk)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(crypto.PKToHexString(&priv.PublicKey))
}

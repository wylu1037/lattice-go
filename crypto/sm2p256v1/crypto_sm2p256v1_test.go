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

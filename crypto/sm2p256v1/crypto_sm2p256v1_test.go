package sm2p256v1

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/stretchr/testify/assert"
	"github.com/wylu1037/lattice-go/common/convert"
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

func TestSm2p256v1Api_Sign(t *testing.T) {
	crypto := New()
	sk, err := crypto.GenerateKeyPair()
	if err != nil {
		t.Error(err)
	}
	fmt.Println(crypto.SKToHexString(sk))
	fmt.Println(crypto.PKToHexString(&sk.PublicKey))

	hash := []byte{1, 2, 3, 4, 5, 6, 7, 8, 1, 2, 3, 4, 5, 6, 7, 8, 1, 2, 3, 4, 5, 6, 7, 8, 1, 2, 3, 4, 5, 6, 7, 8}
	signature, err := crypto.Sign(hash, sk)
	if err != nil {
		t.Error(err)
	} else {
		fmt.Println(hexutil.Encode(signature))
	}
	pass := crypto.Verify(hash, signature, &sk.PublicKey)
	fmt.Println(pass)
}

func TestSm2p256v1Api_PKToAddress(t *testing.T) {
	crypto := New()
	sk, err := crypto.GenerateKeyPair()
	assert.Nil(t, err)
	addr, err := crypto.PKToAddress(&sk.PublicKey)
	assert.Nil(t, err)
	t.Log("ETH Address:", addr.Hex())
	t.Log("ZLTC Address:", convert.AddressToZltc(addr))
	assert.Len(t, addr, 20)
}

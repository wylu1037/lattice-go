package lattice

import (
	"github.com/stretchr/testify/assert"
	"lattice-go/common/convert"
	"lattice-go/crypto"

	"testing"
)

func TestNewLattice(t *testing.T) {
	c := crypto.NewCrypto(crypto.Sm2p256v1)
	sk, err := c.GenerateKeyPair()
	assert.Nil(t, err)
	addr, err := c.PKToAddress(&sk.PublicKey)
	assert.Nil(t, err)
	t.Log("ETH Address:", addr.Hex())
	t.Log("ZLTC Address:", convert.AddressToZltc(addr))
}

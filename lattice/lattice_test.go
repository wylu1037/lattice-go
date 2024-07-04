package lattice

import (
	"context"
	"github.com/stretchr/testify/assert"
	"lattice-go/crypto"
	"testing"
)

func TestLattice_Transfer(t *testing.T) {
	lattice := NewLattice(
		&ChainConfig{ChainId: 1, Curve: crypto.Sm2p256v1},
		&NodeConfig{Ip: "192.168.1.185", HttpPort: 13000},
		&IdentityConfig{AccountAddress: "zltc_Z1pnS94bP4hQSYLs4aP4UwBP9pH8bEvhi", PrivateKey: "0x23d5b2a2eb0a9c8b86d62cbc3955cfd1fb26ec576ecc379f402d0f5d2b27a7bb"},
		&Options{},
	)

	hash, err := lattice.Transfer(context.Background(), "zltc_S5KXbs6gFkEpSnfNpBg3DvZHnB9aasa6Q", "0x10")
	assert.NoError(t, err)
	t.Log(hash.String())
}

func TestLattice_DeployContract(t *testing.T) {

}

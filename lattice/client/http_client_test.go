package client

import (
	"context"
	"github.com/stretchr/testify/assert"
	"lattice-go/crypto"
	"lattice-go/lattice/block"
	"testing"
)

func TestHttpApi_GetLatestBlock(t *testing.T) {
	api := NewHttpApi("http://192.168.1.185:13000", "1")
	latestBlock, err := api.GetLatestBlock(context.Background(), "zltc_Z1pnS94bP4hQSYLs4aP4UwBP9pH8bEvhi")
	assert.Nil(t, err)

	transaction := block.NewTransferTXBuilder().
		SetLatestBlock(latestBlock).
		SetOwner("zltc_Z1pnS94bP4hQSYLs4aP4UwBP9pH8bEvhi").
		SetLinker("zltc_S5KXbs6gFkEpSnfNpBg3DvZHnB9aasa6Q").
		SetPayload("0x02").
		Build()

	err = transaction.SignTX(1, crypto.Sm2p256v1, "0x23d5b2a2eb0a9c8b86d62cbc3955cfd1fb26ec576ecc379f402d0f5d2b27a7bb")
	assert.Nil(t, err)

	hash, err := api.SendSignedTransaction(context.Background(), transaction)
	assert.Nil(t, err)
	t.Log(hash.String())
}

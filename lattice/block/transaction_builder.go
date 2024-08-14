package block

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/wylu1037/lattice-go/common/types"
	"math/big"
	"time"
)

type TransactionBuilder interface {
	Build() *Transaction
	SetLatestBlock(block *types.LatestBlock) TransactionBuilder
	SetOwner(owner string) TransactionBuilder
	SetLinker(linker string) TransactionBuilder
	SetCode(code string) TransactionBuilder
	SetPayload(payload string) TransactionBuilder
	SetAmount(amount uint64) TransactionBuilder
	SetJoule(joule uint64) TransactionBuilder
}

type transactionBuilder struct {
	Transaction *Transaction
}

func NewTransactionBuilder(transactionType TransactionType) TransactionBuilder {
	return &transactionBuilder{
		Transaction: &Transaction{
			Type:      transactionType,
			Timestamp: uint64(time.Now().Unix()),
			Hub:       make([]common.Hash, 0),
		},
	}
}

func (builder *transactionBuilder) Build() *Transaction {
	return builder.Transaction
}

func (builder *transactionBuilder) SetLatestBlock(block *types.LatestBlock) TransactionBuilder {
	builder.Transaction.Height = block.Height + 1
	builder.Transaction.ParentHash = block.Hash
	builder.Transaction.DaemonHash = block.DaemonBlockHash
	return builder
}

func (builder *transactionBuilder) SetOwner(owner string) TransactionBuilder {
	builder.Transaction.Owner = owner
	return builder
}

func (builder *transactionBuilder) SetLinker(linker string) TransactionBuilder {
	builder.Transaction.Linker = linker
	return builder
}

func (builder *transactionBuilder) SetCode(code string) TransactionBuilder {
	builder.Transaction.Code = code
	return builder
}

func (builder *transactionBuilder) SetPayload(payload string) TransactionBuilder {
	builder.Transaction.Payload = payload
	return builder
}

func (builder *transactionBuilder) SetAmount(amount uint64) TransactionBuilder {
	bigAmount := big.NewInt(0)
	bigAmount.SetUint64(amount)
	builder.Transaction.Amount = bigAmount
	return builder
}

func (builder *transactionBuilder) SetJoule(joule uint64) TransactionBuilder {
	bigJoule := big.NewInt(0)
	bigJoule.SetUint64(joule)
	builder.Transaction.Joule = bigJoule
	return builder
}

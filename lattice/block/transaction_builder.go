package block

import (
	"github.com/ethereum/go-ethereum/common"
	"lattice-go/common/types"
	"time"
)

type TransactionBuilder interface {
	Build() *Transaction
	SetLatestBlock(block *types.LatestBlock) TransactionBuilder
	SetOwner(owner string) TransactionBuilder
	SetLinker(linker string) TransactionBuilder
	SetCode(code string) TransactionBuilder
	SetPayload(payload string) TransactionBuilder
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

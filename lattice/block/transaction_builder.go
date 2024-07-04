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

type transferTXBuilder struct {
	Transaction *Transaction
}

func NewTransferTXBuilder() TransactionBuilder {
	return &transferTXBuilder{
		Transaction: &Transaction{
			Type:      TransactionTypeSend,
			Timestamp: uint64(time.Now().Unix()),
			Hub:       make([]common.Hash, 0),
		},
	}
}

func (builder *transferTXBuilder) Build() *Transaction {
	return builder.Transaction
}

func (builder *transferTXBuilder) SetLatestBlock(block *types.LatestBlock) TransactionBuilder {
	builder.Transaction.Height = block.Height + 1
	builder.Transaction.ParentHash = block.Hash
	builder.Transaction.DaemonHash = block.DaemonBlockHash
	return builder
}

func (builder *transferTXBuilder) SetOwner(owner string) TransactionBuilder {
	builder.Transaction.Owner = owner
	return builder
}

func (builder *transferTXBuilder) SetLinker(linker string) TransactionBuilder {
	builder.Transaction.Linker = linker
	return builder
}

func (builder *transferTXBuilder) SetCode(code string) TransactionBuilder {
	builder.Transaction.Code = code
	return builder
}

func (builder *transferTXBuilder) SetPayload(payload string) TransactionBuilder {
	builder.Transaction.Payload = payload
	return builder
}

type deployContractTXBuilder struct {
	Transaction *Transaction
}

func NewDeployContractTXBuilder() TransactionBuilder {
	return &deployContractTXBuilder{
		Transaction: &Transaction{
			Type:      TransactionTypeDeployContract,
			Timestamp: uint64(time.Now().Unix()),
			Hub:       make([]common.Hash, 0),
		},
	}
}

func (builder *deployContractTXBuilder) Build() *Transaction {
	return builder.Transaction
}

func (builder *deployContractTXBuilder) SetLatestBlock(block *types.LatestBlock) TransactionBuilder {
	builder.Transaction.Height = block.Height + 1
	builder.Transaction.ParentHash = block.Hash
	builder.Transaction.DaemonHash = block.DaemonBlockHash
	return builder
}

func (builder *deployContractTXBuilder) SetOwner(owner string) TransactionBuilder {
	builder.Transaction.Owner = owner
	return builder
}

func (builder *deployContractTXBuilder) SetLinker(linker string) TransactionBuilder {
	builder.Transaction.Linker = linker
	return builder
}

func (builder *deployContractTXBuilder) SetCode(code string) TransactionBuilder {
	builder.Transaction.Code = code
	return builder
}

func (builder *deployContractTXBuilder) SetPayload(payload string) TransactionBuilder {
	builder.Transaction.Payload = payload
	return builder
}

type callContractTXBuilder struct {
	Transaction *Transaction
}

func NewCallContractTXBuilder() TransactionBuilder {
	return &callContractTXBuilder{
		Transaction: &Transaction{
			Type:      TransactionTypeCallContract,
			Timestamp: uint64(time.Now().Unix()),
			Hub:       make([]common.Hash, 0),
		},
	}
}

func (builder *callContractTXBuilder) Build() *Transaction {
	return builder.Transaction
}

func (builder *callContractTXBuilder) SetLatestBlock(block *types.LatestBlock) TransactionBuilder {
	builder.Transaction.Height = block.Height + 1
	builder.Transaction.ParentHash = block.Hash
	builder.Transaction.DaemonHash = block.DaemonBlockHash
	return builder
}

func (builder *callContractTXBuilder) SetOwner(owner string) TransactionBuilder {
	builder.Transaction.Owner = owner
	return builder
}

func (builder *callContractTXBuilder) SetLinker(linker string) TransactionBuilder {
	builder.Transaction.Linker = linker
	return builder
}

func (builder *callContractTXBuilder) SetCode(code string) TransactionBuilder {
	builder.Transaction.Code = code
	return builder
}

func (builder *callContractTXBuilder) SetPayload(payload string) TransactionBuilder {
	builder.Transaction.Payload = payload
	return builder
}

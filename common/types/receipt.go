package types

import "github.com/ethereum/go-ethereum/common"

type Receipt struct {
	ConfirmedTimestamp int64       `json:"confirmTime"`
	Success            bool        `json:"success"`
	ReceiptIndex       uint64      `json:"receiptIndex"`
	TBlockHash         common.Hash `json:"tBlockHash"`
	ContractAddress    string      `json:"contractAddress"`
	ContractRet        string      `json:"contractRet"`
	JouleUsed          uint64      `json:"jouleUsed"`
	Events             []*Event    `json:"events"`
	DBlockHash         common.Hash `json:"dBlockHash"`
	DBlockNumber       uint64      `json:"dBlockNumber"`
}

type Event struct {
	Address      string        `json:"address"`
	Topics       []common.Hash `json:"topics"`
	Data         []byte        `json:"data"`
	Index        uint          `json:"logIndex"`
	TBlockHash   common.Hash   `json:"tblockHash"`
	DBlockNumber uint64        `json:"dblockNumber"`
	Removed      bool          `json:"removed"`
	DataHex      string        `json:"dataHex"`
}

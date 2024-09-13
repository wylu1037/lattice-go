package types

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog/log"
	"strconv"
)

type Receipt struct {
	ConfirmedTimestamp string      `json:"confirmTime"`
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

// GetConfirmedTimestamp parse timestamp to int64
func (r *Receipt) GetConfirmedTimestamp() int64 {
	if r.ConfirmedTimestamp == "" {
		return 0
	}
	timestamp, err := strconv.ParseInt(r.ConfirmedTimestamp, 10, 64)
	if err != nil {
		log.Error().Err(err).Msgf("解析回执的ConfirmedTimestamp的值%s为int64失败", r.ConfirmedTimestamp)
		return 0
	}
	return timestamp
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

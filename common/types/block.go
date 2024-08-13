package types

import (
	"github.com/ethereum/go-ethereum/common"
)

type LatestBlock struct {
	Height          uint64      `json:"currentTBlockNumber"`
	Hash            common.Hash `json:"currentTBlockHash"`
	DaemonBlockHash common.Hash `json:"currentDBlockHash"`
}

func (b *LatestBlock) IncrHeight() {
	b.Height++
}

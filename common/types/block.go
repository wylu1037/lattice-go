package types

import (
	"github.com/ethereum/go-ethereum/common"
)

type LatestBlock struct {
	Height          uint64      `json:"currentTBlockNumber"` // 账户最新的高度
	Hash            common.Hash `json:"currentTBlockHash"`   // 账户最新的一笔交易哈希
	DaemonBlockHash common.Hash `json:"currentDBlockHash"`   // 守护区块的哈希
}

// IncrHeight 增长高度
func (b *LatestBlock) IncrHeight() {
	b.Height++
}

// DaemonBlock 守护区块
type DaemonBlock struct {
	Hash       string `json:"hash"`
	ParentHash string `json:"parentHash"`
	Height     uint64 `json:"number"`
	Timestamp  uint64 `json:"timestamp"`
	Size       uint32 `json:"size"`
}

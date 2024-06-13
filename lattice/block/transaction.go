package block

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
	"io"
	"lattice-go/common/types"
	"lattice-go/crypto/sm2p256v1"
	"math/big"
)

type Transaction struct {
	Number      int64         `json:"number"`
	Type        string        `json:"type"`
	ParentHash  common.Hash   `json:"parentHash"`
	Hub         []common.Hash `json:"hub"`
	DaemonHash  common.Hash   `json:"daemonHash"`
	CodeHash    common.Hash   `json:"codeHash"`
	Owner       string        `json:"owner"`
	Linker      string        `json:"linker"`
	Amount      *big.Int      `json:"amount"`
	Joule       int64         `json:"joule"`
	Difficulty  int64         `json:"difficulty"`
	Pow         *big.Int      `json:"pow"`
	ProofOfWork string        `json:"proofOfWork"`
	Payload     string        `json:"payload"`
	Timestamp   uint64        `json:"timestamp"`
	Code        string        `json:"code"`
	Sign        string        `json:"sign"`
	Hash        string        `json:"hash"`
	Hash2       common.Hash   `json:"hash2"`
	Key         string        `json:"key"`
	DataHash    string        `json:"dataHash"`
	ApplyHash   string        `json:"applyHash"`
}

func (tx *Transaction) RlpEncodeHash(chainId *big.Int, curve types.Curve) common.Hash {
	return sm2p256v1.New().EncodeHash(func(writer io.Writer) {
		err := rlp.Encode(writer, []interface{}{
			tx.Number,
			tx.Type,
			tx.ParentHash,
			tx.Hub,
			tx.DaemonHash,
			tx.CodeHash,
			tx.Owner,
			tx.Linker,
			tx.Amount,
			tx.Joule,
			tx.Difficulty,
			tx.ProofOfWork,
			tx.Payload,
			tx.Timestamp,
			chainId,
			uint(0),
			uint(0),
		})
		if err != nil {
			return
		}
	})
}

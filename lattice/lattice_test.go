package lattice

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/stretchr/testify/assert"
	"io"
	"lattice-go/abi"
	"lattice-go/common/types"
	"lattice-go/crypto"
	"lattice-go/lattice/builtin"
	"testing"
	"time"
)

func TestLattice_TransferWaitReceipt(t *testing.T) {
	lattice := NewLattice(
		&ChainConfig{ChainId: 1, Curve: crypto.Sm2p256v1},
		&NodeConfig{Ip: "192.168.1.185", HttpPort: 13000},
		&CredentialConfig{AccountAddress: "zltc_Z1pnS94bP4hQSYLs4aP4UwBP9pH8bEvhi", PrivateKey: "0x23d5b2a2eb0a9c8b86d62cbc3955cfd1fb26ec576ecc379f402d0f5d2b27a7bb"},
		&Options{},
	)

	hash, receipt, err := lattice.TransferWaitReceipt(context.Background(), "zltc_S5KXbs6gFkEpSnfNpBg3DvZHnB9aasa6Q", "0x10", NewBackOffWaitStrategy(10, time.Second))
	assert.NoError(t, err)
	t.Log(hash.String())
	t.Log(receipt)
}

func TestLattice_DeployContractWaitReceipt(t *testing.T) {
	lattice := NewLattice(
		&ChainConfig{ChainId: 1, Curve: crypto.Sm2p256v1},
		&NodeConfig{Ip: "192.168.1.185", HttpPort: 13000},
		&CredentialConfig{AccountAddress: "zltc_Z1pnS94bP4hQSYLs4aP4UwBP9pH8bEvhi", PrivateKey: "0x23d5b2a2eb0a9c8b86d62cbc3955cfd1fb26ec576ecc379f402d0f5d2b27a7bb"},
		&Options{},
	)

	data := "0x60806040526000805534801561001457600080fd5b50610278806100246000396000f3fe608060405234801561001057600080fd5b50600436106100415760003560e01c80635b34b96614610046578063a87d942c14610050578063f5c5ad831461006e575b600080fd5b61004e610078565b005b610058610093565b60405161006591906100d0565b60405180910390f35b61007661009c565b005b600160008082825461008a919061011a565b92505081905550565b60008054905090565b60016000808282546100ae91906101ae565b92505081905550565b6000819050919050565b6100ca816100b7565b82525050565b60006020820190506100e560008301846100c1565b92915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b6000610125826100b7565b9150610130836100b7565b9250817f7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0383136000831215161561016b5761016a6100eb565b5b817f80000000000000000000000000000000000000000000000000000000000000000383126000831216156101a3576101a26100eb565b5b828201905092915050565b60006101b9826100b7565b91506101c4836100b7565b9250827f8000000000000000000000000000000000000000000000000000000000000000018212600084121516156101ff576101fe6100eb565b5b827f7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff018213600084121615610237576102366100eb565b5b82820390509291505056fea2646970667358221220d841351625356129f6266ada896818d690dbc4b0d176774a97d745dfbe2fe50164736f6c634300080b0033"
	hash, receipt, err := lattice.DeployContractWaitReceipt(context.Background(), data, "0x", DefaultFixedWaitStrategy())
	assert.NoError(t, err)
	t.Log(hash.String())
	r, err := json.Marshal(receipt)
	assert.NoError(t, err)
	t.Log(string(r))

	// 获取正在进行的提案
	proposal, err := lattice.HttpApi().GetContractLifecycleProposal(context.Background(), receipt.ContractAddress, types.ProposalStateInitial)
	assert.NoError(t, err)

	voteContract := builtin.NewVoteContract()
	approveData, err := voteContract.Approve(proposal[0].Content.Id)
	assert.Nil(t, err)

	hash, receipt, err = lattice.CallContractWaitReceipt(context.Background(), builtin.VoteBuiltinContract.Address, approveData, "0x", DefaultFixedWaitStrategy())
	assert.NoError(t, err)
	t.Log(hash.String())
	re, _ := json.Marshal(receipt)
	t.Log(string(re))
}

func TestLattice_CallContractWaitReceipt(t *testing.T) {
	lattice := NewLattice(
		&ChainConfig{ChainId: 1, Curve: crypto.Sm2p256v1},
		&NodeConfig{Ip: "192.168.1.185", HttpPort: 13000},
		&CredentialConfig{AccountAddress: "zltc_Z1pnS94bP4hQSYLs4aP4UwBP9pH8bEvhi", PrivateKey: "0x23d5b2a2eb0a9c8b86d62cbc3955cfd1fb26ec576ecc379f402d0f5d2b27a7bb"},
		&Options{},
	)

	myAbi := `[{"inputs":[],"name":"decrementCounter","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[],"name":"getCount","outputs":[{"internalType":"int256","name":"","type":"int256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"incrementCounter","outputs":[],"stateMutability":"nonpayable","type":"function"}]`
	function, err := abi.NewAbi(myAbi).GetLatticeFunction("incrementCounter")
	assert.NoError(t, err)
	data, err := function.Encode()
	assert.NoError(t, err)
	hash, receipt, err := lattice.CallContractWaitReceipt(context.Background(), "zltc_a5jerZSQ2UNMhDbUkHaYuVbiBMxkZWETj", data, "0x", DefaultFixedWaitStrategy())
	assert.NoError(t, err)
	t.Log(hash.String())
	r, err := json.Marshal(receipt)
	assert.NoError(t, err)
	t.Log(string(r))
}

func TestLattice_CallContract(t *testing.T) {
	lattice := NewLattice(
		&ChainConfig{
			ChainId: 1,
			Curve:   crypto.Sm2p256v1,
		},
		&NodeConfig{
			Ip:       "192.168.1.185",
			HttpPort: 13000,
		},
		&CredentialConfig{
			AccountAddress: "zltc_Z1pnS94bP4hQSYLs4aP4UwBP9pH8bEvhi",
			PrivateKey:     "0x23d5b2a2eb0a9c8b86d62cbc3955cfd1fb26ec576ecc379f402d0f5d2b27a7bb",
		},
		&Options{},
	)

	voteContract := builtin.NewVoteContract()
	data, err := voteContract.Approve("0x101")
	assert.NoError(t, err)

	hash, receipt, err := lattice.CallContractWaitReceipt(context.Background(), builtin.VoteBuiltinContract.Address, data, "0x", DefaultFixedWaitStrategy())
	assert.NoError(t, err)
	t.Log(hash.String())
	r, err := json.Marshal(receipt)
	assert.NoError(t, err)
	t.Log(string(r))
}

func TestNewLattice(t *testing.T) {
	salt := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		panic("reading from crypto/rand failed: " + err.Error())
	}
	fmt.Println(hexutil.Encode(salt))
}

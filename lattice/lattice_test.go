package lattice

import (
	"context"
	"encoding/json"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/wylu1037/lattice-go/abi"
	"github.com/wylu1037/lattice-go/common/constant"
	"github.com/wylu1037/lattice-go/common/convert"
	"github.com/wylu1037/lattice-go/common/types"
	"github.com/wylu1037/lattice-go/crypto"
	"github.com/wylu1037/lattice-go/lattice/builtin"
	"github.com/wylu1037/lattice-go/lattice/protobuf"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"
)

const (
	proto      = "syntax = \"proto3\";\n\nmessage Student {\n\tstring id = 1;\n\tstring name = 2;\n}"
	counterAbi = `[{"inputs":[],"name":"decrementCounter","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[],"name":"getCount","outputs":[{"internalType":"int256","name":"","type":"int256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"incrementCounter","outputs":[],"stateMutability":"nonpayable","type":"function"}]`
	chainId    = "2"
)

var latticeApi = NewLattice(
	&ChainConfig{Curve: crypto.Sm2p256v1},
	&ConnectingNodeConfig{Ip: "192.168.1.185", HttpPort: 13800},
	NewMemoryBlockCache(10*time.Second, time.Minute, time.Minute),
	NewAccountLock(),
	&Options{MaxIdleConnsPerHost: 200},
)

var credentials = &Credentials{
	AccountAddress: "zltc_Z1pnS94bP4hQSYLs4aP4UwBP9pH8bEvhi",
	PrivateKey:     "0x23d5b2a2eb0a9c8b86d62cbc3955cfd1fb26ec576ecc379f402d0f5d2b27a7bb",
}

// 创建业务合约地址
func TestNewLattice_CreateBusiness(t *testing.T) {
	contract := builtin.NewCredibilityContract()
	data, err := contract.CreateBusiness()
	assert.NoError(t, err)
	hash, receipt, err := latticeApi.CallContractWaitReceipt(context.Background(), credentials, chainId, contract.GetCreateBusinessContractAddress(), data, constant.ZeroPayload, 0, 0, DefaultFixedRetryStrategy())
	assert.NoError(t, err)
	t.Log(hash.String())
	r, err := json.Marshal(receipt)
	assert.NoError(t, err)
	t.Logf("业务合约地址：%s", receipt.ContractRet)
	t.Log(string(r))
}

// 创建协议
func TestNewLattice_CreateProtocol(t *testing.T) {
	contract := builtin.NewCredibilityContract()
	data, err := contract.CreateProtocol(2, []byte(proto))
	assert.NoError(t, err)
	hash, receipt, err := latticeApi.CallContractWaitReceipt(context.Background(), credentials, chainId, contract.ContractAddress(), data, constant.ZeroPayload, 0, 0, DefaultBackOffRetryStrategy())
	assert.NoError(t, err)
	t.Log(hash.String())
	r, err := json.Marshal(receipt)
	assert.NoError(t, err)
	result, err := abi.DecodeReturn(contract.MyAbi(), "addProtocol", receipt.ContractRet)
	assert.NoError(t, err)
	t.Logf("创建协议返回：%s", result)
	t.Log(string(r))
}

// 批量创建协议
func TestNewLattice_BatchCreateProtocol(t *testing.T) {
	contract := builtin.NewCredibilityContract()
	request := make([]builtin.CreateProtocolRequest, 2)
	request[0] = builtin.CreateProtocolRequest{
		ProtocolSuite: 10,
		Data:          convert.BytesToBytes32Arr([]byte("syntax = \"proto3\";\n\nmessage Student {\n\tstring id = 1;\n\tstring name = 2;\n}")),
	}
	request[1] = builtin.CreateProtocolRequest{
		ProtocolSuite: 10,
		Data:          convert.BytesToBytes32Arr([]byte("syntax = \"proto3\";\n\nmessage Student {\n\tstring id = 1;\n\tstring name = 2;\n}")),
	}
	data, err := contract.BatchCreateProtocol(request)
	assert.NoError(t, err)
	hash, receipt, err := latticeApi.CallContractWaitReceipt(context.Background(), credentials, chainId, contract.ContractAddress(), data, constant.ZeroPayload, 0, 0, DefaultBackOffRetryStrategy())
	assert.NoError(t, err)
	t.Log(hash.String())
	r, err := json.Marshal(receipt)
	assert.NoError(t, err)
	result, err := abi.DecodeReturn(contract.MyAbi(), "addProtocolBatch", receipt.ContractRet)
	assert.NoError(t, err)
	t.Logf("批量创建协议返回：%s", result)
	t.Log(string(r))
}

// 数据存证
func TestNewLattice_Write(t *testing.T) {
	contract := builtin.NewCredibilityContract()
	jsonData := `{"id":"1","name":"jack"}`

	fd := protobuf.MakeFileDescriptor(strings.NewReader(proto))
	dataBytes, err := protobuf.MarshallMessage(fd, jsonData)
	assert.NoError(t, err)
	data, err := contract.Write(&builtin.WriteLedgerRequest{
		ProtocolUri: 42949672968,
		Hash:        "2",
		Data:        convert.BytesToBytes32Arr(dataBytes),
		Address:     common.HexToAddress("0x130c46461ff0e4fe8d6660cf4c8d88afb9ab1daf"),
	})
	assert.NoError(t, err)
	hash, receipt, err := latticeApi.CallContractWaitReceipt(context.Background(), credentials, chainId, contract.ContractAddress(), data, constant.ZeroPayload, 0, 0, DefaultBackOffRetryStrategy())
	assert.NoError(t, err)
	t.Logf("结束数据存证，交易哈希为：%s", hash.String())
	r, err := json.Marshal(receipt)
	assert.NoError(t, err)
	t.Log(string(r))
}

func TestLattice_Transfer(t *testing.T) {
	for i := 0; i < 100; i++ {
		go func() {
			hash, err := latticeApi.Transfer(context.Background(), credentials, chainId, "zltc_S5KXbs6gFkEpSnfNpBg3DvZHnB9aasa6Q", "0x10", 0, 0)
			assert.NoError(t, err)
			t.Log(hash.String())
		}()
	}
	time.Sleep(10 * time.Second)
}

func TestLattice_TransferWaitReceipt(t *testing.T) {
	hash, receipt, err := latticeApi.TransferWaitReceipt(context.Background(), credentials, chainId, "zltc_S5KXbs6gFkEpSnfNpBg3DvZHnB9aasa6Q", "0x10", 0, 0, NewBackOffRetryStrategy(10, time.Second))
	assert.NoError(t, err)
	t.Log(hash.String())
	t.Log(receipt)
}

func TestLattice_DeployContractWaitReceipt(t *testing.T) {
	data := "0x60806040526000805534801561001457600080fd5b50610278806100246000396000f3fe608060405234801561001057600080fd5b50600436106100415760003560e01c80635b34b96614610046578063a87d942c14610050578063f5c5ad831461006e575b600080fd5b61004e610078565b005b610058610093565b60405161006591906100d0565b60405180910390f35b61007661009c565b005b600160008082825461008a919061011a565b92505081905550565b60008054905090565b60016000808282546100ae91906101ae565b92505081905550565b6000819050919050565b6100ca816100b7565b82525050565b60006020820190506100e560008301846100c1565b92915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b6000610125826100b7565b9150610130836100b7565b9250817f7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0383136000831215161561016b5761016a6100eb565b5b817f80000000000000000000000000000000000000000000000000000000000000000383126000831216156101a3576101a26100eb565b5b828201905092915050565b60006101b9826100b7565b91506101c4836100b7565b9250827f8000000000000000000000000000000000000000000000000000000000000000018212600084121516156101ff576101fe6100eb565b5b827f7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff018213600084121615610237576102366100eb565b5b82820390509291505056fea2646970667358221220d841351625356129f6266ada896818d690dbc4b0d176774a97d745dfbe2fe50164736f6c634300080b0033"
	hash, receipt, err := latticeApi.DeployContractWaitReceipt(context.Background(), credentials, chainId, data, constant.ZeroPayload, 0, 0, DefaultFixedRetryStrategy())
	assert.NoError(t, err)
	t.Log(hash.String())
	r, err := json.Marshal(receipt)
	assert.NoError(t, err)
	t.Log(string(r))

	// 获取正在进行的提案
	proposal, err := latticeApi.HttpApi().GetContractLifecycleProposal(context.Background(), chainId, receipt.ContractAddress, types.ProposalStateINITIAL)
	assert.NoError(t, err)

	voteContract := builtin.NewProposalContract()
	approveData, err := voteContract.Approve(proposal[0].Content.Id)
	assert.Nil(t, err)

	hash, receipt, err = latticeApi.CallContractWaitReceipt(context.Background(), credentials, chainId, builtin.ProposalBuiltinContract.Address, approveData, constant.ZeroPayload, 0, 0, DefaultFixedRetryStrategy())
	assert.NoError(t, err)
	t.Log(hash.String())
	re, _ := json.Marshal(receipt)
	t.Log(string(re))
}

func TestLattice_CallContractWaitReceipt(t *testing.T) {
	function, err := abi.NewAbi(counterAbi).GetLatticeFunction("incrementCounter")
	assert.NoError(t, err)
	data, err := function.Encode()
	assert.NoError(t, err)
	hash, receipt, err := latticeApi.CallContractWaitReceipt(context.Background(), credentials, chainId, "zltc_TpPbhQmKH8YoLF6aqLw77CiEoQxFL6SaM", data, constant.ZeroPayload, 0, 0, DefaultFixedRetryStrategy())
	assert.NoError(t, err)
	t.Log(hash.String())
	r, err := json.Marshal(receipt)
	assert.NoError(t, err)
	t.Log(string(r))
}

func TestLattice_PreCallContract(t *testing.T) {
	function, err := abi.NewAbi(counterAbi).GetLatticeFunction("getCount")
	assert.NoError(t, err)
	data, err := function.Encode()
	assert.NoError(t, err)
	receipt, err := latticeApi.PreCallContract(context.Background(), chainId, credentials.AccountAddress, "zltc_ebuyF9qei2hoESDzpFVM2cVFC9ViJXyNn", data, constant.ZeroPayload)
	assert.Nil(t, err)
	receiptBytes, _ := json.Marshal(receipt)
	t.Log(string(receiptBytes))
	assert.NotNil(t, receipt)
}

func TestLattice_CallContract(t *testing.T) {
	startTime := time.Now()

	contract := builtin.NewCredibilityContract()
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(contract builtin.CredibilityContract) {
			defer wg.Done()
			data, err := contract.Write(&builtin.WriteLedgerRequest{
				ProtocolUri: 8589934595,
				Hash:        strconv.FormatInt(int64(i)+1000, 10),
				Data:        convert.BytesToBytes32Arr([]byte{1, 2, 3, 4, 5, 6, 7, 8}),
				Address:     convert.ZltcMustToAddress("zltc_YBomBNykwMqxm719giBL3VtYV4ABT9a8D"),
			})
			assert.NoError(t, err)
			hash, err := latticeApi.CallContract(context.Background(), credentials, chainId, contract.ContractAddress(), data, constant.ZeroPayload, 0, 0)
			assert.NoError(t, err)
			t.Log(hash.String())
		}(contract)
	}
	wg.Wait()

	elapsedTime := time.Since(startTime)
	t.Logf("Elapsed Time: %v", elapsedTime)
}

func TestLattice_JsonRpc(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	proposal := new([]types.Proposal[types.ContractLifecycleProposal])
	err := latticeApi.HttpApi().GetProposal(ctx, "2", "", types.ProposalTypeNone, types.ProposalStateNONE, "zltc_YNHgouP2aJgVnt2HncRtoj2mDqcb34Vwz", "", "20240904", "20240904", proposal)
	assert.NoError(t, err)
	t.Log(proposal)
}

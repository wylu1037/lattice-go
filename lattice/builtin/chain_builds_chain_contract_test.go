package builtin

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/wylu1037/lattice-go/common/convert"
	"math/big"
	"testing"
)

func TestChainBuildsChainContract_NewSubChain(t *testing.T) {
	contract := NewChainBuildsChainContract()
	mem := make([]SubChainMember, 1)
	for i := 0; i < len(mem); i++ {
		mem[i] = SubChainMember{
			Type:    1,
			Address: "zltc_YBomBNykwMqxm719giBL3VtYV4ABT9a8D",
		}
	}
	data, err := contract.NewSubChain(&NewSubChainRequest{
		ChannelId:            big.NewInt(101),
		ChannelName:          "channel101",
		Desc:                 "channel",
		BootStrap:            "createNode",
		Preacher:             "zltc_Z1pnS94bP4hQSYLs4aP4UwBP9pH8bEvhi",
		ChannelMemberGroup:   mem,
		Consensus:            1,
		Tokenless:            true,
		GodAmount:            big.NewInt(100),
		Period:               1000,
		NoEmptyAnchor:        false,
		EmptyAnchorPeriodMul: 3,
		IsContractVote:       true,
		IsDictatorship:       true,
		DeployRule:           1,
		ContractPermission:   true,
		ChainByChainVote:     1,
		ProposalExpireTime:   1000,
	})
	assert.NoError(t, err)
	expectData := "0x10b1efb1000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000001da7b22636f6e73656e737573223a312c22746f6b656e6c657373223a747275652c22676f64416d6f756e74223a3130302c22706572696f64223a313030302c226e6f456d707479416e63686f72223a66616c73652c22656d707479416e63686f72506572696f644d756c223a332c226973436f6e7472616374566f7465223a747275652c2269734469637461746f7273686970223a747275652c226465706c6f7952756c65223a312c226e616d65223a226368616e6e656c313031222c22636861696e4964223a3130312c227072656163686572223a227a6c74635f5a31706e533934625034685153594c733461503455774250397048386245766869222c22626f6f745374726170223a226372656174654e6f6465222c22636861696e4d656d62657247726f7570223a5b7b226d656d62657254797065223a312c226d656d626572223a22307835363137313766373932326132333337323061653338616361613431373463646130626631373636227d5d2c22636f6e74726163745065726d697373696f6e223a747275652c22636861696e4279436861696e566f7465223a312c2270726f706f73616c45787069726554696d65223a313030302c2264657363223a226368616e6e656c222c226578747261223a6e756c6c7d000000000000"
	assert.Equal(t, expectData, data)
}

func TestChainBuildsChainContract_StartSubChain(t *testing.T) {
	contract := NewChainBuildsChainContract()
	data, err := contract.StartSubChain("101")
	assert.NoError(t, err)
	expectData := "0x7b777ddf0000000000000000000000000000000000000000000000000000000000000065"
	assert.Equal(t, expectData, data)

}

func TestChainBuildsChainContract_StopSubChain(t *testing.T) {
	contract := NewChainBuildsChainContract()
	data, err := contract.StopSubChain("101")
	assert.NoError(t, err)
	expectData := "0x27c9d3c80000000000000000000000000000000000000000000000000000000000000065"
	assert.Equal(t, expectData, data)
}

func TestChainBuildsChainContract_DeleteSubChain(t *testing.T) {
	contract := NewChainBuildsChainContract()
	data, err := contract.DeleteSubChain("101")
	assert.NoError(t, err)
	expectData := "0x34084eb10000000000000000000000000000000000000000000000000000000000000065"
	assert.Equal(t, expectData, data)
}

func TestChainBuildsChainContract_JoinSubChain(t *testing.T) {
	contract := NewChainBuildsChainContract()
	mem := make([]common.Address, 1)
	for i := 0; i < len(mem); i++ {
		mem[i], _ = convert.ZltcToAddress("zltc_Z1pnS94bP4hQSYLs4aP4UwBP9pH8bEvhi")
	}
	data, err := contract.JoinSubChain(&JoinSubChainRequest{
		ChannelId:     big.NewInt(101),
		NetworkId:     101,
		NodeInfo:      "zltc_YBomBNykwMqxm719giBL3VtYV4ABT9a8D",
		AccessMembers: mem,
	})
	assert.NoError(t, err)
	expectData := "0x1ca9c48700000000000000000000000000000000000000000000000000000000000000650000000000000000000000000000000000000000000000000000000000000065000000000000000000000000000000000000000000000000000000000000008000000000000000000000000000000000000000000000000000000000000000e000000000000000000000000000000000000000000000000000000000000000267a6c74635f59426f6d424e796b774d71786d3731396769424c33567459563441425439613844000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000010000000000000000000000005f2be9a02b43f748ee460bf36eed24fafa109920"
	assert.Equal(t, expectData, data)
}

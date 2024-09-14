package builtin

import (
	"github.com/stretchr/testify/assert"
	"github.com/wylu1037/lattice-go/common/types"
	"testing"
)

func TestModifyChainConfigurationContract_AddConsensusNodes(t *testing.T) {
	contract := NewModifyChainConfigurationContract()
	data, err := contract.AddConsensusNodes([]string{"zltc_dhdfbm9JEoyDvYoCDVsABiZj52TAo9Ei6"})
	assert.Nil(t, err)
	expect := "0x8bd24adc000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000010000000000000000000000009293c604c644bfac34f498998cc3402f203d4d6b"
	assert.Equal(t, expect, data)
}

func TestModifyChainConfigurationContract_DeleteConsensusNodes(t *testing.T) {
	contract := NewModifyChainConfigurationContract()
	data, err := contract.DeleteConsensusNodes([]string{"0x5f2be9a02b43f748ee460bf36eed24fafa109920"})
	assert.Nil(t, err)
	expect := "0x08ce76a7000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000010000000000000000000000005f2be9a02b43f748ee460bf36eed24fafa109920"
	assert.Equal(t, expect, data)
}

func TestModifyChainConfigurationContract_ReplaceConsensusNodes(t *testing.T) {
	contract := NewModifyChainConfigurationContract()
	data, err := contract.ReplaceConsensusNodes("0x5f2be9a02b43f748ee460bf36eed24fafa109920", "zltc_jF4U7umzNpiE8uU35RCBp9f2qf53H5CZZ")
	assert.Nil(t, err)
	expect := "0x67bc37ed0000000000000000000000005f2be9a02b43f748ee460bf36eed24fafa109920000000000000000000000000cf5e003f56d2b75844b741f491861b9fa6daa7c6"
	assert.Equal(t, expect, data)
}

func TestModifyChainConfigurationContract_EnableContractLifecycleVotingDictatorship(t *testing.T) {
	contract := NewModifyChainConfigurationContract()
	data, err := contract.EnableContractLifecycleVotingDictatorship(true)
	assert.Nil(t, err)
	expect := "0x531770b30000000000000000000000000000000000000000000000000000000000000001"
	assert.Equal(t, expect, data)
}

func TestModifyChainConfigurationContract_EnableContractLifecycle(t *testing.T) {
	contract := NewModifyChainConfigurationContract()
	data, err := contract.EnableContractLifecycle(false)
	assert.Nil(t, err)
	expect := "0xb223a59b0000000000000000000000000000000000000000000000000000000000000000"
	assert.Equal(t, expect, data)
}

func TestModifyChainConfigurationContract_EnableContractManagement(t *testing.T) {
	contract := NewModifyChainConfigurationContract()
	data, err := contract.EnableContractManagement(true)
	assert.Nil(t, err)
	expect := "0x6bbb4f0a0000000000000000000000000000000000000000000000000000000000000001"
	assert.Equal(t, expect, data)
}

func TestModifyChainConfigurationContract_UpdateConsensus(t *testing.T) {
	contract := NewModifyChainConfigurationContract()
	data, err := contract.UpdateConsensus(types.ConsensusPOA)
	assert.Nil(t, err)
	expect := "0x3545cf0600000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000000003506f410000000000000000000000000000000000000000000000000000000000"
	assert.Equal(t, expect, data)
}

func TestModifyChainConfigurationContract_UpdateContractDeploymentVotingRule(t *testing.T) {
	contract := NewModifyChainConfigurationContract()
	data, err := contract.UpdateContractDeploymentVotingRule(types.VotingRuleCONSENSUS)
	assert.Nil(t, err)
	expect := "0x3e719df00000000000000000000000000000000000000000000000000000000000000002"
	assert.Equal(t, expect, data)
}

func TestModifyChainConfigurationContract_EnableNoTxDelayedMinerBlock(t *testing.T) {
	contract := NewModifyChainConfigurationContract()
	data, err := contract.EnableNoTxDelayedMining(true)
	assert.Nil(t, err)
	expect := "0x4864b6ef0000000000000000000000000000000000000000000000000000000000000001"
	assert.Equal(t, expect, data)
}

func TestModifyChainConfigurationContract_UpdateNoTxDelayedMinerBlockPeriodMultiple(t *testing.T) {
	contract := NewModifyChainConfigurationContract()
	data, err := contract.UpdateNoTxDelayedMiningPeriodMultiple(5000)
	assert.Nil(t, err)
	expect := "0x9b6042c90000000000000000000000000000000000000000000000000000000000001388"
	assert.Equal(t, expect, data)
}

func TestModifyChainConfigurationContract_UpdateProposalExpirationDays(t *testing.T) {
	contract := NewModifyChainConfigurationContract()
	data, err := contract.UpdateProposalExpirationDays(14)
	assert.Nil(t, err)
	expect := "0xe56d83dd000000000000000000000000000000000000000000000000000000000000000e"
	assert.Equal(t, expect, data)
}

func TestModifyChainConfigurationContract_UpdateChainByChainVotingRule(t *testing.T) {
	contract := NewModifyChainConfigurationContract()
	data, err := contract.UpdateChainByChainVotingRule(types.VotingRuleCONSENSUS)
	assert.Nil(t, err)
	expect := "0x975f23250000000000000000000000000000000000000000000000000000000000000002"
	assert.Equal(t, expect, data)
}

func TestModifyChainConfigurationContract_UpdatePeriod(t *testing.T) {
	contract := NewModifyChainConfigurationContract()
	data, err := contract.UpdatePeriod(50_000)
	assert.Nil(t, err)
	expect := "0xa4bf1361000000000000000000000000000000000000000000000000000000000000c350"
	assert.Equal(t, expect, data)
}

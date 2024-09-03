package builtin

import (
	"github.com/wylu1037/lattice-go/abi"
	"github.com/wylu1037/lattice-go/common/types"
)

func NewModifyChainConfigurationContract() ModifyChainConfigurationContract {
	return &modifyChainConfigurationContract{
		abi: abi.NewAbi(ModifyChainConfigurationContractBuiltinContract.AbiString),
	}
}

type ModifyChainConfigurationContract interface {
	// ContractAddress 获取修改链配置的合约地址
	//
	// Returns:
	//   - string: 合约地址，zltc_ZwuhH4dudz2Md2h6NFgHc8yrFUhKy2UUZ
	ContractAddress() string
	// UpdatePeriod 更新出块时间
	UpdatePeriod(newPeriod uint32) (string, error)
	// AddConsensusNodes 添加共识节点
	AddConsensusNodes(nodes []string) (string, error)
	// DeleteConsensusNodes 删除共识节点
	DeleteConsensusNodes(nodes []string) (string, error)
	// ReplaceConsensusNodes 替换共识节点
	ReplaceConsensusNodes(oldNode, newNode string) (string, error)
	// EnableContractLifecycleVotingDictatorship 是否开启合约生命周期投票的盟主独裁机制，否则为共识投票
	EnableContractLifecycleVotingDictatorship(enable bool) (string, error)
	// UpdateConsensus 更新链的共识机制
	//
	// Parameters:
	//   - consensus types.Consensus: 仅支持切换为 types.ConsensusPOA 或 types.ConsensusPBFT
	UpdateConsensus(consensus types.Consensus) (string, error)
	// EnableContractLifecycle 是否开启合约生命周期
	EnableContractLifecycle(enable bool) (string, error)
	// EnableContractManagement 是否启用合约内部管理
	EnableContractManagement(enable bool) (string, error)
	// EnableNoTxDelayedMining 无交易时延迟出块
	EnableNoTxDelayedMining(enable bool) (string, error)
	// UpdateNoTxDelayedMiningPeriodMultiple 更新无交易时延迟出块的倍数
	UpdateNoTxDelayedMiningPeriodMultiple(multiple uint64) (string, error)
	// UpdateContractDeploymentVotingRule 设置合约部署的投票规则
	//
	// Parameters:
	//	 - votingRule types.VotingRule
	UpdateContractDeploymentVotingRule(votingRule types.VotingRule) (string, error)
	// UpdateProposalExpirationDays 更新提案的过期天数
	UpdateProposalExpirationDays(expirationDays uint64) (string, error)
	// UpdateChainByChainVotingRule 更新以链建链的投票规则
	UpdateChainByChainVotingRule(votingRule types.VotingRule) (string, error)
}

type modifyChainConfigurationContract struct {
	abi abi.LatticeAbi
}

func (c *modifyChainConfigurationContract) ContractAddress() string {
	return ModifyChainConfigurationContractBuiltinContract.Address
}

func (c *modifyChainConfigurationContract) UpdatePeriod(newPeriod uint32) (string, error) {
	fn, err := c.abi.GetLatticeFunction("changePeriod", newPeriod)
	if err != nil {
		return "", err
	}
	return fn.Encode()
}

func (c *modifyChainConfigurationContract) AddConsensusNodes(nodes []string) (string, error) {
	fn, err := c.abi.GetLatticeFunction("addLatcSaint", nodes)
	if err != nil {
		return "", err
	}
	return fn.Encode()
}

func (c *modifyChainConfigurationContract) DeleteConsensusNodes(nodes []string) (string, error) {
	fn, err := c.abi.GetLatticeFunction("delLatcSaint", nodes)
	if err != nil {
		return "", err
	}
	return fn.Encode()
}

func (c *modifyChainConfigurationContract) ReplaceConsensusNodes(oldNode, newNode string) (string, error) {
	fn, err := c.abi.GetLatticeFunction("replaceLatcSaint", oldNode, newNode)
	if err != nil {
		return "", err
	}
	return fn.Encode()
}

func (c *modifyChainConfigurationContract) EnableContractLifecycleVotingDictatorship(enable bool) (string, error) {
	fn, err := c.abi.GetLatticeFunction("isDictatorship", enable)
	if err != nil {
		return "", err
	}
	return fn.Encode()
}

func (c *modifyChainConfigurationContract) UpdateConsensus(consensus types.Consensus) (string, error) {
	fn, err := c.abi.GetLatticeFunction("switchConsensus", string(consensus))
	if err != nil {
		return "", err
	}
	return fn.Encode()
}

func (c *modifyChainConfigurationContract) EnableContractLifecycle(enable bool) (string, error) {
	fn, err := c.abi.GetLatticeFunction("switchIsContractVote", enable)
	if err != nil {
		return "", err
	}
	return fn.Encode()
}

func (c *modifyChainConfigurationContract) EnableContractManagement(enable bool) (string, error) {
	fn, err := c.abi.GetLatticeFunction("switchContractPermission", enable)
	if err != nil {
		return "", err
	}
	return fn.Encode()
}

func (c *modifyChainConfigurationContract) EnableNoTxDelayedMining(enable bool) (string, error) {
	fn, err := c.abi.GetLatticeFunction("switchNoEmptyAnchor", enable)
	if err != nil {
		return "", err
	}
	return fn.Encode()
}

func (c *modifyChainConfigurationContract) UpdateNoTxDelayedMiningPeriodMultiple(multiple uint64) (string, error) {
	fn, err := c.abi.GetLatticeFunction("changeEmptyAnchorPeriodMul", multiple)
	if err != nil {
		return "", err
	}
	return fn.Encode()
}

func (c *modifyChainConfigurationContract) UpdateContractDeploymentVotingRule(votingRule types.VotingRule) (string, error) {
	fn, err := c.abi.GetLatticeFunction("switchDeployRule", uint8(votingRule))
	if err != nil {
		return "", err
	}
	return fn.Encode()
}

func (c *modifyChainConfigurationContract) UpdateProposalExpirationDays(expirationDays uint64) (string, error) {
	fn, err := c.abi.GetLatticeFunction("changeProposalExpireTime", expirationDays)
	if err != nil {
		return "", err
	}
	return fn.Encode()
}

func (c *modifyChainConfigurationContract) UpdateChainByChainVotingRule(votingRule types.VotingRule) (string, error) {
	fn, err := c.abi.GetLatticeFunction("changeChainByChainVote", uint8(votingRule))
	if err != nil {
		return "", err
	}
	return fn.Encode()
}

var ModifyChainConfigurationContractBuiltinContract = Contract{
	Description: "修改链配置合约",
	Address:     "zltc_ZwuhH4dudz2Md2h6NFgHc8yrFUhKy2UUZ",
	AbiString: `[
		{
			"inputs": [
				{
					"internalType": "address[]",
					"name": "LatcSaint",
					"type": "address[]"
				}
			],
			"name": "addLatcSaint",
			"outputs": [
				{
					"internalType": "bytes",
					"name": "",
					"type": "bytes"
				}
			],
			"stateMutability": "pure",
			"type": "function"
		},
		{
			"inputs": [
				{
					"internalType": "uint256",
					"name": "Period",
					"type": "uint256"
				}
			],
			"name": "changePeriod",
			"outputs": [
				{
					"internalType": "bytes",
					"name": "",
					"type": "bytes"
				}
			],
			"stateMutability": "pure",
			"type": "function"
		},
		{
			"inputs": [
				{
					"internalType": "address[]",
					"name": "LatcSaint",
					"type": "address[]"
				}
			],
			"name": "delLatcSaint",
			"outputs": [
				{
					"internalType": "bytes",
					"name": "",
					"type": "bytes"
				}
			],
			"stateMutability": "pure",
			"type": "function"
		},
		{
			"inputs": [
				{
					"internalType": "address",
					"name": "oldSaint",
					"type": "address"
				},
				{
					"internalType": "address",
					"name": "newSaint",
					"type": "address"
				}
			],
			"name": "replaceLatcSaint",
			"outputs": [
				{
					"internalType": "bytes",
					"name": "",
					"type": "bytes"
				}
			],
			"stateMutability": "pure",
			"type": "function"
		},
		{
			"inputs": [
				{
					"internalType": "bool",
					"name": "IsDictatorship",
					"type": "bool"
				}
			],
			"name": "isDictatorship",
			"outputs": [
				{
					"internalType": "bytes",
					"name": "",
					"type": "bytes"
				}
			],
			"stateMutability": "pure",
			"type": "function"
		},
		{
			"inputs": [
				{
					"internalType": "bool",
					"name": "isContractVote",
					"type": "bool"
				}
			],
			"name": "switchIsContractVote",
			"outputs": [
				{
					"internalType": "bytes",
					"name": "",
					"type": "bytes"
				}
			],
			"stateMutability": "pure",
			"type": "function"
		},
		{
			"inputs": [
				{
					"internalType": "bool",
					"name": "contractPermission",
					"type": "bool"
				}
			],
			"name": "switchContractPermission",
			"outputs": [
				{
					"internalType": "bytes",
					"name": "",
					"type": "bytes"
				}
			],
			"stateMutability": "pure",
			"type": "function"
		},
		{
			"inputs": [
				{
					"internalType": "string",
					"name": "Consensus",
					"type": "string"
				}
			],
			"name": "switchConsensus",
			"outputs": [
				{
					"internalType": "bytes",
					"name": "",
					"type": "bytes"
				}
			],
			"stateMutability": "pure",
			"type": "function"
		},
		{
			"inputs": [
				{
					"internalType": "uint8",
					"name": "deployRule",
					"type": "uint8"
				}
			],
			"name": "switchDeployRule",
			"outputs": [
				{
					"internalType": "bytes",
					"name": "",
					"type": "bytes"
				}
			],
			"stateMutability": "pure",
			"type": "function"
		},
		{
			"inputs": [
				{
					"internalType": "bool",
					"name": "noEmptyAnchor",
					"type": "bool"
				}
			],
			"name": "switchNoEmptyAnchor",
			"outputs": [
				{
					"internalType": "bytes",
					"name": "",
					"type": "bytes"
				}
			],
			"stateMutability": "pure",
			"type": "function"
		},
		{
			"inputs": [
				{
					"internalType": "address",
					"name": "preacher",
					"type": "address"
				}
			],
			"name": "changePreacher",
			"outputs": [
				{
					"internalType": "bytes",
					"name": "",
					"type": "bytes"
				}
			],
			"stateMutability": "pure",
			"type": "function"
		},
		{
			"inputs": [
				{
					"internalType": "uint64",
					"name": "emptyAnchorPeriodMul",
					"type": "uint64"
				}
			],
			"name": "changeEmptyAnchorPeriodMul",
			"outputs": [],
			"stateMutability": "pure",
			"type": "function"
		},
		{
			"inputs": [
				{
					"internalType": "uint64",
					"name": "proposalExpireTime",
					"type": "uint64"
				}
			],
			"name": "changeProposalExpireTime",
			"outputs": [],
			"stateMutability": "pure",
			"type": "function"
		},
		{
			"inputs": [
				{
					"internalType": "uint8",
					"name": "chainByChainVote",
					"type": "uint8"
				}
			],
			"name": "changeChainByChainVote",
			"outputs": [],
			"stateMutability": "pure",
			"type": "function"
		}
	]`,
}

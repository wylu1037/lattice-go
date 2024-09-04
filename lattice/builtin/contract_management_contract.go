package builtin

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/wylu1037/lattice-go/abi"
	"github.com/wylu1037/lattice-go/common/types"
	"strings"
)

type ContractManagementAction string

const (
	ContractManagementActionCREATE ContractManagementAction = "C"
	ContractManagementActionUPDATE ContractManagementAction = "U"
	ContractManagementActionDELETE ContractManagementAction = "D"
)

// SetContractManagementRulesRequest 设置合约管理规则的请求
//   - Address 				   合约地址
//   - ContractManagementRules 合约管理规则
type SetContractManagementRulesRequest struct {
	Address                 common.Address
	ContractManagementRules ContractManagementRules
}

// ContractManagementRules 合约管理规则
//   - PermissionMode 管理模式，types.ContractManagementModeWHITELIST or types.ContractManagementModeBLACKLIST
//   - Threshold 	  阈值
//   - BlackList 	  黑名单
//   - WhiteList 	  白名单
//   - ManagerList 	  权重分配
type ContractManagementRules struct {
	PermissionMode types.ContractManagementMode
	Threshold      uint64
	BlackList      []common.Address
	WhiteList      []common.Address
	ManagerList    []WeightDistribution
}

// WeightDistribution 权重分配
//   - Address 账户地址
//   - Weight  权重，000-255
type WeightDistribution struct {
	Address common.Address
	Weight  uint8
}

func NewContractManagementContract() ContractManagementContract {
	return &contractManagementContract{
		abi: abi.NewAbi(ContractManagementBuiltinContract.AbiString),
	}
}

type ContractManagementContract interface {
	// ContractAddress 获取修改链配置的合约地址
	//
	// Returns:
	//   - string: 合约地址，zltc_ZDdPo8P72X7dtMNTxBeKU8pT7bDXb7NtV
	ContractAddress() string
	// 发起更新合约管理规则
	//
	// Parameters:
	//   - contractAddress string: 合约地址
	//   - operation string: 操作
	//
	// Returns:
	//   - (string, error)
	launch(contractAddress, operation string) (string, error)
	// SetManagementRules 设置合约的管理规则
	SetManagementRules(req *SetContractManagementRulesRequest) (string, error)
	// UpdateVotingThreshold 更新投票阈值
	UpdateVotingThreshold(contractAddress string, threshold uint32) (string, error)
	// UpdateManagementMode 更新合约管理模式
	UpdateManagementMode(contractAddress string, mode types.ContractManagementMode) (string, error)
	// UpdateWhitelist 更新白名单
	UpdateWhitelist(contractAddress string, action ContractManagementAction, addresses []string) (string, error)
	// UpdateBlacklist 更新黑名单
	UpdateBlacklist(contractAddress string, action ContractManagementAction, addresses []string) (string, error)
	// UpdateWeight 更新账户权重
	UpdateWeight(contractAddress string, action ContractManagementAction, weights []WeightDistribution) (string, error)
}

type contractManagementContract struct {
	abi abi.LatticeAbi
}

func (c *contractManagementContract) ContractAddress() string {
	return ContractManagementBuiltinContract.Address
}

func (c *contractManagementContract) launch(contractAddress, operation string) (string, error) {
	fn, err := c.abi.GetLatticeFunction("launch", contractAddress, operation)
	if err != nil {
		return "", err
	}
	return fn.Encode()
}

func (c *contractManagementContract) SetManagementRules(req *SetContractManagementRulesRequest) (string, error) {
	code, err := c.abi.RawAbi().Pack("init", req.Address, req.ContractManagementRules)
	if err != nil {
		return "", err
	}
	return hexutil.Encode(code), nil
}

func (c *contractManagementContract) UpdateVotingThreshold(contractAddress string, threshold uint32) (string, error) {
	return c.launch(contractAddress, fmt.Sprintf("UT%d", threshold))
}

func (c *contractManagementContract) UpdateManagementMode(contractAddress string, mode types.ContractManagementMode) (string, error) {
	return c.launch(contractAddress, fmt.Sprintf("UP%d", mode))
}

func (c *contractManagementContract) UpdateWhitelist(contractAddress string, action ContractManagementAction, addresses []string) (string, error) {
	return c.launch(contractAddress, fmt.Sprintf("%sW%s", action, strings.Join(addresses, "")))
}

func (c *contractManagementContract) UpdateBlacklist(contractAddress string, action ContractManagementAction, addresses []string) (string, error) {
	return c.launch(contractAddress, fmt.Sprintf("%sB%s", action, strings.Join(addresses, "")))
}

func (c *contractManagementContract) UpdateWeight(contractAddress string, action ContractManagementAction, weights []WeightDistribution) (string, error) {
	var builder strings.Builder
	for _, elem := range weights {
		builder.WriteString(elem.Address.String())
		if action != ContractManagementActionDELETE {
			builder.WriteString(fmt.Sprintf("%03d", elem.Weight))
		}
	}
	return c.launch(contractAddress, fmt.Sprintf("%sM%s", action, builder.String()))
}

var ContractManagementBuiltinContract = Contract{
	Description: "合约内部管理合约",
	Address:     "zltc_ZDdPo8P72X7dtMNTxBeKU8pT7bDXb7NtV",
	AbiString: `[
		{
			"inputs": [
				{
					"internalType": "address",
					"name": "Address",
					"type": "address"
				},
				{
					"internalType": "string",
					"name": "operation",
					"type": "string"
				}
			],
			"name": "launch",
			"outputs": [
				{
					"internalType": "string",
					"name": "",
					"type": "string"
				}
			],
			"stateMutability": "pure",
			"type": "function"
		},
		{
			"inputs": [
				{
					"internalType": "address",
					"name": "Address",
					"type": "address"
				},
				{
					"components": [
						{
							"internalType": "uint8",
							"name": "permissionMode",
							"type": "uint8"
						},
						{
							"internalType": "uint64",
							"name": "threshold",
							"type": "uint64"
						},
						{
							"internalType": "address[]",
							"name": "blackList",
							"type": "address[]"
						},
						{
							"internalType": "address[]",
							"name": "whiteList",
							"type": "address[]"
						},
						{
							"components": [
								{
									"internalType": "address",
									"name": "Address",
									"type": "address"
								},
								{
									"internalType": "uint8",
									"name": "weight",
									"type": "uint8"
								}
							],
							"internalType": "struct chainbychain.Manager[]",
							"name": "managerList",
							"type": "tuple[]"
						}
					],
					"internalType": "struct chainbychain.Args",
					"name": "permissionList",
					"type": "tuple"
				}
			],
			"name": "init",
			"outputs": [],
			"stateMutability": "nonpayable",
			"type": "function"
		}
	]`,
}

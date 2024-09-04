package builtin

import (
	"github.com/wylu1037/lattice-go/abi"
)

func NewProposalContract() ProposalContract {
	return &proposalContract{
		abi: abi.NewAbi(ProposalBuiltinContract.AbiString),
	}
}

type ProposalContract interface {
	// ContractAddress 获取修改链配置的合约地址
	//
	// Returns:
	//   - string: 合约地址，zltc_amgWuhifLRUoZc3GSbv9wUUz6YUfTuWy5
	ContractAddress() string

	// Approve 投赞同票
	//
	// Parameters:
	//   - proposalId string: 提案ID
	//
	// Returns:
	//   - string
	//   - error
	Approve(proposalId string) (string, error)

	// Disapprove 投反对票
	//
	// Parameters:
	//   - proposalId string: 提案ID
	//
	// Returns:
	//   - string
	//   - error
	Disapprove(proposalId string) (string, error)
	// Refresh 刷新提案
	Refresh(proposalId string) (string, error)
	// BatchRefresh 批量刷新提案
	BatchRefresh(proposalIds []string) (string, error)
	// Cancel 取消提案
	Cancel(proposalId string) (string, error)
}

type proposalContract struct {
	abi abi.LatticeAbi
}

const (
	disapprove = iota // 反对
	approve           // 同意
)

func (c *proposalContract) vote(proposalId string, approve bool) (string, error) {
	if approve {
		return c.Approve(proposalId)
	} else {
		return c.Disapprove(proposalId)
	}
}

func (c *proposalContract) ContractAddress() string {
	return ProposalBuiltinContract.Address
}

func (c *proposalContract) Approve(proposalId string) (string, error) {
	fn, err := c.abi.GetLatticeFunction("vote", proposalId, approve)
	if err != nil {
		return "", err
	}

	return fn.Encode()
}

func (c *proposalContract) Disapprove(proposalId string) (string, error) {
	fn, err := c.abi.GetLatticeFunction("vote", proposalId, disapprove)
	if err != nil {
		return "", err
	}

	return fn.Encode()
}

func (c *proposalContract) Refresh(proposalId string) (string, error) {
	fn, err := c.abi.GetLatticeFunction("refresh", proposalId)
	if err != nil {
		return "", err
	}

	return fn.Encode()
}

func (c *proposalContract) BatchRefresh(proposalIds []string) (string, error) {
	fn, err := c.abi.GetLatticeFunction("batchRefresh", proposalIds)
	if err != nil {
		return "", err
	}

	return fn.Encode()
}

func (c *proposalContract) Cancel(proposalId string) (string, error) {
	fn, err := c.abi.GetLatticeFunction("cancel", proposalId)
	if err != nil {
		return "", err
	}

	return fn.Encode()
}

var ProposalBuiltinContract = Contract{
	Description: "提案投票合约",
	Address:     "zltc_amgWuhifLRUoZc3GSbv9wUUz6YUfTuWy5",
	AbiString: `[
		{
			"inputs": [
				{
					"internalType": "string",
					"name": "ProposalId",
					"type": "string"
				},
				{
					"internalType": "uint8",
					"name": "VoteSuggestion",
					"type": "uint8"
				}
			],
			"name": "vote",
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
					"name": "ProposalId",
					"type": "string"
				}
			],
			"name": "refresh",
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
					"internalType": "string[]",
					"name": "proposalIds",
					"type": "string[]"
				}
			],
			"name": "batchRefresh",
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
					"name": "proposalId",
					"type": "string"
				}
			],
			"name": "cancel",
			"outputs": [
				{
					"internalType": "bytes",
					"name": "",
					"type": "bytes"
				}
			],
			"stateMutability": "pure",
			"type": "function"
		}
	]`,
}

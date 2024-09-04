package builtin

import (
	"github.com/wylu1037/lattice-go/abi"
)

func NewVoteContract() VoteContract {
	return &voteContract{
		abi: abi.NewAbi(VoteBuiltinContract.AbiString),
	}
}

type VoteContract interface {
	vote(proposalId string, approve bool) (string, error)

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

type voteContract struct {
	abi abi.LatticeAbi
}

const (
	disapprove = iota // 反对
	approve           // 同意
)

func (c *voteContract) vote(proposalId string, approve bool) (string, error) {
	if approve {
		return c.Approve(proposalId)
	} else {
		return c.Disapprove(proposalId)
	}
}

func (c *voteContract) Approve(proposalId string) (string, error) {
	fn, err := c.abi.GetLatticeFunction("vote", proposalId, approve)
	if err != nil {
		return "", err
	}

	return fn.Encode()
}

func (c *voteContract) Disapprove(proposalId string) (string, error) {
	fn, err := c.abi.GetLatticeFunction("vote", proposalId, disapprove)
	if err != nil {
		return "", err
	}

	return fn.Encode()
}

func (c *voteContract) Refresh(proposalId string) (string, error) {
	fn, err := c.abi.GetLatticeFunction("refresh", proposalId)
	if err != nil {
		return "", err
	}

	return fn.Encode()
}

func (c *voteContract) BatchRefresh(proposalIds []string) (string, error) {
	fn, err := c.abi.GetLatticeFunction("batchRefresh", proposalIds)
	if err != nil {
		return "", err
	}

	return fn.Encode()
}

func (c *voteContract) Cancel(proposalId string) (string, error) {
	fn, err := c.abi.GetLatticeFunction("cancel", proposalId)
	if err != nil {
		return "", err
	}

	return fn.Encode()
}

var VoteBuiltinContract = Contract{
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

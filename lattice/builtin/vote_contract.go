package builtin

import (
	"github.com/wylu1037/lattice-go/abi"
)

func NewVoteContract() VoteContract {
	return &voteContract{}
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
}

type voteContract struct{}

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
	fn, err := abi.NewAbi(VoteBuiltinContract.AbiString).GetLatticeFunction("vote", proposalId, approve)
	if err != nil {
		return "", err
	}

	data, err := fn.Encode()
	if err != nil {
		return "", err
	}

	return data, nil
}

func (c *voteContract) Disapprove(proposalId string) (string, error) {
	fn, err := abi.NewAbi(VoteBuiltinContract.AbiString).GetLatticeFunction("vote", proposalId, disapprove)
	if err != nil {
		return "", err
	}

	data, err := fn.Encode()
	if err != nil {
		return "", err
	}

	return data, nil
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

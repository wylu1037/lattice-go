package client

import (
	"context"
	"github.com/wylu1037/lattice-go/common/types"
)

func (api *httpApi) GetContractLifecycleProposal(_ context.Context, chainId, contractAddress string, state types.ProposalState) ([]types.Proposal[types.ContractLifecycleProposal], error) {
	params := map[string]interface{}{
		"proposalType":    types.ProposalTypeContractLifecycle,
		"proposalState":   state,
		"proposalAddress": contractAddress,
	}

	response, err := Post[[]types.Proposal[types.ContractLifecycleProposal]](api.Url, NewJsonRpcBody("wallet_getProposal", params), api.newHeaders(chainId), api.transport)
	if err != nil {
		return nil, err
	}
	if response.Error != nil {
		return nil, response.Error.Error()
	}
	return *response.Result, nil
}

func (api *httpApi) GetVoteById(_ context.Context, chainId, voteId string) (*types.VoteDetails, error) {
	response, err := Post[types.VoteDetails](api.Url, NewJsonRpcBody("wallet_getVoteById", voteId), api.newHeaders(chainId), api.transport)
	if err != nil {
		return nil, err
	}
	if response.Error != nil {
		return nil, response.Error.Error()
	}
	return response.Result, nil
}

// GetProposal get proposal
// Parameters:
//   - ctx: context.Context
//   - chainId(string): query by chain id
//   - proposalId(string): query by ProposalId, can be empty string
//   - proposalType(ProposalType): query by ProposalType, zero represent return all
//   - proposalState(ProposalState): query by ProposalState, zero represent return all
//   - proposalAddress(string): query by proposal address, can be empty string
//   - contractAddress(string): query by contract address, can be empty string
//   - startTime(string): 20240830
//   - endTime(string): 20240830
//   - result([]fusion.Proposal[ContractLifecycleProposalContent|ModifyChainConfigProposalContent]): result is slice
//
// Returns:
//   - error
func (api *httpApi) GetProposal(_ context.Context, chainId, proposalId string, ty types.ProposalType, state types.ProposalState, proposalAddress, contractAddress, startDate, endDate string, result interface{}) error {
	args := map[string]interface{}{"proposalType": ty, "proposalState": state}
	if len(proposalId) != 0 {
		args["proposalId"] = proposalId
	}
	if len(proposalAddress) != 0 {
		args["proposalAddress"] = proposalAddress
	}
	if len(contractAddress) != 0 {
		args["contractAddress"] = contractAddress
	}
	if len(startDate) != 0 {
		args["dateStart"] = startDate
	}
	if len(endDate) != 0 {
		args["dateEnd"] = endDate
	}

	response, err := Post[interface{}](api.Url, NewJsonRpcBody("wallet_getProposal", args), api.newHeaders(chainId), api.transport)
	if err != nil {
		return err
	}
	if response.Error != nil {
		return response.Error.Error()
	}
	result = response.Result
	return nil
}

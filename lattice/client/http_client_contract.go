package client

import (
	"context"
	"github.com/wylu1037/lattice-go/common/types"
)

func (api *httpApi) GetContractInformation(_ context.Context, chainID, contractAddress string) (*types.ContractInformation, error) {
	response, err := Post[types.ContractInformation](api.Url, NewJsonRpcBody("wallet_getContractState", contractAddress), api.newHeaders(chainID), api.transport)
	if err != nil {
		return nil, err
	}
	if response.Error != nil {
		return nil, response.Error.Error()
	}
	return response.Result, nil
}

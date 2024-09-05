package client

import (
	"context"
	"github.com/wylu1037/lattice-go/common/types"
	"math/big"
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

func (api *httpApi) GetContractManagement(_ context.Context, chainID, contractAddress string, daemonBlockHeight *big.Int) (*types.ContractManagement, error) {
	var err error
	var response *JsonRpcResponse[*types.ContractManagement]
	if daemonBlockHeight == nil {
		response, err = Post[types.ContractManagement](api.Url, NewJsonRpcBody("wallet_getPermissionList", contractAddress), api.newHeaders(chainID), api.transport)
	} else {
		response, err = Post[types.ContractManagement](api.Url, NewJsonRpcBody("wallet_getPermissionList", contractAddress, daemonBlockHeight), api.newHeaders(chainID), api.transport)
	}
	if err != nil {
		return nil, err
	}
	if response.Error != nil {
		return nil, response.Error.Error()
	}
	return response.Result, nil
}

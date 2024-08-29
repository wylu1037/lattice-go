package client

import (
	"context"
	"github.com/wylu1037/lattice-go/common/types"
)

func (api *httpApi) GetSubchain(_ context.Context, subchainId string) (*types.Subchain, error) {
	response, err := Post[types.Subchain](api.Url, NewJsonRpcBody("latc_latcInfo"), api.newHeaders(subchainId), api.transport)
	if err != nil {
		return nil, err
	}
	if response.Error != nil {
		return nil, response.Error.Error()
	}
	return response.Result, nil
}

func (api *httpApi) GetCreatedSubchain(_ context.Context) ([]uint64, error) {
	response, err := Post[[]uint64](api.Url, NewJsonRpcBody("cbyc_getCreatedAllChains"), api.newHeaders(emptyChainId), api.transport)
	if err != nil {
		return nil, err
	}
	if response.Error != nil {
		return nil, response.Error.Error()
	}
	return *response.Result, nil
}

func (api *httpApi) GetJoinedSubchain(_ context.Context) ([]uint64, error) {
	response, err := Post[[]uint64](api.Url, NewJsonRpcBody("node_getAllChainId"), api.newHeaders(emptyChainId), api.transport)
	if err != nil {
		return nil, err
	}
	if response.Error != nil {
		return nil, response.Error.Error()
	}
	return *response.Result, nil
}

func (api *httpApi) GetSubchainRunningStatus(_ context.Context, subchainID string) (*types.SubchainRunningStatus, error) {
	response, err := Post[string](api.Url, NewJsonRpcBody("cbyc_getChainStatus"), api.newHeaders(subchainID), api.transport)
	if err != nil {
		return nil, err
	}
	if response.Error != nil {
		return nil, response.Error.Error()
	}
	status := *response.Result
	return &types.SubchainRunningStatus{
		Status:  status,
		Running: status == types.SubchainStatusRUNNING,
	}, nil
}

func (api *httpApi) JoinSubchain(_ context.Context, subchainId, networkId uint64, inode string) error {
	response, err := Post[any](api.Url, NewJsonRpcBody("cbyc_selfJoinChain", subchainId, networkId, inode), api.newHeaders(emptyChainId), api.transport)
	if err != nil {
		return err
	}
	if response.Error != nil {
		return response.Error.Error()
	}
	return nil
}

func (api *httpApi) StartSubchain(_ context.Context, subchainId string) error {
	response, err := Post[any](api.Url, NewJsonRpcBody("cbyc_startSelfChain"), api.newHeaders(subchainId), api.transport)
	if err != nil {
		return err
	}
	if response.Error != nil {
		return response.Error.Error()
	}
	return nil
}

func (api *httpApi) StopSubchain(_ context.Context, subchainId string) error {
	response, err := Post[any](api.Url, NewJsonRpcBody("cbyc_stopSelfChain"), api.newHeaders(subchainId), api.transport)
	if err != nil {
		return err
	}
	if response.Error != nil {
		return response.Error.Error()
	}
	return nil
}

func (api *httpApi) DeleteSubchain(_ context.Context, subchainId string) error {
	response, err := Post[any](api.Url, NewJsonRpcBody("cbyc_delSelfChain"), api.newHeaders(subchainId), api.transport)
	if err != nil {
		return err
	}
	if response.Error != nil {
		return response.Error.Error()
	}
	return nil
}

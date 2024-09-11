package client

import (
	"context"
	"github.com/wylu1037/lattice-go/common/types"
)

func (api *httpApi) GetLatestBlock(ctx context.Context, chainId, accountAddress string) (*types.LatestBlock, error) {
	response, err := Post[types.LatestBlock](ctx, api.Url, NewJsonRpcBody("latc_getCurrentTBDB", accountAddress), api.newHeaders(chainId), api.transport)
	if err != nil {
		return nil, err
	}
	if response.Error != nil {
		return nil, response.Error.Error()
	}
	return response.Result, nil
}

func (api *httpApi) GetLatestBlockWithPending(ctx context.Context, chainId, accountAddress string) (*types.LatestBlock, error) {
	response, err := Post[types.LatestBlock](ctx, api.Url, NewJsonRpcBody("latc_getPendingTBDB", accountAddress), api.newHeaders(chainId), api.transport)
	if err != nil {
		return nil, err
	}
	if response.Error != nil {
		return nil, response.Error.Error()
	}
	return response.Result, nil
}

func (api *httpApi) GetLatestDaemonBlock(ctx context.Context, chainID string) (*types.DaemonBlock, error) {
	response, err := Post[types.DaemonBlock](ctx, api.Url, NewJsonRpcBody("latc_getCurrentDBlock"), api.newHeaders(chainID), api.transport)
	if err != nil {
		return nil, err
	}
	if response.Error != nil {
		return nil, response.Error.Error()
	}
	return response.Result, nil
}

func (api *httpApi) GetDaemonBlockByHash(ctx context.Context, chainId, hash string) (*types.DaemonBlock, error) {
	response, err := Post[types.DaemonBlock](ctx, api.Url, NewJsonRpcBody("latc_getDBlockByHash", hash), api.newHeaders(chainId), api.transport)
	if err != nil {
		return nil, err
	}
	if response.Error != nil {
		return nil, response.Error.Error()
	}
	return response.Result, nil
}

func (api *httpApi) GetTransactionBlockByHash(ctx context.Context, chainId, hash string) (*types.TransactionBlock, error) {
	response, err := Post[types.TransactionBlock](ctx, api.Url, NewJsonRpcBody("latc_getTBlockByHash", hash), api.newHeaders(chainId), api.transport)
	if err != nil {
		return nil, err
	}
	if response.Error != nil {
		return nil, response.Error.Error()
	}
	return response.Result, nil
}

func (api *httpApi) GetTransactionsPagination(ctx context.Context, chainId string, startDaemonBlockHeight uint64, pageSize uint16) (*types.TransactionsPagination, error) {
	response, err := Post[types.TransactionsPagination](ctx, api.Url, NewJsonRpcBody("latc_getTBlockPagesByDNumber", startDaemonBlockHeight, pageSize), api.newHeaders(chainId), api.transport)
	if err != nil {
		return nil, err
	}
	if response.Error != nil {
		return nil, response.Error.Error()
	}
	return response.Result, nil
}

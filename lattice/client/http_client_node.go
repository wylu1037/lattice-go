package client

import (
	"context"
	"github.com/wylu1037/lattice-go/common/types"
	"github.com/wylu1037/lattice-go/wallet"
)

func (api *httpApi) GetNodeInfo(_ context.Context) (*types.NodeInfo, error) {
	response, err := Post[types.NodeInfo](api.Url, NewJsonRpcBody("node_nodeInfo"), api.newHeaders(emptyChainId), api.transport)
	if err != nil {
		return nil, err
	}
	if response.Error != nil {
		return nil, response.Error.Error()
	}
	return response.Result, nil
}

func (api *httpApi) GetConsensusNodesStatus(_ context.Context, chainID string) ([]*types.ConsensusNodeStatus, error) {
	response, err := Post[[]*types.ConsensusNodeStatus](api.Url, NewJsonRpcBody("witness_nodeList"), api.newHeaders(chainID), api.transport)
	if err != nil {
		return nil, err
	}
	if response.Error != nil {
		return nil, response.Error.Error()
	}
	return *response.Result, nil
}

func (api *httpApi) GetGenesisNodeAddress(_ context.Context, chainID string) (string, error) {
	response, err := Post[string](api.Url, NewJsonRpcBody("wallet_getGenesisNode"), api.newHeaders(chainID), api.transport)
	if err != nil {
		return "", err
	}
	if response.Error != nil {
		return "", response.Error.Error()
	}
	return *response.Result, nil
}

func (api *httpApi) GetNodePeers(_ context.Context, chainID string) ([]*types.NodePeer, error) {
	response, err := Post[[]*types.NodePeer](api.Url, NewJsonRpcBody("node_peers"), api.newHeaders(chainID), api.transport)
	if err != nil {
		return nil, err
	}
	if response.Error != nil {
		return nil, response.Error.Error()
	}
	return *response.Result, nil
}

func (api *httpApi) GetNodeConfig(_ context.Context, chainID string) (*types.NodeConfig, error) {
	response, err := Post[types.NodeConfig](api.Url, NewJsonRpcBody("latc_getConfig"), api.newHeaders(chainID), api.transport)
	if err != nil {
		return nil, err
	}
	if response.Error != nil {
		return nil, response.Error.Error()
	}
	return response.Result, nil
}

func (api *httpApi) GetNodeProtocol(_ context.Context, chainId string) (*types.NodeProtocol, error) {
	response, err := Post[types.NodeProtocol](api.Url, NewJsonRpcBody("latc_getProtocols"), api.newHeaders(chainId), api.transport)
	if err != nil {
		return nil, err
	}
	if response.Error != nil {
		return nil, response.Error.Error()
	}
	return response.Result, nil
}

func (api *httpApi) GetNodeConfirmedConfiguration(_ context.Context, chainId string) (*types.NodeConfirmedConfiguration, error) {
	response, err := Post[types.NodeConfirmedConfiguration](api.Url, NewJsonRpcBody("wallet_getConfirmConfig"), api.newHeaders(chainId), api.transport)
	if err != nil {
		return nil, err
	}
	if response.Error != nil {
		return nil, response.Error.Error()
	}
	return response.Result, nil
}

func (api *httpApi) GetNodeVersion(_ context.Context) (*types.NodeVersion, error) {
	response, err := Post[types.NodeVersion](api.Url, NewJsonRpcBody("node_nodeVersion"), api.newHeaders(emptyChainId), api.transport)
	if err != nil {
		return nil, err
	}
	if response.Error != nil {
		return nil, response.Error.Error()
	}
	return response.Result, nil
}

func (api *httpApi) GetNodeSaintKey(_ context.Context) (*wallet.FileKey, error) {
	response, err := Post[wallet.FileKey](api.Url, NewJsonRpcBody("node_getSaintKey"), api.newHeaders(emptyChainId), api.transport)
	if err != nil {
		return nil, err
	}
	if response.Error != nil {
		return nil, response.Error.Error()
	}
	return response.Result, nil
}

func (api *httpApi) GetNodeConfiguration(_ context.Context) (*types.NodeConfiguration, error) {
	response, err := Post[types.NodeConfiguration](api.Url, NewJsonRpcBody("latc_getConfig"), api.newHeaders(emptyChainId), api.transport)
	if err != nil {
		return nil, err
	}
	if response.Error != nil {
		return nil, response.Error.Error()
	}
	return response.Result, nil
}

func (api *httpApi) GetNodeWorkingDirectory(_ context.Context) (string, error) {
	response, err := Post[string](api.Url, NewJsonRpcBody("node_getLocationPath"), api.newHeaders(emptyChainId), api.transport)
	if err != nil {
		return "", err
	}
	if response.Error != nil {
		return "", response.Error.Error()
	}
	return *response.Result, nil
}

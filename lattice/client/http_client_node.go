package client

import (
	"context"
	"github.com/wylu1037/lattice-go/common/types"
)

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

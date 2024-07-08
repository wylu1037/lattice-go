package client

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"io"
	"lattice-go/common/types"
	"lattice-go/lattice/block"
	"net/http"
	"strings"
)

// JsonRpcBody Json-Rpc的请求体结构
type JsonRpcBody struct {
	Id      int           `json:"id,omitempty"`
	JsonRpc string        `json:"jsonrpc,omitempty"`
	Method  string        `json:"method,omitempty"` // 方法名
	Params  []interface{} `json:"params,omitempty"` // 方法参数
}

// JsonRpcResponse Json-Rpc请求的响应结构
type JsonRpcResponse[T any] struct {
	Id      int           `json:"id,omitempty"`
	JsonRpc string        `json:"jsonrpc,omitempty"`
	Result  T             `json:"result,omitempty"`
	Error   *JsonRpcError `json:"error,omitempty"`
}

type JsonRpcError struct {
	Code    uint16 `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

func (e *JsonRpcError) Error() error {
	return fmt.Errorf("%d:%s", e.Code, e.Message)
}

func NewJsonRpcBody(method string, params ...interface{}) *JsonRpcBody {
	return &JsonRpcBody{
		Id:      1,
		JsonRpc: "2.0",
		Method:  method,
		Params:  params,
	}
}

// HttpApiInitParam 初始化HTTP API的参数
type HttpApiInitParam struct {
}

func NewHttpApi(url string, chainId string, transport *http.Transport) HttpApi {
	return &httpApi{
		ChainId:   chainId,
		Url:       url,
		transport: transport,
	}
}

type HttpApi interface {
	// GetLatestBlock 获取当前账户的最新的区块信息
	//
	// Parameters:
	//   - ctx context.Context
	//   - accountAddress string: 账户地址，zltc_Z1pnS94bP4hQSYLs4aP4UwBP9pH8bEvhi
	//
	// Returns:
	//   - types.LatestBlock
	//   - error
	GetLatestBlock(ctx context.Context, accountAddress string) (*types.LatestBlock, error)

	// SendSignedTransaction 发送已签名的交易
	//
	// Parameters:
	//    - ctx context.Context
	//    - signedTX *block.Transaction
	//
	// Returns:
	//    - o
	SendSignedTransaction(ctx context.Context, signedTX *block.Transaction) (*common.Hash, error)

	// GetReceipt 获取交易回执
	//
	// Parameters:
	//    - ctx context.Context
	//    - hash string
	//
	// Returns:
	//    - types.Receipt
	//    - error
	GetReceipt(ctx context.Context, hash string) (*types.Receipt, error)

	// GetContractLifecycleProposal 获取合约生命周期提案
	//
	// Parameters:
	//    - ctx context.Context
	//    - contractAddress string
	//    - state types.ProposalState
	//
	// Returns:
	//    - types.Proposal[types.ContractLifecycleProposal]
	//    - error
	GetContractLifecycleProposal(ctx context.Context, contractAddress string, state types.ProposalState) (*types.Proposal[types.ContractLifecycleProposal], error)
}

type httpApi struct {
	ChainId   string
	Url       string
	transport *http.Transport
}

func (api *httpApi) newHeaders() map[string]string {
	return map[string]string{
		"Content-Type": "application/json",
		"ChainId":      api.ChainId,
	}
}

func (api *httpApi) GetLatestBlock(_ context.Context, accountAddress string) (*types.LatestBlock, error) {
	response, err := Post[types.LatestBlock](api.Url, NewJsonRpcBody("latc_getCurrentTBDB", accountAddress), api.newHeaders(), api.transport)
	if err != nil {
		return nil, err
	}
	if response.Error != nil {
		return nil, response.Error.Error()
	}
	return response.Result, nil
}

func (api *httpApi) SendSignedTransaction(_ context.Context, signedTX *block.Transaction) (*common.Hash, error) {
	response, err := Post[common.Hash](api.Url, NewJsonRpcBody("wallet_sendRawTBlock", signedTX), api.newHeaders(), api.transport)
	if err != nil {
		return nil, err
	}
	if response.Error != nil {
		return nil, response.Error.Error()
	}
	return response.Result, nil
}

func (api *httpApi) GetReceipt(_ context.Context, hash string) (*types.Receipt, error) {
	response, err := Post[types.Receipt](api.Url, NewJsonRpcBody("latc_getReceipt", hash), api.newHeaders(), api.transport)
	if err != nil {
		return nil, err
	}
	if response.Error != nil {
		return nil, response.Error.Error()
	}
	return response.Result, nil
}

func (api *httpApi) GetContractLifecycleProposal(_ context.Context, contractAddress string, state types.ProposalState) (*types.Proposal[types.ContractLifecycleProposal], error) {
	params := map[string]interface{}{
		"proposalType":    types.ProposalTypeContractLifecycle,
		"proposalState":   state,
		"contractAddress": contractAddress,
	}

	response, err := Post[types.Proposal[types.ContractLifecycleProposal]](api.Url, NewJsonRpcBody("wallet_getProposal", params), api.newHeaders(), api.transport)
	if err != nil {
		return nil, err
	}
	if response.Error != nil {
		return nil, response.Error.Error()
	}
	return response.Result, nil
}

// Post send http request use post method
//
// Parameters:
//   - url string: 请求路径，示例：http://192.168.1.20:13000
//   - body sonRpcBody: any, 请求体
//   - headers map[string]string: 请求头
//   - tr http.Transport:
//
// Returns:
//   - []byte: 响应内容
//   - error: 错误
func Post[T any](url string, jsonRpcBody *JsonRpcBody, headers map[string]string, tr *http.Transport) (*JsonRpcResponse[*T], error) {
	bytes, err := json.Marshal(jsonRpcBody)
	if err != nil {
		return nil, err
	}
	body := strings.NewReader(string(bytes))

	request, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		return nil, err
	}

	if headers != nil && len(headers) != 0 {
		for key, value := range headers {
			request.Header.Set(key, value)
		}
	}
	client := &http.Client{Transport: tr}
	request.TransferEncoding = []string{}
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		if err := Body.Close(); err != nil {
			fmt.Println("Failed to close response body")
		}
	}(response.Body)

	if res, err := io.ReadAll(response.Body); err != nil {
		return nil, err
	} else {
		var t JsonRpcResponse[*T]
		if err := json.Unmarshal(res, &t); err != nil {
			return nil, err
		}

		return &t, nil
	}
}

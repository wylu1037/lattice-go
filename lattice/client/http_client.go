package client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/golang-jwt/jwt/v4"
	"io"
	"lattice-go/common/types"
	"lattice-go/lattice/block"
	"net/http"
	"strings"
	"time"
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
	Code    int16  `json:"code,omitempty"`
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

func NewJwt(secret string, expirationDuration time.Duration) Jwt {
	if secret == "" {
		return nil
	}
	return &jwtImpl{
		Secret:             secret,
		Algorithm:          jwt.SigningMethodHS256,
		ExpirationDuration: expirationDuration,
	}
}

// JwtTokenCache jwt token的缓存
type JwtTokenCache struct {
	Token    string
	ExpireAt time.Time
}

// IsValid 验证Token是否有效
//
// Returns:
//   - error
func (cache *JwtTokenCache) IsValid() error {
	if cache.Token == "" {
		return errors.New("token is empty")
	}

	if time.Now().After(cache.ExpireAt) {
		return errors.New("token is expired")
	}

	return nil
}

type Jwt interface {
	GenerateToken() (string, error)
	ParseToken(token string) (*jwt.Token, error)
	GetToken() (string, error)
}

type jwtImpl struct {
	Secret             string            // jwt的secret
	Algorithm          jwt.SigningMethod // jwt.SigningMethodHS256
	ExpirationDuration time.Duration     // token过期时长
	TokenCache         *JwtTokenCache    // token缓存
}

func (j *jwtImpl) GenerateToken() (string, error) {
	now := time.Now()
	expiresAt := now.Add(j.ExpirationDuration).Add(-3 * time.Minute) // 提前3分钟过期
	t := jwt.NewWithClaims(j.Algorithm, jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(expiresAt), // 该项验证
		IssuedAt:  jwt.NewNumericDate(now),       // 该项验证
		NotBefore: jwt.NewNumericDate(now),
		Issuer:    "lattice_go",
		Subject:   "jwt",
		ID:        "1",
		Audience:  []string{"somebody_else"},
	})
	token, err := t.SignedString([]byte(j.Secret))
	if err != nil {
		return "", err
	}

	j.TokenCache.Token = token
	j.TokenCache.ExpireAt = expiresAt

	return token, nil
}

func (j *jwtImpl) ParseToken(token string) (*jwt.Token, error) {
	t, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte("AllYourBase"), nil
	})
	if err != nil {
		return nil, err
	}

	switch {
	case t.Valid:
		return t, nil
	case errors.Is(err, jwt.ErrTokenMalformed):
		return nil, errors.New("that's not even a token")
	case errors.Is(err, jwt.ErrTokenSignatureInvalid):
		return nil, errors.New("invalid signature")
	case errors.Is(err, jwt.ErrTokenExpired) || errors.Is(err, jwt.ErrTokenNotValidYet):
		return nil, errors.New("token is either expired or not active yet")
	default:
		return nil, errors.New("couldn't handle this token")
	}
}

func (j *jwtImpl) GetToken() (string, error) {
	if err := j.TokenCache.IsValid(); err != nil {
		token, err := j.GenerateToken()
		if err != nil {
			return "", err
		}
		return token, nil
	}
	return j.TokenCache.Token, nil
}

// HttpApiInitParam 初始化HTTP API的参数
type HttpApiInitParam struct {
	Url                string          // 节点的URL
	Transport          *http.Transport // tr
	JwtSecret          string          // jwt的secret信息
	ExpirationDuration time.Duration   // jwt token的过期时间
}

func NewHttpApi(args *HttpApiInitParam) HttpApi {
	return &httpApi{
		Url:       args.Url,
		transport: args.Transport,
		jwtApi:    NewJwt(args.JwtSecret, args.ExpirationDuration),
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
	GetLatestBlock(ctx context.Context, chainId, accountAddress string) (*types.LatestBlock, error)

	// SendSignedTransaction 发送已签名的交易
	//
	// Parameters:
	//    - ctx context.Context
	//    - signedTX *block.Transaction
	//
	// Returns:
	//    - o
	SendSignedTransaction(ctx context.Context, chainId string, signedTX *block.Transaction) (*common.Hash, error)

	// PreCallContract 预执行合约
	//
	// Parameters:
	//   - ctx context.Context
	//   - unsignedTX *block.Transaction: 未签名的交易
	//
	// Returns:
	//   - *types.Receipt
	//   - error
	PreCallContract(ctx context.Context, chainId string, unsignedTX *block.Transaction) (*types.Receipt, error)

	// GetReceipt 获取交易回执
	//
	// Parameters:
	//    - ctx context.Context
	//    - hash string
	//
	// Returns:
	//    - types.Receipt
	//    - error
	GetReceipt(ctx context.Context, chainId, hash string) (*types.Receipt, error)

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
	GetContractLifecycleProposal(ctx context.Context, chainId, contractAddress string, state types.ProposalState) ([]types.Proposal[types.ContractLifecycleProposal], error)
}

type httpApi struct {
	Url       string
	transport *http.Transport
	jwtApi    Jwt
}

func (api *httpApi) newHeaders(chainId string) map[string]string {
	return map[string]string{
		"Content-Type": "application/json",
		"ChainId":      chainId,
	}
}

func (api *httpApi) GetLatestBlock(_ context.Context, chainId, accountAddress string) (*types.LatestBlock, error) {
	response, err := Post[types.LatestBlock](api.Url, NewJsonRpcBody("latc_getCurrentTBDB", accountAddress), api.newHeaders(chainId), api.transport)
	if err != nil {
		return nil, err
	}
	if response.Error != nil {
		return nil, response.Error.Error()
	}
	return response.Result, nil
}

func (api *httpApi) SendSignedTransaction(_ context.Context, chainId string, signedTX *block.Transaction) (*common.Hash, error) {
	response, err := Post[common.Hash](api.Url, NewJsonRpcBody("wallet_sendRawTBlock", signedTX), api.newHeaders(chainId), api.transport)
	if err != nil {
		return nil, err
	}
	if response.Error != nil {
		return nil, response.Error.Error()
	}
	return response.Result, nil
}

func (api *httpApi) PreCallContract(ctx context.Context, chainId string, unsignedTX *block.Transaction) (*types.Receipt, error) {
	response, err := Post[types.Receipt](api.Url, NewJsonRpcBody("wallet_preExecuteContract", unsignedTX), api.newHeaders(chainId), api.transport)
	if err != nil {
		return nil, err
	}
	if response.Error != nil {
		return nil, response.Error.Error()
	}
	return response.Result, nil
}

func (api *httpApi) GetReceipt(_ context.Context, chainId, hash string) (*types.Receipt, error) {
	response, err := Post[types.Receipt](api.Url, NewJsonRpcBody("latc_getReceipt", hash), api.newHeaders(chainId), api.transport)
	if err != nil {
		return nil, err
	}
	if response.Error != nil {
		return nil, response.Error.Error()
	}
	return response.Result, nil
}

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

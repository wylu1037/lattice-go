package lattice

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/avast/retry-go"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/samber/lo"
	"lattice-go/common/types"
	"lattice-go/crypto"
	"lattice-go/lattice/block"
	"lattice-go/lattice/client"
	"net/http"
	"strconv"
	"time"
)

const (
	httpProtocol      = "http"
	httpsProtocol     = "https"
	websocketProtocol = "ws"
	zeroAddress       = "zltc_QLbz7JHiBTspS962RLKV8GndWFwjA5K66"
	zeroHash          = "0x0000000000000000000000000000000000000000000000000000000000000000"
)

func NewLattice(chainConfig *ChainConfig, connectingNodeConfig *ConnectingNodeConfig, identityConfig *CredentialConfig, options *Options, blockCache BlockCache) Lattice {
	initHttpClientArgs := &client.HttpApiInitParam{
		Url:                connectingNodeConfig.GetHttpUrl(),
		Transport:          options.GetTransport(),
		JwtSecret:          connectingNodeConfig.JwtSecret,
		ExpirationDuration: connectingNodeConfig.JwtTokenExpirationDuration,
	}
	httpApi := client.NewHttpApi(initHttpClientArgs)

	if blockCache == nil {
		blockCache = newDisabledMemoryBlockCache(httpApi)
	} else {
		blockCache.SetHttpApi(httpApi)
	}

	return &lattice{
		ChainConfig:          chainConfig,
		ConnectingNodeConfig: connectingNodeConfig,
		CredentialConfig:     identityConfig,
		Options:              options,
		httpApi:              httpApi,
		BlockCache:           blockCache,
	}
}

type lattice struct {
	httpApi              client.HttpApi
	ChainConfig          *ChainConfig
	ConnectingNodeConfig *ConnectingNodeConfig
	CredentialConfig     *CredentialConfig
	Options              *Options
	BlockCache           BlockCache
}

// ChainConfig 链配置
type ChainConfig struct {
	Curve     types.Curve // crypto.Secp256k1 or crypto.Sm2p256v1
	TokenLess bool        // false:有通证链，true:无通证链
}

// ConnectingNodeConfig 节点配置
type ConnectingNodeConfig struct {
	Insecure                   bool
	Ip                         string
	HttpPort                   uint16
	WebsocketPort              uint16
	JwtSecret                  string
	JwtTokenExpirationDuration time.Duration
}

// CredentialConfig 凭证配置
type CredentialConfig struct {
	AccountAddress string // 账户地址
	Passphrase     string // 身份密码
	PrivateKey     string // 私钥
}

type Options struct {
	Transport *http.Transport

	InsecureSkipVerify bool

	// MaxIdleConns controls the maximum number of idle (keep-alive)
	// connections across all hosts. Zero means no limit.
	MaxIdleConns int

	// MaxIdleConnsPerHost if non-zero, controls the maximum idle
	// (keep-alive) connections to keep per-host.
	// If zero, DefaultMaxIdleConnsPerHost(2) is used.
	MaxIdleConnsPerHost int
}

func (chain *ChainConfig) GetCurve() types.Curve {
	switch chain.Curve {
	case crypto.Sm2p256v1:
		return crypto.Sm2p256v1
	case crypto.Secp256k1:
		return crypto.Secp256k1
	default:
		return crypto.Sm2p256v1
	}
}

func (options *Options) GetTransport() *http.Transport {
	if options.Transport == nil {
		options.Transport = &http.Transport{
			TLSClientConfig:     &tls.Config{InsecureSkipVerify: options.InsecureSkipVerify},
			MaxIdleConns:        options.MaxIdleConns,
			MaxIdleConnsPerHost: options.MaxIdleConnsPerHost,
		}
	}
	return options.Transport
}

func (identity *CredentialConfig) GetSK() string {
	if identity.PrivateKey == "" {
		// decrypt file key
	}
	return identity.PrivateKey
}

func (node *ConnectingNodeConfig) GetHttpUrl() string {
	return fmt.Sprintf("%s://%s:%d", lo.Ternary(node.Insecure, httpsProtocol, httpProtocol), node.Ip, node.HttpPort)
}

func (node *ConnectingNodeConfig) GetWebsocketUrl() string {
	return fmt.Sprintf("%s://%s:%d", websocketProtocol, node.Ip, node.WebsocketPort)
}

type Strategy string

const (
	BackOff        = "BackOff"
	FixedInterval  = "FixedInterval"
	RandomInterval = "RandomInterval"
)

// RetryStrategy 等待回执策略
type RetryStrategy struct {
	// 具体的策略
	Strategy  Strategy
	Attempts  uint
	Delay     time.Duration
	MaxJitter time.Duration
}

func (strategy *RetryStrategy) GetRetryStrategyOpts() []retry.Option {
	switch strategy.Strategy {
	case BackOff:
		return strategy.BackOffOpts()
	case FixedInterval:
		return strategy.FixedIntervalOpts()
	case RandomInterval:
		return strategy.RandomIntervalOpts()
	default:
		return []retry.Option{}
	}
}

func NewBackOffRetryStrategy(attempts uint, initDelay time.Duration) *RetryStrategy {
	return &RetryStrategy{
		Strategy: BackOff,
		Attempts: attempts,
		Delay:    initDelay,
	}
}

// DefaultBackOffRetryStrategy 创建默认的BackOff等待策略
//
// Parameters:
//
// Returns:
//   - RetryStrategy
func DefaultBackOffRetryStrategy() *RetryStrategy {
	return &RetryStrategy{
		Strategy: BackOff,
		Attempts: 10,
		Delay:    time.Millisecond * 200,
	}
}

func NewFixedRetryStrategy(attempts uint, fixedDelay time.Duration) *RetryStrategy {
	return &RetryStrategy{
		Strategy: FixedInterval,
		Attempts: attempts,
		Delay:    fixedDelay,
	}
}

// DefaultFixedRetryStrategy 创建默认的固定等待策略
//
// Parameters:
//
// Returns:
//   - RetryStrategy
func DefaultFixedRetryStrategy() *RetryStrategy {
	return &RetryStrategy{
		Strategy: FixedInterval,
		Attempts: 10,
		Delay:    time.Millisecond * 150,
	}
}

func NewRandomRetryStrategy(attempts uint, baseDelay time.Duration, maxJitter time.Duration) *RetryStrategy {
	return &RetryStrategy{
		Strategy:  RandomInterval,
		Attempts:  attempts,
		Delay:     baseDelay,
		MaxJitter: maxJitter,
	}
}

// DefaultRandomRetryStrategy 创建默认的随机等待策略
//
// Parameters:
//
// Returns:
//   - RetryStrategy
func DefaultRandomRetryStrategy() *RetryStrategy {
	return &RetryStrategy{
		Strategy:  RandomInterval,
		Attempts:  10,
		Delay:     time.Millisecond * 100,
		MaxJitter: time.Millisecond * 500, // 最大的随机抖动
	}
}

func (strategy *RetryStrategy) BackOffOpts() []retry.Option {
	return []retry.Option{retry.Attempts(strategy.Attempts), retry.Delay(strategy.Delay), retry.DelayType(retry.BackOffDelay)}
}

func (strategy *RetryStrategy) FixedIntervalOpts() []retry.Option {
	return []retry.Option{retry.Attempts(strategy.Attempts), retry.Delay(strategy.Delay), retry.DelayType(retry.FixedDelay)}
}

func (strategy *RetryStrategy) RandomIntervalOpts() []retry.Option {
	return []retry.Option{retry.Attempts(strategy.Attempts), retry.Delay(strategy.Delay), retry.MaxJitter(strategy.MaxJitter), retry.DelayType(retry.RandomDelay)}
}

type Lattice interface {
	// HttpApi return the http api
	//
	// Parameters:
	//
	// Returns:
	//   - client.HttpApi
	HttpApi() client.HttpApi

	// Transfer 发起转账交易
	//
	// Parameters:
	//    - ctx context.Context
	//    - linker string: 转账接收者账户地址
	//    - payload string: 交易备注
	//
	// Returns:
	//    - *common.Hash: 交易哈希
	//    - error
	Transfer(ctx context.Context, chainId, linker, payload string, amount, joule uint64) (*common.Hash, error)

	// DeployContract 发起部署合约交易
	//
	// Parameters:
	//   - ctx context.Context
	//   - data string: 合约数据
	//   - payload string: 交易备注
	//
	// Returns:
	//   - *common.Hash: 交易哈希
	//   - error
	DeployContract(ctx context.Context, chainId, data, payload string, amount, joule uint64) (*common.Hash, error)

	// CallContract 发起调用合约交易
	//
	// Parameters:
	//   - ctx context.Context
	//   - contractAddress string: 合约地址
	//   - data string: 调用的合约数据
	//   - payload string: 交易备注
	//
	// Returns:
	//   - *common.Hash: 交易哈希
	//   - error
	CallContract(ctx context.Context, contractAddress, chainId, data, payload string, amount, joule uint64) (*common.Hash, error)

	// TransferWaitReceipt 发起转账交易并等待回执
	//
	// Parameters:
	//   - ctx context.Context
	//   - linker string
	//   - payload string: 交易备注
	//   - retryStrategy *RetryStrategy: 等待回执策略
	//
	// Returns:
	//   - *common.Hash: 交易哈希
	//   - *types.Receipt: 回执
	//   - error
	TransferWaitReceipt(ctx context.Context, chainId, linker, payload string, amount, joule uint64, retryStrategy *RetryStrategy) (*common.Hash, *types.Receipt, error)

	// DeployContractWaitReceipt 发起部署合约交易并等待回执
	//
	// Parameters:
	//   - ctx context.Context
	//   - data string
	//   - payload string: 交易备注
	//   - retryStrategy *RetryStrategy: 等待回执策略
	//
	// Returns:
	//   - *common.Hash: 交易哈希
	//   - *types.Receipt: 回执
	//   - error
	DeployContractWaitReceipt(ctx context.Context, chainId, data, payload string, amount, joule uint64, retryStrategy *RetryStrategy) (*common.Hash, *types.Receipt, error)

	// CallContractWaitReceipt 发起调用合约交易并等待回执
	//
	// Parameters:
	//   - ctx context.Context
	//   - contractAddress string
	//   - data string
	//   - payload string: 交易备注
	//   - retryStrategy *RetryStrategy: 等待回执策略
	//
	// Returns:
	//   - *common.Hash: 交易哈希
	//   - *types.Receipt: 回执
	//   - error
	CallContractWaitReceipt(ctx context.Context, chainId, contractAddress, data, payload string, amount, joule uint64, retryStrategy *RetryStrategy) (*common.Hash, *types.Receipt, error)

	// PreCallContract 预执行合约，预执行的交易不会上链
	//
	// Parameters:
	//   - ctx context.Context:
	//   - contractAddress string: 合约地址
	//   - data string: 执行的合约代码
	//   - payload string: 交易备注
	//
	// Returns:
	//   - *types.Receipt: 交易回执
	//   - error: 预执行的错误
	PreCallContract(ctx context.Context, chainId, contractAddress, data, payload string) (*types.Receipt, error)
}

func (svc *lattice) HttpApi() client.HttpApi {
	return svc.httpApi
}

func (svc *lattice) Transfer(ctx context.Context, chainId, linker, payload string, amount, joule uint64) (*common.Hash, error) {
	latestBlock, err := svc.BlockCache.GetBlock(chainId, svc.CredentialConfig.AccountAddress)
	if err != nil {
		return nil, err
	}

	transaction := block.NewTransactionBuilder(block.TransactionTypeSend).
		SetLatestBlock(latestBlock).
		SetOwner(svc.CredentialConfig.AccountAddress).
		SetLinker(linker).
		SetPayload(payload).
		SetAmount(amount).
		SetJoule(joule).
		Build()

	chainIdAsInt, err := strconv.Atoi(chainId)
	if err != nil {
		return nil, err
	}
	err = transaction.SignTX(uint64(chainIdAsInt), svc.ChainConfig.GetCurve(), svc.CredentialConfig.GetSK())
	if err != nil {
		return nil, err
	}

	hash, err := svc.httpApi.SendSignedTransaction(ctx, chainId, transaction)
	if err != nil {
		return nil, err
	} else {
		latestBlock.Hash = *hash
		latestBlock.IncrHeight()
		if err := svc.BlockCache.SetBlock(chainId, svc.CredentialConfig.AccountAddress, latestBlock); err != nil {
			fmt.Println(err)
		}
	}
	return hash, nil
}

func (svc *lattice) DeployContract(ctx context.Context, chainId, data, payload string, amount, joule uint64) (*common.Hash, error) {
	latestBlock, err := svc.httpApi.GetLatestBlock(ctx, chainId, svc.CredentialConfig.AccountAddress)
	if err != nil {
		return nil, err
	}

	transaction := block.NewTransactionBuilder(block.TransactionTypeDeployContract).
		SetLatestBlock(latestBlock).
		SetOwner(svc.CredentialConfig.AccountAddress).
		SetLinker(zeroAddress).
		SetCode(data).
		SetPayload(payload).
		SetAmount(amount).
		SetJoule(joule).
		Build()

	cryptoInstance := crypto.NewCrypto(svc.ChainConfig.Curve)
	dataHash := cryptoInstance.Hash(hexutil.MustDecode(data))
	transaction.CodeHash = dataHash

	chainIdAsInt, err := strconv.Atoi(chainId)
	if err != nil {
		return nil, err
	}
	err = transaction.SignTX(uint64(chainIdAsInt), svc.ChainConfig.GetCurve(), svc.CredentialConfig.GetSK())
	if err != nil {
		return nil, err
	}

	hash, err := svc.httpApi.SendSignedTransaction(ctx, chainId, transaction)
	if err != nil {
		return nil, err
	}
	return hash, nil
}

func (svc *lattice) CallContract(ctx context.Context, chainId, contractAddress, data, payload string, amount, joule uint64) (*common.Hash, error) {
	latestBlock, err := svc.httpApi.GetLatestBlock(ctx, chainId, svc.CredentialConfig.AccountAddress)
	if err != nil {
		return nil, err
	}

	transaction := block.NewTransactionBuilder(block.TransactionTypeCallContract).
		SetLatestBlock(latestBlock).
		SetOwner(svc.CredentialConfig.AccountAddress).
		SetLinker(contractAddress).
		SetCode(data).
		SetPayload(payload).
		SetAmount(amount).
		SetJoule(joule).
		Build()

	cryptoInstance := crypto.NewCrypto(svc.ChainConfig.Curve)
	dataHash := cryptoInstance.Hash(hexutil.MustDecode(data))
	transaction.CodeHash = dataHash

	chainIdAsInt, err := strconv.Atoi(chainId)
	if err != nil {
		return nil, err
	}
	err = transaction.SignTX(uint64(chainIdAsInt), svc.ChainConfig.GetCurve(), svc.CredentialConfig.GetSK())
	if err != nil {
		return nil, err
	}

	hash, err := svc.httpApi.SendSignedTransaction(ctx, chainId, transaction)
	if err != nil {
		return nil, err
	}
	return hash, nil
}

func (svc *lattice) waitReceipt(ctx context.Context, chainId string, hash *common.Hash, retryStrategy *RetryStrategy) (*common.Hash, *types.Receipt, error) {
	var err error
	var receipt *types.Receipt
	err = retry.Do(
		func() error {
			receipt, err = svc.httpApi.GetReceipt(ctx, chainId, hash.String())
			if err != nil {
				return err
			}
			return nil
		},
		retryStrategy.GetRetryStrategyOpts()...,
	)

	if err != nil {
		return hash, nil, err
	}
	return hash, receipt, nil
}

func (svc *lattice) TransferWaitReceipt(ctx context.Context, chainId, linker, payload string, amount, joule uint64, retryStrategy *RetryStrategy) (*common.Hash, *types.Receipt, error) {
	hash, err := svc.Transfer(ctx, chainId, linker, payload, amount, joule)
	if err != nil {
		return nil, nil, err
	}

	return svc.waitReceipt(ctx, chainId, hash, retryStrategy)
}

func (svc *lattice) DeployContractWaitReceipt(ctx context.Context, chainId, data, payload string, amount, joule uint64, retryStrategy *RetryStrategy) (*common.Hash, *types.Receipt, error) {
	hash, err := svc.DeployContract(ctx, chainId, data, payload, amount, joule)
	if err != nil {
		return nil, nil, err
	}

	return svc.waitReceipt(ctx, chainId, hash, retryStrategy)
}

func (svc *lattice) CallContractWaitReceipt(ctx context.Context, chainId, contractAddress, data, payload string, amount, joule uint64, retryStrategy *RetryStrategy) (*common.Hash, *types.Receipt, error) {
	hash, err := svc.CallContract(ctx, chainId, contractAddress, data, payload, amount, joule)
	if err != nil {
		return nil, nil, err
	}

	return svc.waitReceipt(ctx, chainId, hash, retryStrategy)
}

func (svc *lattice) PreCallContract(ctx context.Context, chainId, contractAddress, data, payload string) (*types.Receipt, error) {
	transaction := block.NewTransactionBuilder(block.TransactionTypeCallContract).
		SetLatestBlock(
			&types.LatestBlock{
				Height:          0,
				Hash:            common.HexToHash(zeroHash),
				DaemonBlockHash: common.HexToHash(zeroHash),
			}).
		SetOwner(svc.CredentialConfig.AccountAddress).
		SetLinker(contractAddress).
		SetCode(data).
		SetPayload(payload).
		Build()

	receipt, err := svc.httpApi.PreCallContract(ctx, chainId, transaction)
	if err != nil {
		return nil, err
	}
	return receipt, nil
}

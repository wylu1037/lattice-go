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
)

func NewLattice(chainConfig *ChainConfig, nodeConfig *NodeConfig, identityConfig *CredentialConfig, options *Options) Lattice {
	return &lattice{
		Chain:      chainConfig,
		Node:       nodeConfig,
		Credential: identityConfig,
		Options:    options,
		httpApi:    client.NewHttpApi(nodeConfig.GetHttpUrl(), strconv.FormatUint(chainConfig.ChainId, 10), options.GetTransport()),
	}
}

type lattice struct {
	httpApi    client.HttpApi
	Chain      *ChainConfig
	Node       *NodeConfig
	Credential *CredentialConfig
	Options    *Options
}

// ChainConfig 链配置
type ChainConfig struct {
	ChainId uint64
	Curve   types.Curve
}

// NodeConfig 节点配置
type NodeConfig struct {
	Insecure      bool
	Ip            string
	HttpPort      uint16
	WebsocketPort uint16
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

func (node *NodeConfig) GetHttpUrl() string {
	return fmt.Sprintf("%s://%s:%d", lo.Ternary(node.Insecure, httpsProtocol, httpProtocol), node.Ip, node.HttpPort)
}

func (node *NodeConfig) GetWebsocketUrl() string {
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
	Transfer(ctx context.Context, linker, payload string, amount, joule uint64) (*common.Hash, error)

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
	DeployContract(ctx context.Context, data, payload string, amount, joule uint64) (*common.Hash, error)

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
	CallContract(ctx context.Context, contractAddress, data, payload string, amount, joule uint64) (*common.Hash, error)

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
	TransferWaitReceipt(ctx context.Context, linker, payload string, amount, joule uint64, retryStrategy *RetryStrategy) (*common.Hash, *types.Receipt, error)

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
	DeployContractWaitReceipt(ctx context.Context, data, payload string, amount, joule uint64, retryStrategy *RetryStrategy) (*common.Hash, *types.Receipt, error)

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
	CallContractWaitReceipt(ctx context.Context, contractAddress, data, payload string, amount, joule uint64, retryStrategy *RetryStrategy) (*common.Hash, *types.Receipt, error)

	PreCallContract(ctx context.Context, contractAddress, data, payload string) (*types.Receipt, error)
}

func (svc *lattice) HttpApi() client.HttpApi {
	return svc.httpApi
}

func (svc *lattice) Transfer(ctx context.Context, linker, payload string, amount, joule uint64) (*common.Hash, error) {
	latestBlock, err := svc.httpApi.GetLatestBlock(ctx, svc.Credential.AccountAddress)
	if err != nil {
		return nil, err
	}

	transaction := block.NewTransactionBuilder(block.TransactionTypeSend).
		SetLatestBlock(latestBlock).
		SetOwner(svc.Credential.AccountAddress).
		SetLinker(linker).
		SetPayload(payload).
		SetAmount(amount).
		SetJoule(joule).
		Build()

	err = transaction.SignTX(svc.Chain.ChainId, svc.Chain.GetCurve(), svc.Credential.GetSK())
	if err != nil {
		return nil, err
	}

	hash, err := svc.httpApi.SendSignedTransaction(ctx, transaction)
	if err != nil {
		return nil, err
	}
	return hash, nil
}

func (svc *lattice) DeployContract(ctx context.Context, data, payload string, amount, joule uint64) (*common.Hash, error) {
	latestBlock, err := svc.httpApi.GetLatestBlock(ctx, svc.Credential.AccountAddress)
	if err != nil {
		return nil, err
	}

	transaction := block.NewTransactionBuilder(block.TransactionTypeDeployContract).
		SetLatestBlock(latestBlock).
		SetOwner(svc.Credential.AccountAddress).
		SetLinker(zeroAddress).
		SetCode(data).
		SetPayload(payload).
		SetAmount(amount).
		SetJoule(joule).
		Build()

	cryptoInstance := crypto.NewCrypto(svc.Chain.Curve)
	dataHash := cryptoInstance.Hash(hexutil.MustDecode(data))
	transaction.CodeHash = dataHash

	err = transaction.SignTX(svc.Chain.ChainId, svc.Chain.GetCurve(), svc.Credential.GetSK())
	if err != nil {
		return nil, err
	}

	hash, err := svc.httpApi.SendSignedTransaction(ctx, transaction)
	if err != nil {
		return nil, err
	}
	return hash, nil
}

func (svc *lattice) CallContract(ctx context.Context, contractAddress, data, payload string, amount, joule uint64) (*common.Hash, error) {
	latestBlock, err := svc.httpApi.GetLatestBlock(ctx, svc.Credential.AccountAddress)
	if err != nil {
		return nil, err
	}

	transaction := block.NewTransactionBuilder(block.TransactionTypeCallContract).
		SetLatestBlock(latestBlock).
		SetOwner(svc.Credential.AccountAddress).
		SetLinker(contractAddress).
		SetCode(data).
		SetPayload(payload).
		SetAmount(amount).
		SetJoule(joule).
		Build()

	cryptoInstance := crypto.NewCrypto(svc.Chain.Curve)
	dataHash := cryptoInstance.Hash(hexutil.MustDecode(data))
	transaction.CodeHash = dataHash

	err = transaction.SignTX(svc.Chain.ChainId, svc.Chain.GetCurve(), svc.Credential.GetSK())
	if err != nil {
		return nil, err
	}

	hash, err := svc.httpApi.SendSignedTransaction(ctx, transaction)
	if err != nil {
		return nil, err
	}
	return hash, nil
}

func (svc *lattice) waitReceipt(ctx context.Context, hash *common.Hash, retryStrategy *RetryStrategy) (*common.Hash, *types.Receipt, error) {
	var err error
	var receipt *types.Receipt
	err = retry.Do(
		func() error {
			receipt, err = svc.httpApi.GetReceipt(ctx, hash.String())
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

func (svc *lattice) TransferWaitReceipt(ctx context.Context, linker, payload string, amount, joule uint64, retryStrategy *RetryStrategy) (*common.Hash, *types.Receipt, error) {
	hash, err := svc.Transfer(ctx, linker, payload, amount, joule)
	if err != nil {
		return nil, nil, err
	}

	return svc.waitReceipt(ctx, hash, retryStrategy)
}

func (svc *lattice) DeployContractWaitReceipt(ctx context.Context, data, payload string, amount, joule uint64, retryStrategy *RetryStrategy) (*common.Hash, *types.Receipt, error) {
	hash, err := svc.DeployContract(ctx, data, payload, amount, joule)
	if err != nil {
		return nil, nil, err
	}

	return svc.waitReceipt(ctx, hash, retryStrategy)
}

func (svc *lattice) CallContractWaitReceipt(ctx context.Context, contractAddress, data, payload string, amount, joule uint64, retryStrategy *RetryStrategy) (*common.Hash, *types.Receipt, error) {
	hash, err := svc.CallContract(ctx, contractAddress, data, payload, amount, joule)
	if err != nil {
		return nil, nil, err
	}

	return svc.waitReceipt(ctx, hash, retryStrategy)
}

func (svc *lattice) PreCallContract(ctx context.Context, contractAddress, data, payload string) (*types.Receipt, error) {
	transaction := block.NewTransactionBuilder(block.TransactionTypeCallContract).
		SetLatestBlock(
			&types.LatestBlock{
				Height:          0,
				Hash:            common.HexToHash(""),
				DaemonBlockHash: common.HexToHash(""),
			}).
		SetOwner(svc.Credential.AccountAddress).
		SetLinker(contractAddress).
		SetCode(data).
		SetPayload(payload).
		Build()

	//cryptoInstance := crypto.NewCrypto(svc.Chain.Curve)
	//dataHash := cryptoInstance.Hash(hexutil.MustDecode(data))
	//transaction.CodeHash = dataHash

	receipt, err := svc.httpApi.PreCallContract(ctx, transaction)
	if err != nil {
		return nil, err
	}
	return receipt, nil
}

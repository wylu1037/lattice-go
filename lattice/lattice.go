package lattice

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/avast/retry-go"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/rs/zerolog/log"
	"github.com/samber/lo"
	"github.com/wylu1037/lattice-go/common/constant"
	"github.com/wylu1037/lattice-go/common/types"
	"github.com/wylu1037/lattice-go/crypto"
	"github.com/wylu1037/lattice-go/lattice/block"
	"github.com/wylu1037/lattice-go/lattice/client"
	"github.com/wylu1037/lattice-go/wallet"
	"net/http"
	"strconv"
	"time"
)

const (
	httpProtocol      = "http"
	httpsProtocol     = "https"
	websocketProtocol = "ws"
)

func NewLattice(chainConfig *ChainConfig, connectingNodeConfig *ConnectingNodeConfig, blockCache BlockCache, accountLock AccountLock, options *Options) Lattice {
	initHttpClientArgs := &client.HttpApiInitParam{
		Url:                        connectingNodeConfig.GetHttpUrl(),
		Transport:                  options.GetTransport(),
		JwtSecret:                  connectingNodeConfig.JwtSecret,
		JwtTokenExpirationDuration: connectingNodeConfig.JwtTokenExpirationDuration,
	}
	httpApi := client.NewHttpApi(initHttpClientArgs)

	if blockCache == nil {
		blockCache = newDisabledMemoryBlockCache(httpApi)
	} else {
		blockCache.SetHttpApi(httpApi)
	}

	if accountLock == nil {
		accountLock = NewAccountLock()
	}

	return &lattice{
		chainConfig:          chainConfig,
		connectingNodeConfig: connectingNodeConfig,
		options:              options,
		httpApi:              httpApi,
		blockCache:           blockCache,
		accountLock:          accountLock,
	}
}

type lattice struct {
	httpApi              client.HttpApi
	chainConfig          *ChainConfig
	connectingNodeConfig *ConnectingNodeConfig
	// credentials          *Credentials
	blockCache  BlockCache
	accountLock AccountLock
	options     *Options
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

// Credentials 凭证配置
type Credentials struct {
	AccountAddress string // 账户地址
	Passphrase     string // 身份密码
	FileKey        string // FileKey 的json字符串
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

// GetSK 获取私钥的Hex字符串
//
// Returns:
//   - string: 私钥的hex字符串
//   - error
func (credentials *Credentials) GetSK() (string, error) {
	if credentials.PrivateKey == "" {
		fileKey := wallet.NewFileKey(credentials.FileKey)
		privateKey, err := fileKey.Decrypt(credentials.Passphrase)
		if err != nil {
			return "", err
		}

		api := crypto.NewCrypto(lo.Ternary(fileKey.IsGM, crypto.Sm2p256v1, crypto.Secp256k1))
		sk, err := api.SKToHexString(privateKey)
		if err != nil {
			return "", err
		}
		credentials.PrivateKey = sk
		return sk, nil
	}
	return credentials.PrivateKey, nil
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
// Returns:
//   - RetryStrategy
func DefaultBackOffRetryStrategy() *RetryStrategy {
	return &RetryStrategy{
		Strategy: BackOff,
		Attempts: 15,
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
// Returns:
//   - RetryStrategy
func DefaultFixedRetryStrategy() *RetryStrategy {
	return &RetryStrategy{
		Strategy: FixedInterval,
		Attempts: 15,
		Delay:    time.Millisecond * 100,
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
		Attempts:  15,
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
	//    - credentials *Credentials: 发交易的身份凭证
	//   - chainId string
	//    - linker string: 转账接收者账户地址
	//    - payload string: 交易备注
	//
	// Returns:
	//    - *common.Hash: 交易哈希
	//    - error
	Transfer(ctx context.Context, credentials *Credentials, chainId, linker, payload string, amount, joule uint64) (*common.Hash, error)

	// DeployContract 发起部署合约交易
	//
	// Parameters:
	//   - ctx context.Context
	//   - credentials *Credentials: 发交易的身份凭证
	//   - chainId string
	//   - data string: 合约数据
	//   - payload string: 交易备注
	//
	// Returns:
	//   - *common.Hash: 交易哈希
	//   - error
	DeployContract(ctx context.Context, credentials *Credentials, chainId, data, payload string, amount, joule uint64) (*common.Hash, error)

	// CallContract 发起调用合约交易
	//
	// Parameters:
	//   - ctx context.Context
	//   - credentials *Credentials: 发交易的身份凭证
	//   - chainId string
	//   - contractAddress string: 合约地址
	//   - data string: 调用的合约数据
	//   - payload string: 交易备注
	//
	// Returns:
	//   - *common.Hash: 交易哈希
	//   - error
	CallContract(ctx context.Context, credentials *Credentials, chainId, contractAddress, data, payload string, amount, joule uint64) (*common.Hash, error)

	// TransferWaitReceipt 发起转账交易并等待回执
	//
	// Parameters:
	//   - ctx context.Context
	//   - credentials *Credentials: 发交易的身份凭证
	//   - chainId string
	//   - linker string
	//   - payload string: 交易备注
	//   - retryStrategy *RetryStrategy: 等待回执策略
	//
	// Returns:
	//   - *common.Hash: 交易哈希
	//   - *types.Receipt: 回执
	//   - error
	TransferWaitReceipt(ctx context.Context, credentials *Credentials, chainId, linker, payload string, amount, joule uint64, retryStrategy *RetryStrategy) (*common.Hash, *types.Receipt, error)

	// DeployContractWaitReceipt 发起部署合约交易并等待回执
	//
	// Parameters:
	//   - ctx context.Context
	//   - credentials *Credentials: 发交易的身份凭证
	//   - chainId string
	//   - data string
	//   - payload string: 交易备注
	//   - retryStrategy *RetryStrategy: 等待回执策略
	//
	// Returns:
	//   - *common.Hash: 交易哈希
	//   - *types.Receipt: 回执
	//   - error
	DeployContractWaitReceipt(ctx context.Context, credentials *Credentials, chainId, data, payload string, amount, joule uint64, retryStrategy *RetryStrategy) (*common.Hash, *types.Receipt, error)

	// CallContractWaitReceipt 发起调用合约交易并等待回执
	//
	// Parameters:
	//   - ctx context.Context
	//   - credentials *Credentials: 发交易的身份凭证
	//   - chainId string
	//   - contractAddress string
	//   - data string
	//   - payload string: 交易备注
	//   - retryStrategy *RetryStrategy: 等待回执策略
	//
	// Returns:
	//   - *common.Hash: 交易哈希
	//   - *types.Receipt: 回执
	//   - error
	CallContractWaitReceipt(ctx context.Context, credentials *Credentials, chainId, contractAddress, data, payload string, amount, joule uint64, retryStrategy *RetryStrategy) (*common.Hash, *types.Receipt, error)

	// PreCallContract 预执行合约，预执行的交易不会上链
	//
	// Parameters:
	//   - ctx context.Context:
	//   - chainId string
	//   - owner string: 调用者账户地址
	//   - contractAddress string: 合约地址
	//   - data string: 执行的合约代码
	//   - payload string: 交易备注
	//
	// Returns:
	//   - *types.Receipt: 交易回执
	//   - error: 预执行的错误
	PreCallContract(ctx context.Context, chainId, owner, contractAddress, data, payload string) (*types.Receipt, error)
}

func (svc *lattice) HttpApi() client.HttpApi {
	return svc.httpApi
}

func (svc *lattice) Transfer(ctx context.Context, credentials *Credentials, chainId, linker, payload string, amount, joule uint64) (*common.Hash, error) {
	log.Debug().Msgf("开始发起转账交易，chainId: %s, linker: %s, payload: %s, amount: %d, joule: %d", chainId, linker, payload, amount, joule)

	svc.accountLock.Obtain(chainId, credentials.AccountAddress)
	defer svc.accountLock.Unlock(chainId, credentials.AccountAddress)

	latestBlock, err := svc.blockCache.GetBlock(chainId, credentials.AccountAddress)
	if err != nil {
		log.Error().Err(err)
		return nil, err
	}

	transaction := block.NewTransactionBuilder(block.TransactionTypeSend).
		SetLatestBlock(latestBlock).
		SetOwner(credentials.AccountAddress).
		SetLinker(linker).
		SetPayload(payload).
		SetAmount(amount).
		SetJoule(joule).
		Build()

	chainIdAsInt, err := strconv.Atoi(chainId)
	if err != nil {
		log.Error().Err(err)
		return nil, err
	}
	sk, err := credentials.GetSK()
	if err != nil {
		log.Error().Err(err)
		return nil, err
	}
	err = transaction.SignTX(uint64(chainIdAsInt), svc.chainConfig.GetCurve(), sk)
	if err != nil {
		log.Error().Err(err)
		return nil, err
	}

	hash, err := svc.httpApi.SendSignedTransaction(ctx, chainId, transaction)
	if err != nil {
		log.Error().Err(err)
		return nil, err
	} else {
		latestBlock.Hash = *hash
		latestBlock.IncrHeight()
		if err := svc.blockCache.SetBlock(chainId, credentials.AccountAddress, latestBlock); err != nil {
			log.Error().Err(err)
		}
	}
	log.Debug().Msgf("结束转账交易，哈希为：%s", hash.String())
	return hash, nil
}

func (svc *lattice) DeployContract(ctx context.Context, credentials *Credentials, chainId, data, payload string, amount, joule uint64) (*common.Hash, error) {
	log.Debug().Msgf("开始发起部署合约交易，chainId: %s, data: %s, payload: %s, amount: %d, joule: %d", chainId, data, payload, amount, joule)

	svc.accountLock.Obtain(chainId, credentials.AccountAddress)
	defer svc.accountLock.Unlock(chainId, credentials.AccountAddress)

	latestBlock, err := svc.blockCache.GetBlock(chainId, credentials.AccountAddress)
	if err != nil {
		log.Error().Err(err)
		return nil, err
	}

	transaction := block.NewTransactionBuilder(block.TransactionTypeDeployContract).
		SetLatestBlock(latestBlock).
		SetOwner(credentials.AccountAddress).
		SetLinker(constant.ZeroAddress).
		SetCode(data).
		SetPayload(payload).
		SetAmount(amount).
		SetJoule(joule).
		Build()

	cryptoInstance := crypto.NewCrypto(svc.chainConfig.Curve)
	dataHash := cryptoInstance.Hash(hexutil.MustDecode(data))
	transaction.CodeHash = dataHash

	chainIdAsInt, err := strconv.Atoi(chainId)
	if err != nil {
		log.Error().Err(err)
		return nil, err
	}
	sk, err := credentials.GetSK()
	if err != nil {
		log.Error().Err(err)
		return nil, err
	}
	err = transaction.SignTX(uint64(chainIdAsInt), svc.chainConfig.GetCurve(), sk)
	if err != nil {
		log.Error().Err(err)
		return nil, err
	}

	hash, err := svc.httpApi.SendSignedTransaction(ctx, chainId, transaction)
	if err != nil {
		log.Error().Err(err)
		return nil, err
	} else {
		latestBlock.Hash = *hash
		latestBlock.IncrHeight()
		if err := svc.blockCache.SetBlock(chainId, credentials.AccountAddress, latestBlock); err != nil {
			log.Error().Err(err)
		}
	}
	log.Debug().Msgf("结束部署合约，哈希为：%s", hash.String())
	return hash, nil
}

func (svc *lattice) CallContract(ctx context.Context, credentials *Credentials, chainId, contractAddress, data, payload string, amount, joule uint64) (*common.Hash, error) {
	log.Debug().Msgf("开始发起调用合约交易，chainId: %s, contractAddress: %s, data: %s, payload: %s, amount: %d, joule: %d", chainId, contractAddress, data, payload, amount, joule)

	svc.accountLock.Obtain(chainId, credentials.AccountAddress)
	defer svc.accountLock.Unlock(chainId, credentials.AccountAddress)

	latestBlock, err := svc.blockCache.GetBlock(chainId, credentials.AccountAddress)
	if err != nil {
		return nil, err
	}

	transaction := block.NewTransactionBuilder(block.TransactionTypeCallContract).
		SetLatestBlock(latestBlock).
		SetOwner(credentials.AccountAddress).
		SetLinker(contractAddress).
		SetCode(data).
		SetPayload(payload).
		SetAmount(amount).
		SetJoule(joule).
		Build()

	cryptoInstance := crypto.NewCrypto(svc.chainConfig.Curve)
	dataHash := cryptoInstance.Hash(hexutil.MustDecode(data))
	transaction.CodeHash = dataHash

	chainIdAsInt, err := strconv.Atoi(chainId)
	if err != nil {
		log.Error().Err(err)
		return nil, err
	}
	sk, err := credentials.GetSK()
	if err != nil {
		log.Error().Err(err)
		return nil, err
	}
	err = transaction.SignTX(uint64(chainIdAsInt), svc.chainConfig.GetCurve(), sk)
	if err != nil {
		log.Error().Err(err)
		return nil, err
	}

	hash, err := svc.httpApi.SendSignedTransaction(ctx, chainId, transaction)
	if err != nil {
		log.Error().Err(err)
		return nil, err
	} else {
		latestBlock.Hash = *hash
		latestBlock.IncrHeight()
		if err := svc.blockCache.SetBlock(chainId, credentials.AccountAddress, latestBlock); err != nil {
			log.Error().Err(err)
		}
	}
	log.Debug().Msgf("结束调用合约，哈希为：%s", hash.String())
	return hash, nil
}

func (svc *lattice) waitReceipt(ctx context.Context, chainId string, hash *common.Hash, retryStrategy *RetryStrategy) (*common.Hash, *types.Receipt, error) {
	var err error
	var receipt *types.Receipt
	err = retry.Do(
		func() error {
			receipt, err = svc.httpApi.GetReceipt(ctx, chainId, hash.String())
			if err != nil {
				log.Error().Err(err)
				return err
			}
			return nil
		},
		retryStrategy.GetRetryStrategyOpts()...,
	)

	if err != nil {
		log.Error().Err(err)
		return hash, nil, err
	}
	return hash, receipt, nil
}

func (svc *lattice) TransferWaitReceipt(ctx context.Context, credentials *Credentials, chainId, linker, payload string, amount, joule uint64, retryStrategy *RetryStrategy) (*common.Hash, *types.Receipt, error) {
	hash, err := svc.Transfer(ctx, credentials, chainId, linker, payload, amount, joule)
	if err != nil {
		return nil, nil, err
	}

	return svc.waitReceipt(ctx, chainId, hash, retryStrategy)
}

func (svc *lattice) DeployContractWaitReceipt(ctx context.Context, credentials *Credentials, chainId, data, payload string, amount, joule uint64, retryStrategy *RetryStrategy) (*common.Hash, *types.Receipt, error) {
	hash, err := svc.DeployContract(ctx, credentials, chainId, data, payload, amount, joule)
	if err != nil {
		return nil, nil, err
	}

	return svc.waitReceipt(ctx, chainId, hash, retryStrategy)
}

func (svc *lattice) CallContractWaitReceipt(ctx context.Context, credentials *Credentials, chainId, contractAddress, data, payload string, amount, joule uint64, retryStrategy *RetryStrategy) (*common.Hash, *types.Receipt, error) {
	hash, err := svc.CallContract(ctx, credentials, chainId, contractAddress, data, payload, amount, joule)
	if err != nil {
		return nil, nil, err
	}

	return svc.waitReceipt(ctx, chainId, hash, retryStrategy)
}

func (svc *lattice) PreCallContract(ctx context.Context, chainId, owner, contractAddress, data, payload string) (*types.Receipt, error) {
	transaction := block.NewTransactionBuilder(block.TransactionTypeCallContract).
		SetLatestBlock(
			&types.LatestBlock{
				Height:          0,
				Hash:            common.HexToHash(constant.ZeroHash),
				DaemonBlockHash: common.HexToHash(constant.ZeroHash),
			}).
		SetOwner(owner).
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

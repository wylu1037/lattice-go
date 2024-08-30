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
	httpProtocol              = "http"
	httpsProtocol             = "https"
	websocketProtocol         = "ws"
	defaultHttpRequestTimeout = time.Second * 15
)

// NewLattice 初始化LatticeApi
//
// Parameters:
//   - chainConfig *ChainConfig: 链配置信息
//   - connectingNodeConfig *ConnectingNodeConfig: 节点的连接信息
//   - blockCache BlockCache: 区块缓存接口，通过缓存支持账户高并发发交易，为nil时，禁用缓存，或着使用内置的 lattice.NewMemoryBlockCache(10*time.Second, time.Minute, time.Minute)
//   - accountLock AccountLock: 账户锁接口，通过账户锁支持账户高并发发交易，为nil时，默认使用 lattice.NewAccountLock()
//   - options *Options:
//
// Returns:
//   - Lattice
func NewLattice(chainConfig *ChainConfig, connectingNodeConfig *ConnectingNodeConfig, blockCache BlockCache, accountLock AccountLock, options *Options) Lattice {
	if err := chainConfig.validate(); err != nil {
		panic(err)
	}
	if err := connectingNodeConfig.validate(); err != nil {
		panic(err)
	}

	initHttpClientArgs := &client.HttpApiInitParam{
		HttpUrl:                    connectingNodeConfig.GetHttpUrl(),
		GinServerUrl:               connectingNodeConfig.GetGinServerUrl(),
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
	httpApi              client.HttpApi        // 节点的http客户端
	chainConfig          *ChainConfig          // 链信息配置
	connectingNodeConfig *ConnectingNodeConfig // 节点的连接信息配置
	blockCache           BlockCache            // 区块缓存接口
	accountLock          AccountLock           // 账户锁接口
	options              *Options              // 可选配置
}

// ChainConfig 链配置
type ChainConfig struct {
	Curve     types.Curve // crypto.Secp256k1 or crypto.Sm2p256v1
	TokenLess bool        // false:有通证链，true:无通证链
}

// 验证链配置信息是否有效
func (chain *ChainConfig) validate() error {
	if chain.Curve == "" {
		return fmt.Errorf("ChainConfig未指定Curve参数")
	}
	return nil
}

// ConnectingNodeConfig 节点配置
type ConnectingNodeConfig struct {
	Insecure                   bool
	Ip                         string
	HttpPort                   uint16
	WebsocketPort              uint16
	GinHttpPort                uint16
	JwtSecret                  string
	JwtTokenExpirationDuration time.Duration
}

// 验证节点的连接信息是否有效
func (node *ConnectingNodeConfig) validate() error {
	if node.Ip == "" {
		return fmt.Errorf("节点的IP信息不能为空")
	}
	if node.HttpPort == 0 {
		return fmt.Errorf("节点的HttpPort信息不能为空")
	}
	return nil
}

// Credentials 凭证配置
type Credentials struct {
	AccountAddress string // 账户地址
	Passphrase     string // 身份密码
	FileKey        string // FileKey 的json字符串
	PrivateKey     string // 私钥
}

type Options struct {
	Transport *http.Transport // http连接的transport配置

	InsecureSkipVerify bool // 是否跳过https安全验证

	// MaxIdleConns controls the maximum number of idle (keep-alive)
	// connections across all hosts. Zero means no limit.
	MaxIdleConns int

	// MaxIdleConnsPerHost if non-zero, controls the maximum idle
	// (keep-alive) connections to keep per-host.
	// If zero, DefaultMaxIdleConnsPerHost(2) is used.
	MaxIdleConnsPerHost int
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

func (node *ConnectingNodeConfig) GetGinServerUrl() string {
	port := lo.Ternary(node.GinHttpPort == 0, node.HttpPort+2, node.GinHttpPort)
	return fmt.Sprintf("%s://%s:%d", lo.Ternary(node.Insecure, httpsProtocol, httpProtocol), node.Ip, port)
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
		Delay:    time.Millisecond * 150,
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
		Attempts:  15,
		Delay:     time.Millisecond * 150,
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

	// UpgradeContract 发起升级合约交易
	//
	// Parameters:
	//   - ctx context.Context
	//   - credentials *Credentials: 发交易的身份凭证
	//   - chainId string: 链ID
	//   - contractAddress string: 要升级的合约地址
	//   - data string: 升级的合约代码
	//   - payload string: 交易备注
	//
	// Returns:
	//   - *common.Hash: 交易哈希
	//   - error
	UpgradeContract(ctx context.Context, credentials *Credentials, chainId, contractAddress, data, payload string, amount, joule uint64) (*common.Hash, error)

	// UpgradeContractWaitReceipt 发起升级合约交易并等待回执
	//
	// Parameters:
	//   - ctx context.Context
	//   - credentials *Credentials: 发交易的身份凭证
	//   - chainId string: 链ID
	//   - contractAddress string: 要升级的合约地址
	//   - data string：升级的合约代码
	//   - payload string: 交易备注
	//   - retryStrategy *RetryStrategy: 等待回执策略
	//
	// Returns:
	//   - *common.Hash: 交易哈希
	//   - *types.Receipt: 回执
	//   - error
	UpgradeContractWaitReceipt(ctx context.Context, credentials *Credentials, chainId, contractAddress, data, payload string, amount, joule uint64, retryStrategy *RetryStrategy) (*common.Hash, *types.Receipt, error)

	// DeployGoContract 部署GO合约
	//
	// Parameters:
	//   - ctx context.Context
	//   - credentials *Credentials: 部署合约的凭证
	//   - chainId string: 要部署到的链(通道)ID
	//   - data types.DeployMultilingualContractCode
	//   - payload string: 交易备注，16进制带0x前缀的字符串
	//   - amount uint64: 转账额度
	//   - joule uint64: 部署GO合约的手续费
	//
	// Returns:
	//   - *common.Hash: 部署GO合约的交易哈希
	//   - error: 部署GO合约时的错误
	DeployGoContract(ctx context.Context, credentials *Credentials, chainId string, data types.DeployMultilingualContractCode, payload string, amount, joule uint64) (*common.Hash, error)

	UpgradeGoContract(ctx context.Context, credentials *Credentials, chainId, contractAddress string, data types.UpgradeMultilingualContractCode, payload string, amount, joule uint64) (*common.Hash, error)

	CallGoContract(ctx context.Context, credentials *Credentials, chainId, contractAddress string, data types.CallMultilingualContractCode, payload string, amount, joule uint64) (*common.Hash, error)

	// DeployJavaContract 部署JAVA合约
	//
	// Parameters:
	//   - ctx context.Context
	//   - credentials *Credentials: 部署合约的凭证
	//   - chainId string: 要部署到的链(通道)ID
	//   - data types.DeployMultilingualContractCode
	//   - payload string: 交易备注，16进制带0x前缀的字符串
	//   - amount uint64: 转账额度
	//   - joule uint64: 部署JAVA合约的手续费
	//
	// Returns:
	//   - *common.Hash: 部署JAVA合约的交易哈希
	//   - error: 部署JAVA合约时的错误
	DeployJavaContract(ctx context.Context, credentials *Credentials, chainId string, data types.DeployMultilingualContractCode, payload string, amount, joule uint64) (*common.Hash, error)

	UpgradeJavaContract(ctx context.Context, credentials *Credentials, chainId, contractAddress string, data types.UpgradeMultilingualContractCode, payload string, amount, joule uint64) (*common.Hash, error)

	CallJavaContract(ctx context.Context, credentials *Credentials, chainId, contractAddress string, data types.CallMultilingualContractCode, payload string, amount, joule uint64) (*common.Hash, error)

	DeployGoContractWaitReceipt(ctx context.Context, credentials *Credentials, chainId string, data types.DeployMultilingualContractCode, payload string, amount, joule uint64, retryStrategy *RetryStrategy) (*common.Hash, *types.Receipt, error)

	UpgradeGoContractWaitReceipt(ctx context.Context, credentials *Credentials, chainId, contractAddress string, data types.UpgradeMultilingualContractCode, payload string, amount, joule uint64, retryStrategy *RetryStrategy) (*common.Hash, *types.Receipt, error)

	CallGoContractWaitReceipt(ctx context.Context, credentials *Credentials, chainId, contractAddress string, data types.CallMultilingualContractCode, payload string, amount, joule uint64, retryStrategy *RetryStrategy) (*common.Hash, *types.Receipt, error)

	DeployJavaContractWaitReceipt(ctx context.Context, credentials *Credentials, chainId string, data types.DeployMultilingualContractCode, payload string, amount, joule uint64, retryStrategy *RetryStrategy) (*common.Hash, *types.Receipt, error)

	UpgradeJavaContractWaitReceipt(ctx context.Context, credentials *Credentials, chainId, contractAddress string, data types.UpgradeMultilingualContractCode, payload string, amount, joule uint64, retryStrategy *RetryStrategy) (*common.Hash, *types.Receipt, error)

	CallJavaContractWaitReceipt(ctx context.Context, credentials *Credentials, chainId, contractAddress string, data types.CallMultilingualContractCode, payload string, amount, joule uint64, retryStrategy *RetryStrategy) (*common.Hash, *types.Receipt, error)
}

func (svc *lattice) HttpApi() client.HttpApi {
	return svc.httpApi
}

// Start handle transaction, contains
// 1.Sign transaction,
// 2.Send transaction to the chain.
func (svc *lattice) handleTransaction(ctx context.Context, credentials *Credentials, chainId string, transaction *block.Transaction, latestBlock *types.LatestBlock) (*common.Hash, error) {
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
	err = transaction.SignTX(uint64(chainIdAsInt), svc.chainConfig.Curve, sk)
	if err != nil {
		log.Error().Err(err)
		return nil, err
	}

	cancelCtx, cancelFunc := context.WithTimeout(ctx, defaultHttpRequestTimeout)
	defer cancelFunc()
	hash, err := svc.httpApi.SendSignedTransaction(cancelCtx, chainId, transaction)
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
	return hash, nil
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

	hash, err := svc.handleTransaction(ctx, credentials, chainId, transaction, latestBlock)
	if err != nil {
		return nil, err
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

	hash, err := svc.handleTransaction(ctx, credentials, chainId, transaction, latestBlock)
	if err != nil {
		return nil, err
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

	hash, err := svc.handleTransaction(ctx, credentials, chainId, transaction, latestBlock)
	if err != nil {
		return nil, err
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
	log.Debug().Msgf("开始发起预调用合约交易，chainId: %s, owner: %s, contractAddress: %s, data: %s, payload: %s", chainId, owner, contractAddress, data, payload)

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
	log.Debug().Msgf("结束预调用合约，回执为：%+v", receipt)
	return receipt, nil
}

func (svc *lattice) UpgradeContract(ctx context.Context, credentials *Credentials, chainId, contractAddress, data, payload string, amount, joule uint64) (*common.Hash, error) {
	log.Debug().Msgf("开始发起升级合约交易，chainId: %s, contractAddress: %s, data: %s, payload: %s, amount: %d, joule: %d", chainId, contractAddress, data, payload, amount, joule)

	svc.accountLock.Obtain(chainId, credentials.AccountAddress)
	defer svc.accountLock.Unlock(chainId, credentials.AccountAddress)

	latestBlock, err := svc.blockCache.GetBlock(chainId, credentials.AccountAddress)
	if err != nil {
		return nil, err
	}

	transaction := block.NewTransactionBuilder(block.TransactionTypeUpgradeContract).
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

	hash, err := svc.handleTransaction(ctx, credentials, chainId, transaction, latestBlock)
	if err != nil {
		return nil, err
	}
	log.Debug().Msgf("结束升级合约，哈希为：%s", hash.String())
	return hash, nil
}

func (svc *lattice) UpgradeContractWaitReceipt(ctx context.Context, credentials *Credentials, chainId, contractAddress, data, payload string, amount, joule uint64, retryStrategy *RetryStrategy) (*common.Hash, *types.Receipt, error) {
	hash, err := svc.UpgradeContract(ctx, credentials, chainId, contractAddress, data, payload, amount, joule)
	if err != nil {
		return nil, nil, err
	}

	return svc.waitReceipt(ctx, chainId, hash, retryStrategy)
}

func (svc *lattice) deployMultilingualContract(ctx context.Context, credentials *Credentials, chainId string, lang types.ContractLang, data types.DeployMultilingualContractCode, payload string, amount, joule uint64) (*common.Hash, error) {
	log.Debug().Msgf("开始发起部署%s合约交易，chainId: %s, data: %+v, payload: %s, amount: %d, joule: %d", lang, chainId, data, payload, amount, joule)

	svc.accountLock.Obtain(chainId, credentials.AccountAddress)
	defer svc.accountLock.Unlock(chainId, credentials.AccountAddress)

	latestBlock, err := svc.blockCache.GetBlock(chainId, credentials.AccountAddress)
	if err != nil {
		log.Error().Err(err)
		return nil, err
	}

	var transactionType block.TransactionType
	switch lang {
	case types.ContractLangGo:
		transactionType = block.TransactionTypeDeployGoContract
	case types.ContractLangJava:
		transactionType = block.TransactionTypeDeployJavaContract
	default:
	}
	code := data.Encode()
	transaction := block.NewTransactionBuilder(transactionType).
		SetLatestBlock(latestBlock).
		SetOwner(credentials.AccountAddress).
		SetLinker(constant.ZeroAddress).
		SetCode(code).
		SetPayload(payload).
		SetAmount(amount).
		SetJoule(joule).
		Build()

	cryptoInstance := crypto.NewCrypto(svc.chainConfig.Curve)
	dataHash := cryptoInstance.Hash(hexutil.MustDecode(code))
	transaction.CodeHash = dataHash

	hash, err := svc.handleTransaction(ctx, credentials, chainId, transaction, latestBlock)
	if err != nil {
		return nil, err
	}
	log.Debug().Msgf("结束部署%s合约，哈希为：%s", lang, hash.String())
	return hash, nil
}

func (svc *lattice) upgradeMultilingualContract(ctx context.Context, credentials *Credentials, chainId, contractAddress string, lang types.ContractLang, data types.UpgradeMultilingualContractCode, payload string, amount, joule uint64) (*common.Hash, error) {
	log.Debug().Msgf("开始发起升级%s合约交易，chainId: %s, contractAddress: %s, data: %s, payload: %s, amount: %d, joule: %d", lang, chainId, contractAddress, data, payload, amount, joule)

	svc.accountLock.Obtain(chainId, credentials.AccountAddress)
	defer svc.accountLock.Unlock(chainId, credentials.AccountAddress)

	latestBlock, err := svc.blockCache.GetBlock(chainId, credentials.AccountAddress)
	if err != nil {
		return nil, err
	}

	var transactionType block.TransactionType
	switch lang {
	case types.ContractLangGo:
		transactionType = block.TransactionTypeUpgradeGoContract
	case types.ContractLangJava:
		transactionType = block.TransactionTypeUpgradeJavaContract
	default:
	}
	code := data.Encode()
	transaction := block.NewTransactionBuilder(transactionType).
		SetLatestBlock(latestBlock).
		SetOwner(credentials.AccountAddress).
		SetLinker(contractAddress).
		SetCode(code).
		SetPayload(payload).
		SetAmount(amount).
		SetJoule(joule).
		Build()

	cryptoInstance := crypto.NewCrypto(svc.chainConfig.Curve)
	dataHash := cryptoInstance.Hash(hexutil.MustDecode(code))
	transaction.CodeHash = dataHash

	hash, err := svc.handleTransaction(ctx, credentials, chainId, transaction, latestBlock)
	if err != nil {
		return nil, err
	}
	log.Debug().Msgf("结束升级%s合约，哈希为：%s", lang, hash.String())
	return hash, nil
}

func (svc *lattice) callMultilingualContract(ctx context.Context, credentials *Credentials, chainId, contractAddress string, lang types.ContractLang, data types.CallMultilingualContractCode, payload string, amount, joule uint64) (*common.Hash, error) {
	log.Debug().Msgf("开始发起调用%s合约交易，chainId: %s, contractAddress: %s, data: %s, payload: %s, amount: %d, joule: %d", lang, chainId, contractAddress, data, payload, amount, joule)

	svc.accountLock.Obtain(chainId, credentials.AccountAddress)
	defer svc.accountLock.Unlock(chainId, credentials.AccountAddress)

	latestBlock, err := svc.blockCache.GetBlock(chainId, credentials.AccountAddress)
	if err != nil {
		return nil, err
	}

	var transactionType block.TransactionType
	switch lang {
	case types.ContractLangGo:
		transactionType = block.TransactionTypeCallGoContract
	case types.ContractLangJava:
		transactionType = block.TransactionTypeCallJavaContract
	default:
	}
	code := data.Encode()
	transaction := block.NewTransactionBuilder(transactionType).
		SetLatestBlock(latestBlock).
		SetOwner(credentials.AccountAddress).
		SetLinker(contractAddress).
		SetCode(code).
		SetPayload(payload).
		SetAmount(amount).
		SetJoule(joule).
		Build()

	cryptoInstance := crypto.NewCrypto(svc.chainConfig.Curve)
	dataHash := cryptoInstance.Hash(hexutil.MustDecode(code))
	transaction.CodeHash = dataHash

	hash, err := svc.handleTransaction(ctx, credentials, chainId, transaction, latestBlock)
	if err != nil {
		return nil, err
	}
	log.Debug().Msgf("结束调用%s合约，哈希为：%s", lang, hash.String())
	return hash, nil
}

func (svc *lattice) DeployGoContract(ctx context.Context, credentials *Credentials, chainId string, data types.DeployMultilingualContractCode, payload string, amount, joule uint64) (*common.Hash, error) {
	return svc.deployMultilingualContract(ctx, credentials, chainId, types.ContractLangGo, data, payload, amount, joule)
}

func (svc *lattice) UpgradeGoContract(ctx context.Context, credentials *Credentials, chainId, contractAddress string, data types.UpgradeMultilingualContractCode, payload string, amount, joule uint64) (*common.Hash, error) {
	return svc.upgradeMultilingualContract(ctx, credentials, chainId, contractAddress, types.ContractLangGo, data, payload, amount, joule)
}

func (svc *lattice) CallGoContract(ctx context.Context, credentials *Credentials, chainId, contractAddress string, data types.CallMultilingualContractCode, payload string, amount, joule uint64) (*common.Hash, error) {
	return svc.callMultilingualContract(ctx, credentials, chainId, contractAddress, types.ContractLangGo, data, payload, amount, joule)
}

func (svc *lattice) DeployJavaContract(ctx context.Context, credentials *Credentials, chainId string, data types.DeployMultilingualContractCode, payload string, amount, joule uint64) (*common.Hash, error) {
	return svc.deployMultilingualContract(ctx, credentials, chainId, types.ContractLangJava, data, payload, amount, joule)
}

func (svc *lattice) UpgradeJavaContract(ctx context.Context, credentials *Credentials, chainId, contractAddress string, data types.UpgradeMultilingualContractCode, payload string, amount, joule uint64) (*common.Hash, error) {
	return svc.upgradeMultilingualContract(ctx, credentials, chainId, contractAddress, types.ContractLangJava, data, payload, amount, joule)
}

func (svc *lattice) CallJavaContract(ctx context.Context, credentials *Credentials, chainId, contractAddress string, data types.CallMultilingualContractCode, payload string, amount, joule uint64) (*common.Hash, error) {
	return svc.callMultilingualContract(ctx, credentials, chainId, contractAddress, types.ContractLangJava, data, payload, amount, joule)
}

func (svc *lattice) DeployGoContractWaitReceipt(ctx context.Context, credentials *Credentials, chainId string, data types.DeployMultilingualContractCode, payload string, amount, joule uint64, retryStrategy *RetryStrategy) (*common.Hash, *types.Receipt, error) {
	hash, err := svc.DeployGoContract(ctx, credentials, chainId, data, payload, amount, joule)
	if err != nil {
		return nil, nil, err
	}

	return svc.waitReceipt(ctx, chainId, hash, retryStrategy)
}

func (svc *lattice) UpgradeGoContractWaitReceipt(ctx context.Context, credentials *Credentials, chainId, contractAddress string, data types.UpgradeMultilingualContractCode, payload string, amount, joule uint64, retryStrategy *RetryStrategy) (*common.Hash, *types.Receipt, error) {
	hash, err := svc.UpgradeGoContract(ctx, credentials, chainId, contractAddress, data, payload, amount, joule)
	if err != nil {
		return nil, nil, err
	}

	return svc.waitReceipt(ctx, chainId, hash, retryStrategy)
}

func (svc *lattice) CallGoContractWaitReceipt(ctx context.Context, credentials *Credentials, chainId, contractAddress string, data types.CallMultilingualContractCode, payload string, amount, joule uint64, retryStrategy *RetryStrategy) (*common.Hash, *types.Receipt, error) {
	hash, err := svc.CallGoContract(ctx, credentials, chainId, contractAddress, data, payload, amount, joule)
	if err != nil {
		return nil, nil, err
	}

	return svc.waitReceipt(ctx, chainId, hash, retryStrategy)
}

func (svc *lattice) DeployJavaContractWaitReceipt(ctx context.Context, credentials *Credentials, chainId string, data types.DeployMultilingualContractCode, payload string, amount, joule uint64, retryStrategy *RetryStrategy) (*common.Hash, *types.Receipt, error) {
	hash, err := svc.DeployJavaContract(ctx, credentials, chainId, data, payload, amount, joule)
	if err != nil {
		return nil, nil, err
	}

	return svc.waitReceipt(ctx, chainId, hash, retryStrategy)
}

func (svc *lattice) UpgradeJavaContractWaitReceipt(ctx context.Context, credentials *Credentials, chainId, contractAddress string, data types.UpgradeMultilingualContractCode, payload string, amount, joule uint64, retryStrategy *RetryStrategy) (*common.Hash, *types.Receipt, error) {
	hash, err := svc.UpgradeJavaContract(ctx, credentials, chainId, contractAddress, data, payload, amount, joule)
	if err != nil {
		return nil, nil, err
	}

	return svc.waitReceipt(ctx, chainId, hash, retryStrategy)
}

func (svc *lattice) CallJavaContractWaitReceipt(ctx context.Context, credentials *Credentials, chainId, contractAddress string, data types.CallMultilingualContractCode, payload string, amount, joule uint64, retryStrategy *RetryStrategy) (*common.Hash, *types.Receipt, error) {
	hash, err := svc.CallJavaContract(ctx, credentials, chainId, contractAddress, data, payload, amount, joule)
	if err != nil {
		return nil, nil, err
	}

	return svc.waitReceipt(ctx, chainId, hash, retryStrategy)
}

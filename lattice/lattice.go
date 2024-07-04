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
)

func NewLattice(chainConfig *ChainConfig, nodeConfig *NodeConfig, identityConfig *IdentityConfig, options *Options) Lattice {
	return &lattice{
		ChainConfig:    chainConfig,
		NodeConfig:     nodeConfig,
		IdentityConfig: identityConfig,
		Options:        options,
		httpApi:        client.NewHttpApi(nodeConfig.GetHttpUrl(), strconv.FormatUint(chainConfig.ChainId, 10), options.GetTransport()),
	}
}

type lattice struct {
	httpApi        client.HttpApi
	ChainConfig    *ChainConfig
	NodeConfig     *NodeConfig
	IdentityConfig *IdentityConfig
	Options        *Options
}

type ChainConfig struct {
	ChainId uint64
	Curve   types.Curve
}

type NodeConfig struct {
	Insecure      bool
	Ip            string
	HttpPort      uint16
	WebsocketPort uint16
}

type IdentityConfig struct {
	AccountAddress string
	Passphrase     string
	PrivateKey     string
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

func (identity *IdentityConfig) GetSK() string {
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

// WaitStrategy 等待回执策略
type WaitStrategy struct {
	// 具体的策略
	Strategy  Strategy
	Attempts  uint
	Delay     time.Duration
	MaxJitter time.Duration
}

func (strategy *WaitStrategy) GetWaitStrategyOpts() []retry.Option {
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

func NewBackOffWaitStrategy(attempts uint, initDelay time.Duration) *WaitStrategy {
	return &WaitStrategy{
		Strategy: BackOff,
		Attempts: attempts,
		Delay:    initDelay,
	}
}

func NewFixedWaitStrategy(attempts uint, fixedDelay time.Duration) *WaitStrategy {
	return &WaitStrategy{
		Strategy: FixedInterval,
		Attempts: attempts,
		Delay:    fixedDelay,
	}
}

func NewRandomWaitStrategy(attempts uint, baseDelay time.Duration, maxJitter time.Duration) *WaitStrategy {
	return &WaitStrategy{
		Strategy:  RandomInterval,
		Attempts:  attempts,
		Delay:     baseDelay,
		MaxJitter: maxJitter,
	}
}

func (strategy *WaitStrategy) BackOffOpts() []retry.Option {
	return []retry.Option{retry.Attempts(strategy.Attempts), retry.Delay(strategy.Delay), retry.DelayType(retry.BackOffDelay)}
}

func (strategy *WaitStrategy) FixedIntervalOpts() []retry.Option {
	return []retry.Option{retry.Attempts(strategy.Attempts), retry.Delay(strategy.Delay), retry.DelayType(retry.FixedDelay)}
}

func (strategy *WaitStrategy) RandomIntervalOpts() []retry.Option {
	return []retry.Option{retry.Attempts(strategy.Attempts), retry.Delay(strategy.Delay), retry.MaxJitter(strategy.MaxJitter), retry.DelayType(retry.RandomDelay)}
}

type Lattice interface {
	// Transfer 发起转账交易
	//
	// Parameters:
	//    - ctx context.Context
	//    - linker string: 转账接收者账户地址
	//    - payload string: 交易备注
	//
	// Returns:
	//    - common.Hash: 交易哈希
	//    - error
	Transfer(ctx context.Context, linker, payload string) (*common.Hash, error)
	DeployContract(ctx context.Context, data, payload string) (*common.Hash, error)
	CallContract(ctx context.Context, contractAddress, data, payload string) (*common.Hash, error)
	TransferWaitReceipt(ctx context.Context, linker, payload string, waitStrategy *WaitStrategy) (*common.Hash, *types.Receipt, error)
	DeployContractWaitReceipt(ctx context.Context, data, payload string, waitStrategy *WaitStrategy) (*common.Hash, *types.Receipt, error)
	CallContractWaitReceipt(ctx context.Context, contractAddress, data, payload string, waitStrategy *WaitStrategy) (*common.Hash, *types.Receipt, error)
}

func (svc *lattice) Transfer(ctx context.Context, linker, payload string) (*common.Hash, error) {
	latestBlock, err := svc.httpApi.GetLatestBlock(ctx, svc.IdentityConfig.AccountAddress)
	if err != nil {
		return nil, err
	}

	transaction := block.NewTransactionBuilder(block.TransactionTypeSend).
		SetLatestBlock(latestBlock).
		SetOwner(svc.IdentityConfig.AccountAddress).
		SetLinker(linker).
		SetPayload(payload).
		Build()

	err = transaction.SignTX(svc.ChainConfig.ChainId, svc.ChainConfig.GetCurve(), svc.IdentityConfig.GetSK())
	if err != nil {
		return nil, err
	}

	hash, err := svc.httpApi.SendSignedTransaction(ctx, transaction)
	if err != nil {
		return nil, err
	}
	return hash, nil
}

func (svc *lattice) DeployContract(ctx context.Context, data, payload string) (*common.Hash, error) {
	latestBlock, err := svc.httpApi.GetLatestBlock(ctx, svc.IdentityConfig.AccountAddress)
	if err != nil {
		return nil, err
	}

	transaction := block.NewTransactionBuilder(block.TransactionTypeDeployContract).
		SetLatestBlock(latestBlock).
		SetOwner(svc.IdentityConfig.AccountAddress).
		SetLinker("zltc_QLbz7JHiBTspS962RLKV8GndWFwjA5K66").
		SetCode(data).
		SetPayload(payload).
		Build()

	cryptoInstance := crypto.NewCrypto(svc.ChainConfig.Curve)
	dataHash := cryptoInstance.Hash(hexutil.MustDecode(data))
	transaction.CodeHash = dataHash

	err = transaction.SignTX(svc.ChainConfig.ChainId, svc.ChainConfig.GetCurve(), svc.IdentityConfig.GetSK())
	if err != nil {
		return nil, err
	}

	hash, err := svc.httpApi.SendSignedTransaction(ctx, transaction)
	if err != nil {
		return nil, err
	}
	return hash, nil
}

func (svc *lattice) CallContract(ctx context.Context, contractAddress, data, payload string) (*common.Hash, error) {
	latestBlock, err := svc.httpApi.GetLatestBlock(ctx, svc.IdentityConfig.AccountAddress)
	if err != nil {
		return nil, err
	}

	transaction := block.NewTransactionBuilder(block.TransactionTypeCallContract).
		SetLatestBlock(latestBlock).
		SetOwner(svc.IdentityConfig.AccountAddress).
		SetLinker(contractAddress).
		SetCode(data).
		SetPayload(payload).
		Build()

	cryptoInstance := crypto.NewCrypto(svc.ChainConfig.Curve)
	dataHash := cryptoInstance.Hash(hexutil.MustDecode(data))
	transaction.CodeHash = dataHash

	err = transaction.SignTX(svc.ChainConfig.ChainId, svc.ChainConfig.GetCurve(), svc.IdentityConfig.GetSK())
	if err != nil {
		return nil, err
	}

	hash, err := svc.httpApi.SendSignedTransaction(ctx, transaction)
	if err != nil {
		return nil, err
	}
	return hash, nil
}

func (svc *lattice) waitReceipt(ctx context.Context, hash *common.Hash, waitStrategy *WaitStrategy) (*common.Hash, *types.Receipt, error) {
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
		waitStrategy.GetWaitStrategyOpts()...,
	)

	if err != nil {
		return hash, nil, err
	}
	return hash, receipt, nil
}

func (svc *lattice) TransferWaitReceipt(ctx context.Context, linker, payload string, waitStrategy *WaitStrategy) (*common.Hash, *types.Receipt, error) {
	hash, err := svc.Transfer(ctx, linker, payload)
	if err != nil {
		return nil, nil, err
	}

	return svc.waitReceipt(ctx, hash, waitStrategy)
}

func (svc *lattice) DeployContractWaitReceipt(ctx context.Context, data, payload string, waitStrategy *WaitStrategy) (*common.Hash, *types.Receipt, error) {
	hash, err := svc.DeployContract(ctx, data, payload)
	if err != nil {
		return nil, nil, err
	}

	return svc.waitReceipt(ctx, hash, waitStrategy)
}

func (svc *lattice) CallContractWaitReceipt(ctx context.Context, contractAddress, data, payload string, waitStrategy *WaitStrategy) (*common.Hash, *types.Receipt, error) {
	hash, err := svc.CallContract(ctx, contractAddress, data, payload)
	if err != nil {
		return nil, nil, err
	}

	return svc.waitReceipt(ctx, hash, waitStrategy)
}

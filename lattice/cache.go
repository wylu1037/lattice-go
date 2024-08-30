package lattice

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/allegro/bigcache/v3"
	"github.com/rs/zerolog/log"
	"github.com/wylu1037/lattice-go/common/types"
	"github.com/wylu1037/lattice-go/lattice/client"
	"sync"
	"time"
)

// NewMemoryBlockCache 初始化一个内存缓存
//
// Parameters:
//   - enable bool: 是否启用缓存
//   - httpApi client.HttpApi
//   - daemonHashExpirationDuration time.Duration: 守护区块哈希的过期时长
//   - lifeDuration time.Duration: 缓存的存活时长
//   - cleanInterval time.Duration: 过期缓存的清理间隔时长
//
// Returns:
//   - BlockCache
func NewMemoryBlockCache(daemonHashExpirationDuration time.Duration, lifeDuration time.Duration, cleanInterval time.Duration) BlockCache {
	memoryCacheApi, err := bigcache.New(context.Background(), newMemoryCacheConfig(lifeDuration, cleanInterval))
	if err != nil {
		panic(err)
	}
	return &memoryBlockCache{
		enable:                       true,
		memoryCacheApi:               memoryCacheApi,
		daemonHashExpirationDuration: daemonHashExpirationDuration,
	}
}

func newDisabledMemoryBlockCache(httpApi client.HttpApi) BlockCache {
	return &memoryBlockCache{
		enable:  false,
		httpApi: httpApi,
	}
}

func newMemoryCacheConfig(lifeDuration time.Duration, cleanInterval time.Duration) bigcache.Config {
	config := bigcache.Config{
		// number of shards (must be a power of 2)
		Shards: 1024,

		// time after which entry can be evicted
		// 缓存的过期时间
		LifeWindow: lifeDuration,

		// Interval between removing expired entries (clean up).
		// If set to <= 0 then no action is performed.
		// Setting to < 1 second is counterproductive — bigcache has a one second resolution.
		// 清理过期条目的时间间隔
		CleanWindow: cleanInterval,

		// rps * lifeWindow, used only in initial memory allocation
		// 在 LifeWindow 时间内可能的最大条目数，支持1024个账户的缓存
		MaxEntriesInWindow: 1024,

		// max entry size in bytes, used only in initial memory allocation
		// 每个条目的最大字节数，512byte = 0.5KB
		MaxEntrySize: 512,

		// prints information about additional memory allocation
		// 是否打印内存分配的详细信息
		Verbose: true,

		// cache will not allocate more memory than this limit, value in MB
		// if value is reached then the oldest entries can be overridden for the new ones
		// 0 value means no size limit
		// 缓存系统的最大内存限制，以MB为单位，达到设置的上限时，新条目会覆盖旧条目，最大8192
		HardMaxCacheSize: 512,

		// callback fired when the oldest entry is removed because of its expiration time or no space left
		// for the new entry, or because delete was called. A bitmask representing the reason will be returned.
		// Default value is nil which means no callback and it prevents from unwrapping the oldest entry.
		// 当最旧的条目被移除时触发的回调函数
		OnRemove: nil,

		// OnRemoveWithReason is a callback fired when the oldest entry is removed because of its expiration time or no space left
		// for the new entry, or because delete was called. A constant representing the reason will be passed through.
		// Default value is nil which means no callback and it prevents from unwrapping the oldest entry.
		// Ignored if OnRemove is specified.
		// 与 OnRemove 类似，但此回调函数会带有移除条目的原因
		OnRemoveWithReason: nil,
	}
	return config
}

type BlockCache interface {
	SetHttpApi(httpApi client.HttpApi)

	// SetBlock 设置区块缓存
	//
	// Parameters:
	//   - key string: 缓存的Key
	//   - block *types.LatestBlock: 缓存的区块
	//
	// Returns:
	//   - error
	SetBlock(chainId, address string, block *types.LatestBlock) error

	// GetBlock 获取区块缓存
	//
	// Parameters:
	//   - key string: 缓存的Key
	//
	// Returns:
	//   - *types.LatestBlock: 缓存的区块信息
	//   - error
	GetBlock(chainId, address string) (*types.LatestBlock, error)
}

// type redisBlockCache struct{}

type memoryBlockCache struct {
	enable                       bool               // 是否启用缓存
	httpApi                      client.HttpApi     // 节点的http客户端
	memoryCacheApi               *bigcache.BigCache // big cache
	daemonHashExpireAtMap        sync.Map           // 守护哈希的过期时间，每个链维护一个
	daemonHashExpirationDuration time.Duration      // 守护区块哈希的过期时长
}

func (c *memoryBlockCache) SetHttpApi(httpApi client.HttpApi) {
	c.httpApi = httpApi
}

func (c *memoryBlockCache) SetBlock(chainId, address string, block *types.LatestBlock) error {
	if !c.enable {
		return nil
	}
	bytes, err := json.Marshal(block)
	if err != nil {
		log.Error().Err(err).Msgf("json序列化block失败，chainId: %s, accountAddress: %s", chainId, address)
		return err
	}
	if err := c.memoryCacheApi.Set(fmt.Sprintf("%s_%s", chainId, address), bytes); err != nil {
		log.Error().Err(err).Msgf("设置区块缓存信息失败，chainId: %s, accountAddress: %s", chainId, address)
		return err
	}

	_, ok := c.daemonHashExpireAtMap.Load(chainId)
	if !ok {
		c.daemonHashExpireAtMap.Store(chainId, time.Now().Add(c.daemonHashExpirationDuration))
	}

	return nil
}

func (c *memoryBlockCache) GetBlock(chainId, address string) (*types.LatestBlock, error) {
	if !c.enable {
		return c.httpApi.GetLatestBlock(context.Background(), chainId, address)
	}
	// load cached block from memory
	cacheBlockBytes, err := c.memoryCacheApi.Get(fmt.Sprintf("%s_%s", chainId, address))
	if err != nil {
		if errors.Is(err, bigcache.ErrEntryNotFound) {
			return c.httpApi.GetLatestBlock(context.Background(), chainId, address)
		}
		log.Error().Err(err).Msgf("获取区块缓存信息失败，chainId: %s, accountAddress: %s", chainId, address)
		return nil, err
	}
	cacheBlock := new(types.LatestBlock)
	if err := json.Unmarshal(cacheBlockBytes, cacheBlock); err != nil {
		log.Error().Err(err).Msgf("json序列化block失败，chainId: %s, accountAddress: %s", chainId, address)
		return nil, err
	}
	// judge daemon hash expiration time
	daemonHashExpireAt, ok := c.daemonHashExpireAtMap.Load(chainId)
	if !ok {
		daemonHashExpireAt = time.Now().Add(c.daemonHashExpirationDuration)
		c.daemonHashExpireAtMap.LoadOrStore(chainId, daemonHashExpireAt)
	}
	if time.Now().After(daemonHashExpireAt.(time.Time)) {
		log.Debug().Msgf("守护区块哈希已过期，开始更新守护区块哈希，chainId: %s, accountAddress: %s", chainId, address)
		block, err := c.httpApi.GetLatestBlock(context.Background(), chainId, address)
		if err != nil {
			log.Error().Err(err).Msgf("请求节点获取最新区块信息失败，chainId: %s, accountAddress: %s", chainId, address)
			return nil, err
		}
		c.daemonHashExpireAtMap.Store(chainId, time.Now().Add(c.daemonHashExpirationDuration))
		cacheBlock.DaemonBlockHash = block.DaemonBlockHash
	}

	return cacheBlock, nil
}

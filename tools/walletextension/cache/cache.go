package cache

import (
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/log"
)

type Cache interface {
	// EvictShortLiving - notify the cache that all short living elements cached before the events should be considered as evicted.
	EvictShortLiving()

	// DisableShortLiving disables the caching of short-living elements.
	DisableShortLiving()

	// IsEvicted - based on the eviction event and the time of caching, calculates whether the key was evicted
	IsEvicted(key []byte, originalTTL time.Duration) bool

	Set(key []byte, value any, ttl time.Duration) bool
	Get(key []byte) (value any, ok bool)
	Remove(key []byte)
}

func NewCache(nrElems int, logger log.Logger) (Cache, error) {
	return NewRistrettoCacheWithEviction(nrElems, logger)
}

type Strategy uint8

const (
	NoCache     Strategy = iota
	LatestBatch Strategy = iota
	LongLiving  Strategy = iota

	longCacheTTL  = 5 * time.Hour
	shortCacheTTL = 1 * time.Minute
)

type Cfg struct {
	Type        Strategy
	DynamicType func() Strategy
}

func WithCache[R any](cache Cache, cfg *Cfg, cacheKey []byte, onCacheMiss func() (*R, error)) (*R, error) {
	if cfg == nil {
		return onCacheMiss()
	}

	cacheType := cfg.Type
	if cfg.DynamicType != nil {
		cacheType = cfg.DynamicType()
	}

	if cacheType == NoCache {
		return onCacheMiss()
	}

	// we implement a custom cache eviction logic for the cache strategy of type LatestBatch.
	// when a new batch is created, all entries with "LatestBatch" are considered evicted.
	// elements not cached for a specific batch are not evicted
	isEvicted := false
	ttl := longCacheTTL
	if cacheType == LatestBatch {
		ttl = shortCacheTTL
		isEvicted = cache.IsEvicted(cacheKey, ttl)
	}

	if !isEvicted {
		cachedValue, foundInCache := cache.Get(cacheKey)
		if foundInCache {
			returnValue, ok := cachedValue.(*R)
			if !ok {
				return nil, fmt.Errorf("unexpected error. Invalid format cached. %v", cachedValue)
			}
			return returnValue, nil
		}
	}

	result, err := onCacheMiss()

	// cache only non-nil values
	if err == nil && result != nil {
		cache.Set(cacheKey, result, ttl)
	}

	return result, err
}

type noOpCache struct{}

func NewNoOpCache() Cache {
	return &noOpCache{}
}

func (c *noOpCache) EvictShortLiving() {}

func (c *noOpCache) DisableShortLiving() {}

func (c *noOpCache) IsEvicted(key []byte, originalTTL time.Duration) bool {
	return false
}

func (c *noOpCache) Set(key []byte, value any, ttl time.Duration) bool {
	return false
}

func (c *noOpCache) Get(key []byte) (value any, ok bool) {
	return nil, false
}

func (c *noOpCache) Remove(key []byte) {}

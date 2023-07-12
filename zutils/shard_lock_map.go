package zutils

import "sync"

const (
	ShardCount = 32
	Prime      = 16777619
	HashVal    = 2166136261
)

type ShardLockMaps struct {
	shards       []*SingleShardMap
	shardKeyFunc func(key interface{}) string
}

type SingleShardMap struct {
	items map[string]interface{}
	sync.RWMutex
}

func NewShardLockMaps(shardKeyFunc func(key interface{}) string) ShardLockMaps {
	slm := ShardLockMaps{
		shards:       make([]*SingleShardMap, ShardCount),
		shardKeyFunc: shardKeyFunc,
	}

	for i := 0; i < ShardCount; i++ {
		slm.shards[i] = &SingleShardMap{items: make(map[string]interface{})}
	}
	return slm
}

func fnv32(key string) uint32 {
	hashVal := uint32(HashVal)
	prime := uint32(Prime)
	keyLength := len(key)
	for i := 0; i < keyLength; i++ {
		hashVal *= prime
		hashVal ^= uint32(key[i])
	}
	return hashVal
}

func (slm ShardLockMaps) GetShard(key string) *SingleShardMap {
	return slm.shards[fnv32(key)%ShardCount]
}

func (slm ShardLockMaps) Count() int {
	count := 0
	for i := 0; i < ShardCount; i++ {
		shard := slm.shards[i]
		shard.RLock()
		count += len(shard.items)
		shard.RUnlock()
	}
	return count
}

func (slm ShardLockMaps) Get(key string) (interface{}, bool) {
	shard := slm.GetShard(key)
	shard.RLock()
	val, ok := shard.items[key]
	shard.RUnlock()
	return val, ok
}

func (slm ShardLockMaps) Set(key string, value interface{}) {
	shard := slm.GetShard(key)
	shard.Lock()
	shard.items[key] = value
	shard.Unlock()
}

func (slm ShardLockMaps) MSet(data map[string]interface{}) {
	for key, value := range data {
		shard := slm.GetShard(key)
		shard.Lock()
		shard.items[key] = value
		shard.Unlock()
	}
}

func (slm ShardLockMaps) Has(key string) bool {
	shard := slm.GetShard(key)
	shard.RLock()
	_, ok := shard.items[key]
	shard.RUnlock()
	return ok
}

func (slm ShardLockMaps) Remove(key string) {
	shard := slm.GetShard(key)
	shard.Lock()
	delete(shard.items, key)
	shard.Unlock()
}

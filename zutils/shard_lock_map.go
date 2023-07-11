package zutils

import "sync"

const ShardCount = 32

type ShardLockMaps struct {
	shards       []*SingleShardMap
	shardKeyFunc func(key interface{}) uint32
}

type SingleShardMap struct {
	data map[interface{}]interface{}
	sync.RWMutex
}

func NewShardLockMaps(shardKeyFunc func(key interface{}) uint32) ShardLockMaps {
	slm := ShardLockMaps{
		shards:       make([]*SingleShardMap, ShardCount),
		shardKeyFunc: shardKeyFunc,
	}

	for i := 0; i < ShardCount; i++ {
		slm.shards[i] = &SingleShardMap{data: make(map[interface{}]interface{})}
	}
	return slm
}

func (slm ShardLockMaps) GetShard(key uint64) *SingleShardMap {
	return slm.shards[slm.shardKeyFunc(key)%ShardCount]
}

func (slm ShardLockMaps) Count() int {
	count := 0
	for i := 0; i < ShardCount; i++ {
		shard := slm.shards[i]
		shard.RLock()
		count += len(shard.data)
		shard.RUnlock()
	}
	return count
}

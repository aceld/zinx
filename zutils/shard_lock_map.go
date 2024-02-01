package zutils

import (
	"encoding/json"
	"sync"
)

var ShardCount = 32

// ShardLockMaps A "thread" safe map of type string:Anything.
// To avoid lock bottlenecks this map is dived to several (ShardCount) map shards.
type ShardLockMaps struct {
	shards []*SingleShardMap
	hash   IHash
}

// SingleShardMap A "thread" safe string to anything map.
type SingleShardMap struct {
	items map[string]interface{}
	sync.RWMutex
}

// createShardLockMaps Creates a new concurrent map.
func createShardLockMaps(hash IHash) ShardLockMaps {
	slm := ShardLockMaps{
		shards: make([]*SingleShardMap, ShardCount),
		hash:   hash,
	}
	for i := 0; i < ShardCount; i++ {
		slm.shards[i] = &SingleShardMap{items: make(map[string]interface{})}
	}
	return slm
}

// NewShardLockMaps Creates a new ShardLockMaps.
func NewShardLockMaps() ShardLockMaps {
	return createShardLockMaps(DefaultHash())
}

func NewWithCustomHash(hash IHash) ShardLockMaps {
	return createShardLockMaps(hash)
}

// GetShard returns shard under given key
func (slm ShardLockMaps) GetShard(key string) *SingleShardMap {
	return slm.shards[slm.hash.Sum(key)%uint32(ShardCount)]
}

// Count returns the number of elements within the map.
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

// Get retrieves an element from map under given key.
func (slm ShardLockMaps) Get(key string) (interface{}, bool) {
	shard := slm.GetShard(key)
	shard.RLock()
	val, ok := shard.items[key]
	shard.RUnlock()
	return val, ok
}

// Set Sets the given value under the specified key.
func (slm ShardLockMaps) Set(key string, value interface{}) {
	shard := slm.GetShard(key)
	shard.Lock()
	shard.items[key] = value
	shard.Unlock()
}

// SetNX Sets the given value under the specified key if no value was associated with it.
func (slm ShardLockMaps) SetNX(key string, value interface{}) bool {
	shard := slm.GetShard(key)
	shard.Lock()
	_, ok := shard.items[key]
	if !ok {
		shard.items[key] = value
	}
	shard.Unlock()
	return !ok
}

// MSet Sets the given value under the specified key.
func (slm ShardLockMaps) MSet(data map[string]interface{}) {
	for key, value := range data {
		shard := slm.GetShard(key)
		shard.Lock()
		shard.items[key] = value
		shard.Unlock()
	}
}

// Has Looks up an item under specified key
func (slm ShardLockMaps) Has(key string) bool {
	shard := slm.GetShard(key)
	shard.RLock()
	_, ok := shard.items[key]
	shard.RUnlock()
	return ok
}

// Remove removes an element from the map.
func (slm ShardLockMaps) Remove(key string) {
	shard := slm.GetShard(key)
	shard.Lock()
	delete(shard.items, key)
	shard.Unlock()
}

// RemoveCb is a callback executed in a map.RemoveCb() call, while Lock is held
// If returns true, the element will be removed from the map
type RemoveCb func(key string, v interface{}, exists bool) bool

// RemoveCb locks the shard containing the key, retrieves its current value and calls the callback with those params
// If callback returns true and element exists, it will remove it from the map
// Returns the value returned by the callback (even if element was not present in the map)
func (slm ShardLockMaps) RemoveCb(key string, cb RemoveCb) bool {

	shard := slm.GetShard(key)
	shard.Lock()
	v, ok := shard.items[key]
	remove := cb(key, v, ok)
	if remove && ok {
		delete(shard.items, key)
	}
	shard.Unlock()
	return remove
}

// Pop removes an element from the map and returns it
func (slm ShardLockMaps) Pop(key string) (v interface{}, exists bool) {
	shard := slm.GetShard(key)
	shard.Lock()
	v, exists = shard.items[key]
	delete(shard.items, key)
	shard.Unlock()
	return v, exists
}

// Clear removes all items from map.
func (slm ShardLockMaps) Clear() {
	for item := range slm.IterBuffered() {
		slm.Remove(item.Key)
	}
}

// IsEmpty checks if map is empty.
func (slm ShardLockMaps) IsEmpty() bool {
	return slm.Count() == 0
}

// Tuple Used by the IterBuffered functions to wrap two variables together over a channel,
type Tuple struct {
	Key string
	Val interface{}
}

// Returns a array of channels that contains elements in each shard,
// which likely takes a snapshot of `slm`.
// It returns once the size of each buffered channel is determined,
// before all the channels are populated using goroutines.
func snapshot(slm ShardLockMaps) (chanList []chan Tuple) {
	chanList = make([]chan Tuple, ShardCount)
	wg := sync.WaitGroup{}
	wg.Add(ShardCount)
	for index, shard := range slm.shards {
		go func(index int, shard *SingleShardMap) {
			shard.RLock()
			chanList[index] = make(chan Tuple, len(shard.items))
			wg.Done()
			for key, val := range shard.items {
				chanList[index] <- Tuple{key, val}
			}
			shard.RUnlock()
			close(chanList[index])
		}(index, shard)
	}
	wg.Wait()
	return chanList
}

// fanIn reads elements from channels `chanList` into channel `out`
func fanIn(chanList []chan Tuple, out chan Tuple) {
	wg := sync.WaitGroup{}
	wg.Add(len(chanList))
	for _, ch := range chanList {
		go func(ch chan Tuple) {
			for t := range ch {
				out <- t
			}
			wg.Done()
		}(ch)
	}
	wg.Wait()
	close(out)
}

// IterBuffered returns a buffered iterator which could be used in a for range loop.
func (slm ShardLockMaps) IterBuffered() <-chan Tuple {
	chanList := snapshot(slm)
	total := 0
	for _, c := range chanList {
		total += cap(c)
	}
	ch := make(chan Tuple, total)
	go fanIn(chanList, ch)
	return ch
}

// Items returns all items as map[string]interface{}
func (slm ShardLockMaps) Items() map[string]interface{} {
	tmp := make(map[string]interface{})

	for item := range slm.IterBuffered() {
		tmp[item.Key] = item.Val
	}

	return tmp
}

// Keys returns all keys as []string
func (slm ShardLockMaps) Keys() []string {
	count := slm.Count()
	ch := make(chan string, count)
	go func() {
		wg := sync.WaitGroup{}
		wg.Add(ShardCount)
		for _, shard := range slm.shards {
			go func(shard *SingleShardMap) {
				shard.RLock()
				for key := range shard.items {
					ch <- key
				}
				shard.RUnlock()
				wg.Done()
			}(shard)
		}
		wg.Wait()
		close(ch)
	}()

	keys := make([]string, 0, count)
	for k := range ch {
		keys = append(keys, k)
	}
	return keys
}

// IterCb Iterator callback,called for every key,value found in maps.
// RLock is held for all calls for a given shard
// therefore callback sess consistent view of a shard,
// but not across the shards
type IterCb func(key string, v interface{})

// IterCb Callback based iterator, cheapest way to read
// all elements in a map.
func (slm ShardLockMaps) IterCb(fn IterCb) {
	for idx := range slm.shards {
		shard := (slm.shards)[idx]
		shard.RLock()
		for key, value := range shard.items {
			fn(key, value)
		}
		shard.RUnlock()
	}
}

// MarshalJSON Reviles ConcurrentMap "private" variables to json marshal.
func (slm ShardLockMaps) MarshalJSON() ([]byte, error) {
	tmp := make(map[string]interface{})

	for item := range slm.IterBuffered() {
		tmp[item.Key] = item.Val
	}
	return json.Marshal(tmp)
}

// UnmarshalJSON Reverse process of Marshal.
func (slm ShardLockMaps) UnmarshalJSON(b []byte) (err error) {
	tmp := make(map[string]interface{})

	if err := json.Unmarshal(b, &tmp); err != nil {
		return err
	}

	for key, val := range tmp {
		slm.Set(key, val)
	}
	return nil
}

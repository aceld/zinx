package zutils

import (
	"fmt"
	"sync"
	"time"
)

/*
	workerID 的大小为 10 位，因此可以生成的最大 workerID 数量为 2^10 = 1024。
	sequenceBits 的大小为 12 位，因此对于每个工作器，可以在同一毫秒内生成的最大序列号数为 2^12 = 4096。
	因此，在同一毫秒内，每个工作者最多可以生成 4096 个分布式 ID。
	由于时间戳使用的是毫秒级别的时间戳，因此每个工作者在一秒钟内最多可以生成 4096 * 1000 = 4,096,000 个分布式 ID。
	总体上，如果 workerID 保持唯一，这个算法可以生成多达 1024 x 4,096,000 = 4,194,304,000 个分布式 ID。
*/

const (
	workerBits     uint8 = 10
	maxWorker      int64 = -1 ^ (-1 << workerBits)
	sequenceBits   uint8 = 12
	sequenceMask   int64 = -1 ^ (-1 << sequenceBits)
	workerShift    uint8 = sequenceBits
	timestampShift uint8 = sequenceBits + workerBits
)

type IDWorker struct {
	sequence      int64
	lastTimestamp int64
	workerId      int64
	mutex         sync.Mutex
}

func NewIDWorker(workerId int64) (*IDWorker, error) {
	if workerId < 0 || workerId > maxWorker {
		return nil, fmt.Errorf("worker ID can't be greater than %d or less than 0", maxWorker)
	}

	return &IDWorker{
		workerId:      workerId,
		lastTimestamp: -1,
		sequence:      0,
	}, nil
}

func (w *IDWorker) NextID() (int64, error) {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	timestamp := time.Now().UnixNano() / 1000000

	if timestamp < w.lastTimestamp {
		return 0, fmt.Errorf("clock moved backwards")
	}

	if timestamp == w.lastTimestamp {
		w.sequence = (w.sequence + 1) & sequenceMask
		if w.sequence == 0 {
			timestamp = w.nextMillisecond(timestamp)
		}
	} else {
		w.sequence = 0
	}

	w.lastTimestamp = timestamp

	return (timestamp << timestampShift) | (w.workerId << workerShift) | w.sequence, nil
}

func (w *IDWorker) nextMillisecond(currentTimestamp int64) int64 {
	for currentTimestamp <= w.lastTimestamp {
		currentTimestamp = time.Now().UnixNano() / 1000000
	}
	return currentTimestamp
}

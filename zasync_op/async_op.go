/*
Package zasync_op
@Author：14March
@File：async_op.go
*/
package zasync_op

import (
	"sync"
)

// 异步worker组
var asyncWorkerArray = [2048]*AsyncWorker{}
var initAsyncWorkerLocker = &sync.Mutex{}

func Process(opId int, asyncOp func()) {
	if asyncOp == nil {
		return
	}

	curWorker := getCurWorker(opId)

	if curWorker != nil {
		curWorker.process(asyncOp)
	}

}

func getCurWorker(opId int) *AsyncWorker {
	if opId < 0 {
		opId = -opId
	}

	workerIndex := opId % len(asyncWorkerArray)
	curWorker := asyncWorkerArray[workerIndex]

	if nil != curWorker {
		return curWorker
	}

	// 初始化
	initAsyncWorkerLocker.Lock()
	defer initAsyncWorkerLocker.Unlock()

	// 重新拿到这个干活的工人
	curWorker = asyncWorkerArray[workerIndex]

	// 并重新进行空指针判断
	if curWorker != nil {
		return curWorker
	}

	// 双重检查之后还是空值,进行初始化操作
	curWorker = &AsyncWorker{
		taskQ: make(chan func(), 2048),
	}

	asyncWorkerArray[workerIndex] = curWorker
	go curWorker.loopExecTask()

	return curWorker
}

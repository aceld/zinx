/*
Package zasync_op
@Author：14March
@File：async_op.go
*/
package zasync_op

import (
	"sync"
)

/*
异步IO模块：
1.业务线程执行业务操作，发送一个IO请求，由IO线程来完成写库，如果写完库之后，还有其他操作呢？
a.接下来的逻辑就在 IO 线程里执行了；
b.回到不是原来的业务线程，而是另一个业务线程执行；
这2种情况，就相当于一部分业务逻辑在 A 线程里，一部分业务逻辑在 B 线程了；两个线程同时操作一块内存区域，会出现脏读写问题。

2.因此，必须回到原本所属的业务线程里执行
意思就是说，
业务逻辑原先是由谁来执行的，那么 IO 操作完成之后，继续交还给原来的人去执行。

3.使用：
a.调用 Process 选择一个异步worker进行异步IO操作逻辑；
b.在异步IO逻辑中设置需要共享的变量，及异步返回结果：asyncResult.SetReturnedObj
c.注册设置异步回调，即回到原本的业务线程里继续进行后续的操作：asyncResult.OnComplete
*/

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

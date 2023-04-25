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
	<异步IO模块简介>

1.业务线程执行业务操作，发送一个IO请求，由IO线程来完成写库，如果写完库之后，还有其他操作呢？
	a.接下来的逻辑就在 IO 线程里执行了；
	b.回到不是原来的业务线程，而是另一个业务线程执行；
	这2种情况，就相当于一部分业务逻辑在 A 线程里，一部分业务逻辑在 B 线程了；两个线程同时操作一块内存区域，会出现脏读写问题。

2.因此，必须回到原本所属的业务线程里执行,意思就是说，业务逻辑原先是由谁来执行的，那么 IO 操作完成之后，继续交还给原来的人去执行。

3.使用：
	a.调用 Process 选择一个异步worker进行异步IO操作逻辑；
	b.在异步IO逻辑中设置需要共享的变量，及异步返回结果：asyncResult.SetReturnedObj
	c.注册设置异步回调，即回到原本的业务线程里继续进行后续的操作：asyncResult.OnComplete
*/

/*
	<Asynchronous IO Module Introduction>

1.Business threads execute business operations, send an IO request, and IO threads complete the write to the database. What if there are other operations after the write is complete?
	a. The next logic will be executed in the IO thread;
	b. Back to a different business thread instead of the original one.
    These two situations mean that some business logic is in thread A and some in thread B. When two threads operate on the same memory area, dirty reads and writes occur.

2.Therefore, the logic must return to the original business thread for execution, which means that the business logic was originally executed by whom, and after the IO operation is completed, it is returned to the original person to continue execution.

3.Usage:
	a. Call Process to select an asynchronous worker for asynchronous IO operation logic;
	b. Set the variables that need to be shared in the asynchronous IO logic and the asynchronous return result: asyncResult.SetReturnedObj
	c. Register and set the asynchronous callback, that is, return to the original business thread to continue subsequent operations: asyncResult.OnComplete
*/

// Asynchronous worker group (异步worker组)
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

	// Initialization (初始化)
	initAsyncWorkerLocker.Lock()
	defer initAsyncWorkerLocker.Unlock()

	// Re-acquire the worker who does the job (重新拿到这个干活的工人)
	curWorker = asyncWorkerArray[workerIndex]

	// And recheck for nil value (并重新进行空指针判断)
	if curWorker != nil {
		return curWorker
	}

	// If it is still nil after double checking, perform initialization operation
	// (双重检查之后还是空值,进行初始化操作)
	curWorker = &AsyncWorker{
		taskQ: make(chan func(), 2048),
	}

	asyncWorkerArray[workerIndex] = curWorker
	go curWorker.loopExecTask()

	return curWorker
}

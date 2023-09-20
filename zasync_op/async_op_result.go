/*
	Package zasync_op
	@Author：14March
	@File：async_op_result.go
*/

package zasync_op

import (
	"sync/atomic"

	"github.com/aceld/zinx/ziface"
	"github.com/aceld/zinx/znet"
)

type AsyncOpResult struct {
	// Player connection (玩家链接)
	conn ziface.IConnection
	// Returned object (已返回对象)
	returnedObj interface{}
	// Completion callback function(完成回调函数)
	completeFunc func()
	// Whether the return value has been set(是否已有返回值)
	hasReturnedObj int32
	// Whether the completion callback function has been set(是否已有回调函数)
	hasCompleteFunc int32
	// Whether the completion function has already been called(是否已经调用过完成函数)
	completeFuncHasAlreadyBeenCalled int32 // Default value = 0, not yet called (默认值 = 0, 没被调用过)
}

// GetReturnedObj returns the return value (获取返回值)
func (aor *AsyncOpResult) GetReturnedObj() interface{} {
	return aor.returnedObj
}

// SetReturnedObj sets the return value(设置返回值)
func (aor *AsyncOpResult) SetReturnedObj(val interface{}) {
	if atomic.CompareAndSwapInt32(&aor.hasReturnedObj, 0, 1) {
		aor.returnedObj = val
		// **** 防止未调用回调函数问题: 设置处理结果时，直接调用回调函数：1.回调函数未绑定，调用空；2.回调函数已绑定,立马调用 ****
		// **** Prevent the problem of not calling the completion callback function:
		// Call the callback function directly when setting the processing result:
		// 1. If the callback function is not bound, it will be called with nil;
		// 2. If the callback function is bound, it will be called immediately. ****
		aor.doComplete()
	}
}

// OnComplete sets the completion callback function(完成回调函数)
func (aor *AsyncOpResult) OnComplete(val func()) {
	if atomic.CompareAndSwapInt32(&aor.hasCompleteFunc, 0, 1) {
		aor.completeFunc = val

		// **** 防止未调用回调函数问题:设置回调函数时，发现已经有处理结果了，直接调用 ****
		// **** Prevent the problem of not calling the completion callback function:
		// If a processing result already exists when setting the callback function,
		// call the callback function directly. ****
		if atomic.LoadInt32(&aor.hasReturnedObj) == 1 {
			aor.doComplete()
		}
	}
}

// doComplete executes the completion callback function, double-checking when setting the processing result and callback function to prevent the possibility of not calling the callback function
// (执行完成回调,在设置处理结果和回调函数的时候双重检查，杜绝未调用回调函数的可能性)
func (aor *AsyncOpResult) doComplete() {
	if aor.completeFunc == nil {
		return
	}

	// Prevent re-entry problems (防止可重入问题)
	if atomic.CompareAndSwapInt32(&aor.completeFuncHasAlreadyBeenCalled, 0, 1) {
		// Prevent cross-thread calling problems,
		// throw it to the corresponding business thread to execute
		// (防止跨线程调用问题,扔到所属业务线程里去执行)
		request := znet.NewFuncRequest(aor.conn, aor.completeFunc)
		aor.conn.GetMsgHandler().SendMsgToTaskQueue(request)
	}
}

// NewAsyncOpResult creates a new asynchronous result(新建异步结果)
func NewAsyncOpResult(conn ziface.IConnection) *AsyncOpResult {
	result := &AsyncOpResult{}
	result.conn = conn
	return result
}

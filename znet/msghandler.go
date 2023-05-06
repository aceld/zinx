package znet

import (
	"encoding/hex"

	"github.com/aceld/zinx/zconf"
	"github.com/aceld/zinx/ziface"
	"github.com/aceld/zinx/zlog"
)

const (
	// If the Worker goroutine pool is not started, a virtual WorkerID is assigned to the MsgHandler, which is 0, for metric counting
	// After starting the Worker goroutine pool, the ID of each worker is 0,1,2,3...
	// (如果不启动Worker协程池，则会给MsgHandler分配一个虚拟的WorkerID，这个workerID为0, 便于指标统计
	// 启动了Worker协程池后，每个worker的ID为0,1,2,3...)
	WorkerIDWithoutWorkerPool int = 0
)

// MsgHandle is the module for handling message processing callbacks
// (对消息的处理回调模块)
type MsgHandle struct {
	// The number of worker goroutines in the business work Worker pool
	// (业务工作Worker池的数量)
	WorkerPoolSize uint32

	// A message queue for workers to take tasks
	// (Worker负责取任务的消息队列)
	TaskQueue []chan ziface.IRequest

	// Chain builder for the responsibility chain
	// (责任链构造器)
	builder      *chainBuilder
	RouterSlices *RouterSlices
}

// newMsgHandle creates MsgHandle
// zinxRole: IServer/IClient
func newMsgHandle() *MsgHandle {
	handle := &MsgHandle{
		RouterSlices:   NewRouterSlices(),
		WorkerPoolSize: zconf.GlobalObject.WorkerPoolSize,
		// One worker corresponds to one queue (一个worker对应一个queue)
		TaskQueue: make([]chan ziface.IRequest, zconf.GlobalObject.WorkerPoolSize),
		builder:   newChainBuilder(),
	}

	// It is necessary to add the MsgHandle to the responsibility chain here, and it is the last link in the responsibility chain. After decoding in the MsgHandle, data distribution is done by router
	// (此处必须把 msghandler 添加到责任链中，并且是责任链最后一环，在msghandler中进行解码后由router做数据分发)
	handle.builder.Tail(handle)
	return handle
}

// Data processing interceptor that is necessary by default in Zinx
// (Zinx默认必经的数据处理拦截器)
func (mh *MsgHandle) Intercept(chain ziface.IChain) ziface.IcResp {
	request := chain.Request()
	if request != nil {
		switch request.(type) {
		case ziface.IRequest:
			iRequest := request.(ziface.IRequest)
			if zconf.GlobalObject.WorkerPoolSize > 0 {
				// If the worker pool mechanism has been started, hand over the message to the worker for processing
				// (已经启动工作池机制，将消息交给Worker处理)
				mh.SendMsgToTaskQueue(iRequest)
			} else {

				// Execute the corresponding Handle method from the bound message and its corresponding processing method
				// (从绑定好的消息和对应的处理方法中执行对应的Handle方法)
				go mh.doMsgHandlerSlices(iRequest, WorkerIDWithoutWorkerPool)
			}
		}
	}

	return chain.Proceed(chain.Request())
}

func (mh *MsgHandle) AddInterceptor(interceptor ziface.IInterceptor) {
	if mh.builder != nil {
		mh.builder.AddInterceptor(interceptor)
	}
}

func (mh *MsgHandle) GetTaskQueueWorkerId(request ziface.IRequest) uint64 {
	// Assign the worker responsible for processing the current connection based on the ConnID
	// Using a round-robin average allocation rule to get the workerID that needs to process this connection
	// (根据ConnID来分配当前的连接应该由哪个worker负责处理
	// 轮询的平均分配法则
	// 得到需要处理此条连接的workerID)
	workerID := request.GetConnection().GetConnID() % uint64(mh.WorkerPoolSize)
	return workerID
}

// SendMsgToTaskQueue sends the message to the TaskQueue for processing by the worker
// (将消息交给TaskQueue,由worker进行处理)
func (mh *MsgHandle) SendMsgToTaskQueue(request ziface.IRequest) {
	workerID := mh.GetTaskQueueWorkerId(request)
	// zlog.Ins().DebugF("Add ConnID=%d request msgID=%d to workerID=%d", request.GetConnection().GetConnID(), request.GetMsgID(), workerID)
	// Send the request message to the task queue
	mh.TaskQueue[workerID] <- request
	zlog.Ins().DebugF("SendMsgToTaskQueue-->%s", hex.EncodeToString(request.GetData()))
}

// doFuncHandler handles functional requests (执行函数式请求)
func (mh *MsgHandle) doFuncHandler(request ziface.IFuncRequest, workerId int) {
	defer func() {
		if err := recover(); err != nil {
			zlog.Ins().ErrorF("doFuncRequest panic: %v", err)
		}
	}()
	// Execute the functional request (执行函数式请求)
	request.CallFunc()
}

func (mh *MsgHandle) Execute(request ziface.IRequest) {
	// Pass the message to the responsibility chain to handle it through interceptors layer by layer and pass it on layer by layer.
	// (将消息丢到责任链，通过责任链里拦截器层层处理层层传递)
	mh.builder.Execute(request)
}

// AddRouter adds specific processing logic for messages
// (为消息添加具体的处理逻辑)
func (mh *MsgHandle) AddRouter(msgID uint32, router ziface.IRouter) {
	mh.AddRouterSlices(msgID, router.PreHandle, router.Handle, router.PostHandle)
}

// AddRouterSlices adds router handlers using slices
// (切片路由添加)
func (mh *MsgHandle) AddRouterSlices(msgId uint32, handler ...ziface.RouterHandler) ziface.IRouterSlices {
	mh.RouterSlices.AddHandler(msgId, handler...)
	return mh.RouterSlices
}

// Group routes into a group (路由分组)
func (mh *MsgHandle) Group(start, end uint32, Handlers ...ziface.RouterHandler) ziface.IGroupRouterSlices {
	return NewGroup(start, end, mh.RouterSlices, Handlers...)
}
func (mh *MsgHandle) Use(Handlers ...ziface.RouterHandler) ziface.IRouterSlices {
	mh.RouterSlices.Use(Handlers...)
	return mh.RouterSlices
}

func (mh *MsgHandle) doMsgHandlerSlices(request ziface.IRequest, workerID int) {
	defer func() {
		if err := recover(); err != nil {
			zlog.Ins().ErrorF("doMsgHandler panic: %v", err)
		}
	}()

	msgId := request.GetMsgID()
	handlers, ok := mh.RouterSlices.GetHandlers(msgId)
	if !ok {
		zlog.Ins().ErrorF("api msgID = %d is not FOUND!", request.GetMsgID())
		return
	}

	request.BindRouterSlices(handlers)
	request.RouterSlicesNext()

}

// StartOneWorker starts a worker workflow
// (启动一个Worker工作流程)
func (mh *MsgHandle) StartOneWorker(workerID int, taskQueue chan ziface.IRequest) {
	zlog.Ins().InfoF("Worker ID = %d is started.", workerID)
	// Continuously wait for messages in the queue
	// (不断地等待队列中的消息)
	for {
		select {
		// If there is a message, take out the Request from the queue and execute the bound business method
		// (有消息则取出队列的Request，并执行绑定的业务方法)
		case request := <-taskQueue:

			switch req := request.(type) {

			case ziface.IFuncRequest:
				// Internal function call request (内部函数调用request)

				mh.doFuncHandler(req, workerID)

			case ziface.IRequest: // Client message request

				mh.doMsgHandlerSlices(req, workerID)

			}
		}
	}
}

// StartWorkerPool starts the worker pool
func (mh *MsgHandle) StartWorkerPool() {
	// Iterate through the required number of workers and start them one by one
	// (遍历需要启动worker的数量，依此启动)
	for i := 0; i < int(mh.WorkerPoolSize); i++ {
		// A worker is started
		// Allocate space for the corresponding task queue for the current worker
		// (给当前worker对应的任务队列开辟空间)
		mh.TaskQueue[i] = make(chan ziface.IRequest, zconf.GlobalObject.MaxWorkerTaskLen)

		// Start the current worker, blocking and waiting for messages to be passed in the corresponding task queue
		// (启动当前Worker，阻塞的等待对应的任务队列是否有消息传递进来)
		go mh.StartOneWorker(i, mh.TaskQueue[i])
	}
}

package znet

import (
	"encoding/hex"
	"fmt"
	"sync"

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
	// A map property that stores the processing methods for each MsgID
	// (存放每个MsgID 所对应的处理方法的map属性)
	Apis map[uint32]ziface.IRouter

	// The number of worker goroutines in the business work Worker pool
	// (业务工作Worker池的数量)
	WorkerPoolSize uint32

	// A collection of idle workers, used for zconf.WorkerModeBind
	// 空闲worker集合，用于zconf.WorkerModeBind
	freeWorkers  map[uint32]struct{}
	freeWorkerMu sync.Mutex

	// A message queue for workers to take tasks
	// (Worker负责取任务的消息队列)
	TaskQueue []chan ziface.IRequest

	// A collection of extra workers, used for zconf.WorkerModeDynamicBind
	// (池里的工作线程不够用的时候, 可临时额外分配workerID集合, 用于zconf.WorkerModeDynamicBind)
	extraFreeWorkers  map[uint32]struct{}
	extraFreeWorkerMu sync.Mutex

	// Chain builder for the responsibility chain
	// (责任链构造器)
	builder      *chainBuilder
	RouterSlices *RouterSlices
}

// newMsgHandle creates MsgHandle
// zinxRole: IServer
func newMsgHandle() *MsgHandle {
	var freeWorkers map[uint32]struct{}
	var extraFreeWorkers map[uint32]struct{}

	if zconf.GlobalObject.WorkerMode == zconf.WorkerModeBind {
		// Assign a workder to each link, avoid interactions when multiple links are processed by the same worker
		// MaxWorkerTaskLen can also be reduced, for example, 50
		// 为每个连接分配一个workder，避免同一worker处理多个连接时的互相影响
		// 同时可以减小MaxWorkerTaskLen，比如50，因为每个worker的负担减轻了
		zconf.GlobalObject.WorkerPoolSize = uint32(zconf.GlobalObject.MaxConn)
		freeWorkers = make(map[uint32]struct{}, zconf.GlobalObject.WorkerPoolSize)
		for i := uint32(0); i < zconf.GlobalObject.WorkerPoolSize; i++ {
			freeWorkers[i] = struct{}{}
		}
	}

	TaskQueueLen := zconf.GlobalObject.WorkerPoolSize

	if zconf.GlobalObject.WorkerMode == zconf.WorkerModeDynamicBind {
		zlog.Ins().DebugF("WorkerMode = %s", zconf.WorkerModeDynamicBind)
		freeWorkers = make(map[uint32]struct{}, zconf.GlobalObject.WorkerPoolSize)
		for i := uint32(0); i < zconf.GlobalObject.WorkerPoolSize; i++ {
			freeWorkers[i] = struct{}{}
		}

		extraFreeWorkers = make(map[uint32]struct{}, zconf.GlobalObject.MaxConn-int(zconf.GlobalObject.WorkerPoolSize))
		for i := zconf.GlobalObject.WorkerPoolSize; i < uint32(zconf.GlobalObject.MaxConn); i++ {
			extraFreeWorkers[i] = struct{}{}
		}
		TaskQueueLen = uint32(zconf.GlobalObject.MaxConn)
	}

	handle := &MsgHandle{
		Apis:         make(map[uint32]ziface.IRouter),
		RouterSlices: NewRouterSlices(),
		freeWorkers:  freeWorkers,
		builder:      newChainBuilder(),
		// 可额外临时分配的workerID集合
		extraFreeWorkers: extraFreeWorkers,
	}

	// server
	handle.WorkerPoolSize = zconf.GlobalObject.WorkerPoolSize
	// One worker corresponds to one queue (一个worker对应一个queue)
	handle.TaskQueue = make([]chan ziface.IRequest, TaskQueueLen)

	// It is necessary to add the MsgHandle to the responsibility chain here, and it is the last link in the responsibility chain. After decoding in the MsgHandle, data distribution is done by router
	// (此处必须把 msghandler 添加到责任链中，并且是责任链最后一环，在msghandler中进行解码后由router做数据分发)
	handle.builder.Tail(handle)
	return handle
}

// newCliMsgHandle creates MsgHandle
// zinxRole: IClient
func newCliMsgHandle() *MsgHandle {
	var freeWorkers map[uint32]struct{}
	var extraFreeWorkers map[uint32]struct{}

	if zconf.GlobalObject.WorkerMode == zconf.WorkerModeBind {
		// Assign a workder to each link, avoid interactions when multiple links are processed by the same worker
		// MaxWorkerTaskLen can also be reduced, for example, 50
		// 为每个连接分配一个workder，避免同一worker处理多个连接时的互相影响
		// 同时可以减小MaxWorkerTaskLen，比如50，因为每个worker的负担减轻了
		zconf.GlobalObject.WorkerPoolSize = uint32(zconf.GlobalObject.MaxConn)
		freeWorkers = make(map[uint32]struct{}, zconf.GlobalObject.WorkerPoolSize)
		for i := uint32(0); i < zconf.GlobalObject.WorkerPoolSize; i++ {
			freeWorkers[i] = struct{}{}
		}
	}

	TaskQueueLen := zconf.GlobalObject.WorkerPoolSize

	if zconf.GlobalObject.WorkerMode == zconf.WorkerModeDynamicBind {
		zlog.Ins().DebugF("WorkerMode = %s", zconf.WorkerModeDynamicBind)
		freeWorkers = make(map[uint32]struct{}, zconf.GlobalObject.WorkerPoolSize)
		for i := uint32(0); i < zconf.GlobalObject.WorkerPoolSize; i++ {
			freeWorkers[i] = struct{}{}
		}

		extraFreeWorkers = make(map[uint32]struct{}, zconf.GlobalObject.MaxConn-int(zconf.GlobalObject.WorkerPoolSize))
		for i := zconf.GlobalObject.WorkerPoolSize; i < uint32(zconf.GlobalObject.MaxConn); i++ {
			extraFreeWorkers[i] = struct{}{}
		}
		TaskQueueLen = uint32(zconf.GlobalObject.MaxConn)
	}

	handle := &MsgHandle{
		Apis:         make(map[uint32]ziface.IRouter),
		RouterSlices: NewRouterSlices(),
		freeWorkers:  freeWorkers,
		builder:      newChainBuilder(),
		// 可额外临时分配的workerID集合
		extraFreeWorkers: extraFreeWorkers,
	}

	// client: Set worker pool size to 0 to turn off the worker pool in the client (客户端将协程池关闭)
	handle.WorkerPoolSize = 0
	// One worker corresponds to one queue (一个worker对应一个queue)
	handle.TaskQueue = make([]chan ziface.IRequest, TaskQueueLen)

	// It is necessary to add the MsgHandle to the responsibility chain here, and it is the last link in the responsibility chain. After decoding in the MsgHandle, data distribution is done by router
	// (此处必须把 msghandler 添加到责任链中，并且是责任链最后一环，在msghandler中进行解码后由router做数据分发)
	handle.builder.Tail(handle)
	return handle
}

// Use worker ID
// 占用workerID
func useWorker(conn ziface.IConnection) uint32 {
	var workerId uint32

	mh, _ := conn.GetMsgHandler().(*MsgHandle)
	if mh == nil {
		zlog.Ins().ErrorF("useWorker failed, mh is nil")
		return 0
	}

	if zconf.GlobalObject.WorkerMode == zconf.WorkerModeBind {
		mh.freeWorkerMu.Lock()
		defer mh.freeWorkerMu.Unlock()

		for k := range mh.freeWorkers {
			delete(mh.freeWorkers, k)
			return k
		}
	}

	if zconf.GlobalObject.WorkerMode == zconf.WorkerModeDynamicBind {
		mh.freeWorkerMu.Lock()
		// try to get workerID from workerPool first
		// 首先尝试从工作线程池里获取一个空闲的workerID
		for workerID := range mh.freeWorkers {
			delete(mh.freeWorkers, workerID)
			mh.freeWorkerMu.Unlock()
			return workerID
		}
		mh.freeWorkerMu.Unlock()

		// 工作池的worker用完了，临时从extraFreeWorkers取一个额外的workerID, 并相应启动一个临时的worker
		mh.extraFreeWorkerMu.Lock()
		defer mh.extraFreeWorkerMu.Unlock()
		for workerID := range mh.extraFreeWorkers {
			zlog.Ins().DebugF("start extra worker, workerID=%d", workerID)
			mh.TaskQueue[workerID] = make(chan ziface.IRequest, zconf.GlobalObject.MaxWorkerTaskLen)
			go mh.StartOneWorker(int(workerID), mh.TaskQueue[workerID])
			return workerID
		}
	}

	//Compatible with the situation where the client has no worker, and solve the situation divide 0
	//(兼容client没有worker情况，解决除0的情况)
	if mh.WorkerPoolSize == 0 {
		workerId = 0
	} else {
		// Assign the worker responsible for processing the current connection based on the ConnID
		// Using a round-robin average allocation rule to get the workerID that needs to process this connection
		// (根据ConnID来分配当前的连接应该由哪个worker负责处理
		// 轮询的平均分配法则
		// 得到需要处理此条连接的workerID)
		workerId = uint32(conn.GetConnID() % uint64(mh.WorkerPoolSize))
	}

	return workerId
}

// Free worker ID
// 释放workerid
func freeWorker(conn ziface.IConnection) {
	mh, _ := conn.GetMsgHandler().(*MsgHandle)
	if mh == nil {
		zlog.Ins().ErrorF("useWorker failed, mh is nil")
		return
	}

	if zconf.GlobalObject.WorkerMode == zconf.WorkerModeBind {
		mh.freeWorkerMu.Lock()
		defer mh.freeWorkerMu.Unlock()

		mh.freeWorkers[conn.GetWorkerID()] = struct{}{}
	}

	if zconf.GlobalObject.WorkerMode == zconf.WorkerModeDynamicBind {
		workerID := conn.GetWorkerID()
		if workerID < mh.WorkerPoolSize {
			// 说明这个是工作线程池里的workerID，回收这个workerID, workerID对应的worker不需要销毁
			mh.freeWorkerMu.Lock()
			mh.freeWorkers[workerID] = struct{}{}
			mh.freeWorkerMu.Unlock()
		} else {
			// 说明这个worker是一个临时的worker，需要销毁这个worker
			mh.StopOneWorker(int(workerID))
			// 回收workerID, 放回额外workerID池里
			mh.extraFreeWorkerMu.Lock()
			mh.extraFreeWorkers[workerID] = struct{}{}
			mh.extraFreeWorkerMu.Unlock()
		}
	}
}

// Data processing interceptor that is necessary by default in Zinx
// (Zinx默认必经的数据处理拦截器)
func (mh *MsgHandle) Intercept(chain ziface.IChain) ziface.IcResp {
	request := chain.Request()
	if request != nil {
		switch request.(type) {
		case ziface.IRequest:
			iRequest := request.(ziface.IRequest)
			if mh.WorkerPoolSize > 0 {
				// If the worker pool mechanism has been started, hand over the message to the worker for processing
				// (已经启动工作池机制，将消息交给Worker处理)
				mh.SendMsgToTaskQueue(iRequest)
			} else {

				// Execute the corresponding Handle method from the bound message and its corresponding processing method
				// (从绑定好的消息和对应的处理方法中执行对应的Handle方法)
				if !zconf.GlobalObject.RouterSlicesMode {
					go mh.doMsgHandler(iRequest, WorkerIDWithoutWorkerPool)
				} else if zconf.GlobalObject.RouterSlicesMode {
					go mh.doMsgHandlerSlices(iRequest, WorkerIDWithoutWorkerPool)
				}

			}
		}
	}

	return chain.Proceed(chain.Request())
}

// SetHeadInterceptor sets the head interceptor of the responsibility chain, which is the first interceptor to be executed
// (SetHeadInterceptor 设置责任链的头拦截器，也就是第一个要执行的拦截器)
// will replace the default head interceptor
func (mh *MsgHandle) SetHeadInterceptor(interceptor ziface.IInterceptor) {
	if mh.builder != nil {
		mh.builder.Head(interceptor)
	}
}

func (mh *MsgHandle) AddInterceptor(interceptor ziface.IInterceptor) {
	if mh.builder != nil {
		mh.builder.AddInterceptor(interceptor)
	}
}

// SendMsgToTaskQueue sends the message to the TaskQueue for processing by the worker
// (将消息交给TaskQueue,由worker进行处理)
func (mh *MsgHandle) SendMsgToTaskQueue(request ziface.IRequest) {
	workerID := request.GetConnection().GetWorkerID()
	// zlog.Ins().DebugF("Add ConnID=%d request msgID=%d to workerID=%d", request.GetConnection().GetConnID(), request.GetMsgID(), workerID)
	// Send the request message to the task queue
	mh.TaskQueue[workerID] <- request
	if zlog.Ins().DebugEnabled() {
		zlog.Ins().DebugF("SendMsgToTaskQueue-->%s", hex.EncodeToString(request.GetData()))
	}
}

// doFuncHandler handles functional requests (执行函数式请求)
func (mh *MsgHandle) doFuncHandler(request ziface.IFuncRequest, workerID int) {
	defer func() {
		if err := recover(); err != nil {
			zlog.Ins().ErrorF("workerID: %d doFuncRequest panic: %v", workerID, err)
		}
	}()
	// Execute the functional request (执行函数式请求)
	request.CallFunc()
}

// doMsgHandler immediately handles messages in a non-blocking manner
// (立即以非阻塞方式处理消息)
func (mh *MsgHandle) doMsgHandler(request ziface.IRequest, workerID int) {
	defer func() {
		if err := recover(); err != nil {
			zlog.Ins().ErrorF("workerID: %d doMsgHandler panic: %v", workerID, err)
		}
	}()

	msgId := request.GetMsgID()
	handler, ok := mh.Apis[msgId]

	if !ok {
		zlog.Ins().ErrorF("api msgID = %d is not FOUND!", request.GetMsgID())
		return
	}

	// Bind the Request request to the corresponding Router relationship
	// (Request请求绑定Router对应关系)
	request.BindRouter(handler)

	// Execute the corresponding processing method
	request.Call()

	// 执行完成后回收 Request 对象回对象池
	PutRequest(request)
}

func (mh *MsgHandle) Execute(request ziface.IRequest) {
	// Pass the message to the responsibility chain to handle it through interceptors layer by layer and pass it on layer by layer.
	// (将消息丢到责任链，通过责任链里拦截器层层处理层层传递)
	mh.builder.Execute(request)
}

// AddRouter adds specific processing logic for messages
// (为消息添加具体的处理逻辑)
func (mh *MsgHandle) AddRouter(msgID uint32, router ziface.IRouter) {
	// 1. Check whether the current API processing method bound to the msgID already exists
	// (判断当前msg绑定的API处理方法是否已经存在)
	if _, ok := mh.Apis[msgID]; ok {
		msgErr := fmt.Sprintf("repeated api , msgID = %+v\n", msgID)
		panic(msgErr)
	}
	// 2. Add the binding relationship between msg and API
	// (添加msg与api的绑定关系)
	mh.Apis[msgID] = router
	zlog.Ins().InfoF("Add Router msgID = %d", msgID)
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
			zlog.Ins().ErrorF("workerID: %d doMsgHandler panic: %v", workerID, err)
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
	// 执行完成后回收 Request 对象回对象池
	PutRequest(request)
}

func (mh *MsgHandle) StopOneWorker(workerID int) {
	zlog.Ins().DebugF("stop Worker ID = %d ", workerID)
	// Stop the worker by closing the corresponding taskQueue
	// (停止一个Worker，通过关闭对应的taskQueue)
	close(mh.TaskQueue[workerID])
}

// StartOneWorker starts a worker workflow
// (启动一个Worker工作流程)
func (mh *MsgHandle) StartOneWorker(workerID int, taskQueue chan ziface.IRequest) {
	zlog.Ins().DebugF("Worker ID = %d is started.", workerID)
	// Continuously wait for messages in the queue
	// (不断地等待队列中的消息)
	for {
		select {
		// If there is a message, take out the Request from the queue and execute the bound business method
		// (有消息则取出队列的Request，并执行绑定的业务方法)
		case request, ok := <-taskQueue:
			if !ok {
				// DynamicBind Mode, destroy current worker by close the taskQueue
				// (DynamicBind模式下，临时创建的worker, 是通过关闭taskQueue 来销毁当前worker)
				zlog.Ins().ErrorF(" taskQueue is closed, Worker ID = %d quit", workerID)
				return
			}
			switch req := request.(type) {

			case ziface.IFuncRequest:
				// Internal function call request (内部函数调用request)

				mh.doFuncHandler(req, workerID)

			case ziface.IRequest: // Client message request

				if !zconf.GlobalObject.RouterSlicesMode {
					mh.doMsgHandler(req, workerID)
				} else if zconf.GlobalObject.RouterSlicesMode {
					mh.doMsgHandlerSlices(req, workerID)
				}
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

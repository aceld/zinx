package znet

import (
	"encoding/hex"
	"fmt"
	"strconv"
	"time"

	"github.com/aceld/zinx/zmetrics"

	"github.com/aceld/zinx/zconf"
	"github.com/aceld/zinx/ziface"
	"github.com/aceld/zinx/zlog"
)

const (
	//如果不启动Worker协程池，则会给MsgHandler分配一个虚拟的WorkerID，这个workerID为0, 便于指标统计
	//启动了Worker协程池后，每个worker的ID为0,1,2,3...
	WorkerIDWithoutWorkerPool int = 0
)

// MsgHandle 对消息的处理回调模块
type MsgHandle struct {
	Apis           map[uint32]ziface.IRouter // 存放每个MsgID 所对应的处理方法的map属性
	WorkerPoolSize uint32                    // 业务工作Worker池的数量
	TaskQueue      []chan ziface.IRequest    // Worker负责取任务的消息队列
	builder        *chainBuilder             // 责任链构造器
	RouterSlices   *RouterSlices
}

// NewMsgHandle 创建MsgHandle
// zinxRole: IServer/IClient
func newMsgHandle() *MsgHandle {
	handle := &MsgHandle{
		Apis:           make(map[uint32]ziface.IRouter),
		RouterSlices:   NewRouterSlices(),
		WorkerPoolSize: zconf.GlobalObject.WorkerPoolSize,
		// 一个worker对应一个queue
		TaskQueue: make([]chan ziface.IRequest, zconf.GlobalObject.WorkerPoolSize),
		builder:   newChainBuilder(),
	}
	// 此处必须把 msghandler 添加到责任链中，并且是责任链最后一环，在msghandler中进行解码后由router做数据分发
	handle.builder.Tail(handle)
	return handle
}

// Zinx默认必经的数据处理拦截器
func (mh *MsgHandle) Intercept(chain ziface.IChain) ziface.IcResp {
	request := chain.Request()
	if request != nil {
		switch request.(type) {
		case ziface.IRequest:
			iRequest := request.(ziface.IRequest)
			if zconf.GlobalObject.WorkerPoolSize > 0 {
				// 已经启动工作池机制，将消息交给Worker处理
				mh.SendMsgToTaskQueue(iRequest)
			} else {

				// 从绑定好的消息和对应的处理方法中执行对应的Handle方法
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

func (mh *MsgHandle) AddInterceptor(interceptor ziface.IInterceptor) {
	if mh.builder != nil {
		mh.builder.AddInterceptor(interceptor)
	}
}

func (mh *MsgHandle) GetTaskQueueWorkerId(request ziface.IRequest) uint64 {
	// 根据ConnID来分配当前的连接应该由哪个worker负责处理
	// 轮询的平均分配法则
	// 得到需要处理此条连接的workerID
	workerID := request.GetConnection().GetConnID() % uint64(mh.WorkerPoolSize)
	return workerID
}

// SendMsgToTaskQueue 将消息交给TaskQueue,由worker进行处理
func (mh *MsgHandle) SendMsgToTaskQueue(request ziface.IRequest) {
	workerID := mh.GetTaskQueueWorkerId(request)
	// zlog.Ins().DebugF("Add ConnID=%d request msgID=%d to workerID=%d", request.GetConnection().GetConnID(), request.GetMsgID(), workerID)
	// 将请求消息发送给任务队列
	mh.TaskQueue[workerID] <- request
	zlog.Ins().DebugF("SendMsgToTaskQueue-->%s", hex.EncodeToString(request.GetData()))
}

func (mh *MsgHandle) StatisticsMetrics(request ziface.IRequest, workerId int, msgId uint32, timeNow time.Time) {

	conn := request.GetConnection()

	//统计MsgID被调度的路由次数
	zmetrics.Metrics().IncRouterSchedule(conn.LocalAddrString(), conn.GetName(), strconv.Itoa(workerId), strconv.Itoa(int(msgId)))

	//统计Router和MsgID业务调度的耗时
	zmetrics.Metrics().ObserveRouterScheduleDuration(conn.LocalAddrString(), conn.GetName(), strconv.Itoa(workerId), strconv.Itoa(int(msgId)), time.Since(timeNow))
}

// DoMsgHandler 立即以非阻塞方式处理消息
func (mh *MsgHandle) doMsgHandler(request ziface.IRequest, workerID int) {
	defer func() {
		if err := recover(); err != nil {
			zlog.Ins().ErrorF("doMsgHandler panic: %v", err)
		}
	}()

	var timeNow time.Time
	if zmetrics.Metrics().IsEnable() {
		timeNow = time.Now()
	}

	msgId := request.GetMsgID()
	handler, ok := mh.Apis[msgId]

	if !ok {
		zlog.Ins().ErrorF("api msgID = %d is not FOUND!", request.GetMsgID())
		return
	}

	// Request请求绑定Router对应关系
	request.BindRouter(handler)
	// 执行对应处理方法
	request.Call()

	//统计Router调度指标数据
	mh.StatisticsMetrics(request, workerID, msgId, timeNow)
}

func (mh *MsgHandle) Execute(request ziface.IRequest) {
	mh.builder.Execute(request) // 将消息丢到责任链，通过责任链里拦截器层层处理层层传递
}

// AddRouter 为消息添加具体的处理逻辑
func (mh *MsgHandle) AddRouter(msgID uint32, router ziface.IRouter) {
	// 1 判断当前msg绑定的API处理方法是否已经存在
	if _, ok := mh.Apis[msgID]; ok {
		msgErr := fmt.Sprintf("repeated api , msgID = %+v\n", msgID)
		panic(msgErr)
	}
	// 2 添加msg与api的绑定关系
	mh.Apis[msgID] = router
	zlog.Ins().InfoF("Add Router msgID = %d", msgID)
}

// 切片路由添加
func (mh *MsgHandle) AddRouterSlices(msgId uint32, handler ...ziface.RouterHandler) ziface.IRouterSlices {
	mh.RouterSlices.AddHandler(msgId, handler...)
	return mh.RouterSlices
}

// 路由分组
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

	var timeNow time.Time
	if zmetrics.Metrics().IsEnable() {
		timeNow = time.Now()
	}

	msgId := request.GetMsgID()
	handlers, ok := mh.RouterSlices.GetHandlers(msgId)
	if !ok {
		zlog.Ins().ErrorF("api msgID = %d is not FOUND!", request.GetMsgID())
		return
	}

	request.BindRouterSlices(handlers)
	request.RouterSlicesNext()

	//统计MsgID被调度的路由次数
	mh.StatisticsMetrics(request, workerID, msgId, timeNow)
}

// StartOneWorker 启动一个Worker工作流程
func (mh *MsgHandle) StartOneWorker(workerID int, taskQueue chan ziface.IRequest) {
	zlog.Ins().InfoF("Worker ID = %d is started.", workerID)
	// 不断地等待队列中的消息
	for {
		select {
		// 有消息则取出队列的Request，并执行绑定的业务方法
		case request := <-taskQueue:

			switch req := request.(type) {

			case ziface.IFuncRequest: // 内部函数调用request

				req.CallFunc()

			case ziface.IRequest: // 客户端消息request

				if !zconf.GlobalObject.RouterSlicesMode {
					mh.doMsgHandler(req, workerID)
				} else if zconf.GlobalObject.RouterSlicesMode {
					mh.doMsgHandlerSlices(req, workerID)
				}

				// Metrics统计，每次处理完一个请求，当前WorkId处理的任务数量+1
				conn := request.GetConnection()
				zmetrics.Metrics().IncTask(conn.LocalAddrString(), conn.GetName(), strconv.Itoa(workerID))
			}
		}
	}
}

// StartWorkerPool 启动worker工作池
func (mh *MsgHandle) StartWorkerPool() {
	// 遍历需要启动worker的数量，依此启动
	for i := 0; i < int(mh.WorkerPoolSize); i++ {
		// 一个worker被启动
		// 给当前worker对应的任务队列开辟空间
		mh.TaskQueue[i] = make(chan ziface.IRequest, zconf.GlobalObject.MaxWorkerTaskLen)
		// 启动当前Worker，阻塞的等待对应的任务队列是否有消息传递进来
		go mh.StartOneWorker(i, mh.TaskQueue[i])
	}
}

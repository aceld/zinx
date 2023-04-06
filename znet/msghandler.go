package znet

import (
	"encoding/hex"
	"fmt"

	"github.com/aceld/zinx/zconf"
	"github.com/aceld/zinx/ziface"
	"github.com/aceld/zinx/zinterceptor"
	"github.com/aceld/zinx/zlog"
)

// MsgHandle 对消息的处理回调模块
type MsgHandle struct {
	Apis           map[uint32]ziface.IRouter // 存放每个MsgID 所对应的处理方法的map属性
	WorkerPoolSize uint32                    // 业务工作Worker池的数量
	TaskQueue      []chan ziface.IRequest    // Worker负责取任务的消息队列
	builder        ziface.IBuilder           // 责任链构造器
}

// NewMsgHandle 创建MsgHandle
func NewMsgHandle() *MsgHandle {
	handle := &MsgHandle{
		Apis:           make(map[uint32]ziface.IRouter),
		WorkerPoolSize: zconf.GlobalObject.WorkerPoolSize,
		// 一个worker对应一个queue
		TaskQueue: make([]chan ziface.IRequest, zconf.GlobalObject.WorkerPoolSize),
		builder:   zinterceptor.NewBuilder(),
	}
	// 此处必须把 msghandler 添加到责任链中，并且是责任链最后一环，在msghandler中进行解码后由router做数据分发
	handle.builder.Tail(handle)
	return handle
}

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
				go mh.doMsgHandler(iRequest)
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

// SendMsgToTaskQueue 将消息交给TaskQueue,由worker进行处理
func (mh *MsgHandle) SendMsgToTaskQueue(request ziface.IRequest) {
	// 根据ConnID来分配当前的连接应该由哪个worker负责处理
	// 轮询的平均分配法则

	// 得到需要处理此条连接的workerID
	workerID := request.GetConnection().GetConnID() % uint64(mh.WorkerPoolSize)
	// zlog.Ins().DebugF("Add ConnID=%d request msgID=%d to workerID=%d", request.GetConnection().GetConnID(), request.GetMsgID(), workerID)
	// 将请求消息发送给任务队列
	mh.TaskQueue[workerID] <- request
	zlog.Ins().DebugF("SendMsgToTaskQueue-->%s", hex.EncodeToString(request.GetData()))
}

// DoMsgHandler 马上以非阻塞方式处理消息
func (mh *MsgHandle) doMsgHandler(request ziface.IRequest) {
	defer func() {
		if err := recover(); err != nil {
			zlog.Ins().ErrorF("doMsgHandler panic: %v", err)
		}
	}()

	handler, ok := mh.Apis[request.GetMsgID()]
	if !ok {
		zlog.Ins().ErrorF("api msgID = %d is not FOUND!", request.GetMsgID())
		return
	}

	// Request请求绑定Router对应关系
	request.BindRouter(handler)
	// 执行对应处理方法
	request.Call()
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

// StartOneWorker 启动一个Worker工作流程
func (mh *MsgHandle) StartOneWorker(workerID int, taskQueue chan ziface.IRequest) {
	zlog.Ins().InfoF("Worker ID = %d is started.", workerID)
	// 不断的等待队列中的消息
	for {
		select {
		// 有消息则取出队列的Request，并执行绑定的业务方法
		case request := <-taskQueue:
			mh.doMsgHandler(request)
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

package router

import (
	"encoding/json"
	"time"

	"github.com/aceld/zinx/examples/zinx_async_op/async_op_apis"
	"github.com/aceld/zinx/examples/zinx_async_op/db_model"
	"github.com/aceld/zinx/examples/zinx_async_op/msg_struct"
	"github.com/aceld/zinx/ziface"
	"github.com/aceld/zinx/zlog"
	"github.com/aceld/zinx/znet"
)

type LoginRouter struct {
	znet.BaseRouter
}

func (hr *LoginRouter) Handle(request ziface.IRequest) {
	zlog.Debug("AsyncOpRouter Handle IN ===>111")

	asyncResult := async_op_apis.AsyncUserSaveData(request) // 测试DB异步操作

	// 测试：执行了一大推业务逻辑... 才设置回调函数
	time.Sleep(1 * time.Second)

	// 异步回调
	asyncResult.OnComplete(func() {
		zlog.Debug("OnComplete IN===>333")
		returnedObj := asyncResult.GetReturnedObj()
		if returnedObj == nil {
			zlog.Debug("注册回调函数时，还未设置异步结果")
			return
		}

		user := returnedObj.(*db_model.UserModel)

		userLoginRsp := &msg_struct.MsgLoginResponse{
			UserId:    user.UserId,
			UserName:  user.Name,
			ErrorCode: 0,
		}

		marshalData, marErr := json.Marshal(userLoginRsp)
		if marErr != nil {
			zlog.Error("LoginRouter marErr", marErr.Error())
			return
		}

		// 回包客户端
		conn := request.GetConnection()
		if sendErr := conn.SendMsg(1, marshalData); sendErr != nil {
			zlog.Error("LoginRouter sendErr", sendErr.Error())
			return
		}
		zlog.Debug("OnComplete OUT===>333")
	})

	// 测试：
	// 原来所属的线程阻塞3秒，回调函数因为是回到原来所属的线程里执行的，所以一定在3秒后执行
	time.Sleep(time.Second * 3)

	zlog.Debug("AsyncOpRouter Handle OUT ===>111")
}

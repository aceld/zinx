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

	asyncResult := async_op_apis.AsyncUserSaveData(request) // // Test DB asynchronous operation(测试DB异步操作)

	// 测试：执行了一大推业务逻辑,才设置回调函数
	// Test: A lot of business logic is executed before setting the callback function
	time.Sleep(1 * time.Second)

	// Asynchronous callback (异步回调)
	asyncResult.OnComplete(func() {
		zlog.Debug("OnComplete IN===>333")
		returnedObj := asyncResult.GetReturnedObj()
		if returnedObj == nil {
			zlog.Debug("The asynchronous result has not been set when registering the callback function.")
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

		// Send response to the client
		conn := request.GetConnection()
		if sendErr := conn.SendMsg(1, marshalData); sendErr != nil {
			zlog.Error("LoginRouter sendErr", sendErr.Error())
			return
		}
		zlog.Debug("OnComplete OUT===>333")

		// Test actively throwing an exception (测试主动异常)
		/*
			a := 0
			b := 1
			c := b / a
			fmt.Println(c)
		*/
	})

	// Test: The original thread is blocked for 3 seconds, and the callback function is executed in the original thread,
	//       so it will be executed after 3 seconds
	// 测试：原来所属的线程阻塞3秒，回调函数因为是回到原来所属的线程里执行的，所以一定在3秒后执行.
	time.Sleep(time.Second * 3)

	zlog.Debug("AsyncOpRouter Handle OUT ===>111")
}

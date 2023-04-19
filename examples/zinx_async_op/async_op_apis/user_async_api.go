package async_op_apis

import (
	"github.com/aceld/zinx/examples/zinx_async_op/db_model"
	"github.com/aceld/zinx/zasync_op"
	"github.com/aceld/zinx/ziface"
)

func AsyncUserSaveData(request ziface.IRequest) *zasync_op.AsyncOpResult {

	opId := 1 // 玩家的唯一标识Id
	asyncResult := zasync_op.NewAsyncOpResult(request.GetConnection())

	zasync_op.Process(
		int(opId),
		func() {
			// 	执行db操作
			user := db_model.SaveUserData()

			// 设置异步返回结果
			asyncResult.SetReturnedObj(user)
		},
	)

	return asyncResult
}

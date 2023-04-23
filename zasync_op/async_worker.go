/*
Package zasync_op
@Author：14March
@File：async_worker.go
*/
package zasync_op

import "github.com/aceld/zinx/zlog"

type AsyncWorker struct {
	taskQ chan func()
}

func (aw *AsyncWorker) process(asyncOp func()) {
	if asyncOp == nil {
		zlog.Error("异步操作为空")
		return
	}

	if aw.taskQ == nil {
		zlog.Error("任务队列尚未初始化")
		return
	}

	aw.taskQ <- func() {
		defer func() {
			if err := recover(); err != nil {
				zlog.Ins().ErrorF("async process panic: %v", err)
			}
		}()

		// 执行异步操作
		asyncOp()

	}
}

func (aw *AsyncWorker) loopExecTask() {
	if aw.taskQ == nil {
		zlog.Error("任务队列尚未初始化")
		return
	}

	for {
		task := <-aw.taskQ
		if task != nil {
			task()
		}
	}
}

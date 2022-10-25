package znet

import "github.com/chnkenc/zinx-xiaoan/ziface"

var logger ziface.ILog

// SetLogger 设置logger
func SetLogger(logHandler ziface.ILog) {
	if logHandler == nil {
		// TODO: 加载默认日志组件
	} else {
		logger = logHandler
	}
}

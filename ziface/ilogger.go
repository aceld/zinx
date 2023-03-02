package ziface

import "context"

type ILogger interface {
	//没有context的日志接口
	InfoF(format string, v ...interface{})
	ErrorF(format string, v ...interface{})
	DebugF(format string, v ...interface{})

	//携带context的日志接口
	InfoFX(ctx context.Context, format string, v ...interface{})
	ErrorFX(ctx context.Context, format string, v ...interface{})
	DebugFX(ctx context.Context, format string, v ...interface{})
}

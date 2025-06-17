package ziface

import "context"

type ILogger interface {
	//without context
	InfoF(format string, v ...interface{})
	ErrorF(format string, v ...interface{})
	DebugF(format string, v ...interface{})

	//with context
	InfoFX(ctx context.Context, format string, v ...interface{})
	ErrorFX(ctx context.Context, format string, v ...interface{})
	DebugFX(ctx context.Context, format string, v ...interface{})

	// 此处增加了 interface 定义，可能会导致原有第三方实现无法通过编译
	DebugEnabled() bool
}

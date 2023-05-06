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
}

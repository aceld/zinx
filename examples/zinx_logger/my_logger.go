package main

import (
	"context"
	"fmt"
)

// User-defined logging method
// The internal engine logging method of zinx can be reset by the logging method of its own business.
// In this example, fmt.Println is used.
// 用户自定义日志方式，
// 可以通过自身业务的日志方式，来重置zinx内部引擎的日志打印方式
// 本例以fmt.Println为例
type MyLogger struct{}

// Without context logging interface
func (l *MyLogger) InfoF(format string, v ...interface{}) {
	fmt.Printf(format, v...)
}

func (l *MyLogger) ErrorF(format string, v ...interface{}) {
	fmt.Printf(format, v...)
}

func (l *MyLogger) DebugF(format string, v ...interface{}) {
	fmt.Printf(format, v...)
}

// Logging interface with context
func (l *MyLogger) InfoFX(ctx context.Context, format string, v ...interface{}) {
	fmt.Println(ctx)
	fmt.Printf(format, v...)
}

func (l *MyLogger) ErrorFX(ctx context.Context, format string, v ...interface{}) {
	fmt.Println(ctx)
	fmt.Printf(format, v...)
}

func (l *MyLogger) DebugFX(ctx context.Context, format string, v ...interface{}) {
	fmt.Println(ctx)
	fmt.Printf(format, v...)
}

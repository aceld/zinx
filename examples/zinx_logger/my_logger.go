package main

import (
	"context"
	"fmt"
)

//用户自定义日志方式，
//可以通过自身业务的日志方式，来重置zinx内部引擎的日志打印方式
//本例以fmt.Println为例
type MyLogger struct{}

//没有context的日志接口
func (l *MyLogger) InfoF(format string, v ...interface{}) {
	fmt.Printf(format, v...)
}

func (l *MyLogger) ErrorF(format string, v ...interface{}) {
	fmt.Printf(format, v...)
}

func (l *MyLogger) DebugF(format string, v ...interface{}) {
	fmt.Printf(format, v...)
}

//携带context的日志接口
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

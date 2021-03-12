// Package zlog 主要提供zinx相关日志记录接口
// 包括:
//		stdzlog模块， 提供全局日志方法
//		zlogger模块,  日志内部定义协议，均为对象类方法
//
// 当前文件描述:
// @Title  stdzlog.go
// @Description    包裹zlogger日志方法，提供全局方法
// @Author  Aceld - Thu Mar 11 10:32:29 CST 2019
package zlog

/*
   全局默认提供一个Log对外句柄，可以直接使用API系列调用
   全局日志对象 StdZinxLog
*/

import "os"

//StdZinxLog 创建全局log
var StdZinxLog = NewZinxLog(os.Stderr, "", BitDefault)

//Flags 获取StdZinxLog 标记位
func Flags() int {
	return StdZinxLog.Flags()
}

//ResetFlags 设置StdZinxLog标记位
func ResetFlags(flag int) {
	StdZinxLog.ResetFlags(flag)
}

//AddFlag 添加flag标记
func AddFlag(flag int) {
	StdZinxLog.AddFlag(flag)
}

//SetPrefix 设置StdZinxLog 日志头前缀
func SetPrefix(prefix string) {
	StdZinxLog.SetPrefix(prefix)
}

//SetLogFile 设置StdZinxLog绑定的日志文件
func SetLogFile(fileDir string, fileName string) {
	StdZinxLog.SetLogFile(fileDir, fileName)
}

//CloseDebug 设置关闭debug
func CloseDebug() {
	StdZinxLog.CloseDebug()
}

//OpenDebug 设置打开debug
func OpenDebug() {
	StdZinxLog.OpenDebug()
}

//Debugf ====> Debug <====
func Debugf(format string, v ...interface{}) {
	StdZinxLog.Debugf(format, v...)
}

//Debug Debug
func Debug(v ...interface{}) {
	StdZinxLog.Debug(v...)
}

//Infof ====> Info <====
func Infof(format string, v ...interface{}) {
	StdZinxLog.Infof(format, v...)
}

//Info -
func Info(v ...interface{}) {
	StdZinxLog.Info(v...)
}

// ====> Warn <====
func Warnf(format string, v ...interface{}) {
	StdZinxLog.Warnf(format, v...)
}

func Warn(v ...interface{}) {
	StdZinxLog.Warn(v...)
}

// ====> Error <====
func Errorf(format string, v ...interface{}) {
	StdZinxLog.Errorf(format, v...)
}

func Error(v ...interface{}) {
	StdZinxLog.Error(v...)
}

// ====> Fatal 需要终止程序 <====
func Fatalf(format string, v ...interface{}) {
	StdZinxLog.Fatalf(format, v...)
}

func Fatal(v ...interface{}) {
	StdZinxLog.Fatal(v...)
}

// ====> Panic  <====
func Panicf(format string, v ...interface{}) {
	StdZinxLog.Panicf(format, v...)
}

func Panic(v ...interface{}) {
	StdZinxLog.Panic(v...)
}

// ====> Stack  <====
func Stack(v ...interface{}) {
	StdZinxLog.Stack(v...)
}

func init() {
	//因为StdZinxLog对象 对所有输出方法做了一层包裹，所以在打印调用函数的时候，比正常的logger对象多一层调用
	//一般的zinxLogger对象 calldDepth=2, StdZinxLog的calldDepth=3
	StdZinxLog.calldDepth = 3
}

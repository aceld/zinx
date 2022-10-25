package ziface

// ILog 日志接口
type ILog interface {
	Info(args ...interface{})  // 消息日志
	Warn(args ...interface{})  // 警告日志
	Error(args ...interface{}) // 普通错误日志
	Fatal(args ...interface{}) // 致命错误日志
	Panic(args ...interface{}) // 崩溃日志
	Debug(args ...interface{}) // Debug日志
}

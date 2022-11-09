package ziface

// ILogger 日志接口
type ILogger interface {
	Infof(format string, args ...interface{})  // 消息日志
	Warnf(format string, args ...interface{})  // 警告日志
	Errorf(format string, args ...interface{}) // 普通错误日志
	Fatalf(format string, args ...interface{}) // 致命错误日志
	Panicf(format string, args ...interface{}) // 崩溃日志
	Debugf(format string, args ...interface{}) // Debug日志
}

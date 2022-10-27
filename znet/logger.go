package znet

import (
	"github.com/chnkenc/zinx-xiaoan/ziface"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	logger          ziface.ILogger
	defaultLogLevel LogLevel
)

type LogLevel uint8

// SetLogger 设置日志处理器
func SetLogger(logHandler ziface.ILogger, level LogLevel) {
	if logHandler == nil {
		logger = LoadZapLogger(level)
	} else {
		logger = logHandler
	}
}

// GetLogger 获取日志处理器
func GetLogger() ziface.ILogger {
	return logger
}

// LoadZapLogger 加载zap日志处理器
// 默认使用zap sugar
func LoadZapLogger(level LogLevel) ziface.ILogger {
	zapConfig := zap.NewDevelopmentConfig()
	zapLogLevel := zapcore.Level(level)
	zapConfig.Level = zap.NewAtomicLevelAt(zapLogLevel)
	zapConfig.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("[2006-01-02 15:04:05.000]")
	zapLogger, _ := zapConfig.Build()
	return zapLogger.Sugar()
}

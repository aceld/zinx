package zlog_test

import (
	"testing"

	"github.com/aceld/zinx/zlog"
)

func TestStdZLog(t *testing.T) {

	//测试 默认debug输出
	zlog.Debug("zinx debug content1")
	zlog.Debug("zinx debug content2")

	zlog.Debugf(" zinx debug a = %d\n", 10)

	//设置log标记位，加上长文件名称 和 微秒 标记
	zlog.ResetFlags(zlog.BitDate | zlog.BitLongFile | zlog.BitLevel)
	zlog.Info("zinx info content")

	//设置日志前缀，主要标记当前日志模块
	zlog.SetPrefix("MODULE")
	zlog.Error("zinx error content")

	//添加标记位
	zlog.AddFlag(zlog.BitShortFile | zlog.BitTime)
	zlog.Stack(" Zinx Stack! ")

	//设置日志写入文件
	zlog.SetLogFile("./log", "testfile.log")
	zlog.Debug("===> zinx debug content ~~666")
	zlog.Debug("===> zinx debug content ~~888")
	zlog.Error("===> zinx Error!!!! ~~~555~~~")

	//调试隔离级别
	zlog.Debug("=================================>")
	//1.debug
	zlog.SetLogLevel(zlog.LogInfo)
	zlog.Debug("===> 调试Debug：debug不应该出现")
	zlog.Info("===> 调试Debug：info应该出现")
	zlog.Warn("===> 调试Debug：warn应该出现")
	zlog.Error("===> 调试Debug：error应该出现")
	//2.info
	zlog.SetLogLevel(zlog.LogWarn)
	zlog.Debug("===> 调试Info：debug不应该出现")
	zlog.Info("===> 调试Info：info不应该出现")
	zlog.Warn("===> 调试Info：warn应该出现")
	zlog.Error("===> 调试Info：error应该出现")
	//3.warn
	zlog.SetLogLevel(zlog.LogError)
	zlog.Debug("===> 调试Warn：debug不应该出现")
	zlog.Info("===> 调试Warn：info不应该出现")
	zlog.Warn("===> 调试Warn：warn不应该出现")
	zlog.Error("===> 调试Warn：error应该出现")
	//4.error
	zlog.SetLogLevel(zlog.LogPanic)
	zlog.Debug("===> 调试Error：debug不应该出现")
	zlog.Info("===> 调试Error：info不应该出现")
	zlog.Warn("===> 调试Error：warn不应该出现")
	zlog.Error("===> 调试Error：error不应该出现")
}

func TestZLogger(t *testing.T) {
}

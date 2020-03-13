package zlog

import (
	"testing"
)

func TestStdZLog(t *testing.T) {

	//测试 默认debug输出
	Debug("zinx debug content1")
	Debug("zinx debug content2")

	Debugf(" zinx debug a = %d\n", 10)

	//设置log标记位，加上长文件名称 和 微秒 标记
	ResetFlags(BitDate | BitLongFile | BitLevel)
	Info("zinx info content")

	//设置日志前缀，主要标记当前日志模块
	SetPrefix("MODULE")
	Error("zinx error content")

	//添加标记位
	AddFlag(BitShortFile | BitTime)
	Stack(" Zinx Stack! ")

	//设置日志写入文件
	SetLogFile("./log", "testfile.log")
	Debug("===> zinx debug content ~~666")
	Debug("===> zinx debug content ~~888")
	Error("===> zinx Error!!!! ~~~555~~~")

	//关闭debug调试
	CloseDebug()
	Debug("===> 我不应该出现~！")
	Debug("===> 我不应该出现~！")
	Error("===> zinx Error  after debug close !!!!")

}

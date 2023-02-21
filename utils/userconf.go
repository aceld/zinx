package utils

import (
	"github.com/aceld/zinx/ziface"
)

type Config struct {

	//Server
	TcpServer  ziface.IServer //当前全局的Server对象
	Host       string         //当前服务器主机监听的IP
	TcpPort    int            //当前服务器监听的端口
	Name       string         //当前服务器的名称
	TcpVersion string         //tcp版本

	//服务器可选配置
	Version          string //版本
	MaxConn          int    //最大连接数量
	MaxPacketSize    uint32 //当前框架数据包的最大尺寸
	WorkerPoolSize   uint32 //业务工作Worker池的数量
	MaxWorkerTaskLen uint32 //业务工作Worker对应负责的任务队列最大任务存储数量
	MaxMsgChanLen    uint32 //SendBuffMsg发送消息的缓冲最大长度

	/*
		logger
	*/
	LogDir        string //日志所在文件夹 默认"./log"
	LogFile       string //日志文件名称   默认""  --如果没有设置日志文件，打印信息将打印至stderr
	LogDebugClose bool   //是否关闭Debug日志级别调试信息 默认false  -- 默认打开debug信息
}

//注意如果使用UserConf应该调用方法同步至 GlobalConfObject 因为其他参数是调用的此结构体参数
func UserConfToGlobal(config *Config) {

	//Server配置
	GlobalObject.Name = config.Name
	GlobalObject.Host = config.Host
	GlobalObject.TCPPort = config.TcpPort

	//Zinx配置项设置
	GlobalObject.Version = config.Version
	GlobalObject.WorkerPoolSize = config.WorkerPoolSize
	GlobalObject.MaxConn = config.MaxConn
	GlobalObject.MaxPacketSize = config.MaxPacketSize
	GlobalObject.MaxMsgChanLen = config.MaxWorkerTaskLen
	GlobalObject.MaxMsgChanLen = config.MaxMsgChanLen

	//日志配置项目
	//默认就是False config没有初始化即使用默认配置
	GlobalObject.LogDebugClose = config.LogDebugClose
	//不同于上方必填项 日志目前如果没配置应该使用默认配置
	if config.LogDir != "" {
		GlobalObject.LogDir = config.LogDir
	}
	if config.LogFile != "" {
		GlobalObject.LogFile = config.LogFile
	}

}

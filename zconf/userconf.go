package zconf

import "github.com/aceld/zinx/zlog"

// UserConfToGlobal, Note that if UserConf is used,
// the method should be called to synchronize with GlobalConfObject
// because other parameters are called from this structure parameter.
// (注意如果使用UserConf应该调用方法同步至 GlobalConfObject 因为其他参数是调用的此结构体参数)
func UserConfToGlobal(config *Config) {

	// Server
	if config.Name != "" {
		GlobalObject.Name = config.Name
	}
	if config.Host != "" {
		GlobalObject.Host = config.Host
	}
	if config.TCPPort != 0 {
		GlobalObject.TCPPort = config.TCPPort
	}

	// Zinx
	if config.Version != "" {
		GlobalObject.Version = config.Version
	}
	if config.MaxPacketSize != 0 {
		GlobalObject.MaxPacketSize = config.MaxPacketSize
	}
	if config.MaxConn != 0 {
		GlobalObject.MaxConn = config.MaxConn
	}
	if config.WorkerPoolSize != 0 {
		GlobalObject.WorkerPoolSize = config.WorkerPoolSize
	}
	if config.MaxWorkerTaskLen != 0 {
		GlobalObject.MaxWorkerTaskLen = config.MaxWorkerTaskLen
	}
	if config.WorkerMode != "" {
		GlobalObject.WorkerMode = config.WorkerMode
	}

	if config.MaxMsgChanLen != 0 {
		GlobalObject.MaxMsgChanLen = config.MaxMsgChanLen
	}
	if config.IOReadBuffSize != 0 {
		GlobalObject.IOReadBuffSize = config.IOReadBuffSize
	}

	// logger
	// By default, it is False. If the config is not initialized, the default configuration will be used.
	// (默认是False, config没有初始化即使用默认配置)
	GlobalObject.LogIsolationLevel = config.LogIsolationLevel
	if GlobalObject.LogIsolationLevel > zlog.LogDebug {
		zlog.SetLogLevel(GlobalObject.LogIsolationLevel)
	}

	// Different from the required fields mentioned above, the logging module should use the default configuration if it is not configured.
	// (不同于上方必填项 日志目前如果没配置应该使用默认配置)
	if config.LogDir != "" {
		GlobalObject.LogDir = config.LogDir
	}

	if config.LogFile != "" {
		GlobalObject.LogFile = config.LogFile
		zlog.SetLogFile(GlobalObject.LogDir, GlobalObject.LogFile)
	}

	// Keepalive
	if config.HeartbeatMax != 0 {
		GlobalObject.HeartbeatMax = config.HeartbeatMax
	}

	// TLS
	if config.CertFile != "" {
		GlobalObject.CertFile = config.CertFile
	}
	if config.PrivateKeyFile != "" {
		GlobalObject.PrivateKeyFile = config.PrivateKeyFile
	}

	if config.Mode != "" {
		GlobalObject.Mode = config.Mode
	}
	if config.WsPort != 0 {
		GlobalObject.WsPort = config.WsPort
	}

	if config.RouterSlicesMode {
		GlobalObject.RouterSlicesMode = config.RouterSlicesMode
	}

	if config.KcpPort != 0 {
		GlobalObject.KcpPort = config.KcpPort
	}

	if config.KcpACKNoDelay {
		GlobalObject.KcpACKNoDelay = config.KcpACKNoDelay
	}

	if !config.KcpStreamMode {
		GlobalObject.KcpStreamMode = config.KcpStreamMode
	}

	if config.KcpNoDelay != 0 {
		GlobalObject.KcpNoDelay = config.KcpNoDelay
	}

	if config.KcpInterval != 0 {
		GlobalObject.KcpInterval = config.KcpInterval
	}

	if config.KcpResend != 0 {
		GlobalObject.KcpResend = config.KcpResend
	}

	if config.KcpNc != 0 {
		GlobalObject.KcpNc = config.KcpNc
	}

	if config.KcpSendWindow != 0 {
		GlobalObject.KcpSendWindow = config.KcpSendWindow
	}

	if config.KcpRecvWindow != 0 {
		GlobalObject.KcpRecvWindow = config.KcpRecvWindow
	}

}

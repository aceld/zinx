package utils

const (
	ZD_CONN_BUFSIZE     = 1024 //一次连接读取缓冲大小
	ZD_CONN_READTIMEOUT = 6000 //读取超时时间(毫秒)

	//zinx unit 集群单元的状态
	ZINX_UNIT_STATUS_ALIVE = true
	ZINX_UNIT_STATUS_DEAD  = false

	//zinx 集群中的角色
	ZINX_ROLE_SERVER  = 0
	ZINX_ROLE_CLIENT  = 1
	ZINX_ROLE_MONITOR = 2

	//API处理HTTP指令端口
	ZINX_API_RET_SUCC = "SUCC"
	ZINX_API_RET_FAIL = "FAIL"

	//API指令CMD
	ZINX_CMD_ID_NODE_ADD    = 100
	ZINX_CMD_ID_NODE_REMOVE = 110

	//集群端口定义
	ZINX_API_PORT  = 17770 //集群对外API端口
	ZINX_SYNC_PORT = 17771 //集群数据同步端口
	ZINX_RAFT_PORT = 17772 //集群raft协商协议端口

	//网卡名称
	//TODO 从配置文件中读取
	INTERNET_DEVICE_NAME = "en5"

	//zinx 分布式版本号
	ZINX_DISTRIBUTED_VERSION = "1.0"
)

func NodeRoleName(role int32) string {
	switch role {
	case ZINX_ROLE_SERVER:
		return "server"
	case ZINX_ROLE_CLIENT:
		return "client"
	case ZINX_ROLE_MONITOR:
		return "monitor"
	}

	return "UNKNOW"
}

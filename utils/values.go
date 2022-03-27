package utils

const (
	ZD_CONN_BUFSIZE     = 1024 //一次连接读取缓冲大小
	ZD_CONN_READTIMEOUT = 6000 //读取超时时间(毫秒)

	//zinx unit 集群单元的状态
	ZINX_UNIT_STATUS_ALIVE = true
	ZINX_UNIT_STATUS_DEAD  = false

	//zinx 集群中的角色
	ZINX_ROLE_SERVER = 0
	ZINX_ROLE_CLIENT = 1

	//raft 集群角色
	ZINX_RAFT_LEADER    = 0
	ZINX_RAFT_CONDIDATE = 1
	ZINX_RAFT_FOLLOWER  = 2

	//API处理HTTP指令端口
	ZINX_API_RET_SUCC     = "SUCC"
	ZINX_API_RET_FAIL     = "FAIL"
	ZINX_API_RETCODE_OK   = 1
	ZINX_API_RETCODE_FAIL = 0

	//API指令CMD
	ZINX_CMD_ID_NODE_SYNC_ACK = 100 //回执消息
	ZINX_CMD_ID_NODE_ADD      = 110
	ZINX_CMD_ID_NODE_REMOVE   = 111

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

//将集群角色标识转换为名称:server/client
func NodeRoleName(role int32) string {
	switch role {
	case ZINX_ROLE_SERVER:
		return "server"
	case ZINX_ROLE_CLIENT:
		return "client"
	}

	return "UNKNOW"
}

//将集群raft角色标识转化为名称:leader/condidate/follower
func NodeRaftName(raftRole int32) string {
	switch raftRole {
	case ZINX_RAFT_LEADER:
		return "leader"
	case ZINX_RAFT_CONDIDATE:
		return "condidate"
	case ZINX_RAFT_FOLLOWER:
		return "follower"
	}

	return "UNKNOW"
}

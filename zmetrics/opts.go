package zmetrics

const (
	GANGEVEC_ZINX_CONNECTION_TOTAL_NAME string = "zinx_connection_online_total"
	GANGEVEC_ZINX_CONNECTION_TOTAL_HELP string = "Zinx All Online Connections Group By (Address, Name) (zinx 不同Server的链接总数,根据(Address, Name) 分组)"

	GANGEVEC_ZINX_TASK_TOTAL_NAME string = "zinx_task_total"
	GANGEVEC_ZINX_TASK_TOTAL_HELP string = "Zinx All Task Total Group By (Address, Name, WorkerID) (zinx 已经处理的数据任务总数,根据(Address, Name, WorkerID)分组)"

	GANGEVEC_ZINX_ROUTER_SCHEDULE_TOTAL_NAME string = "zinx_router_schedule_total"
	GANGEVEC_ZINX_ROUTER_SCHEDULE_TOTAL_HELP string = "Zinx Router Schedule Total Group By (Address, Name, WorkerID, MsgID) (zinx 路由调度的Handler总数,根据(Address, Name, WorkerID, MsgID)分组)"
)

const (
	LABEL_ADDRESS   string = "address"
	LABEL_NAME      string = "name"
	LABEL_WORKER_ID string = "worker_id"
	LABEL_MSG_ID    string = "msg_id"
)

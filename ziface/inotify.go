package ziface

type Inotify interface {
	// Whether there is a connection with this id
	// (是否有这个id)
	HasIdConn(id uint64) bool

	// Get the number of connections stored
	// (存储的map长度)
	ConnNums() int

	// Add a connection
	// (添加链接)
	SetNotifyID(Id uint64, conn IConnection)

	// Get a connection by id
	// (得到某个链接)
	GetNotifyByID(Id uint64) (IConnection, error)

	// Delete a connection by id
	// (删除某个链接)
	DelNotifyByID(Id uint64)

	// Notify a connection with the given id
	// (通知某个id的方法)
	NotifyToConnByID(Id uint64, MsgId uint32, data []byte) error

	// Notify all connections
	// (通知所有人)
	NotifyAll(MsgId uint32, data []byte) error

	// Notify a connection with the given id using a buffer queue
	// (通过缓冲队列通知某个id的方法)
	NotifyBuffToConnByID(Id uint64, MsgId uint32, data []byte) error

	// Notify all connections using a buffer queue
	// (缓冲队列通知所有人)
	NotifyBuffAll(MsgId uint32, data []byte) error
}

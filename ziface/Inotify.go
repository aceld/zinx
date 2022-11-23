package ziface

type Inotify interface {
	HasIdConn(id uint64) bool
	//通知某个id的方法
	NotifyToConnByID(Id uint64, MsgId uint32, data []byte) error
	//通知所有人
	NotifyAll(MsgId uint32, data []byte) error

	//通过缓冲队列通知某个id的方法
	NotifyBuffToConnByID(Id uint64, MsgId uint32, data []byte) error
	//缓冲队列通知所有人
	NotifyBuffAll(MsgId uint32, data []byte) error
}

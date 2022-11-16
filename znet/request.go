package znet

import "github.com/aceld/zinx/ziface"

//Request 请求
type Request struct {
	conn   ziface.IConnection //已经和客户端建立好的 链接
	msg    ziface.IMessage    //客户端请求的数据
	router ziface.IRouter     //请求处理的函数
	index  int8               //用来控制路由函数执行
}

//GetConnection 获取请求连接信息
func (r *Request) GetConnection() ziface.IConnection {
	return r.conn
}

//GetData 获取请求消息的数据
func (r *Request) GetData() []byte {
	return r.msg.GetData()
}

//GetMsgID 获取请求的消息的ID
func (r *Request) GetMsgID() uint32 {
	return r.msg.GetMsgID()
}

func (r *Request) BindRouter(router ziface.IRouter) {
	r.router = router
}

func (r *Request) Next() {
	r.index++
	for r.index < 4 {
		switch r.index {
		case 1:
			r.router.PreHandle(r)
		case 2:
			r.router.Handle(r)
		case 3:
			r.router.PostHandle(r)
		}
		r.index++
	}
}

func (r *Request) Abort() {
	r.index = 4
}

package zdnet

import (
	"fmt"
	"net"
)

/*
   主要用来针对ZINX_SYNC_PORT 集群数据同步端口的业务消息处理
*/

type ZdTcpServer struct {
	//服务器IP
	Ip string
	//服务器端口
	Port int

	//对于每个链接的处理业务
	Handler func(*ZDConn)
}

func NewZdTcpServer(port int, handler func(*ZDConn)) *ZdTcpServer {
	s := &ZdTcpServer{
		Ip:      "127.0.0.1",
		Port:    port,
		Handler: handler,
	}

	return s
}

func (s *ZdTcpServer) Start() {

	fmt.Printf("[START] ZdTcpServer listenner at IP: %s, Port %d is starting\n", s.Ip, s.Port)

	//开启一个go去做服务端Linster业务
	//1 获取一个TCP的Addr
	addr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", s.Ip, s.Port))
	if err != nil {
		fmt.Println("resolve tcp addr err: ", err)
		return
	}

	//2 监听服务器地址
	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		fmt.Println("listen tcp", "err", err)
		return
	}

	//已经监听成功
	fmt.Println("start Zinx SyncServer  succ, now listenning...")

	//3 启动server网络连接业务
	for {
		//3.1 阻塞等待客户端建立连接请求
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Accept err ", err)
			continue
		}

		zdConn := MakeZDConn(conn)

		go s.Handler(zdConn)
	}
}

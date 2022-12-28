# <img width="80px" src="https://s2.ax1x.com/2019/10/09/u4yHo9.png" />
English | [简体中文](README-CN.md)

[![License](https://img.shields.io/badge/License-GPL%203.0-black.svg)](LICENSE)
[![Discord](https://img.shields.io/badge/zinx-Discord-blue.svg)](https://discord.gg/X7BUn6bT)
[![Gitter](https://img.shields.io/badge/zinx-Gitter-green.svg)](https://gitter.im/zinx_go/community) 
[![zinx tutorial](https://img.shields.io/badge/ZinxTutorial-YuQue-red.svg)](https://www.yuque.com/aceld/npyr8s/bgftov) 
[![Original Book of Zinx](https://img.shields.io/badge/OriginalBook-YuQue-black.svg)](https://www.yuque.com/aceld)


Zinx is a lightweight concurrent server framework based on Golang.

Website:http://zinx.me

> **ps**:   
> Zinx has been developed and used in many enterprises: Service of message transfer, Persistent Connection TCP/IP Server, The middleware of Web Service and so on. 
> Zinx is positioned for code simplicity, Developers can use Zinx to redevelop a module suitable for their own enterprise scenarios.

---
## Source of Zinx
### Github
Git: https://github.com/aceld/zinx

### Gitee(China)
Git: https://gitee.com/Aceld/zinx

---

## Online Tutorial

| platform | online video | 
| ---- | ---- | 
| <img src="https://s1.ax1x.com/2022/09/22/xFePUK.png" width = "100" height = "100" alt="" align=center />| [![zinx-BiliBili](https://s2.ax1x.com/2019/10/13/uv340S.jpg)](https://www.bilibili.com/video/av71067087)| 
| <img src="https://s1.ax1x.com/2022/09/22/xFeRVx.png" width = "100" height = "100" alt="" align=center />  | [![zinx-BiliBili](https://s2.ax1x.com/2019/10/13/uv340S.jpg)](https://www.douyin.com/video/6983301202939333891) |
| <img src="https://s1.ax1x.com/2022/09/23/xkQcng.png" width = "100" height = "100" alt="" align=center />| [![zinx-youtube](https://s2.ax1x.com/2019/10/14/KSurCR.jpg)](https://www.youtube.com/watch?v=U95iF-HMWsU&list=PL_GrAPKmuajzeNI8HBTi-k5NQO1g0rM-A)| 


## The Document of Zinx

### PC terminal

[YuQue - Zinx Framework tutorial-Lightweight server based on Golang》](https://www.yuque.com/aceld)


### Mobile terminal(WeChat)
![gongzhonghao](https://s1.ax1x.com/2020/07/07/UFyUdx.th.jpg)


## I. One word that has been said before

Why are we doing Zinx? There are a lot of frameworks for servers in The Go Programing Language, but very few lightweight enterprise frameworks for games or other persistent connection TCP/IP Server domains.

Zinx is designed so that developers can use the Zinx framework to understand the overall outline of writing a TCP server based on Golang, Let more Gopher can learn and understand this field in a simple way.

The Zinx framework projects are done in parallel with the coding and learning tutorials, bringing all of the progressive and iterative thinking of development into the tutorials, rather than giving everyone a very complete framework to learn at once, leaving many people confused about how to learn.

The tutorial will iterate from release to release, with minor additions to each release, giving a small, curvewise approach to the domain of the server framework.

Of course, I hope that more people will join Zinx and give us valuable suggestions, so that Zinx can become a real solution server framework for enterprises! Thank you for your attention!

### Reply from chatGPT(AI)
![what-is-zinx](https://user-images.githubusercontent.com/7778936/209745848-acfc14eb-74cd-4513-b386-8bc6e0bcc09f.png)

![compare-zinx](https://user-images.githubusercontent.com/7778936/209745864-7d8984b0-bd73-4109-b4ec-aec152f8f8e8.png)


### The honor of zinx
#### GVP Most Valuable Open Source Project of the Year at OSCHINA

![GVP-zinx](https://s2.ax1x.com/2019/10/13/uvYVBV.jpg)



#### Stargazers over time

[![Stargazers over time](https://api.star-history.com/svg?repos=aceld/zinx&type=Date)](#zinx)


## II. Quick start

**Version**
Golang 1.16+

```bash
# clone from git
$ git clone https://github.com/aceld/zinx.git

# cd the dir of Demo
$ cd ./zinx/examples/zinx_server

# Build
$ make build

# Build for docker image
$ make image

# start and run
$ make run 

# cd the dir of Demo Client
$ cd ../zinx_client

# run 
$ go run main.go 

```


## III. Zinx architecture

![1-Zinx框架.png](https://camo.githubusercontent.com/903d1431358fa6f4634ebaae3b49a28d97e23d77/68747470733a2f2f75706c6f61642d696d616765732e6a69616e7368752e696f2f75706c6f61645f696d616765732f31313039333230352d633735666636383232333362323533362e706e673f696d6167654d6f6772322f6175746f2d6f7269656e742f7374726970253743696d61676556696577322f322f772f31323430)

![zinx-start](https://user-images.githubusercontent.com/7778936/126594039-98dddd10-ec6a-4881-9e06-a09ec34f1af7.gif)



## IV. Zinx development API documentation

### (1) Quick start

#### A. Demo
1. Compile Demo example, in dir `zinx/example/zinx_server`, we get `server`, in `zinx/example/zinx_client`, we get`client`.
```bash
$ cd zinx/
$ make
```
2. run Demo server (don't close the terminal)
```bash
$ cd example/zinx_server
$ ./server 
                                        
              ██                        
              ▀▀                        
 ████████   ████     ██▄████▄  ▀██  ██▀ 
     ▄█▀      ██     ██▀   ██    ████   
   ▄█▀        ██     ██    ██    ▄██▄   
 ▄██▄▄▄▄▄  ▄▄▄██▄▄▄  ██    ██   ▄█▀▀█▄  
 ▀▀▀▀▀▀▀▀  ▀▀▀▀▀▀▀▀  ▀▀    ▀▀  ▀▀▀  ▀▀▀ 
                                        
┌───────────────────────────────────────────────────┐
│ [Github] https://github.com/aceld                 │
│ [tutorial] https://www.kancloud.cn/aceld/zinx     │
└───────────────────────────────────────────────────┘
[Zinx] Version: V0.11, MaxConn: 3, MaxPacketSize: 4096
Add api msgId =  0
Add api msgId =  1
[START] Server name: zinx server Demo,listenner at IP: 127.0.0.1, Port 8999 is starting
Worker ID =  0  is started.
Worker ID =  1  is started.
Worker ID =  2  is started.
Worker ID =  3  is started.
Worker ID =  4  is started.
Worker ID =  7  is started.
Worker ID =  6  is started.
Worker ID =  8  is started.
Worker ID =  9  is started.
Worker ID =  5  is started.
start Zinx server   zinx server Demo  succ, now listenning...
...
```

3. Then open the new terminal and start the Client Demo to test communication
```bash
$ cd example/zinx_client
$ ./client
==> Test Router:[Ping] Recv Msg: ID= 2 , len= 21 , data= DoConnection BEGIN... ==> Test Router:[Ping] Recv Msg: ID= 0 , len= 18 , data= ping...ping...ping 
==> Test Router:[Ping] Recv Msg: ID= 0 , len= 18 , data= ping...ping...ping
==> Test Router:[Ping] Recv Msg: ID= 0 , len= 18 , data= ping...ping...ping
...
t

```

#### B. server

In the server application developed based on Zinx framework, the main function steps are relatively simple and only need 3 steps at most.

1. Create the server object
2. Configure user-defined routes and services
3. Start the service

```go
func main() {
	//1 Create the server object
	s := znet.NewServer()

	//2 Configure user-defined routes and services
	s.AddRouter(0, &PingRouter{})

	//3 Start the service
	s.Serve()
}
```

The custom route and service configuration methods are as follows：
```go
import (
	"fmt"
	"github.com/aceld/zinx/ziface"
	"github.com/aceld/zinx/znet"
)

//ping test custom route
type PingRouter struct {
	znet.BaseRouter
}

//Ping Handle
func (this *PingRouter) Handle(request ziface.IRequest) {
	//Read the data from the client first
	fmt.Println("recv from client : msgId=", request.GetMsgID(), ", data=", string(request.GetData()))
	
	//To go back to write  "ping...ping...ping"
	err := request.GetConnection().SendBuffMsg(0, []byte("ping...ping...ping"))
	if err != nil {
		fmt.Println(err)
	}
}
```

#### C. client
Zinx's message packet format processing uses `[MsgLength]|[MsgID]|[Data]` .
```go
package main

import (
	"fmt"
	"io"
	"net"
	"time"
	"github.com/aceld/zinx/znet"
)

/*
	Simulation Client
*/
func main() {

	fmt.Println("Client Test ... start")
	//A test request is made after 3 seconds, giving the server a chance to start the service
	time.Sleep(3 * time.Second)

	conn,err := net.Dial("tcp", "127.0.0.1:7777")
	if err != nil {
		fmt.Println("client start err, exit!")
		return
	}

	for n := 3; n >= 0; n-- {
		//Send a packet message
		dp := znet.NewDataPack()
		msg, _ := dp.Pack(znet.NewMsgPackage(0,[]byte("Zinx Client Test Message")))
		_, err := conn.Write(msg)
		if err !=nil {
			fmt.Println("write error err ", err)
			return
		}

		//Read the head part of the stream first
		headData := make([]byte, dp.GetHeadLen())
		_, err = io.ReadFull(conn, headData) //ReadFull 会把msg填充满为止
		if err != nil {
			fmt.Println("read head error")
			break
		}
		//Unpack the headData byte stream into MSG
		msgHead, err := dp.Unpack(headData)
		if err != nil {
			fmt.Println("server unpack err:", err)
			return
		}

		if msgHead.GetDataLen() > 0 {
			//msg has data data, which needs to be read again
			msg := msgHead.(*znet.Message)
			msg.Data = make([]byte, msg.GetDataLen())

			//Read the byte stream from the IO according to dataLen
			_, err := io.ReadFull(conn, msg.Data)
			if err != nil {
				fmt.Println("server unpack data err:", err)
				return
			}

			fmt.Println("==> Recv Msg: ID=", msg.Id, ", len=", msg.DataLen, ", data=", string(msg.Data))
		}

		time.Sleep(1*time.Second)
	}
}
```

### (2) Zinx configuration file
```json
{
  "Name":"zinx v-0.10 demoApp",
  "Host":"127.0.0.1",
  "TcpPort":7777,
  "MaxConn":3,
  "WorkerPoolSize":10,
  "LogDir": "./mylog",
  "LogFile":"zinx.log"
}
```

`Name`:Server Application Name

`Host`:Server IP

`TcpPort`:Server listening port

`MaxConn`:Maximum number of client links allowed

`WorkerPoolSize`:Maximum number of working Goroutines in the work task pool

`LogDir`: Log folder

`LogFile`: Log file name (if not provided, log information is printed to Stderr)


### (3) Server Module 
```go
  func NewServer () ziface.IServer 
```
Create a Zinx server object that serves as the primary hub for the current server application, including the following functions:

#### A. Start the Server
```go
  func (s *Server) Start()
```
#### B. Stop the Server
```go
  func (s *Server) Stop()
```
#### C. Run the Server
```go
  func (s *Server) Serve()
```
#### D. Registered router
```go
func (s *Server) AddRouter (msgId uint32, router ziface.IRouter) 
```
#### E. Register the link to create the Hook function
```go
func (s *Server) SetOnConnStart(hookFunc func (ziface.IConnection))
```
#### F. Register the link destruction Hook function
```go
func (s *Server) SetOnConnStop(hookFunc func (ziface.IConnection))
```

### (4) Router Module
```go
//When you implement Router, you embed the base class and then override the methods of the base class as needed.
type BaseRouter struct {}

//The BaseRouter's methods are null because some Router does not want to
//have PreHandle or PostHandle. 
//The Router inherits all BaseRouter's methods because PreHandle and PostHandle can be instantiated 
//without implementing them
func (br *BaseRouter)PreHandle(req ziface.IRequest){}
func (br *BaseRouter)Handle(req ziface.IRequest){}
func (br *BaseRouter)PostHandle(req ziface.IRequest){}
```


### (5) Connection Module
#### A. Get the socket net.TCPConn
```go
  func (c *Connection) GetTCPConnection() *net.TCPConn 
```
#### B. Get the Connection ID
```go
  func (c *Connection) GetConnID() uint32 
```
#### C. Get the remote client address 
```go
  func (c *Connection) RemoteAddr() net.Addr 
```
#### D. send message
```go
  func (c *Connection) SendMsg(msgId uint32, data []byte) error 
  func (c *Connection) SendBuffMsg(msgId uint32, data []byte) error
```
#### E. Connection attributes
```go
//Setting connection attributes
func (c *Connection) SetProperty(key string, value interface{})

//Getting connection attributes
func (c *Connection) GetProperty(key string) (interface{}, error)

//remove connection attributes
func (c *Connection) RemoveProperty(key string) 
```

---

#### Developers
-   刘丹冰([@aceld](https://github.com/aceld))
-   张超([@zhngcho](https://github.com/zhngcho))
-   高智辉Roger([@adsian](https://github.com/adsian))
-   胡贵建([@huguijian](https://github.com/huguijian))
-   张继瑀([@kstwoak](https://github.com/huguijian))


---
[zinx(with C++)](https://github.com/marklion/zinx)
#### Developers
-  刘洋([@marklion](https://github.com/marklion))


---
[zinx(with Lua)](https://github.com/huqitt/zinx-lua)
#### Developers
-  胡琪([@huqitt](https://github.com/huqitt))

---
[zinx(for websocket)](https://github.com/aceld/zinx/tree/wsserver)
#### Developers
-  胡贵建([@huguijian](https://github.com/huguijian))

---

Thanks to all the developers who contributed to Zinx!

<a href="https://github.com/aceld/zinx/graphs/contributors">
  <img src="https://contrib.rocks/image?repo=aceld/zinx" />
</a>    


---
### About the author

`name`：`Aceld(刘丹冰)`

`mail`:
[danbing.at@gmail.com](mailto:danbing.at@gmail.com)

`github`:
[https://github.com/aceld](https://github.com/aceld)

`original work`:
[https://www.yuque.com/aceld](https://www.yuque.com/aceld)

### Zinx Technical Discussion Community
|  **WeChat**   | **WeChat Public Account**  | **QQ Group**  |
|  ----  | ----  | ----  |
| <img src="https://s1.ax1x.com/2020/07/07/UF6rNV.png" width = "150" height = "180" alt="" align=center />  | <img src="https://s1.ax1x.com/2020/07/07/UFyUdx.th.jpg" height = "150"  alt="" align=center /> | <img src="https://s1.ax1x.com/2020/07/07/UF6Y9S.th.png" width = "150" height = "150" alt="" align=center /> |

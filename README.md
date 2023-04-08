# <img width="80px" src="https://s2.ax1x.com/2019/10/09/u4yHo9.png" />
English | [简体中文](README-CN.md)

[![License](https://img.shields.io/badge/License-GPL%203.0-black.svg)](LICENSE)
[![Discord](https://img.shields.io/badge/zinx-Discord-blue.svg)](https://discord.gg/rM2sw5uv)
[![Gitter](https://img.shields.io/badge/zinx-Gitter-green.svg)](https://gitter.im/zinx_go/community) 
[![zinx tutorial](https://img.shields.io/badge/ZinxTutorial-YuQue-red.svg)](https://www.yuque.com/aceld/npyr8s/bgftov) 
[![Original Book of Zinx](https://img.shields.io/badge/OriginalBook-YuQue-black.svg)](https://www.yuque.com/aceld)


Zinx is a lightweight concurrent server framework based on Golang.

##  Document 

[《Zinx Documentation》](https://www.yuque.com/aceld/tsgooa/sbvzgczh3hqz8q3l)


> **ps**:   
> Zinx has been developed and used in many enterprises: Service of message transfer, Persistent Connection TCP/IP Server, The middleware of Web Service and so on. 
> Zinx is positioned for code simplicity, Developers can use Zinx to redevelop a module suitable for their own enterprise scenarios.

---
## Source of Zinx
### Github
Git: https://github.com/aceld/zinx

### Gitee(China)
Git: https://gitee.com/Aceld/zinx

### Website
http://zinx.me

---

## Online Tutorial

[YuQue - 《Zinx Framework tutorial-Lightweight server based on Golang》](https://www.yuque.com/aceld)

| platform | online video | 
| ---- | ---- | 
| <img src="https://s1.ax1x.com/2022/09/22/xFePUK.png" width = "100" height = "100" alt="" align=center />| [![zinx-BiliBili](https://s2.ax1x.com/2019/10/13/uv340S.jpg)](https://www.bilibili.com/video/av71067087)| 
| <img src="https://s1.ax1x.com/2022/09/22/xFeRVx.png" width = "100" height = "100" alt="" align=center />  | [![zinx-BiliBili](https://s2.ax1x.com/2019/10/13/uv340S.jpg)](https://www.douyin.com/video/6983301202939333891) |
| <img src="https://s1.ax1x.com/2022/09/23/xkQcng.png" width = "100" height = "100" alt="" align=center />| [![zinx-youtube](https://s2.ax1x.com/2019/10/14/KSurCR.jpg)](https://www.youtube.com/watch?v=U95iF-HMWsU&list=PL_GrAPKmuajzeNI8HBTi-k5NQO1g0rM-A)| 


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


## II. Zinx architecture
![Zinx框架](https://user-images.githubusercontent.com/7778936/220058446-0ad45112-2225-4b71-b0d8-69a7f3cee5ca.jpg)

![流程图](https://raw.githubusercontent.com/wenyoufu/testaaaaaa/master/%E6%B5%81%E7%A8%8B%E5%9B%BE-en.jpg)
![zinx-start](https://user-images.githubusercontent.com/7778936/126594039-98dddd10-ec6a-4881-9e06-a09ec34f1af7.gif)



## III. Zinx development API documentation


## The Document of Zinx

[YuQue - 《Zinx Documentation》](https://www.yuque.com/aceld/tsgooa/sbvzgczh3hqz8q3l)

### (1) QuickStart

[<Zinx's TCP Debugging Tool>](https://github.com/xxl6097/tcptest)

DownLoad zinx Source

```bash
$go get github.com/aceld/zinx
```

> note: Golang Version 1.16+

#### Zinx-Server
```go
package main

import (
	"fmt"
	"github.com/aceld/zinx/ziface"
	"github.com/aceld/zinx/znet"
)

// PingRouter MsgId=1 
type PingRouter struct {
	znet.BaseRouter
}

//Ping Handle MsgId=1
func (r *PingRouter) Handle(request ziface.IRequest) {
	//read client data
	fmt.Println("recv from client : msgId=", request.GetMsgID(), ", data=", string(request.GetData()))
}

func main() {
	//1 Create a server service
	s := znet.NewServer()

	//2 configure routing
	s.AddRouter(1, &PingRouter{})

	//3 start service
	s.Serve()
}
```

Run Server

```bash
$ go run server.go
```

```bash
                                        
              ██                        
              ▀▀                        
 ████████   ████     ██▄████▄  ▀██  ██▀ 
     ▄█▀      ██     ██▀   ██    ████   
   ▄█▀        ██     ██    ██    ▄██▄   
 ▄██▄▄▄▄▄  ▄▄▄██▄▄▄  ██    ██   ▄█▀▀█▄  
 ▀▀▀▀▀▀▀▀  ▀▀▀▀▀▀▀▀  ▀▀    ▀▀  ▀▀▀  ▀▀▀ 
                                        
┌──────────────────────────────────────────────────────┐
│ [Github] https://github.com/aceld                    │
│ [tutorial] https://www.yuque.com/aceld/npyr8s/bgftov │
└──────────────────────────────────────────────────────┘
[Zinx] Version: V1.0, MaxConn: 12000, MaxPacketSize: 4096
===== Zinx Global Config =====
TCPServer: <nil>
Host: 0.0.0.0
TCPPort: 8999
Name: ZinxServerApp
Version: V1.0
MaxPacketSize: 4096
MaxConn: 12000
WorkerPoolSize: 10
MaxWorkerTaskLen: 1024
MaxMsgChanLen: 1024
ConfFilePath: /Users/Aceld/go/src/zinx-usage/quick_start/conf/zinx.json
LogDir: /Users/Aceld/go/src/zinx-usage/quick_start/log
LogFile: 
LogDebugClose: false
HeartbeatMax: 10
==============================
2023/03/09 18:39:49 [INFO]msghandler.go:61: Add api msgID = 1
2023/03/09 18:39:49 [INFO]server.go:112: [START] Server name: ZinxServerApp,listenner at IP: 0.0.0.0, Port 8999 is starting
2023/03/09 18:39:49 [INFO]msghandler.go:66: Worker ID = 0 is started.
2023/03/09 18:39:49 [INFO]msghandler.go:66: Worker ID = 1 is started.
2023/03/09 18:39:49 [INFO]msghandler.go:66: Worker ID = 3 is started.
2023/03/09 18:39:49 [INFO]msghandler.go:66: Worker ID = 2 is started.
2023/03/09 18:39:49 [INFO]msghandler.go:66: Worker ID = 4 is started.
2023/03/09 18:39:49 [INFO]msghandler.go:66: Worker ID = 6 is started.
2023/03/09 18:39:49 [INFO]msghandler.go:66: Worker ID = 7 is started.
2023/03/09 18:39:49 [INFO]msghandler.go:66: Worker ID = 8 is started.
2023/03/09 18:39:49 [INFO]msghandler.go:66: Worker ID = 9 is started.
2023/03/09 18:39:49 [INFO]msghandler.go:66: Worker ID = 5 is started.
2023/03/09 18:39:49 [INFO]server.go:134: [START] start Zinx server  ZinxServerApp succ, now listenning...

```



#### Zinx-Client

```go
package main

import (
	"fmt"
	"github.com/aceld/zinx/ziface"
	"github.com/aceld/zinx/znet"
	"time"
)

//Client custom business
func pingLoop(conn ziface.IConnection) {
	for {
		err := conn.SendMsg(1, []byte("Ping...Ping...Ping...[FromClient]"))
		if err != nil {
			fmt.Println(err)
			break
		}

		time.Sleep(1 * time.Second)
	}
}

//Executed when a connection is created
func onClientStart(conn ziface.IConnection) {
	fmt.Println("onClientStart is Called ... ")
	go pingLoop(conn)
}

func main() {
	//Create a client client
	client := znet.NewClient("127.0.0.1", 8999)

	//Set the hook function after the link is successfully established
	client.SetOnConnStart(onClientStart)

	//start the client
	client.Start()

	//Prevent the process from exiting, waiting for an interrupt signal
	select {}
}
```

Run Client

```bash
$ go run client.go 
2023/03/09 19:04:54 [INFO]client.go:73: [START] Zinx Client LocalAddr: 127.0.0.1:55294, RemoteAddr: 127.0.0.1:8999
2023/03/09 19:04:54 [INFO]connection.go:354: ZINX CallOnConnStart....
```

Terminal of Zinx Print:
```bash
recv from client : msgId= 1 , data= Ping...Ping...Ping...[FromClient]
recv from client : msgId= 1 , data= Ping...Ping...Ping...[FromClient]
recv from client : msgId= 1 , data= Ping...Ping...Ping...[FromClient]
recv from client : msgId= 1 , data= Ping...Ping...Ping...[FromClient]
recv from client : msgId= 1 , data= Ping...Ping...Ping...[FromClient]
recv from client : msgId= 1 , data= Ping...Ping...Ping...[FromClient]
...
```


### (2) Zinx configuration file
```json
{
  "Name":"zinx v-0.10 demoApp",
  "Host":"127.0.0.1",
  "TCPPort":7777,
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


---

#### Developers

|  **Zinx**   | **Authors**  |
|  ----  | ----  | 
|[zinx](https://github.com/aceld/zinx)|刘丹冰([@aceld](https://github.com/aceld)) 张超([@zhngcho](https://github.com/zhngcho)) 高智辉Roger([@adsian](https://github.com/adsian)) 胡贵建([@huguijian](https://github.com/huguijian)) 张继瑀([@kstwoak](https://github.com/kstwoak)) 夏小力([@xxl6097](https://github.com/xxl6097)) 李志成([@clukboy](https://github.com/clukboy)）|
|[zinx(C++)](https://github.com/marklion/zinx) |刘洋([@marklion](https://github.com/marklion))|
|[zinx(Lua)](https://github.com/huqitt/zinx-lua)|胡琪([@huqitt](https://github.com/huqitt))|
|[ginx(Java)](https://github.com/ModuleCode/ginx)|ModuleCode([@ModuleCode](https://github.com/ModuleCode))|

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

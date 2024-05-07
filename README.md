# <img width="80px" src="https://s2.ax1x.com/2019/10/09/u4yHo9.png" />
English | [简体中文](README-CN.md)

[![License](https://img.shields.io/badge/License-MIT-black.svg)](LICENSE)
[![Discord](https://img.shields.io/badge/zinx-Discord-blue.svg)](https://discord.gg/xQ8Xxfyfcz)
[![Gitter](https://img.shields.io/badge/zinx-Gitter-green.svg)](https://gitter.im/zinx_go/community) 
[![zinx tutorial](https://img.shields.io/badge/ZinxTutorial-YuQue-red.svg)](https://www.yuque.com/aceld/npyr8s/bgftov) 
[![Original Book of Zinx](https://img.shields.io/badge/OriginalBook-YuQue-black.svg)](https://www.yuque.com/aceld)

Zinx is a lightweight concurrent server framework based on Golang.

##  Document 

[ < Zinx Wiki : English > ](https://github.com/aceld/zinx/wiki)

[ < Zinx 文档 : 简体中文> ](https://www.yuque.com/aceld/tsgooa/sbvzgczh3hqz8q3l)


> **Note**: 
> Zinx has been widely used in many enterprises for development purposes, including message forwarding for backend modules, long-linked game servers, and message handling plugins for web frameworks. 
> Zinx is positioned as a framework with concise code that allows developers to quickly understand the internal details of the framework and easily customize it based on their own enterprise scenarios.

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

| platform | Entry | 
| ---- | ---- | 
| <img src="https://user-images.githubusercontent.com/7778936/236784004-b6d99e26-b1ab-4bc3-988e-7a46108b85fe.png" width = "100" height = "100" alt="" align=center />| [Zinx Framework tutorial-Lightweight server based on Golang](https://dev.to/aceld/1building-basic-services-with-zinx-framework-296e)| 
|<img src="https://user-images.githubusercontent.com/7778936/236784168-6528a9b8-d37b-4b02-a37c-b9988d7508d8.jpeg" width = "100" height = "100" alt="" align=center />|[《Golang轻量级并发服务器框架zinx》](https://www.yuque.com/aceld)|


## Online Tutorial Video

| platform | online video | 
| ---- | ---- | 
| <img src="https://s1.ax1x.com/2022/09/22/xFePUK.png" width = "100" height = "100" alt="" align=center />| [![zinx-BiliBili](https://s2.ax1x.com/2019/10/13/uv340S.jpg)](https://www.bilibili.com/video/av71067087)| 
| <img src="https://s1.ax1x.com/2022/09/22/xFeRVx.png" width = "100" height = "100" alt="" align=center />  | [![zinx-BiliBili](https://s2.ax1x.com/2019/10/13/uv340S.jpg)](https://www.douyin.com/video/6983301202939333891) |
| <img src="https://s1.ax1x.com/2022/09/23/xkQcng.png" width = "100" height = "100" alt="" align=center />| [![zinx-youtube](https://s2.ax1x.com/2019/10/14/KSurCR.jpg)](https://www.youtube.com/watch?v=U95iF-HMWsU&list=PL_GrAPKmuajzeNI8HBTi-k5NQO1g0rM-A)| 


## I. One word that has been said before

Why did we create Zinx? Although there are many Golang application frameworks for servers, there are few lightweight enterprise frameworks applied in the gaming or other long-linked fields.

The purpose of designing Zinx is to provide a complete outline of how to write a TCP server based on Golang, so that more Golang enthusiasts can learn and understand this field in a straightforward manner.

The development of the Zinx framework project is synchronized with the creation of learning tutorials, and all the incremental and iterative thinking involved in the development process is incorporated into the tutorials. This approach avoids overwhelming beginners with a complete framework that they may find difficult to grasp all at once.

The tutorials will be iterated version by version, with each version adding small increments of functionality, allowing a beginner to gradually and comprehensively learn about the field of server frameworks.

Of course, we hope that more people will join Zinx and provide us with valuable feedback, enabling Zinx to become a truly enterprise-level server framework. Thank you for your attention!


### Reply from chatGPT(AI)
![what-is-zinx](https://user-images.githubusercontent.com/7778936/209745848-acfc14eb-74cd-4513-b386-8bc6e0bcc09f.png)

![compare-zinx](https://user-images.githubusercontent.com/7778936/209745864-7d8984b0-bd73-4109-b4ec-aec152f8f8e8.png)


### The honor of zinx
#### GVP Most Valuable Open Source Project of the Year at OSCHINA

![GVP-zinx](https://s2.ax1x.com/2019/10/13/uvYVBV.jpg)




## II. Zinx architecture
![Zinx框架](https://user-images.githubusercontent.com/7778936/220058446-0ad45112-2225-4b71-b0d8-69a7f3cee5ca.jpg)

![流程图](https://raw.githubusercontent.com/wenyoufu/testaaaaaa/master/%E6%B5%81%E7%A8%8B%E5%9B%BE-en.jpg)
![zinx-start](https://user-images.githubusercontent.com/7778936/126594039-98dddd10-ec6a-4881-9e06-a09ec34f1af7.gif)



## III. Zinx development API documentation


### (1) QuickStart

[<Zinx's TCP Debugging Tool>](https://github.com/xxl6097/tcptest)

DownLoad zinx Source

```bash
$go get github.com/aceld/zinx
```

> note: Golang Version 1.17+

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
LogIsolationLevel: 0
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
  "Host":"0.0.0.0",
  "TCPPort":9090,
  "MaxConn":3,
  "WorkerPoolSize":10,
  "LogDir": "./mylog",
  "LogFile":"app.log",
  "LogSaveDays":15,
  "LogCons": true,
  "LogIsolationLevel":0
}
```

`Name`:Server Application Name

`Host`:Server IP

`TcpPort`:Server listening port

`MaxConn`:Maximum number of client links allowed

`WorkerPoolSize`:Maximum number of working Goroutines in the work task pool

`LogDir`: Log folder

`LogFile`: Log file name (if not provided, log information is printed to Stderr)

`LogIsolationLevel`: Log Isolation Level -0: Full On 1: Off debug 2: Off debug/info 3: Off debug/info/warn

---


#### Developers

| **Zinx**                                                       | **Authors**  |
|----------------------------------------------------------------| ----  | 
| [zinx](https://github.com/aceld/zinx)                          |刘丹冰([@aceld](https://github.com/aceld)) 张超([@zhngcho](https://github.com/zhngcho)) 高智辉Roger([@adsian](https://github.com/adsian)) 胡贵建([@huguijian](https://github.com/huguijian)) 张继瑀([@kstwoak](https://github.com/kstwoak)) 夏小力([@xxl6097](https://github.com/xxl6097)) 李志成([@clukboy](https://github.com/clukboy)）姚承政([@hcraM41](https://github.com/hcraM41)）李国杰([@LI-GUOJIE](https://github.com/LI-GUOJIE)）余喆宁([@YanHeDoki](https://github.com/YanHeDoki)）|
| [moke-kit(Microservices)](https://github.com/GStones/moke-kit) |GStones([@GStones](https://github.com/GStones))|
| [zinx(C++)](https://github.com/marklion/zinx)                  |刘洋([@marklion](https://github.com/marklion))|
| [zinx(Lua)](https://github.com/huqitt/zinx-lua)                |胡琪([@huqitt](https://github.com/huqitt))|
| [ginx(Java)](https://github.com/ModuleCode/ginx)               |ModuleCode([@ModuleCode](https://github.com/ModuleCode))|

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

### Join the Zinx community 

| platform | Entry | 
| ---- | ---- | 
| <img src="https://user-images.githubusercontent.com/7778936/236775008-6bd488e3-249a-4d43-8885-7e3889e11e2d.png" width = "100" height = "100" alt="" align=center />| https://discord.gg/xQ8Xxfyfcz| 
| <img src="https://user-images.githubusercontent.com/7778936/236775137-5381f8a6-f534-49c4-8628-e52bf245c3bc.jpeg" width = "100" height = "100" alt="" align=center />  | 加微信: `ace_ld`  或扫二维码，备注`zinx`即可。</br><img src="https://user-images.githubusercontent.com/7778936/236781258-2f0371bd-5797-49e8-a74c-680e9f15843d.png" width = "150" height = "150" alt="" align=center /> |
|<img src="https://user-images.githubusercontent.com/7778936/236778547-9cdadfb6-0f62-48ac-851a-b940389038d0.jpeg" width = "100" height = "100" alt="" align=center />|<img src="https://s1.ax1x.com/2020/07/07/UFyUdx.th.jpg" height = "150"  alt="" align=center /> **WeChat Public Account** |
|<img src="https://user-images.githubusercontent.com/7778936/236779000-70f16c8f-0eec-4b5f-9faa-e1d5229a43e0.png" width = "100" height = "100" alt="" align=center />|<img src="https://github.com/aceld/zinx/assets/7778936/461b409f-6337-48a8-826b-a7a746aaee31" width = "150" height = "150" alt="" align=center /> **QQ Group** |

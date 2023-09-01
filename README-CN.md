# <img width="80px" src="https://s2.ax1x.com/2019/10/09/u4yHo9.png" /> 

[English](README.md) | 简体中文

[![License](https://img.shields.io/badge/License-GPL%203.0-black.svg)](LICENSE)
[![Discord](https://img.shields.io/badge/zinx-Discord在线社区-blue.svg)](https://discord.gg/xQ8Xxfyfcz)
[![Gitter](https://img.shields.io/badge/zinx-Gitter在线交流-green.svg)](https://gitter.im/zinx_go/community)
[![zinx tutorial](https://img.shields.io/badge/Zinx教程-YuQue-red.svg)](https://www.yuque.com/aceld/npyr8s/bgftov)
[![Original Book of Zinx](https://img.shields.io/badge/原创书籍-YuQue-black.svg)](https://www.yuque.com/aceld)

Zinx 是一个基于Golang的轻量级并发服务器框架

## 开发者文档

[ < Zinx Wiki : English > ](https://github.com/aceld/zinx/wiki)

[ < Zinx 文档 : 简体中文> ](https://www.yuque.com/aceld/tsgooa/sbvzgczh3hqz8q3l)

> **说明**:目前zinx已经在很多企业进行开发使用，具体使用领域包括:后端模块的消息中转、长链接游戏服务器、Web框架中的消息处理插件等。zinx的定位是代码简洁，让更多的开发者迅速的了解框架的内脏细节并且可以快速基于zinx DIY(二次开发)一款适合自己企业场景的模块。

---
## zinx源码地址
### Github
Git: https://github.com/aceld/zinx

### 码云(Gitee)
Git: https://gitee.com/Aceld/zinx

### 官网
http://zinx.me

---

## 在线开发教程

### 文字教程

| platform | Entry | 
| ---- | ---- | 
| <img src="https://user-images.githubusercontent.com/7778936/236784004-b6d99e26-b1ab-4bc3-988e-7a46108b85fe.png" width = "100" height = "100" alt="" align=center />| [Zinx Framework tutorial-Lightweight server based on Golang](https://dev.to/aceld/1building-basic-services-with-zinx-framework-296e)| 
|<img src="https://user-images.githubusercontent.com/7778936/236784168-6528a9b8-d37b-4b02-a37c-b9988d7508d8.jpeg" width = "100" height = "100" alt="" align=center />|[《Golang轻量级并发服务器框架zinx》](https://www.yuque.com/aceld)|


### 视频教程

| platform | online video | 
| ---- | ---- | 
| <img src="https://s1.ax1x.com/2022/09/22/xFePUK.png" width = "100" height = "100" alt="" align=center />| [![zinx-BiliBili](https://s2.ax1x.com/2019/10/13/uv340S.jpg)](https://www.bilibili.com/video/av71067087)| 
| <img src="https://s1.ax1x.com/2022/09/22/xFesxJ.png" width = "100" height = "80" alt="" align=center />  | [![zinx-BiliBili](https://s2.ax1x.com/2019/10/13/uv340S.jpg)](https://www.douyin.com/video/6983301202939333891) |
| <img src="https://s1.ax1x.com/2022/09/23/xkQcng.png" width = "100" height = "100" alt="" align=center />| [![zinx-youtube](https://s2.ax1x.com/2019/10/14/KSurCR.jpg)](https://www.youtube.com/watch?v=U95iF-HMWsU&list=PL_GrAPKmuajzeNI8HBTi-k5NQO1g0rM-A)| 

    
## 一、写在前面

我们为什么要做Zinx，Golang目前在服务器的应用框架很多，但是应用在游戏领域或者其他长链接的领域的轻量级企业框架甚少。

设计Zinx的目的是我们可以通过Zinx框架来了解基于Golang编写一个TCP服务器的整体轮廓，让更多的Golang爱好者能深入浅出的去学习和认识这个领域。

Zinx框架的项目制作采用编码和学习教程同步进行，将开发的全部递进和迭代思维带入教程中，而不是一下子给大家一个非常完整的框架去学习，让很多人一头雾水，不知道该如何学起。

教程会一个版本一个版本迭代，每个版本的添加功能都是微小的，让一个服务框架小白，循序渐进的曲线方式了解服务器框架的领域。

当然，最后希望Zinx会有更多的人加入，给我们提出宝贵的意见，让Zinx成为真正的解决企业的服务器框架！在此感谢您的关注！

### 来自chatGPT(AI)的回复
![什么是zinx](https://user-images.githubusercontent.com/7778936/209745655-7463be0d-1450-4a70-b201-6d9279935aff.jpg)
![zinx和其他库对比](https://user-images.githubusercontent.com/7778936/209745668-e6938534-113d-4465-a949-58328c4dca5c.jpg)

### zinx荣誉
#### 开源中国GVP年度最有价值开源项目
![GVP-zinx](https://s2.ax1x.com/2019/10/13/uvYVBV.jpg)


#### Stargazers over time

[![Stargazers over time](https://api.star-history.com/svg?repos=aceld/zinx&type=Date)](#zinx)



## 二、初探Zinx架构

![1-Zinx框架.png](https://camo.githubusercontent.com/903d1431358fa6f4634ebaae3b49a28d97e23d77/68747470733a2f2f75706c6f61642d696d616765732e6a69616e7368752e696f2f75706c6f61645f696d616765732f31313039333230352d633735666636383232333362323533362e706e673f696d6167654d6f6772322f6175746f2d6f7269656e742f7374726970253743696d61676556696577322f322f772f31323430)
![流程图](https://github.com/wenyoufu/testaaaaaa/blob/abc8a50078a86aed37e8af6082d1d867bc165c32/%E6%9C%AA%E5%91%BD%E5%90%8D%E6%B5%81%E7%A8%8B%E5%9B%BE%20(1).jpg?raw=true)
![zinx-start](https://user-images.githubusercontent.com/7778936/126594039-98dddd10-ec6a-4881-9e06-a09ec34f1af7.gif)



## 三、Zinx开发接口文档


### （1）快速开始

[<Zinx的Tcp调试工具>](https://github.com/xxl6097/tcptest)

**版本**
Golang 1.17+

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

// PingRouter MsgId=1的路由
type PingRouter struct {
	znet.BaseRouter
}

//Ping Handle MsgId=1的路由处理方法
func (r *PingRouter) Handle(request ziface.IRequest) {
	//读取客户端的数据
	fmt.Println("recv from client : msgId=", request.GetMsgID(), ", data=", string(request.GetData()))
}

func main() {
	//1 创建一个server服务
	s := znet.NewServer()

	//2 配置路由
	s.AddRouter(1, &PingRouter{})

	//3 启动服务
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

//客户端自定义业务
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

//创建连接的时候执行
func onClientStart(conn ziface.IConnection) {
	fmt.Println("onClientStart is Called ... ")
	go pingLoop(conn)
}

func main() {
	//创建Client客户端
	client := znet.NewClient("127.0.0.1", 8999)

	//设置链接建立成功后的钩子函数
	client.SetOnConnStart(onClientStart)

	//启动客户端
	client.Start()

	//防止进程退出，等待中断信号
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


### （2）Zinx配置文件
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

`Name`:服务器应用名称

`Host`:服务器IP

`TcpPort`:服务器监听端口

`MaxConn`:允许的客户端链接最大数量

`WorkerPoolSize`:工作任务池最大工作Goroutine数量

`LogDir`: 日志文件夹

`LogFile`: 日志文件名称(如果不提供，则日志信息打印到Stderr)

`LogIsolationLevel`: 日志隔离级别 0：全开, 1：关debug, 2：关debug/info, 3：关debug/info/warn 

---

#### 开发者
|  **Zinx**   | **开发者**  |
|  ----  | ----  | 
|[zinx](https://github.com/aceld/zinx)|刘丹冰([@aceld](https://github.com/aceld)) 张超([@zhngcho](https://github.com/zhngcho)) 高智辉Roger([@adsian](https://github.com/adsian)) 胡贵建([@huguijian](https://github.com/huguijian)) 张继瑀([@kstwoak](https://github.com/kstwoak)) 夏小力([@xxl6097](https://github.com/xxl6097)) 李志成([@clukboy](https://github.com/clukboy)）姚承政([@hcraM41](https://github.com/hcraM41)）李国杰([@LI-GUOJIE](https://github.com/LI-GUOJIE)）|
|[zinx(Lua)](https://github.com/huqitt/zinx-lua)|胡琪([@huqitt](https://github.com/huqitt))|
|[ginx(Java)](https://github.com/ModuleCode/ginx)|ModuleCode([@ModuleCode](https://github.com/ModuleCode))|

---

感谢所有为zinx贡献的开发者


<a href="https://github.com/aceld/zinx/graphs/contributors">
  <img src="https://contrib.rocks/image?repo=aceld/zinx" />
</a>    


---
### 关于作者：

作者：`Aceld(刘丹冰)`

`mail`:
[danbing.at@gmail.com](mailto:danbing.at@gmail.com)

`github`:
[https://github.com/aceld](https://github.com/aceld)

`原创书籍`:
[https://www.yuque.com/aceld](https://www.yuque.com/aceld)


### 加入Zinx技术社区

| platform | Entry | 
| ---- | ---- | 
| <img src="https://user-images.githubusercontent.com/7778936/236775008-6bd488e3-249a-4d43-8885-7e3889e11e2d.png" width = "100" height = "100" alt="" align=center />| https://discord.gg/xQ8Xxfyfcz| 
| <img src="https://user-images.githubusercontent.com/7778936/236775137-5381f8a6-f534-49c4-8628-e52bf245c3bc.jpeg" width = "100" height = "100" alt="" align=center />  | 加微信: `ace_ld`  或扫二维码，备注`zinx`即可。</br><img src="https://user-images.githubusercontent.com/7778936/236781258-2f0371bd-5797-49e8-a74c-680e9f15843d.png" width = "150" height = "150" alt="" align=center /> |
|<img src="https://user-images.githubusercontent.com/7778936/236778547-9cdadfb6-0f62-48ac-851a-b940389038d0.jpeg" width = "100" height = "100" alt="" align=center />|<img src="https://s1.ax1x.com/2020/07/07/UFyUdx.th.jpg" height = "150"  alt="" align=center /> **WeChat Public Account** |
|<img src="https://user-images.githubusercontent.com/7778936/236779000-70f16c8f-0eec-4b5f-9faa-e1d5229a43e0.png" width = "100" height = "100" alt="" align=center />|<img src="https://s1.ax1x.com/2020/07/07/UF6Y9S.th.png" width = "150" height = "150" alt="" align=center /> **QQ Group** |

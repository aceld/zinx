# <img width="80px" src="https://s2.ax1x.com/2019/10/09/u4yHo9.png" /> 
[![License](https://img.shields.io/badge/License-GPL%203.0-blue.svg)](LICENSE) [![Gitter](https://img.shields.io/badge/在线交流-Gitter-green.svg)](https://gitter.im/zinx_go/community) [![zinx详细教程](https://img.shields.io/badge/zinx详细教程-简书-red.svg)](https://www.kancloud.cn/aceld/zinx) [![zinx原创书籍下载](https://img.shields.io/badge/原创书籍下载-KanCloud-black.svg)](https://www.kancloud.cn/@aceld)

Zinx 是一个基于Golang的轻量级并发服务器框架

官网:http://zinx.me

> **说明**:目前zinx已经在很多企业进行开发使用，具体使用领域包括:后端模块的消息中转、长链接游戏服务器、Web框架中的消息处理插件等。zinx的定位是代码简洁，让更多的开发者迅速的了解框架的内脏细节并且可以快速基于zinx DIY(二次开发)一款适合自己企业场景的模块。

---
## zinx源码地址
### Github
Git: https://github.com/aceld/zinx

### 码云(Gitee)
Git: https://gitee.com/Aceld/zinx

---

## 在线开发教程
### [B站]

[![zinx-视频教程B站](https://s2.ax1x.com/2019/10/13/uv340S.jpg)](https://www.bilibili.com/video/av71067087)


### [YouTube]

[![zinx-youtube](https://s2.ax1x.com/2019/10/14/KSurCR.jpg)](https://www.youtube.com/watch?v=U95iF-HMWsU&list=PL_GrAPKmuajzeNI8HBTi-k5NQO1g0rM-A)

> 注意：修正Github打不开图片方法 windows: https://zhuanlan.zhihu.com/p/107196957 , mac/Linux: 修改/etc/hosts文件即可 

## Zinx详细教程及文档

### PC端文档

[看云-《Zinx框架教程-基于Golang的轻量级并发服务器》](https://www.kancloud.cn/aceld/zinx)

[简书-《Zinx框架教程-基于Golang的轻量级并发服务器》](https://www.jianshu.com/p/23d07c0a28e5)

### 移动端文档
![gongzhonghao](https://s1.ax1x.com/2020/07/07/UFyUdx.th.jpg)

    
## 一、写在前面

我们为什么要做Zinx，Golang目前在服务器的应用框架很多，但是应用在游戏领域或者其他长链接的领域的轻量级企业框架甚少。

设计Zinx的目的是我们可以通过Zinx框架来了解基于Golang编写一个TCP服务器的整体轮廓，让更多的Golang爱好者能深入浅出的去学习和认识这个领域。

Zinx框架的项目制作采用编码和学习教程同步进行，将开发的全部递进和迭代思维带入教程中，而不是一下子给大家一个非常完整的框架去学习，让很多人一头雾水，不知道该如何学起。

教程会一个版本一个版本迭代，每个版本的添加功能都是微小的，让一个服务框架小白，循序渐进的曲线方式了解服务器框架的领域。

当然，最后希望Zinx会有更多的人加入，给我们提出宝贵的意见，让Zinx成为真正的解决企业的服务器框架！在此感谢您的关注！

**zinx荣誉**
#### 开源中国GVP年度最有价值开源项目

![GVP-zinx](https://s2.ax1x.com/2019/10/13/uvYVBV.jpg)


#### Stargazers over time

[![Stargazers over time](https://starchart.cc/aceld/zinx.svg)](https://starchart.cc/aceld/zinx)


## 二、快速启动

```bash
# 克隆项目
$ git clone https://github.com/aceld/zinx.git

# 进入服务端样例目录
$ cd ./zinx/examples/zinx_server

# 服务端编译
$ make build

# 服务端容器化
$ make image

# 服务端启动
$ make run 

# 进入客户端样例目录
$ cd ../zinx_client

# 启动客户端进行测试
$ go run main.go 

```


## 三、初探Zinx架构

![1-Zinx框架.png](https://camo.githubusercontent.com/903d1431358fa6f4634ebaae3b49a28d97e23d77/68747470733a2f2f75706c6f61642d696d616765732e6a69616e7368752e696f2f75706c6f61645f696d616765732f31313039333230352d633735666636383232333362323533362e706e673f696d6167654d6f6772322f6175746f2d6f7269656e742f7374726970253743696d61676556696577322f322f772f31323430)

![zinx-start](https://user-images.githubusercontent.com/7778936/126594039-98dddd10-ec6a-4881-9e06-a09ec34f1af7.gif)



## 三、Zinx开发API文档

### 快速开始

#### 示例演示
1. 编译demo示例，会在`zinx/example/zinx_server`下得到`server`, 在`zinx/example/zinx_client`下得到`client`.
```bash
$ cd zinx/
$ make
```
2. 启动Demo server, 该终端不要关闭
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

3. 再打开新终端，启动client Demo测试通信
```bash
$ cd example/zinx_client
$ ./client
==> Test Router:[Ping] Recv Msg: ID= 2 , len= 21 , data= DoConnection BEGIN... ==> Test Router:[Ping] Recv Msg: ID= 0 , len= 18 , data= ping...ping...ping 
==> Test Router:[Ping] Recv Msg: ID= 0 , len= 18 , data= ping...ping...ping
==> Test Router:[Ping] Recv Msg: ID= 0 , len= 18 , data= ping...ping...ping
...
t

```

#### server
基于Zinx框架开发的服务器应用，主函数步骤比较精简，最多只需要3步即可。
1. 创建server句柄
2. 配置自定义路由及业务
3. 启动服务

```go
func main() {
	//1 创建一个server句柄
	s := znet.NewServer()

	//2 配置路由
	s.AddRouter(0, &PingRouter{})

	//3 开启服务
	s.Serve()
}
```

其中自定义路由及业务配置方式如下：
```go
import (
	"fmt"
	"github.com/aceld/zinx/ziface"
	"github.com/aceld/zinx/znet"
)

//ping test 自定义路由
type PingRouter struct {
	znet.BaseRouter
}

//Ping Handle
func (this *PingRouter) Handle(request ziface.IRequest) {
	//先读取客户端的数据
	fmt.Println("recv from client : msgId=", request.GetMsgID(), ", data=", string(request.GetData()))

    //再回写ping...ping...ping
	err := request.GetConnection().SendBuffMsg(0, []byte("ping...ping...ping"))
	if err != nil {
		fmt.Println(err)
	}
}
```

#### client
Zinx的消息处理采用，`[MsgLength]|[MsgID]|[Data]`的封包格式
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
	模拟客户端
 */
func main() {

	fmt.Println("Client Test ... start")
	//3秒之后发起测试请求，给服务端开启服务的机会
	time.Sleep(3 * time.Second)

	conn,err := net.Dial("tcp", "127.0.0.1:7777")
	if err != nil {
		fmt.Println("client start err, exit!")
		return
	}

	for n := 3; n >= 0; n-- {
		//发封包message消息
		dp := znet.NewDataPack()
		msg, _ := dp.Pack(znet.NewMsgPackage(0,[]byte("Zinx Client Test Message")))
		_, err := conn.Write(msg)
		if err !=nil {
			fmt.Println("write error err ", err)
			return
		}

		//先读出流中的head部分
		headData := make([]byte, dp.GetHeadLen())
		_, err = io.ReadFull(conn, headData) //ReadFull 会把msg填充满为止
		if err != nil {
			fmt.Println("read head error")
			break
		}
		//将headData字节流 拆包到msg中
		msgHead, err := dp.Unpack(headData)
		if err != nil {
			fmt.Println("server unpack err:", err)
			return
		}

		if msgHead.GetDataLen() > 0 {
			//msg 是有data数据的，需要再次读取data数据
			msg := msgHead.(*znet.Message)
			msg.Data = make([]byte, msg.GetDataLen())

			//根据dataLen从io中读取字节流
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

### Zinx配置文件
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

`Name`:服务器应用名称

`Host`:服务器IP

`TcpPort`:服务器监听端口

`MaxConn`:允许的客户端链接最大数量

`WorkerPoolSize`:工作任务池最大工作Goroutine数量

`LogDir`: 日志文件夹

`LogFile`: 日志文件名称(如果不提供，则日志信息打印到Stderr)


### I.服务器模块Server
```go
  func NewServer () ziface.IServer 
```
创建一个Zinx服务器句柄，该句柄作为当前服务器应用程序的主枢纽，包括如下功能：

#### 1)开启服务
```go
  func (s *Server) Start()
```
#### 2)停止服务
```go
  func (s *Server) Stop()
```
#### 3)运行服务
```go
  func (s *Server) Serve()
```
#### 4)注册路由
```go
func (s *Server) AddRouter (msgId uint32, router ziface.IRouter) 
```
#### 5)注册链接创建Hook函数
```go
func (s *Server) SetOnConnStart(hookFunc func (ziface.IConnection))
```
#### 6)注册链接销毁Hook函数
```go
func (s *Server) SetOnConnStop(hookFunc func (ziface.IConnection))
```
### II.路由模块

```go
//实现router时，先嵌入这个基类，然后根据需要对这个基类的方法进行重写
type BaseRouter struct {}

//这里之所以BaseRouter的方法都为空，
// 是因为有的Router不希望有PreHandle或PostHandle
// 所以Router全部继承BaseRouter的好处是，不需要实现PreHandle和PostHandle也可以实例化
func (br *BaseRouter)PreHandle(req ziface.IRequest){}
func (br *BaseRouter)Handle(req ziface.IRequest){}
func (br *BaseRouter)PostHandle(req ziface.IRequest){}
```


### III.链接模块
#### 1)获取原始的socket TCPConn
```go
  func (c *Connection) GetTCPConnection() *net.TCPConn 
```
#### 2)获取链接ID
```go
  func (c *Connection) GetConnID() uint32 
```
#### 3)获取远程客户端地址信息
```go
  func (c *Connection) RemoteAddr() net.Addr 
```
#### 4)发送消息
```go
  func (c *Connection) SendMsg(msgId uint32, data []byte) error 
  func (c *Connection) SendBuffMsg(msgId uint32, data []byte) error
```
#### 5)链接属性
```go
//设置链接属性
func (c *Connection) SetProperty(key string, value interface{})

//获取链接属性
func (c *Connection) GetProperty(key string) (interface{}, error)

//移除链接属性
func (c *Connection) RemoveProperty(key string) 
```

---

#### 开发者
-   刘丹冰([@aceld](https://github.com/aceld))
-   张超([@zhngcho](https://github.com/zhngcho))
-   高智辉Roger([@adsian](https://github.com/adsian))
-   胡贵建([@huguijian](https://github.com/huguijian))
-   张继瑀([@kstwoak](https://github.com/huguijian))


---
[zinx(C++版本)](https://github.com/marklion/zinx)
#### 开发者
-  刘洋([@marklion](https://github.com/marklion))


---
[zinx(Lua版本)](https://github.com/huqitt/zinx-lua)
#### 开发者
-  胡琪([@huqitt](https://github.com/huqitt))

---
[zinx(websocket版本)](https://github.com/aceld/zinx/tree/wsserver)
#### 开发者
-  胡贵建([@huguijian](https://github.com/huguijian))

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
[https://www.kancloud.cn/@aceld](https://www.kancloud.cn/@aceld)

### Zinx技术讨论社区

### **微信公众号**
![gongzhonghao](https://s1.ax1x.com/2020/07/07/UFyUdx.th.jpg)


### **QQ资源分享**

![qqqun](https://s1.ax1x.com/2020/07/07/UF6Y9S.th.png)

### **微信群**
<img src="https://s1.ax1x.com/2020/07/07/UF6rNV.png" width = "200" height = "230" alt="weixin" />
欢迎大家加入，一起学习分享

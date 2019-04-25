# Zinx
[![License](https://img.shields.io/badge/License-GPL%203.0-blue.svg)](LICENSE) [![Gitter](https://img.shields.io/badge/在线交流-Gitter-green.svg)](https://gitter.im/zinx_go/community) [![zinx详细教程](https://img.shields.io/badge/zinx详细教程-简书-red.svg)](https://www.jianshu.com/p/23d07c0a28e5) [![zinx原创书籍下载](https://img.shields.io/badge/原创书籍下载-Gitbook-black.svg)](https://legacy.gitbook.com/book/aceld/zinx/details)

Zinx 是一个基于Golang的轻量级并发服务器框架

### 开发者

-   刘丹冰([@aceld](https://github.com/aceld))
-   张超([@zhngcho](https://github.com/zhngcho))

---
   [zinx(C++版本)](https://github.com/marklion/Game_zinx)
### 开发者
-  刘洋([@marklion](https://github.com/marklion))
---
    
## 一、写在前面

我们为什么要做Zinx，Golang目前在服务器的应用框架很多，但是应用在游戏领域或者其他长链接的领域的轻量级企业框架甚少。

设计Zinx的目的是我们可以通过Zinx框架来了解基于Golang编写一个TCP服务器的整体轮廓，让更多的Golang爱好者能深入浅出的去学习和认识这个领域。

Zinx框架的项目制作采用编码和学习教程同步进行，将开发的全部递进和迭代思维带入教程中，而不是一下子给大家一个非常完整的框架去学习，让很多人一头雾水，不知道该如何学起。

教程会一个版本一个版本迭代，每个版本的添加功能都是微小的，让一个服务框架小白，循序渐进的曲线方式了解服务器框架的领域。

当然，最后希望Zinx会有更多的人加入，给我们提出宝贵的意见，让Zinx成为真正的解决企业的服务器框架！在此感谢您的关注！

## 二、初探Zinx架构

![1-Zinx框架.png](https://upload-images.jianshu.io/upload_images/11093205-c75ff682233b2536.png?imageMogr2/auto-orient/strip%7CimageView2/2/w/1240)

![zinx-start.gif](https://upload-images.jianshu.io/upload_images/11093205-490b49098ac63af6.gif?imageMogr2/auto-orient/strip)

## 三、Zinx详细教程(代码教程同步更新)
[《Zinx框架教程-基于Golang的轻量级并发服务器》](https://www.jianshu.com/p/23d07c0a28e5)

## 四、Zinx开发API文档

### 快速开始

#### server
基于Zinx框架开发的服务器应用，主函数步骤比较精简，最多主需要3步即可。
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
	"zinx/ziface"
	"zinx/znet"
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
	"zinx/znet"
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
### 关于作者：

作者：`Aceld(刘丹冰)`
简书号：`IT无崖子`

`mail`:
[danbing.at@gmail.com](mailto:danbing.at@gmail.com)
`github`:
[https://github.com/aceld](https://github.com/aceld)
`原创书籍gitbook`:
[http://legacy.gitbook.com/@aceld](http://legacy.gitbook.com/@aceld)

### Zinx技术讨论社区

QQ技术讨论群:

![gopool5.jpeg](https://upload-images.jianshu.io/upload_images/11093205-6cdfd381e8ffa127.jpeg?imageMogr2/auto-orient/strip%7CimageView2/2/w/1240)

欢迎大家加入，获取更多相关学习资料

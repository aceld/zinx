/**
* @Author: Aceld
* @Date: 2023/03/02
* @Mail: danbing.at@gmail.com
*    zinx client demo
 */
package main

import (
	"fmt"
	"github.com/aceld/zinx/examples/zinx_client/c_router"
	"github.com/aceld/zinx/ziface"
	"github.com/aceld/zinx/zlog"
	"github.com/aceld/zinx/znet"
	"github.com/aceld/zinx/zpack"
	"io"
	"net"
	"os"
	"os/signal"
	"time"
)

// 客户端自定义业务
func business(conn ziface.IConnection) {

	for {
		err := conn.SendMsg(100, []byte("Ping...[FromClient]"))
		if err != nil {
			fmt.Println(err)
			zlog.Error(err)
			break
		}

		time.Sleep(1 * time.Second)
	}
}

// 创建连接的时候执行
func DoClientConnectedBegin(conn ziface.IConnection) {
	zlog.Debug("DoConnecionBegin is Called ... ")

	//设置两个链接属性，在连接创建之后
	conn.SetProperty("Name", "刘丹冰")
	conn.SetProperty("Home", "https://yuque.com/aceld")

	go business(conn)
}

// 连接断开的时候执行
func DoClientConnectedLost(conn ziface.IConnection) {
	//在连接销毁之前，查询conn的Name，Home属性
	if name, err := conn.GetProperty("Name"); err == nil {
		zlog.Debug("Conn Property Name = ", name)
	}

	if home, err := conn.GetProperty("Home"); err == nil {
		zlog.Debug("Conn Property Home = ", home)
	}

	zlog.Debug("DoClientConnectedLost is Called ... ")
}

func main() {
	//创建一个Client句柄，使用Zinx的API
	client := znet.NewClient("127.0.0.1", 8999)

	//添加首次建立链接时的业务
	client.SetOnConnStart(DoClientConnectedBegin)
	client.SetOnConnStop(DoClientConnectedLost)

	//注册收到服务器消息业务路由
	client.AddRouter(2, &c_router.PingRouter{})
	client.AddRouter(3, &c_router.HelloRouter{})

	//启动客户端client
	client.Start()

	// close
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	sig := <-c
	fmt.Println("===exit===", sig)
	// 清理客户端
	client.Stop()
	time.Sleep(time.Second * 2)
}

/*
模拟客户端, 不使用client模块方式
*/
func main_old() {
	conn, err := net.Dial("tcp", "127.0.0.1:8999")
	if err != nil {
		fmt.Println("client start err, exit!", err)
		return
	}

	for {
		//发封包message消息
		dp := zpack.NewDataPack()
		msg, _ := dp.Pack(zpack.NewMsgPackage(0, []byte("Zinx client Demo Test MsgID=0, [Ping]")))
		_, err := conn.Write(msg)
		if err != nil {
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
			msg := msgHead.(*zpack.Message)
			msg.Data = make([]byte, msg.GetDataLen())

			//根据dataLen从io中读取字节流
			_, err := io.ReadFull(conn, msg.Data)
			if err != nil {
				fmt.Println("server unpack data err:", err)
				return
			}

			fmt.Println("==> Test Router:[Ping] Recv Msg: ID=", msg.ID, ", len=", msg.DataLen, ", data=", string(msg.Data))
		}

		time.Sleep(1 * time.Second)
	}
}

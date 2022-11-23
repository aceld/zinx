package Notify

import (
	"fmt"
	"github.com/aceld/zinx/utils"
	"github.com/aceld/zinx/ziface"
	"github.com/aceld/zinx/znet"
	"github.com/aceld/zinx/zpack"
	"net"
	"strconv"
	"testing"
	"time"
)

var nt = NewNotify()

type router struct {
	znet.BaseRouter
}

func (r *router) Handle(req ziface.IRequest) {
	id, _ := strconv.Atoi(string(req.GetData()))
	nt.SetNotifyID(uint64(id), req.GetConnection())
}

func Server() {
	s := znet.NewUserConfServer(&utils.Config{
		Host:             "127.0.0.1",
		TcpPort:          9991,
		Name:             "NtTest",
		TcpVersion:       "tcp",
		Version:          "1",
		MaxConn:          10000,
		MaxPacketSize:    4096,
		WorkerPoolSize:   10,
		MaxWorkerTaskLen: 10,
		MaxMsgChanLen:    10,
	})

	s.AddRouter(1, &router{})
	s.Serve()
}

func Clinet() {
	//conf.ConfigInit()
	//1创建直接链接
	for i := 0; i < 9000; i++ {
		go func(i int) {
			conn, err := net.Dial("tcp", "127.0.0.1:9991")
			if err != nil {
				fmt.Println("net dial err:", err)
				return
			}
			defer conn.Close()
			//链接调用write方法写入数据
			id := strconv.Itoa(i)
			dp := zpack.NewDataPack()
			msg, err := dp.Pack(zpack.NewMsgPackage(1, []byte(id)))
			if err != nil {
				return
			}
			_, err = conn.Write(msg)

			if err != nil {
				return
			}
			select {}
		}(i)
	}
}

func init() {
	go Server()
	go Clinet()
	go ClinetJoin()
}

func ClinetJoin() {
	t := time.NewTicker(50 * time.Millisecond)
	i := 10000
	for {
		select {
		case <-t.C:
			go func(i int) {
				conn, err := net.Dial("tcp", "127.0.0.1:9991")
				if err != nil {
					fmt.Println("net dial err:", err)
					return
				}
				defer conn.Close()
				//链接调用write方法写入数据
				id := strconv.Itoa(i)
				dp := zpack.NewDataPack()
				msg, err := dp.Pack(zpack.NewMsgPackage(1, []byte(id)))
				if err != nil {
					return
				}
				_, err = conn.Write(msg)

				if err != nil {
					return
				}
				select {}
			}(i)
			i++
		}
	}

}

func TestAA(t *testing.T) {
	time.AfterFunc(5*time.Second, func() {
		fmt.Println(len(nt.cimap))
	})
	time.Sleep(6 * time.Second)
}

func BenchmarkNotify(b *testing.B) {
	time.Sleep(5 * time.Second)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		nt.NotifyAll(1, []byte("雪下的是盐"))
	}
	fmt.Println(len(nt.cimap))
}

package znotify

import (
	"fmt"
	"net"
	"strconv"
	"testing"
	"time"

	"github.com/aceld/zinx/zconf"
	"github.com/aceld/zinx/ziface"
	"github.com/aceld/zinx/znet"
	"github.com/aceld/zinx/zpack"
)

var nt = NewZNotify()

type router struct {
	znet.BaseRouter
}

func (r *router) Handle(req ziface.IRequest) {
	id, _ := strconv.Atoi(string(req.GetData()))
	nt.SetNotifyID(uint64(id), req.GetConnection())
}

func Server() {
	s := znet.NewUserConfServer(&zconf.Config{
		Host:             "127.0.0.1",
		TCPPort:          9991,
		Name:             "NtTest",
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

func ClientJoin() {
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
	})
	time.Sleep(6 * time.Second)
	nt.ConnNums()
}

func BenchmarkNotify(b *testing.B) {
	fmt.Println("Begin BenchmarkNotify")
	time.Sleep(60 * time.Second)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		nt.NotifyAll(1, []byte("雪下的是盐"))
	}
	nt.ConnNums()
}

func init() {
	go Server()
	go Clinet()
	go ClientJoin()
}

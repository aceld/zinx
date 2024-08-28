package main

import (
	"encoding/hex"
	"github.com/aceld/zinx/zconf"
	"github.com/aceld/zinx/ziface"
	"github.com/aceld/zinx/zinterceptor"
	"github.com/aceld/zinx/zlog"
	"github.com/aceld/zinx/znet"
)

type Test struct {
	znet.BaseRouter
}

func (tlv *Test) Intercept(chain ziface.IChain) ziface.IcResp {
	//1. Get the IMessage of zinx
	iMessage := chain.GetIMessage()
	if iMessage == nil {
		// Go to the next layer in the chain of responsibility
		return chain.ProceedWithIMessage(iMessage, nil)
	}

	iMessage.SetMsgID(1000)

	//6. Pass the decoded data to the next layer.
	// (将解码后的数据进入下一层)
	return chain.ProceedWithIMessage(iMessage, chain.Request())
}
func (this *Test) Handle(request ziface.IRequest) {
	_response := request.GetResponse()
	if _response != nil {
		switch _response.(type) {
		case *znet.Request:
			_data := _response.(*znet.Request)
			zlog.Ins().InfoF("body:%s\n", string(_data.GetData()))
			break
		}
	} else {
		zlog.Ins().InfoF("Message0x15 Handle %s \n", hex.EncodeToString(request.GetMessage().GetData()))
	}
}

func main() {
	conf := &zconf.Config{}
	conf.TCPPort = 8045
	conf.Name = "clink-go-tcp-yehoo-app"
	conf.Host = "0.0.0.0"
	conf.WorkerPoolSize = 10
	conf.MaxConn = 1000000

	s := znet.NewUserConfServer(conf)
	//注册链接hook回调函数
	s.SetOnConnStart(func(connection ziface.IConnection) {
		zlog.Ins().InfoF("DoConnectionBegin is Called ...")
	})

	s.SetOnConnStop(func(connection ziface.IConnection) {
		zlog.Ins().InfoF("Conn is Lost")
	})

	frameDecoder := zinterceptor.NewDelimiterBasedFrameDecoder(10240, false, true, []byte("_$"))

	s.SetFrameDecoder(frameDecoder)
	test := Test{}
	s.AddInterceptor(&test)
	s.AddRouter(1000, &test)

	//开启服务
	s.Serve()
}

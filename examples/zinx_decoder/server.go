package main

import (
	"github.com/aceld/zinx/examples/zinx_decoder/router"
	"github.com/aceld/zinx/zdecoder"
	"github.com/aceld/zinx/ziface"
	"github.com/aceld/zinx/zlog"
	"github.com/aceld/zinx/znet"
)

func DoConnectionBegin(conn ziface.IConnection) {
	zlog.Ins().InfoF("DoConnectionBegin is Called ...")
}

func DoConnectionLost(conn ziface.IConnection) {
	zlog.Ins().InfoF("Conn is Lost")
}

func main() {
	//创建一个server句柄
	s := znet.NewServer()

	//注册链接hook回调函数
	s.SetOnConnStart(DoConnectionBegin)
	s.SetOnConnStop(DoConnectionLost)

	//s.AddRouter(0x00000001, &router.TLVBusinessRouter{}) //TLV协议对应业务功能
	//处理HTLVCRC协议数据
	s.SetDecoder(zdecoder.NewHTLVCRCDecoder())
	s.AddRouter(0x10, &router.HtlvCrcBusinessRouter{}) //TLV协议对应业务功能，因为client.go中模拟数据funcode字段为0x10
	s.AddRouter(0x13, &router.HtlvCrcBusinessRouter{}) //TLV协议对应业务功能，因为client.go中模拟数据funcode字段为0x13

	//开启服务
	s.Serve()
}

package main

import (
	"github.com/aceld/zinx/examples/zinx_decoder/decode"
	"github.com/aceld/zinx/examples/zinx_decoder/router"
	"github.com/aceld/zinx/ziface"
	"github.com/aceld/zinx/zlog"
	"github.com/aceld/zinx/znet"
)

func DoConnectionBegin(conn ziface.IConnection) {
	zlog.Ins().InfoF("DoConnecionBegin is Called ...")
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

	//处理TLV协议数据
	s.AddInterceptor(&decode.TLVDecoder{})               //TVL协议解码器
	s.AddRouter(0x00000001, &router.TLVBusinessRouter{}) //TLV协议对应业务功能

	//处理HTLVCRC协议数据
	s.AddInterceptor(&decode.HtlvCrcDecoder{})         //TVL协议解码器
	s.AddRouter(0x10, &router.HtlvCrcBusinessRouter{}) //TLV协议对应业务功能，因为client.go中模拟数据funcode字段为0x10

	//开启服务
	s.Serve()
}

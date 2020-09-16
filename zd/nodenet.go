package zd

import (
	"errors"
	"fmt"

	"github.com/aceld/zinx/utils"

	"github.com/aceld/zinx/zdnet"
)

/*
本文件主要是node的基本server服务的接口, 不负责处理网络通信实现，
而是处理服务通信基础业务和消息的协议封装和解析
*/

/*
向集群中的Leader Node发送消息
*/
func (node *Node) SendToLeader(cmdId int, port int, data []byte) ([]byte, error) {
	//1. 获取当前集群中的Leader
	leader := node.GetLeader()
	if leader == nil {
		return nil, errors.New(fmt.Sprintf("Leader not Found, Please Try Again, cmdId = %d", cmdId))
	}

	//2. 创建短连接
	conn := zdnet.NewZDConn(leader.Ip, port)
	if conn == nil {
		return nil, errors.New("Can not connected to leader... " + leader.Ip)
	}

	defer conn.Close()

	//3. 发送消息给Leader节点
	if err := node.SendToNode(conn, cmdId, data); err != nil {
		return nil, errors.New(fmt.Sprintf("sending To Node error, cmdId=%d", cmdId))
	} else {
		//4.读取远程Node返回的消息
		msg := node.RecvFromNode(conn)
		if msg != nil && msg.CmdId != utils.ZINX_CMD_ID_NODE_SYNC_ACK {
			return nil, errors.New(fmt.Sprintf("sync/raft deal response error, cmdid = %d", msg.CmdId))
		}
		return msg.Data, nil
	}

	return nil, errors.New("unknown error")
}

/*
	一个node发送消息给另一个node
*/
func (node *Node) SendToNode(conn *zdnet.ZDConn, cmdId int, data []byte) error {
	unit := node.GetZinxUnit()
	unit.Ip = conn.Ip

	//将消息打包成Message
	msg := NewZdMessage(cmdId, unit, data)

	//将msg序列化并发送
	binaryBuf, err := ZdMsgPack(msg)
	if err != nil {
		fmt.Println("SendToNode, message pack error")
		return err
	}

	return conn.Send(binaryBuf)
}

/*
	从一个Node读取消息, 需要先调用SendToNode，然后再处理返回的数据
*/
//TODO 设置超时时间的接口
func (node *Node) RecvFromNode(conn *zdnet.ZDConn) *ZdMessage {
	binaryData := conn.Recv()

	//读到了数据
	if binaryData != nil && len(binaryData) > 0 {
		msg, err := ZdMsgUnpack(binaryData)
		if err != nil {
			fmt.Println("ZdMsgUnpack error")
			return nil
		}

		return msg
	}

	return nil
}

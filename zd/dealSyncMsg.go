package zd

import (
	"encoding/json"
	"fmt"

	"github.com/aceld/zinx/utils"
	"github.com/aceld/zinx/zdnet"
)

/*
	处理API的Message消息
*/
func (node *Node) DealSyncMsg(conn *zdnet.ZDConn) {
	msg := node.RecvFromNode(conn)

	if msg == nil {
		conn.Close()
		return
	}

	switch msg.CmdId {
	case utils.ZINX_CMD_ID_NODE_ADD:
		node.dealNodeAdd(conn, msg)
	case utils.ZINX_CMD_ID_NODE_REMOVE:
		node.dealNodeRemove(conn, msg)
	}

	//再次读取，如果读取不到数据，则会关闭链接，防止链接过多超出限制
	node.DealSyncMsg(conn)
}

//新增node节点
func (node *Node) dealNodeAdd(conn *zdnet.ZDConn, msg *ZdMessage) {
	//得到 需要添加的 node信息
	nodelist := []string{}
	if json.Unmarshal(msg.Data, &nodelist) != nil {
		fmt.Println("dealNodeAdd json unmarshal error")
		return
	}

	if nodelist != nil && len(nodelist) > 0 {
		for _, ip := range nodelist {

			if _, ok := node.Peers[ip]; ok {
				//key存在，不做任何操作
				continue
			}

			//更新当前节点信息
			node.AddPeersUnit(&ZinxUnit{Id: ip, Ip: ip})
		}
	}

	//回执对端消息
	node.SendToNode(conn, utils.ZINX_CMD_ID_NODE_SYNC_ACK, nil)
}

//删除node节点
func (node *Node) dealNodeRemove(conn *zdnet.ZDConn, msg *ZdMessage) {
	//得到需要删除的node信息
	nodelist := []string{}
	if json.Unmarshal(msg.Data, &nodelist) != nil {
		fmt.Println("dealNodeRemove json unmarshal error")
		return
	}

	if nodelist != nil && len(nodelist) > 0 {
		for _, ip := range nodelist {
			for _, v := range node.Peers {
				unit := v.(*ZinxUnit)
				if ip == unit.Ip {
					fmt.Printf("delete peer node: %s, ip:%d\n", unit.Name, unit.Ip)
					node.peersLock.Lock()
					delete(node.Peers, unit.Id)
					node.peersLock.Unlock()
				}
			}
		}
	}

	//回执对端消息
	node.SendToNode(conn, utils.ZINX_CMD_ID_NODE_SYNC_ACK, nil)
}

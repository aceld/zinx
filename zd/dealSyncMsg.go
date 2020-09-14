package zd

import "github.com/aceld/zinx/zdnet"

func (node *Node) DealSyncMsg(conn *zdnet.ZDConn) {
	msg := node.RecvFromNode(conn)

	if msg == nil {
		conn.Close()
	}

	switch msg.CmdId {

	}
}

package znet

import (
	"errors"
	"github.com/aceld/zinx/zlog"
	"sync"

	"github.com/aceld/zinx/ziface"
)

type ConnManager struct {
	connections map[uint64]ziface.IConnection
	connLock    sync.RWMutex
}

func newConnManager() *ConnManager {
	return &ConnManager{
		connections: make(map[uint64]ziface.IConnection),
	}
}

func (connMgr *ConnManager) Add(conn ziface.IConnection) {

	connMgr.connLock.Lock()
	connMgr.connections[conn.GetConnID()] = conn //将conn连接添加到ConnMananger中
	connMgr.connLock.Unlock()

	zlog.Ins().InfoF("connection add to ConnManager successfully: conn num = %d", connMgr.Len())
}

func (connMgr *ConnManager) Remove(conn ziface.IConnection) {

	connMgr.connLock.Lock()
	delete(connMgr.connections, conn.GetConnID()) //删除连接信息
	connMgr.connLock.Unlock()

	zlog.Ins().InfoF("connection Remove ConnID=%d successfully: conn num = %d", conn.GetConnID(), connMgr.Len())
}

func (connMgr *ConnManager) Get(connID uint64) (ziface.IConnection, error) {
	connMgr.connLock.RLock()
	defer connMgr.connLock.RUnlock()

	if conn, ok := connMgr.connections[connID]; ok {
		return conn, nil
	}

	return nil, errors.New("connection not found")
}

func (connMgr *ConnManager) Len() int {

	connMgr.connLock.RLock()
	length := len(connMgr.connections)
	connMgr.connLock.RUnlock()

	return length
}

func (connMgr *ConnManager) ClearConn() {
	connMgr.connLock.Lock()

	// Stop and delete all connection information
	for connID, conn := range connMgr.connections {
		//停止
		conn.Stop()
		delete(connMgr.connections, connID)
	}
	connMgr.connLock.Unlock()

	zlog.Ins().InfoF("Clear All Connections successfully: conn num = %d", connMgr.Len())
}

func (connMgr *ConnManager) GetAllConnID() []uint64 {
	ids := make([]uint64, 0, len(connMgr.connections))

	connMgr.connLock.RLock()
	defer connMgr.connLock.RUnlock()

	for id := range connMgr.connections {
		ids = append(ids, id)
	}

	return ids
}

func (connMgr *ConnManager) Range(cb func(uint64, ziface.IConnection, interface{}) error, args interface{}) (err error) {

	for connID, conn := range connMgr.connections {
		err = cb(connID, conn, args)
	}

	return err
}

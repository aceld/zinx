package znet

import (
	"errors"
	"github.com/aceld/zinx/zlog"
	"sync"

	"github.com/aceld/zinx/ziface"
)

//ConnManager 连接管理模块
type ConnManager struct {
	//主链接结合
	connections map[uint64]ziface.IConnection
	//只读的链接集合
	connectionsReadOnly map[uint64]ziface.IConnection
	connLock            sync.RWMutex
}

//NewConnManager 创建一个链接管理
func NewConnManager() *ConnManager {
	return &ConnManager{
		connections:         make(map[uint64]ziface.IConnection),
		connectionsReadOnly: make(map[uint64]ziface.IConnection),
	}
}

//Add 添加链接
func (connMgr *ConnManager) Add(conn ziface.IConnection) {

	connMgr.connLock.Lock()
	connMgr.connections[conn.GetConnID()] = conn //将conn连接添加到ConnMananger中
	connMgr.connectionsReadOnly[conn.GetConnID()] = conn
	connMgr.connLock.Unlock()

	zlog.Ins().InfoF("connection add to ConnManager successfully: conn num = %d", connMgr.Len())
}

//Remove 删除连接
func (connMgr *ConnManager) Remove(conn ziface.IConnection) {

	connMgr.connLock.Lock()
	delete(connMgr.connections, conn.GetConnID()) //删除连接信息
	delete(connMgr.connectionsReadOnly, conn.GetConnID())
	connMgr.connLock.Unlock()

	zlog.Ins().InfoF("connection Remove ConnID=%d successfully: conn num = %d", conn.GetConnID(), connMgr.Len())
}

//Get 利用ConnID获取链接
func (connMgr *ConnManager) Get(connID uint64) (ziface.IConnection, error) {
	connMgr.connLock.RLock()
	defer connMgr.connLock.RUnlock()

	if conn, ok := connMgr.connections[connID]; ok {
		return conn, nil
	}

	return nil, errors.New("connection not found")
}

//Len 获取当前连接
func (connMgr *ConnManager) Len() int {

	connMgr.connLock.RLock()
	length := len(connMgr.connections)
	connMgr.connLock.RUnlock()

	return length
}

//ClearConn 清除并停止所有连接
func (connMgr *ConnManager) ClearConn() {
	connMgr.connLock.Lock()

	//停止并删除全部的连接信息
	for connID, conn := range connMgr.connections {
		//停止
		conn.Stop()
		delete(connMgr.connections, connID)
		delete(connMgr.connectionsReadOnly, connID)
	}
	connMgr.connLock.Unlock()

	zlog.Ins().InfoF("Clear All Connections successfully: conn num = %d", connMgr.Len())
}

// GetAllConnID 获取所有连接的ID
func (connMgr *ConnManager) GetAllConnID() []uint64 {
	ids := make([]uint64, 0, len(connMgr.connectionsReadOnly))

	for id := range connMgr.connectionsReadOnly {
		ids = append(ids, id)
	}

	return ids
}

// Range 遍历所有连接
func (connMgr *ConnManager) Range(cb func(uint64, ziface.IConnection, interface{}) error, args interface{}) (err error) {

	for connID, conn := range connMgr.connections {
		err = cb(connID, conn, args)
	}

	return err
}

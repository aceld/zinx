package server

import (
	"errors"
	"fmt"
	"sync"
	"wsserver/iserverface"
)
/*
	连接管理模块
*/
type ConnManager struct {
	connections map[uint64]iserverface.IConnection //管理的连接信息
	connLock    sync.RWMutex                  //读写连接的读写锁
}

/*
	创建一个链接管理
*/
func NewConnManager() *ConnManager {
	return &ConnManager{
		connections: make(map[uint64]iserverface.IConnection),
	}
}

//添加链接
func (connMgr *ConnManager) Add(conn iserverface.IConnection) {
	//保护共享资源Map 加写锁
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()

	//将conn连接添加到ConnMananger中
	connMgr.connections[conn.GetConnID()] = conn

	fmt.Println("connection add to ConnManager successfully: conn num = ", connMgr.Len())
}

//删除连接
func (connMgr *ConnManager) Remove(conn iserverface.IConnection) {
	//保护共享资源Map 加写锁
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()

	//删除连接信息
	delete(connMgr.connections, conn.GetConnID())

	fmt.Println("connection Remove ConnID=", conn.GetConnID(), " successfully: conn num = ", connMgr.Len())
}

//利用ConnID获取链接
func (connMgr *ConnManager) Get(connID uint64) (iserverface.IConnection, error) {
	//保护共享资源Map 加读锁
	connMgr.connLock.RLock()
	defer connMgr.connLock.RUnlock()

	if conn, ok := connMgr.connections[connID]; ok {
		return conn, nil
	} else {
		return nil, errors.New("connection not found")
	}
}

//获取当前连接
func (connMgr *ConnManager) Len() int {
	return len(connMgr.connections)
}

//清除并停止所有连接
func (connMgr *ConnManager) ClearConn() {
	//保护共享资源Map 加写锁
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()

	//停止并删除全部的连接信息
	for connID, conn := range connMgr.connections {
		//停止
		conn.Close()
		//删除
		delete(connMgr.connections, connID)
	}

	fmt.Println("Clear All connections successfully: conn num = ", connMgr.Len())
}

func (connMgr *ConnManager) PushAll(msg []byte) {
	for _, conn := range connMgr.connections {
		conn.SendMessage(1,msg)
	}
}

func (connMgr *ConnManager) GetConnByProName(key string,val string) (iserverface.IConnection, error) {
	//保护共享资源Map 加读锁
	connMgr.connLock.RLock()
	defer connMgr.connLock.RUnlock()
	var ConnID uint64
	for connId, conn := range connMgr.connections {

		if proVal,err :=conn.GetProperty("parkid"); err!=nil {
			fmt.Println("获取属性出错",err)
		}else{

			if proVal==val {
				ConnID = connId
				break;
			}
		}
	}
	if conn, ok := connMgr.connections[ConnID]; ok {
		return conn, nil
	} else {
		return nil, errors.New("connection not found")
	}
}

package znotify

import (
	"errors"
	"fmt"
	"sync"

	"github.com/aceld/zinx/ziface"
	"github.com/aceld/zinx/zlog"
)

// ConnIDMap Establish a structure that maps user-defined IDs to connections
// Map will have concurrent access issues, as well as looping through large amounts of data
// Currently, a map structure is used to store the data, but it may not be the best choice
// (建立一个用户自定义ID和连接映射的结构
// map会存在 并发问题，大量数据循环读取问题
// 暂时先用map结构存储，但是应该不是最好的选择，抛砖引玉)
type ConnIDMap map[uint64]ziface.IConnection

type notify struct {
	cimap ConnIDMap
	sync.RWMutex
}

func NewZNotify() ziface.Inotify {
	return &notify{
		cimap: make(map[uint64]ziface.IConnection, 5000),
	}
}

func (n *notify) ConnNums() int {
	return len(n.cimap)
}

func (n *notify) HasIdConn(Id uint64) bool {
	n.RLock()
	defer n.RUnlock()
	_, ok := n.cimap[Id]
	return ok
}

func (n *notify) SetNotifyID(Id uint64, conn ziface.IConnection) {
	n.Lock()
	defer n.Unlock()
	n.cimap[Id] = conn
}

func (n *notify) GetNotifyByID(Id uint64) (ziface.IConnection, error) {
	n.RLock()
	defer n.RUnlock()
	Conn, ok := n.cimap[Id]
	if !ok {
		return nil, errors.New(" Not Find UserId")
	}
	return Conn, nil
}

func (n *notify) DelNotifyByID(Id uint64) {
	n.RLock()
	defer n.RUnlock()
	delete(n.cimap, Id)
}

func (n *notify) NotifyToConnByID(Id uint64, MsgId uint32, data []byte) error {
	Conn, err := n.GetNotifyByID(Id)
	if err != nil {
		return err
	}
	err = Conn.SendMsg(MsgId, data)
	if err != nil {
		fmt.Printf("Notify to %d err:%s \n", Id, err)
		return err
	}
	return nil
}

func (n *notify) NotifyAll(MsgId uint32, data []byte) error {
	n.RLock()
	defer n.RUnlock()
	for Id, v := range n.cimap {
		err := v.SendMsg(MsgId, data)
		if err != nil {
			zlog.Ins().ErrorF("Notify to %d err:%s \n", Id, err)
		}
	}
	return nil
}

func (n *notify) notifyAll(MsgId uint32, data []byte) error {
	n.RLock()
	defer n.RUnlock()
	var err error
	for Id, v := range n.cimap {
		er := v.SendMsg(MsgId, data)
		if er != nil {
			zlog.Ins().ErrorF("Notify to %d err:%s \n", Id, er)
			err = er
		}
	}
	return err
}

// In extreme cases where many people are joining and sending messages at the same time, and the map needs to be released as soon as possible,
// but there are currently many problems and it is not used
// (极端情况 同时加入和发送的人很多需要尽快释放map的情况， 目前问题很多不采用)
func (n *notify) notifyAll2(MsgId uint32, data []byte) error {
	conns := make([]ziface.IConnection, 0, len(n.cimap))
	n.RLock()
	for _, v := range n.cimap {
		conns = append(conns, v)
	}
	n.RUnlock()

	var err error
	for i := 0; i < len(conns); i++ {
		if er := conns[i].SendMsg(MsgId, data); er != nil {
			err = er
		}
	}
	return err
}

func (n *notify) NotifyBuffToConnByID(Id uint64, MsgId uint32, data []byte) error {
	Conn, err := n.GetNotifyByID(Id)
	if err != nil {
		return err
	}
	err = Conn.SendBuffMsg(MsgId, data)
	if err != nil {
		zlog.Ins().ErrorF("Notify to %d err:%s \n", Id, err)
		return err
	}
	return nil
}

func (n *notify) NotifyBuffAll(MsgId uint32, data []byte) error {
	n.RLock()
	defer n.RUnlock()
	for Id, v := range n.cimap {
		err := v.SendBuffMsg(MsgId, data)
		if err != nil {
			zlog.Ins().ErrorF("Notify to %d err:%s \n", Id, err)
		}
	}
	return nil
}

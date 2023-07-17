package znotify

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/aceld/zinx/ziface"
	"github.com/aceld/zinx/zlog"
	"github.com/aceld/zinx/zutils"
)

// ConnIDMap Establish a structure that maps user-defined IDs to connections
// Map will have concurrent access issues, as well as looping through large amounts of data
// Use the map structure of shard and lock storage to minimize lock granularity and lock holding time
// 建立一个用户自定义ID和连接映射的结构
// map会存在并发问题，大量数据循环读取问题
// 使用分片加锁的map结构存储，尽量减少锁的粒度和锁的持有时间

type notify struct {
	connIdMap zutils.ShardLockMaps
}

func NewZNotify() ziface.Inotify {
	return &notify{
		connIdMap: zutils.NewShardLockMaps(),
	}
}

func (n *notify) genConnStrId(connID uint64) string {
	strConnId := strconv.FormatUint(connID, 10)
	return strConnId
}

func (n *notify) ConnNums() int {
	return n.connIdMap.Count()
}

func (n *notify) HasIdConn(Id uint64) bool {
	strId := n.genConnStrId(Id)
	return n.connIdMap.Has(strId)
}

func (n *notify) SetNotifyID(Id uint64, conn ziface.IConnection) {
	strId := n.genConnStrId(Id)
	n.connIdMap.Set(strId, conn)
}

func (n *notify) GetNotifyByID(Id uint64) (ziface.IConnection, error) {

	strId := n.genConnStrId(Id)
	Conn, ok := n.connIdMap.Get(strId)
	if !ok {
		return nil, errors.New(" Not Find UserId")
	}
	return Conn.(ziface.IConnection), nil
}

func (n *notify) DelNotifyByID(Id uint64) {
	strId := n.genConnStrId(Id)
	n.connIdMap.Remove(strId)
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

	n.connIdMap.IterCb(func(key string, v interface{}) {
		conn, _ := v.(ziface.IConnection)
		err := conn.SendMsg(MsgId, data)
		if err != nil {
			zlog.Ins().ErrorF("Notify to %s err:%s \n", key, err)
		}
	})

	return nil
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

	n.connIdMap.IterCb(func(key string, v interface{}) {
		conn, _ := v.(ziface.IConnection)
		err := conn.SendBuffMsg(MsgId, data)
		if err != nil {
			zlog.Ins().ErrorF("Notify to %s err:%s \n", key, err)
		}
	})

	return nil
}

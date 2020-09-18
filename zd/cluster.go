package zd

/*
	处理Node的集群管理相关实现
*/

import (
	"fmt"

	"github.com/aceld/zinx/utils"
)

/* 一个Node 在集群中的基础信息 */
type ZinxUnit struct {
	Name    string `json:"name"`
	Group   string `json:"group"`
	Id      string `json:"id"`
	Ip      string `json:"ip"`
	Status  bool   `json:"status"`
	Role    int32  `json:"role"`
	Version string `json:"version"`
}

func (node *Node) GetZinxUnit() *ZinxUnit {
	return &ZinxUnit{
		Name:    node.Name,
		Group:   node.Group,
		Id:      node.Id,
		Ip:      node.Ip,
		Status:  utils.ZINX_UNIT_STATUS_ALIVE,
		Role:    node.Role,
		Version: utils.ZINX_DISTRIBUTED_VERSION,
	}
}

//获取当前node所在集群的Leader节点信息
func (node *Node) GetLeader() *ZinxUnit {
	node.mutex.RLock()
	leader := node.Leader
	node.mutex.RUnlock()

	return leader
}

//设置Leader节点信息
func (node *Node) SetLeader(unit *ZinxUnit) {
	node.mutex.Lock()
	defer node.mutex.Unlock()

	if node.Leader != nil {
		if node.Leader.Id == unit.Id {
			//当前node就是leader
			return
		} else {
			fmt.Printf("Leader Changed %s --> %s\n", node.Leader.Name, unit.Name)
		}
	} else {
		fmt.Printf("Leader  nil ---> %s\n", unit.Name)
	}

	node.Leader = unit
}

//更新节点的zinxUnit信息
func (node *Node) AddPeersUnit(unit *ZinxUnit) {
	//添加节点
	node.peersLock.Lock()
	node.Peers[unit.Id] = unit
	node.peersLock.Unlock()

	//leader更新判断
	leaderUnit := node.GetLeader()
	if leaderUnit != nil && leaderUnit.Id == unit.Id {
		node.SetLeader(unit)
	}

	//去掉初始化时写入的IP作为Key值的记录
	node.peersLock.Lock()
	if _, ok := node.Peers[unit.Ip]; ok {
		if unit.Id != unit.Ip {
			delete(node.Peers, unit.Ip)
		}
	}
	node.peersLock.Unlock()
}

//查询全部的Peers节点
func (node *Node) GetPeers() []ZinxUnit {
	peers := make([]ZinxUnit, 0)

	//添加当前node自身
	peers = append(peers, *node.GetZinxUnit())

	node.peersLock.RLock()
	for _, unit := range node.Peers {
		peers = append(peers, *unit.(*ZinxUnit))
	}
	node.peersLock.RUnlock()

	return peers
}

func (node *Node) ElectionLeader() {
	if len(node.Peers) > 1 {
		//TODO 集群在两个以上
	} else {
		//集群目前仅有1个节点，就是自己,那么Leader也是自己
		node.SetLeader(node.GetZinxUnit())
	}
}

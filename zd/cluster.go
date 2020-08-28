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
	//	node.mutex.Lock()
	//defer node.mutex.Unlock()

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

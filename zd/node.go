/*
	zinx 分布式 node节点信息
*/
package zd

import (
	"fmt"
	"os"
	"sync"

	"github.com/aceld/zinx/zdnet"

	"github.com/aceld/zinx/utils"
	uuid "github.com/aceld/zinx/utils/go.uuid"
)

/* node 节点以及属性 */
type Node struct {
	//当前Node的集群通信连接
	Conn *zdnet.ZDConn

	/* 基本信息 */
	Group string //集群名称
	Id    string //节点ID
	Name  string //节点主机名称
	Ip    string //主机节点ip

	/* 配置文件信息 */

	/* 集群与选举信息 */
	Leader    *ZinxUnit              //Leader节点信息(zinx集群中的一个节点单元)
	Peers     map[string]interface{} //集群所有节点信息
	peersLock sync.RWMutex           //防止竞争Peers的读写锁
	Role      int32                  //集群角色

	/* 锁 */
	mutex sync.RWMutex //保护Node对象竞争访问的锁
}

func NewNode() *Node {
	//主机名称
	hostname, err := os.Hostname()
	if err != nil {
		fmt.Println("get hostname error,", err)
		return nil
	}

	node := &Node{
		Id:    NodeId(),
		Ip:    "127.0.0.1",
		Name:  hostname,
		Group: "default.group.zinx",
		Role:  utils.ZINX_ROLE_SERVER,
		Peers: make(map[string]interface{}),
	}

	//TODO 从配置文件中读取
	node.Ip = utils.GetInternalIP(utils.INTERNET_DEVICE_NAME)
	if node.Ip == "" {
		fmt.Println("error get Ip")
		return nil
	}


	return node
}

//获取NodeID
func NodeId() string {
	u1 := uuid.Must(uuid.NewV4()).String()
	fmt.Println("NodeId=", u1)
	// NewV4用来返回一个随机生成的UUID，以及生成过程中遇到的错误
	// Must将返回结果为(UUID, error)类型的值进行包裹，如果error非nil，程序会panic。在生产环境下，不建议使用。
	// String会将UUID类型转化成string类型

	/*
		u2, err := uuid.FromString("6ba7b810-9dad-11d1-80b4-00c04fd430c8")
		if err != nil {
			fmt.Println("failed to parse uuid: ", err)
			return ""
		}
	*/

	return u1
}

//启动服务
func (node *Node) Start() {
	//解析命令行参数
	if ParseCommand() == true {
		//有单独指令执行，处理指令，然后退出，不走业务
		os.Exit(0)
	}

	//读取配置文件

	//打印node信息
	fmt.Println("------------------------")
	fmt.Println("NodeID		:", node.Id)
	fmt.Println("NodeIP		:", node.Ip)
	fmt.Println("NodeGroup	:", node.Group)
	fmt.Println("NodeName	:", node.Name)
	fmt.Println("Role		:", utils.NodeRoleName(node.Role))
	fmt.Println("------------------------")

	//启动消息数据同步的TCP Server
	go zdnet.NewZdTcpServer(utils.ZINX_SYNC_PORT, node.DealSyncMsg).Start()

	//启动API服务
	go ApiRun(node)

	//初始化选举Leader
	node.ElectionLeader()
}

package zd

import (
	"flag"

	"github.com/aceld/zinx/utils"
)

//初始化命令行参数
func init() {
	var name string
	flag.StringVar(&name, "name", "defaultName", "设置当前node节点名称")

	var group string
	flag.StringVar(&group, "group", "default.group.zinx", "设置当前node节点所在的组")

	var role int
	flag.IntVar(&role, "role", utils.ZINX_ROLE_SERVER, "设置当前role角色: 0-server, 1-client")
}

//通过命令行配置加载参数
func (node *Node) LoadWithCommand() {

}

//通过配置文件加载参数
func (node *Node) LoadWithConfig() {

}

package zd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/aceld/zinx/utils"
)

type NodeApiHandler struct {
	param interface{}
}

func (api *NodeApiHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	retStr := utils.ZINX_API_RET_SUCC

	if r.Method == "POST" {
		//API: -addnode  添加节点
		data, _ := ioutil.ReadAll(r.Body)

		//得到 需要添加的 node信息
		nodelist := []string{}
		if json.Unmarshal(data, &nodelist) != nil {
			fmt.Println("addnode format error")
			retStr = utils.ZINX_API_RET_FAIL
		}

		//依次添加node到集群中
		node := api.param.(*Node)

		fmt.Printf("get data = %+v\n", data)

		_, err := node.SendToLeader(utils.ZINX_CMD_ID_NODE_ADD, utils.ZINX_SYNC_PORT, data)
		if err != nil {
			fmt.Println(err)
			retStr = utils.ZINX_API_RET_FAIL
		}

	} else if r.Method == "GET" {
		//API: -nodes  获取节点集合

	} else if r.Method == "DELETE" {
		//API: -delnode 删除节点
	}

	w.Write([]byte(retStr))
}

func ApiRun(node *Node) {
	//绑定API指令
	http.Handle("/node", &NodeApiHandler{param: node})
	http.ListenAndServe(fmt.Sprintf(":%d", utils.ZINX_API_PORT), nil)
}

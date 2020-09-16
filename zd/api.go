package zd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/aceld/zinx/utils"
)

type ApiResponse struct {
	RetCode int         `json:"retcode"`
	RetStr  string      `json:"retstr"`
	Data    interface{} `json:"data"`
}

/*
	处理 /addnode 添加node指令 API业务
*/
type ApiAddNode struct {
	param interface{}
}

func (api *ApiAddNode) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	retCode := utils.ZINX_API_RETCODE_OK

	if r.Method == "POST" {
		//API: -addnode  添加节点
		data, _ := ioutil.ReadAll(r.Body)

		//依次添加node到集群中
		node := api.param.(*Node)

		_, err := node.SendToLeader(utils.ZINX_CMD_ID_NODE_ADD, utils.ZINX_SYNC_PORT, data)
		if err != nil {
			fmt.Println(err)
			retCode = utils.ZINX_API_RETCODE_FAIL
		}
	} else {
		retCode = utils.ZINX_API_RETCODE_FAIL
	}

	//发送返回数据json
	resp := &ApiResponse{
		RetCode: retCode,
		RetStr:  retCode2Str(retCode),
		Data:    nil,
	}

	jsonData, _ := json.Marshal(resp)

	//回写给客户端(Command命令行端)
	w.Write(jsonData)
}

/*
	处理 /nodes 获取全部node信息 API业务
*/
type ApiGetNodes struct {
	param interface{}
}

func (api *ApiGetNodes) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	retStr := utils.ZINX_API_RET_SUCC

	if r.Method == "GET" {
		//API: -addnode  添加节点
		data, _ := ioutil.ReadAll(r.Body)

		//依次添加node到集群中
		node := api.param.(*Node)

		_, err := node.SendToLeader(utils.ZINX_CMD_ID_NODE_ADD, utils.ZINX_SYNC_PORT, data)
		if err != nil {
			fmt.Println(err)
			retStr = utils.ZINX_API_RET_FAIL
		}
	} else {
		retStr = utils.ZINX_API_RET_FAIL
	}

	w.Write([]byte(retStr))
}

func ApiRun(node *Node) {
	//绑定API指令
	http.Handle("/addnode", &ApiAddNode{param: node})
	http.Handle("/nodes", &ApiGetNodes{param: node})
	http.ListenAndServe(fmt.Sprintf(":%d", utils.ZINX_API_PORT), nil)
}

func retCode2Str(retCode int) string {
	if retCode == utils.ZINX_API_RETCODE_OK {
		return utils.ZINX_API_RET_SUCC
	} else if retCode == utils.ZINX_API_RETCODE_FAIL {
		return utils.ZINX_API_RET_FAIL
	}

	return utils.ZINX_API_RET_FAIL
}

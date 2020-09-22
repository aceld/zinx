package zd

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/aceld/zinx/utils"
)

type ZinxCmd struct {
	nodes    bool     //查看所有节点
	addNodes []string //要添加的节点
	delNodes []string //要删除的节点 (,分隔多个)

	services    bool     //查看所有服务
	addService  string   //要添加的服务
	delServices []string //要删除的服务
}

//全局唯一命令行结果
var zinxCmd ZinxCmd

type sliceValue []string

func newSliceValue(vals []string, p *[]string) *sliceValue {
	*p = vals
	return (*sliceValue)(p)
}
func (s *sliceValue) Set(val string) error {
	*s = sliceValue(strings.Split(val, ","))
	return nil
}
func (s *sliceValue) Get() interface{} { return []string(*s) }

func (s *sliceValue) String() string { return strings.Join([]string(*s), ",") }

//初始化命令行参数
func init() {
	/* node 配置*/
	flag.Var(newSliceValue([]string{}, &zinxCmd.addNodes), "addnode", "添加node节点(ip1,ip2,ip3...)")
	flag.BoolVar(&zinxCmd.nodes, "nodes", false, "查看全部的node信息")
	flag.Var(newSliceValue([]string{}, &zinxCmd.delNodes), "delnode", "删除节点(ip:port, ip:port, ...)")

	/* service 配置*/
	flag.StringVar(&zinxCmd.addService, "addservice", "unknow", "添加一个服务")
	flag.BoolVar(&zinxCmd.services, "services", false, "查看全部的service信息")
	flag.Var(newSliceValue([]string{}, &zinxCmd.delServices), "delservice", "删除服务")
}

/*合法命令行 返回true，否则返回false*/
func ParseCommand() bool {
	//解析命令行
	flag.Parse()

	//	fmt.Println("zinxCmd = ", zinxCmd)

	//处理命令行结果
	if zinxCmd.nodes == true {
		//显示全部节点信息
		fmt.Println("===> CMD nodes")
		CommandNodes()
		return true
	}

	if len(zinxCmd.addNodes) != 0 {
		//添加新节点
		fmt.Println("===> CMD addnode, node = ", zinxCmd.addNodes)
		CommandAddNode()
		return true
	}

	if len(zinxCmd.delNodes) != 0 {
		//删除节点
		fmt.Println("====> CMD delnode, node = ", zinxCmd.delNodes)
		CommandDelNode()
		return true
	}

	if zinxCmd.services == true {
		//显示全部节点信息
		fmt.Println("===> CMD services")

		return true
	}

	if zinxCmd.addService != "unknow" {
		//添加一个新节点
		fmt.Println("===> CMD addservice, service = ", zinxCmd.addService)

		return true
	}

	if len(zinxCmd.delServices) != 0 {
		//删除节点
		fmt.Println("====> CMD delservices, node = ", zinxCmd.delServices)

		return true
	}

	return false
}

/*
	查询集群中当前在线的node节点信息
	usage: [zinx -nodes]
*/
func CommandNodes() {
	//发送http请求,给API层
	url := fmt.Sprintf("http://127.0.0.1:%d/nodes", utils.ZINX_API_PORT)
	response, err := http.Get(url)
	if err != nil {
		fmt.Printf("http GET %s error: %s", url, err)
		return
	}

	defer response.Body.Close()

	//处理API回执内容
	if response.StatusCode != 200 {
		fmt.Println("command ERROR, StatusCode = ", response.StatusCode)
		return
	}

	retData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("read Body error, ", err)
		return
	}

	//解析json文件
	resp := ApiResponse{}

	fmt.Println("json: ")
	fmt.Printf("%s\n", retData)

	if err := json.Unmarshal(retData, &resp); err != nil {
		fmt.Println("unmarshal API json retData error ,", err)
		return
	}

	if resp.RetCode == utils.ZINX_API_RETCODE_OK {
		if resp.Data != nil {
			nodelist := resp.Data.([]interface{})
			for _, v := range nodelist {
				node := v.(map[string]interface{})
				fmt.Println("------- Node -------")
				fmt.Printf("Id:\t%s\n", node["id"])
				fmt.Printf("Name:\t%s\n", node["name"])
				fmt.Printf("Group:\t%s\n", node["group"])
				fmt.Printf("Ip:\t%s\n", node["ip"])
				fmt.Printf("Role:\t%v\n", node["role"])
				fmt.Printf("Status:\t%v\n", node["status"])
				fmt.Printf("Version:\t%s\n", node["version"])
			}
		}
		fmt.Println("[GetNodes Success!]")
	} else {
		fmt.Println("[GetNodes Fail!]")
	}
}

/*
	添加一个node节点到集群中
	usage: [zinx -addnode IP1,IP2,IP3...]
*/
func CommandAddNode() {
	data, err := json.Marshal(zinxCmd.addNodes)
	if err != nil {
		fmt.Println("json marshal addNodes err:", err)
		return
	}

	//fmt.Printf("请求数据: %s\n", data)

	//发送http请求,给API层
	url := fmt.Sprintf("http://127.0.0.1:%d/addnode", utils.ZINX_API_PORT)
	response, err := http.Post(url,
		//"application/x-www-form-urlencoded",
		"",
		bytes.NewReader(data))

	if err != nil {
		fmt.Printf("http POST %s error: %s", url, err)
		return
	}

	defer response.Body.Close()

	//处理API回执内容
	if response.StatusCode != 200 {
		fmt.Println("command ERROR, StatusCode = ", response.StatusCode)
		return
	}

	retData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("read Body error, ", err)
		return
	}

	fmt.Println("json: ")
	fmt.Printf("%s\n", retData)

	//解析json文件
	resp := ApiResponse{}

	if err := json.Unmarshal(retData, &resp); err != nil {
		fmt.Println("unmarshal API json retData error ,", err)
		return
	}

	if resp.RetCode == utils.ZINX_API_RETCODE_OK {
		fmt.Println("[AddNode Success!]")
	} else {
		fmt.Println("[AddNode Fail!]")
	}
}

/*
	删除node节点从集群中
	usage: [zinx -delnode IP1,IP2,IP3...]
*/
func CommandDelNode() {
	data, err := json.Marshal(zinxCmd.delNodes)
	if err != nil {
		fmt.Println("json marshal delNodes err:", err)
		return
	}

	//发送http请求,给API层
	url := fmt.Sprintf("http://127.0.0.1:%d/delnode", utils.ZINX_API_PORT)
	response, err := http.Post(url,
		//"application/x-www-form-urlencoded",
		"",
		bytes.NewReader(data))

	if err != nil {
		fmt.Printf("http POST %s error: %s", url, err)
		return
	}

	defer response.Body.Close()

	//处理API回执内容
	if response.StatusCode != 200 {
		fmt.Println("command ERROR, StatusCode = ", response.StatusCode)
		return
	}

	retData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("read Body error, ", err)
		return
	}

	fmt.Println("json: ")
	fmt.Printf("%s\n", retData)

	//解析json文件
	resp := ApiResponse{}

	if err := json.Unmarshal(retData, &resp); err != nil {
		fmt.Println("unmarshal API json retData error ,", err)
		return
	}

	if resp.RetCode == utils.ZINX_API_RETCODE_OK {
		fmt.Println("[DelNode Success!]")
	} else {
		fmt.Println("[DelNode Fail!]")
	}
}

// Package utils 提供zinx相关工具类函数
// 包括:
//
//	全局配置
//	配置文件加载
//
// 当前文件描述:
// @Title  globalobj.go
// @Description  相关配置文件定义及加载方式
// @Author  Aceld - Thu Mar 11 10:32:29 CST 2019
package zconf

import (
	"encoding/json"
	"fmt"
	"github.com/aceld/zinx/utils/commandline/args"
	"github.com/aceld/zinx/utils/commandline/uflag"
	"github.com/aceld/zinx/zlog"
	"io/ioutil"
	"os"
	"reflect"
	"testing"
	"time"
)

/*
存储一切有关Zinx框架的全局参数，供其他模块使用
一些参数也可以通过 用户根据 zinx.json来配置
*/
type Config struct {
	/*
		Server
	*/
	Host    string //当前服务器主机IP
	TCPPort int    //当前服务器主机监听端口号
	Name    string //当前服务器名称

	/*
		Zinx
	*/
	Version          string //当前Zinx版本号
	MaxPacketSize    uint32 //读写数据包的最大值
	MaxConn          int    //当前服务器主机允许的最大链接个数
	WorkerPoolSize   uint32 //业务工作Worker池的数量
	MaxWorkerTaskLen uint32 //业务工作Worker对应负责的任务队列最大任务存储数量
	MaxMsgChanLen    uint32 //SendBuffMsg发送消息的缓冲最大长度
	IOReadBuffSize   uint32 //每次IO最大的读取长度

	/*
		logger
	*/
	LogDir        string //日志所在文件夹 默认"./log"
	LogFile       string //日志文件名称   默认""  --如果没有设置日志文件，打印信息将打印至stderr
	LogDebugClose bool   //是否关闭Debug日志级别调试信息 默认false  -- 默认打开debug信息

	/*
		Keepalive
	*/
	HeartbeatMax int //最长心跳检测间隔时间(单位：秒),超过改时间间隔，则认为超时，从配置文件读取

	/*
		TLS
	*/
	CertFile       string // 证书文件名称 默认""
	PrivateKeyFile string // 私钥文件名称 默认"" --如果没有设置证书和私钥文件，则不启用TLS加密
}

/*
定义一个全局的对象
*/
var GlobalObject *Config

// PathExists 判断一个文件是否存在
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// Reload 读取用户的配置文件
func (g *Config) Reload() {
	confFilePath := args.Args.ConfigFile
	if confFileExists, _ := PathExists(confFilePath); confFileExists != true {
		zlog.Ins().ErrorF("Config File %s is not exist!!", confFilePath)
		return
	}

	data, err := ioutil.ReadFile(confFilePath)
	if err != nil {
		panic(err)
	}
	//将json数据解析到struct中
	err = json.Unmarshal(data, g)
	if err != nil {
		panic(err)
	}

	//Logger 设置
	if g.LogFile != "" {
		zlog.SetLogFile(g.LogDir, g.LogFile)
	}
	if g.LogDebugClose == true {
		zlog.CloseDebug()
	}
}

// 提示详细
func (g *Config) Show() {
	//提示当前配置信息
	objVal := reflect.ValueOf(g).Elem()
	objType := reflect.TypeOf(*g)

	fmt.Println("===== Zinx Global Config =====")
	for i := 0; i < objVal.NumField(); i++ {
		field := objVal.Field(i)
		typeField := objType.Field(i)

		fmt.Printf("%s: %v\n", typeField.Name, field.Interface())
	}
	fmt.Println("==============================")
}

func (g *Config) HeartbeatMaxDuration() time.Duration {
	return time.Duration(g.HeartbeatMax) * time.Second
}

/*
提供init方法，默认加载
*/
func init() {
	pwd, err := os.Getwd()
	if err != nil {
		pwd = "."
	}

	// 初始化配置模块flag
	args.InitConfigFlag(pwd+"/conf/zinx.json", "配置文件，如果没有设置，则默认为<exeDir>/conf/zinx.json")
	// 初始化日志模块flag TODO

	// 解析
	testing.Init() //防止 go test 出现"flag provided but not defined: -test.paniconexit0"等错误
	uflag.Parse()

	// 解析之后的操作
	args.FlagHandle()

	//初始化GlobalObject变量，设置一些默认值
	GlobalObject = &Config{
		Name:             "ZinxServerApp",
		Version:          "V1.0",
		TCPPort:          8999,
		Host:             "0.0.0.0",
		MaxConn:          12000,
		MaxPacketSize:    4096,
		WorkerPoolSize:   10,
		MaxWorkerTaskLen: 1024,
		MaxMsgChanLen:    1024,
		LogDir:           pwd + "/log",
		LogFile:          "",
		LogDebugClose:    true,
		HeartbeatMax:     10, //默认心跳检测最长间隔为10秒
		IOReadBuffSize:   1024,
		CertFile:         "",
		PrivateKeyFile:   "",
	}
	//NOTE: 从配置文件中加载一些用户配置的参数
	GlobalObject.Reload()
}

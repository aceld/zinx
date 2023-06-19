// @Title  globalobj.go
// @Description  相关配置文件定义及加载方式
// defines a configuration structure named "Config" along with its methods.
// The package is named "zconf", and the file is named "globalobj.go".
// @Author  Aceld - Thu Mar 11 10:32:29 CST 2019
package zconf

import (
	"fmt"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/aceld/zinx/zlog"
	"github.com/aceld/zinx/zutils/commandline/args"
	"github.com/aceld/zinx/zutils/commandline/uflag"
	"github.com/ilyakaznacheev/cleanenv"
)

const (
	ServerModeTcp       = "tcp"
	ServerModeWebsocket = "websocket"
)

const (
	WorkerModeHash = "Hash" //By default, the round-robin average allocation rule is used.(默认使用取余的方式)
	WorkerModeBind = "Bind" //Bind a worker to each connection.(为每个连接分配一个worker)
)

/*
	   Store all global parameters related to the Zinx framework for use by other modules.
	   Some parameters can also be configured by the user based on the zinx.json file.
		(存储一切有关Zinx框架的全局参数，供其他模块使用
		一些参数也可以通过 用户根据 zinx.json来配置)
*/
type Config struct {
	/*
		Server
	*/
	Host    string `json:"Host" env:"PORT" env-default:"0.0.0.0"`       //The IP address of the current server. (当前服务器主机IP)
	TCPPort int    `json:"TCPPort" env:"TCP_PORT" env-default:"8999"`   //The port number on which the server listens for TCP connections.(当前服务器主机监听端口号)
	WsPort  int    `json:"WsPort" env:"WS_PORT" env-default:"9000"`     //The port number on which the server listens for WebSocket connections.(当前服务器主机websocket监听端口)
	Name    string `json:"Name" env:"NAME" env-default:"ZinxServerApp"` //The name of the current server.(当前服务器名称)

	/*
		Zinx
	*/
	Version          string `json:"Version" env:"VERSION" env-default:"V1.0.0"`                    //The version of the Zinx framework.(当前Zinx版本号)
	MaxPacketSize    uint32 `json:"MaxPacketSize" env:"MAX_PACKET_SIZE" env-default:"4096"`        //The maximum size of the packets that can be sent or received.(读写数据包的最大值)
	MaxConn          int    `json:"MaxConn" env:"MAX_CONN" env-default:"12000"`                    //The maximum number of connections that the server can handle.(当前服务器主机允许的最大链接个数)
	WorkerPoolSize   uint32 `json:"WorkerPoolSize" env:"WORKER_POOL_SIZE" env-default:"10"`        //The number of worker pools in the business logic.(业务工作Worker池的数量)
	MaxWorkerTaskLen uint32 `json:"MaxWorkerTaskLen" env:"MAX_WORKER_TASK_LEN" env-default:"1024"` //The maximum number of tasks that a worker pool can handle.(业务工作Worker对应负责的任务队列最大任务存储数量)
	WorkerMode       string `json:"WorkerMode" env:"WORKER_MODE" env-default:"Hash" `              //The way to assign workers to connections.(为链接分配worker的方式)
	MaxMsgChanLen    uint32 `json:"MaxMsgChanLen" env:"MAX_MSG_CHAN_LEN" env-default:"1024"`       //The maximum length of the send buffer message queue.(SendBuffMsg发送消息的缓冲最大长度)
	IOReadBuffSize   uint32 `json:"IOReadBuffSize" env:"IO_READ_BUFF_SIZE" env-default:"1024"`     //The maximum size of the read buffer for each IO operation.(每次IO最大的读取长度)

	//The server mode, which can be "tcp" or "websocket". If it is empty, both modes are enabled.
	//"tcp":tcp监听, "websocket":websocket 监听 为空时同时开启
	Mode string `json:"Mode" env:"MODE" env-default:"tcp"`

	// A boolean value that indicates whether the new or old version of the router is used. The default value is false.
	//路由模式 false为旧版本路由，true为启用新版本的路由 默认使用旧版本
	RouterSlicesMode bool `json:"RouterSlicesMode" env:"ROUTER_SLICES_MODE" env-default:"false"`

	/*
		logger
	*/
	LogDir string `json:"LogDir" env:"LOG_DIR" env-default:"./log"` //The directory where log files are stored. The default value is "./log".(日志所在文件夹 默认"./log")

	// The name of the log file. If it is empty, the log information will be printed to stderr.
	// (日志文件名称   默认""  --如果没有设置日志文件，打印信息将打印至stderr)
	LogFile string `json:"LogFile" env:"LOG_FILE" env-default:""`

	LogSaveDays int   `json:"LogSaveDays" env:"LOG_SAVE_DAYS" env-default:"7" ` // 日志最大保留天数
	LogFileSize int64 `json:"LogFileSize" env:"LOG_FILE_SIZE" `                 // 日志单个日志最大容量 默认 64MB,单位：字节，记得一定要换算成MB（1024 * 1024）
	LogCons     bool  `json:"LogCons" env:"LOG_CONS" env-default:"false"`       // 日志标准输出  默认 false

	// The level of log isolation. The values can be 0 (all open), 1 (debug off), 2 (debug/info off), 3 (debug/info/warn off), and so on.
	//日志隔离级别  -- 0：全开 1：关debug 2：关debug/info 3：关debug/info/warn ...
	LogIsolationLevel int `json:"LogIsolationLevel" env:"LOG_ISOLATION_LEVEL" env-default:"0"`

	/*
		Keepalive
	*/
	// The maximum interval for heartbeat detection in seconds.
	// 最长心跳检测间隔时间(单位：秒),超过改时间间隔，则认为超时，从配置文件读取, 默认10秒
	HeartbeatMax int `json:"HeartbeatMax" env:"HEARTBEAT_MAX" env-default:"10"`

	/*
		TLS
	*/
	CertFile       string `json:"CertFile" env:"CERT_FILE" env-default:""`              // The name of the certificate file. If it is empty, TLS encryption is not enabled.(证书文件名称 默认"")
	PrivateKeyFile string `json:"PrivateKeyFile" env:"PRIVATE_KEY_FILE" env-default:""` // The name of the private key file. If it is empty, TLS encryption is not enabled.(私钥文件名称 默认"" --如果没有设置证书和私钥文件，则不启用TLS加密)
}

/*
Define a global object.(定义一个全局的对象)
*/
var GlobalObject *Config

// PathExists Check if a file exists.(判断一个文件是否存在)
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
// This method is used to reload the configuration file.
// It reads the configuration file specified in the command-line arguments,
// and updates the fields of the "Config" structure accordingly.
// If the configuration file does not exist, it prints an error message to the log and returns.
func (g *Config) Reload() error {
	confFilePath := args.Args.ConfigFile

	if confFileExists, _ := PathExists(confFilePath); !confFileExists {
		zlog.Ins().ErrorF("Config File %s is not exist!!", confFilePath)
	} else {
		if err := cleanenv.ReadConfig(confFilePath, g); err != nil {
			return fmt.Errorf("read config error: %w", err)
		}
	}
	if err := cleanenv.ReadEnv(g); err != nil {
		return fmt.Errorf("read evn config error: %w", err)
	}
	g.InitLogConfig()
	return nil
}

// Show Zinx Config Info
func (g *Config) Show() {
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

func (g *Config) InitLogConfig() {
	if g.LogFile != "" {
		zlog.SetLogFile(g.LogDir, g.LogFile)
		zlog.SetCons(g.LogCons)
	}
	if g.LogSaveDays > 0 {
		zlog.SetMaxAge(g.LogSaveDays)
	}
	if g.LogFileSize > 0 {
		zlog.SetMaxSize(g.LogFileSize)
	}
	if g.LogIsolationLevel > zlog.LogDebug {
		zlog.SetLogLevel(g.LogIsolationLevel)
	}
}

/*
init, set default value
*/
func init() {
	pwd, err := os.Getwd()
	if err != nil {
		pwd = "."
	}

	args.InitConfigFlag(pwd+"/conf/zinx.json", "The configuration file defaults to <exeDir>/conf/zinx.json if it is not set.")

	// Note: Prevent errors like "flag provided but not defined: -test.paniconexit0" from occurring in go test.
	// (防止 go test 出现"flag provided but not defined: -test.paniconexit0"等错误)
	testing.Init()
	uflag.Parse()

	// after parsing
	args.FlagHandle()

	// Initialize the GlobalObject variable and set some default values.
	// (初始化GlobalObject变量，设置一些默认值)
	GlobalObject = &Config{}

	// Note: Load some user-configured parameters from the configuration file.
	// (从配置文件中加载一些用户配置的参数)
	if err := GlobalObject.Reload(); err != nil {
		panic(err)
	}
}

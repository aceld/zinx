package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/lestrrat/go-file-rotatelogs"
	"github.com/pkg/errors"
	"github.com/rifflock/lfshook"
	log "github.com/sirupsen/logrus"
	"os"
	"path"
	"time"
	"wsserver/common"
	"wsserver/configs"
	"wsserver/iserverface"
	"wsserver/server"
	"wsserver/zlog"
)

var (
	configFile string

)
func ConfigLocalFilesystemLogger(logPath string, logFileName string, maxAge time.Duration, rotationTime time.Duration) {
	baseLogPaht := path.Join(logPath, logFileName)
	writer, err := rotatelogs.New(
		baseLogPaht+".%Y%m%d%H%M",
		rotatelogs.WithLinkName(baseLogPaht),      // 生成软链，指向最新日志文件
		rotatelogs.WithMaxAge(maxAge),             // 文件最大保存时间
		rotatelogs.WithRotationTime(rotationTime), // 日志切割时间间隔
	)

	src, err := os.OpenFile(os.DevNull, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		fmt.Println("err", err)
	}
	writer1 := bufio.NewWriter(src)
	log.SetOutput(writer1)

	if err != nil {
		log.Errorf("config local file system logger error. %+v", errors.WithStack(err))
	}
	lfHook := lfshook.NewHook(lfshook.WriterMap{
		log.DebugLevel: writer,
		log.InfoLevel:  writer,
		log.WarnLevel:  writer,
		log.ErrorLevel: writer,
		log.FatalLevel: writer,
		log.PanicLevel: writer,
	}, &log.JSONFormatter{})
	log.AddHook(lfHook)
}

func initCmd() {
	flag.StringVar(&configFile, "config", "./config.json", "where load config json")
	flag.Parse()
}

// 	WebSocket服务端
//ping test 自定义路由
type PingRouter struct {
	server.BaseRouter
}

//Ping Handle
func (this *PingRouter) Handle(request iserverface.IRequest) {

	err := request.GetConnection().SendMessage(request.GetMsgType(), []byte("ping...ping...ping"))
	if err != nil {
		zlog.Error(err)
	}
}



func main() {
	initCmd()
	ConfigLocalFilesystemLogger("./", "face-service.log", time.Duration(86400)*time.Second, time.Duration(604800)*time.Second)
	var err error = nil
	bindAddress := ""
	if err = configs.LoadConfig(configFile); err != nil {
		fmt.Println("Load config json error:",err)
	}
	common.InitRedis()
	server.GWServer = server.NewServer()

	//配置路由
	server.GWServer.AddRouter("ping", &PingRouter{})


	bindAddress = fmt.Sprintf("%s:%d", configs.GConf.Ip, configs.GConf.Port)
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.GET("/", server.WsHandler)
	r.Run(bindAddress)
}

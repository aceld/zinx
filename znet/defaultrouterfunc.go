package znet

import (
	"fmt"
	"github.com/aceld/zinx/ziface"
	"github.com/aceld/zinx/zlog"
	"path"
	"runtime"
	"strings"
	"time"
)

//用来存放一些RouterSlicesMode下的路由可用的默认中间件

// RouterRecovery 如果使用NewDefaultRouterSlicesServer方法初始化的获得的server将自带这个函数
// 作用是接收业务执行上产生的panic并且尝试记录现场信息
func RouterRecovery(request ziface.IRequest) {
	defer func() {
		if err := recover(); err != nil {
			funcname, filename, lineNo := getInfo(3)
			panicInfo := fmt.Sprintf("MsgId:%d  funcname:%s filename:%s LineNo:%d", request.GetMsgID(), funcname, filename, lineNo)

			//记录错误
			zlog.Ins().ErrorF("Handler panic: info:%s: err:%v  ", panicInfo, err)

			//fmt.Printf("Handler panic: info: %s  err: %v ", panicInfo, err)

			//应该回传一个错误的
			//request.GetConnection().SendMsg()
		}

	}()
	request.RouterSlicesNext()
}

// RouterTime 简单累计所有路由组的耗时，不启用
func RouterTime(request ziface.IRequest) {
	now := time.Now()
	request.RouterSlicesNext()
	duration := time.Since(now)
	fmt.Println(duration.String())
}

func getInfo(ship int) (funcname, filename string, lineNo int) {
	pc, file, lineNo, ok := runtime.Caller(ship)
	if !ok {
		zlog.Ins().ErrorF("runtime.caller() err")
		return
	}
	funcname = runtime.FuncForPC(pc).Name()
	filename = path.Base(file)
	funcname = strings.Split(funcname, ".")[1]
	return funcname, filename, lineNo

}

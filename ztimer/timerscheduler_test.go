package ztimer

/**
* @Author: Aceld(刘丹冰)
* @Date: 2019/5/9 10:14
* @Mail: danbing.at@gmail.com
*
*  时间轮定时器调度器单元测试
 */

import (
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/aceld/zinx/zlog"
)

//触发函数
func foo(args ...interface{}) {
	fmt.Printf("I am No. %d function, delay %d ms\n", args[0].(int), args[1].(int))
}

//手动创建调度器运转时间轮
func TestNewTimerScheduler(t *testing.T) {
	timerScheduler := NewTimerScheduler()
	timerScheduler.Start()

	//在scheduler中添加timer
	for i := 1; i < 2000; i++ {
		f := NewDelayFunc(foo, []interface{}{i, i * 3})
		tID, err := timerScheduler.CreateTimerAfter(f, time.Duration(3*i)*time.Millisecond)
		if err != nil {
			zlog.Error("create timer error", tID, err)
			break
		}
	}

	//执行调度器触发函数
	go func() {
		delayFuncChan := timerScheduler.GetTriggerChan()
		for df := range delayFuncChan {
			df.Call()
		}
	}()

	//阻塞等待
	select {}
}

//采用自动调度器运转时间轮
func TestNewAutoExecTimerScheduler(t *testing.T) {
	autoTS := NewAutoExecTimerScheduler()

	//给调度器添加Timer
	for i := 0; i < 2000; i++ {
		f := NewDelayFunc(foo, []interface{}{i, i * 3})
		tID, err := autoTS.CreateTimerAfter(f, time.Duration(3*i)*time.Millisecond)
		if err != nil {
			zlog.Error("create timer error", tID, err)
			break
		}
	}

	//阻塞等待
	select {}
}

//测试取消一个定时器
func TestCancelTimerScheduler(t *testing.T) {
	Scheduler := NewAutoExecTimerScheduler()
	f1 := NewDelayFunc(foo, []interface{}{3, 3})
	f2 := NewDelayFunc(foo, []interface{}{5, 5})
	timerID1, err := Scheduler.CreateTimerAfter(f1, time.Duration(3)*time.Second)
	if nil != err {
		t.Log("Scheduler.CreateTimerAfter(f1, time.Duration(3)*time.Second)", "err：", err)
	}
	timerID2, err := Scheduler.CreateTimerAfter(f2, time.Duration(5)*time.Second)
	if nil != err {
		t.Log("Scheduler.CreateTimerAfter(f1, time.Duration(3)*time.Second)", "err：", err)
	}
	log.Printf("timerID1=%d ,timerID2=%d\n", timerID1, timerID2)
	Scheduler.CancelTimer(timerID1) //删除timerID1

	//阻塞等待
	select {}
}

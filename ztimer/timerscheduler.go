package ztimer

/**
* @Author: Aceld
* @Date: 2019/5/8 17:43
* @Mail: danbing.at@gmail.com
*
*  时间轮调度器
*   依赖模块，delayfunc.go  timer.go timewheel.go
 */

import (
	"math"
	"sync"
	"time"

	"github.com/aceld/zinx/zlog"
)

const (
	//MaxChanBuff 默认缓冲触发函数队列大小
	MaxChanBuff = 2048
	//MaxTimeDelay 默认最大误差时间
	MaxTimeDelay = 100
)

// TimerScheduler 计时器调度器
type TimerScheduler struct {
	//当前调度器的最高级时间轮
	tw *TimeWheel
	//定时器编号累加器
	IDGen uint32
	//已经触发定时器的channel
	triggerChan chan *DelayFunc
	//互斥锁
	sync.RWMutex
}

// NewTimerScheduler 返回一个定时器调度器 ，主要创建分层定时器，并做关联，并依次启动
func NewTimerScheduler() *TimerScheduler {

	//创建秒级时间轮
	secondTw := NewTimeWheel(SecondName, SecondInterval, SecondScales, TimersMaxCap)
	//创建分钟级时间轮
	minuteTw := NewTimeWheel(MinuteName, MinuteInterval, MinuteScales, TimersMaxCap)
	//创建小时级时间轮
	hourTw := NewTimeWheel(HourName, HourInterval, HourScales, TimersMaxCap)

	//将分层时间轮做关联
	hourTw.AddTimeWheel(minuteTw)
	minuteTw.AddTimeWheel(secondTw)

	//时间轮运行
	secondTw.Run()
	minuteTw.Run()
	hourTw.Run()

	return &TimerScheduler{
		tw:          hourTw,
		triggerChan: make(chan *DelayFunc, MaxChanBuff),
	}
}

// CreateTimerAt 创建一个定点Timer 并将Timer添加到分层时间轮中， 返回Timer的tID
func (ts *TimerScheduler) CreateTimerAt(df *DelayFunc, unixNano int64) (uint32, error) {
	ts.Lock()
	defer ts.Unlock()

	ts.IDGen++
	return ts.IDGen, ts.tw.AddTimer(ts.IDGen, NewTimerAt(df, unixNano))
}

// CreateTimerAfter 创建一个延迟Timer 并将Timer添加到分层时间轮中， 返回Timer的tID
func (ts *TimerScheduler) CreateTimerAfter(df *DelayFunc, duration time.Duration) (uint32, error) {
	ts.Lock()
	defer ts.Unlock()

	ts.IDGen++
	return ts.IDGen, ts.tw.AddTimer(ts.IDGen, NewTimerAfter(df, duration))
}

// CancelTimer 删除timer
func (ts *TimerScheduler) CancelTimer(tID uint32) {
	ts.Lock()
	defer ts.Unlock()

	tw := ts.tw
	for tw != nil {
		tw.RemoveTimer(tID)
		tw = tw.nextTimeWheel
	}
}

// GetTriggerChan 获取计时结束的延迟执行函数通道
func (ts *TimerScheduler) GetTriggerChan() chan *DelayFunc {
	return ts.triggerChan
}

// Start 非阻塞的方式启动timerSchedule
func (ts *TimerScheduler) Start() {
	go func() {
		for {
			//当前时间
			now := UnixMilli()
			//获取最近MaxTimeDelay 毫秒的超时定时器集合
			timerList := ts.tw.GetTimerWithIn(MaxTimeDelay * time.Millisecond)
			for _, timer := range timerList {
				if math.Abs(float64(now-timer.unixts)) > MaxTimeDelay {
					//已经超时的定时器，报警
					zlog.Error("want call at ", timer.unixts, "; real call at", now, "; delay ", now-timer.unixts)
				}
				ts.triggerChan <- timer.delayFunc
			}
			time.Sleep(MaxTimeDelay / 2 * time.Millisecond)
		}
	}()
}

// NewAutoExecTimerScheduler 时间轮定时器 自动调度
func NewAutoExecTimerScheduler() *TimerScheduler {
	//创建一个调度器
	autoExecScheduler := NewTimerScheduler()
	//启动调度器
	autoExecScheduler.Start()

	//永久从调度器中获取超时 触发的函数 并执行
	go func() {
		delayFuncChan := autoExecScheduler.GetTriggerChan()
		for df := range delayFuncChan {
			go df.Call()
		}
	}()

	return autoExecScheduler
}

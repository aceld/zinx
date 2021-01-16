package ztimer

/**
* @Author: Aceld
* @Date: 2019/4/30 17:42
* @Mail: danbing.at@gmail.com
 */

import (
	"time"
)

const (
	//HourName 小时
	HourName = "HOUR"
	//HourInterval 小时间隔ms为精度
	HourInterval = 60 * 60 * 1e3
	//HourScales  12小时制
	HourScales = 12

	//MinuteName 分钟
	MinuteName = "MINUTE"
	//MinuteInterval 每分钟时间间隔
	MinuteInterval = 60 * 1e3
	//MinuteScales 60分钟
	MinuteScales = 60

	//SecondName  秒
	SecondName = "SECOND"
	//SecondInterval 秒的间隔
	SecondInterval = 1e3
	//SecondScales  60秒
	SecondScales = 60
	//TimersMaxCap //每个时间轮刻度挂载定时器的最大个数
	TimersMaxCap = 2048
)

/*
   注意：
    有关时间的几个换算
   	time.Second(秒) = time.Millisecond * 1e3
	time.Millisecond(毫秒) = time.Microsecond * 1e3
	time.Microsecond(微秒) = time.Nanosecond * 1e3

	time.Now().UnixNano() ==> time.Nanosecond (纳秒)
*/

//Timer 定时器实现
type Timer struct {
	//延迟调用函数
	delayFunc *DelayFunc
	//调用时间(unix 时间， 单位ms)
	unixts int64
}

//UnixMilli 返回1970-1-1至今经历的毫秒数
func UnixMilli() int64 {
	return time.Now().UnixNano() / 1e6
}

//NewTimerAt   创建一个定时器,在指定的时间触发 定时器方法 df: DelayFunc类型的延迟调用函数类型；unixNano: unix计算机从1970-1-1至今经历的纳秒数
func NewTimerAt(df *DelayFunc, unixNano int64) *Timer {
	return &Timer{
		delayFunc: df,
		unixts:    unixNano / 1e6, //将纳秒转换成对应的毫秒 ms ，定时器以ms为最小精度
	}
}

//NewTimerAfter 创建一个定时器，在当前时间延迟duration之后触发 定时器方法
func NewTimerAfter(df *DelayFunc, duration time.Duration) *Timer {
	return NewTimerAt(df, time.Now().UnixNano()+int64(duration))
}

//Run 启动定时器，用一个go承载
func (t *Timer) Run() {
	go func() {
		now := UnixMilli()
		//设置的定时器是否在当前时间之后
		if t.unixts > now {
			//睡眠，直至时间超时,已微秒为单位进行睡眠
			time.Sleep(time.Duration(t.unixts-now) * time.Millisecond)
		}

		//调用事先注册好的超时延迟方法
		t.delayFunc.Call()
	}()
}

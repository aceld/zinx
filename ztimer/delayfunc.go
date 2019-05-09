/**
* @Author: Aceld
* @Date: 2019/4/30 11:57
* @Mail: danbing.at@gmail.com
*/
package ztimer

import (
	"fmt"
	"reflect"
	"zinx/zlog"
)

/*
   定义一个延迟调用函数
	延迟调用函数就是 时间定时器超时的时候，触发的事先注册好的
	回调函数
*/
type DelayFunc struct {
	f    func(...interface{}) //f : 延迟函数调用原型
	args []interface{}        //args: 延迟调用函数传递的形参
}

/*
	创建一个延迟调用函数
*/
func NewDelayFunc(f func(v ...interface{}), args []interface{}) *DelayFunc {
	return &DelayFunc{
		f:f,
		args:args,
	}
}

//打印当前延迟函数的信息，用于日志记录
func (df *DelayFunc) String() string {
	return fmt.Sprintf("{DelayFun:%s, args:%v}", reflect.TypeOf(df.f).Name(), df.args)
}



/*
	执行延迟函数---如果执行失败，抛出异常
 */
func (df *DelayFunc) Call() {
	defer func() {
		if err := recover(); err != nil {
			zlog.Error(df.String(), "Call err: ", err)
		}
	}()

	//调用定时器超时函数
	df.f(df.args...)
}
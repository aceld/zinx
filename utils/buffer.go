// Package utils 提供zinx相关工具类函数
// 包括:
//		全局配置
//		配置文件加载
//      buffer池
//
// 当前文件描述:
// @Title  buffer.go
// @Description  采用sync.Pool优化buffer的申请释放
// @Author  sirodeneko 2021/05/23
package utils

import (
	"bytes"
	"sync"
)

var bufferPool = sync.Pool{
	New: func() interface{} {
		return new(bytes.Buffer)
	},
}

// NewBuffer 从池中获取新 bytes.Buffer
func NewBuffer() *bytes.Buffer {
	return bufferPool.Get().(*bytes.Buffer)
}

// PutBuffer 将 Buffer放入池中
func PutBuffer(buf *bytes.Buffer) {
	// See https://golang.org/issue/23199
	const maxSize = 1 << 16
	if buf.Cap() < maxSize { // 对于大Buffer直接丢弃
		buf.Reset()
		bufferPool.Put(buf)
	}
}

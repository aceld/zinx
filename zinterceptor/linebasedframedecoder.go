package zinterceptor

import (
	"bytes"
	"fmt"
	"github.com/aceld/zinx/ziface"
	"sync"
)

// LineBasedFrameDecoder https://blog.csdn.net/h_sn9999/article/details/106492570
type LineBasedFrameDecoder struct {
	maxLength      int32
	failFast       bool //快速失败
	stripDelimiter bool
	discarding     bool
	discardedBytes int32
	offset         int32
	lock           sync.Mutex
	in             []byte
}

func NewLineBasedFrameDecoder(maxLength int32, stripDelimiter, failFast bool) ziface.IFrameDecoder {
	return LineBasedFrameDecoder{
		maxLength:      maxLength,
		failFast:       failFast,
		stripDelimiter: stripDelimiter,
		in:             make([]byte, 0),
	}
}
func (l LineBasedFrameDecoder) fail(length int) {
	panic(fmt.Sprintf("fail length: %d", length))
}

func (l LineBasedFrameDecoder) Decode(buff []byte) [][]byte {
	l.lock.Lock()
	defer l.lock.Unlock()

	l.in = append(l.in, buff...)
	resp := make([][]byte, 0)

	for {
		arr, skipCount := l.decode(l.in)
		if arr != nil {
			//证明已经解析出一个完整包
			resp = append(resp, arr)
			_size := len(arr) + skipCount
			if _size > 0 {
				l.in = l.in[_size:]
			}
		} else {
			return resp
		}
	}
}

func (l LineBasedFrameDecoder) decode(buf []byte) ([]byte, int) {
	in := bytes.NewBuffer(buf)
	eol := l.findEndOfLine(in)
	var length int
	if !l.discarding {
		if eol >= 0 {
			length = eol - 0
			var delimLength int
			if in.Bytes()[eol] == 13 {
				delimLength = 2
			} else {
				delimLength = 1
			}
			if length > int(l.maxLength) {
				bb := make([]byte, eol+delimLength)
				in.Read(bb)
				l.fail(length)
				return nil, 0
			} else {
				var newBytes []byte
				var skipCount int
				if l.stripDelimiter {
					index := 0
					newBytes = in.Bytes()
					newBytes = newBytes[index:length]
					skipCount = delimLength
				} else {
					index := 0 //in.Cap() - in.Len()
					newBytes = in.Bytes()
					newBytes = newBytes[index : length+delimLength]
					skipCount = 0
				}
				return newBytes, skipCount
			}
		} else {
			length = in.Len()
			if length > int(l.maxLength) {
				l.discardedBytes = int32(length)
				l.discarding = true
				l.offset = 0
				if l.failFast {
					l.fail(int(l.discardedBytes))
				}
			}
		}
	} else {
		if eol >= 0 {
			length = int(l.discardedBytes) + eol - (in.Cap() - in.Len())
			if in.Bytes()[eol] == 13 {
				length = 2
			} else {
				length = 1
			}
			bb := make([]byte, eol+length)
			in.Read(bb)
			l.discardedBytes = 0
			l.discarding = false
			if !l.failFast {
				l.fail(length)
			}
		} else {
			l.discardedBytes += int32(in.Len())
		}
	}
	return nil, 0
}

func (l LineBasedFrameDecoder) findEndOfLine(buffer *bytes.Buffer) int {
	totalLength := buffer.Len()
	readerIndex := 0 //buffer.Cap() - buffer.Len()
	bs := buffer.Bytes()
	start := (readerIndex) + int(l.offset)
	end := (totalLength) - int(l.offset)
	i := bytes.IndexByte(bs[start:end], 10)
	if i >= 0 {
		l.offset = 0
		if i > 0 && bs[i-1] == 13 {
			i = i - 1
		}
	} else {
		l.offset = int32(totalLength)
	}
	return i
}

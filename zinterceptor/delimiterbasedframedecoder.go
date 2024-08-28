package zinterceptor

import (
	"bytes"
	"fmt"
	"github.com/aceld/zinx/ziface"
	"math"
	"sync"
)

type DelimiterBasedFrameDecoder struct {
	delimiters             [][]byte
	maxFrameLength         int32
	stripDelimiter         bool
	failFast               bool
	discardingTooLongFrame bool
	tooLongFrameLength     int32
	lineBasedDecoder       *LineBasedFrameDecoder
	lock                   sync.Mutex
	in                     []byte
}

func validateMaxFrameLength(maxFrameLength int32) {
	if maxFrameLength <= 0 {
		panic("maxFrameLength must be greater than zero")
	}
}
func validateDelimiter(delimiter []byte) {
	if delimiter == nil {
		panic("delimiter must not be nil")
	} else if len(delimiter) <= 0 {
		panic("delimiter must be greater than zero")
	}
}
func isLineBased(delimiters [][]byte) bool {
	if len(delimiters) != 2 {
		return false
	} else {
		a := delimiters[0]
		b := delimiters[1]
		if len(a) < len(b) {
			a = delimiters[1]
			b = delimiters[0]
		}
		return len(a) == 2 && len(b) == 1 && a[0] == 13 && a[1] == 10 && b[0] == 10
	}
}
func NewDelimiterBasedFrameDecoder(maxFrameLength int32, stripDelimiter, failFast bool, delimiters ...[]byte) ziface.IFrameDecoder {
	validateMaxFrameLength(maxFrameLength)
	if delimiters == nil {
		panic("delimiter is nil")
	} else if len(delimiters) == 0 {
		panic("delimiter is empty")
	} else {
		d := DelimiterBasedFrameDecoder{}
		if isLineBased(delimiters) {
			d.lineBasedDecoder = &LineBasedFrameDecoder{
				maxLength:      maxFrameLength,
				stripDelimiter: stripDelimiter,
				failFast:       failFast,
			}
			d.delimiters = nil
		} else {
			d.delimiters = make([][]byte, 0)
			for _, delimiter := range delimiters {
				d.delimiters = append(d.delimiters, delimiter)
			}
			d.lineBasedDecoder = nil
		}
		d.maxFrameLength = maxFrameLength
		d.stripDelimiter = stripDelimiter
		d.failFast = failFast
		return &d
	}
}

func (this *DelimiterBasedFrameDecoder) Decode(buff []byte) [][]byte {
	this.lock.Lock()
	defer this.lock.Unlock()

	this.in = append(this.in, buff...)
	resp := make([][]byte, 0)

	for {
		arr, skipCount := this.decode(this.in)
		if arr != nil {
			//证明已经解析出一个完整包
			resp = append(resp, arr)
			_size := len(arr) + skipCount
			if _size > 0 {
				this.in = this.in[_size:]
			}
		} else {
			return resp
		}
	}
}

func (this *DelimiterBasedFrameDecoder) decode(buffer []byte) ([]byte, int) {
	in := bytes.NewBuffer(buffer)
	if this.lineBasedDecoder != nil {
		return this.lineBasedDecoder.decode(buffer)
	} else {
		minFrameLength := math.MaxInt32
		var skipCount int
		var minDelim []byte
		for _, delimiter := range this.delimiters {
			frameLength := indexOf(buffer, delimiter)
			if frameLength >= 0 && frameLength < int32(minFrameLength) {
				minFrameLength = int(frameLength)
				minDelim = delimiter
			}
		}
		if minDelim != nil {
			minDelimLength := len(minDelim)
			var newBytes []byte
			if this.discardingTooLongFrame {
				this.discardingTooLongFrame = false
				//in.Next(minFrameLength + minDelimLength)
				skipCount = minFrameLength + minDelimLength

				tooLongFrameLength := this.tooLongFrameLength
				this.tooLongFrameLength = 0
				if !this.failFast {
					this.fail(tooLongFrameLength)
				}
				return nil, skipCount
			}

			if minFrameLength > int(this.maxFrameLength) {
				//in.Next(minFrameLength + minDelimLength)
				skipCount = minFrameLength + minDelimLength
				this.fail(int32(minFrameLength))
				return nil, skipCount
			}

			if this.stripDelimiter {
				index := 0
				newBytes = in.Bytes()
				newBytes = newBytes[index:minFrameLength]
				skipCount = minDelimLength
			} else {
				index := 0 //in.Cap() - in.Len()
				newBytes = in.Bytes()
				newBytes = newBytes[index : minFrameLength+minDelimLength]
				skipCount = 0
			}
			return newBytes, skipCount
		} else {
			if !this.discardingTooLongFrame {
				if int32(in.Len()) > this.maxFrameLength {
					this.tooLongFrameLength = int32(in.Len())
					skipCount = in.Len()
					this.discardingTooLongFrame = true
					if this.failFast {
						this.fail(this.tooLongFrameLength)
					}
				}
			} else {
				this.tooLongFrameLength += int32(in.Len())
				skipCount = in.Len()
			}
			return nil, skipCount
		}
	}
}

func indexOf(haystack, needle []byte) int32 {
	for i := 0; i < len(haystack); i++ {
		haystackIndex := i
		var needleIndex int
		for needleIndex = 0; needleIndex < len(needle); needleIndex++ {
			if haystack[haystackIndex] != needle[needleIndex] {
				break
			} else {
				haystackIndex++
				if haystackIndex == len(haystack) && needleIndex != (len(needle)-1) {
					return -1
				}
			}
		}

		if needleIndex == len(needle) {
			return int32(i)
		}
	}
	return -1
}

func (this *DelimiterBasedFrameDecoder) fail(length int32) {
	panic(fmt.Sprintf("fail length: %d", length))
}

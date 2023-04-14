/**
 * @author uuxia
 * @date 15:58 2023/3/10
 * @description 通过拦截，处理数据，任务向下传递
 **/

package zinterceptor

// 暂时不用

/*
// Interceptor 基于LengthField规则的拦截器
type Interceptor struct {
	frameDecoder ziface.IFrameDecoder
}

func NewInterceptor(maxFrameLength uint64,
	lengthFieldOffset, lengthFieldLength, lengthAdjustment, initialBytesToStrip int) ziface.IInterceptor {
	return &Interceptor{
		frameDecoder: NewFrameDecoderByParams(maxFrameLength, lengthFieldOffset, lengthFieldLength, lengthAdjustment, initialBytesToStrip),
	}
}

func (l *Interceptor) Intercept(chain ziface.IChain) ziface.IcResp {
	req := chain.Request()

	if req == nil || l.frameDecoder == nil {
		goto END
	}

	switch req.(type) {
	case ziface.IRequest:
		iRequest := req.(ziface.IRequest)
		iMessage := iRequest.GetMessage()

		if iMessage == nil {
			break
		}

		data := iMessage.GetData()

		bytebuffers := l.frameDecoder.Decode(data)
		size := len(bytebuffers)
		if size == 0 { //半包，或者其他情况，任务就不要往下再传递了
			return nil
		}

		for i := 0; i < size; i++ {
			buffer := bytebuffers[i]
			if buffer == nil {
				continue
			}
			bufferSize := len(buffer)
			iMessage.SetData(buffer)
			iMessage.SetDataLen(uint32(bufferSize))

			if i < size-1 {
				chain.Proceed(chain.Request())
			}
		}
	}

END:
	return chain.Proceed(chain.Request())
}
*/

package ziface

import "time"

type MsgSendOptionObj struct {
	Timeout time.Duration
}

type MsgSendOption func(opt *MsgSendOptionObj)

func WithSendMsgTimeout(timeout time.Duration) MsgSendOption {
	return func(opt *MsgSendOptionObj) {
		opt.Timeout = timeout
	}
}

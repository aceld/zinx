package znet

import (
	"time"
)

const (
	maxDelay = 1 * time.Second
)

var AcceptDelay *acceptDelay

func init() {
	AcceptDelay = &acceptDelay{duration: 0}
}

type acceptDelay struct {
	duration time.Duration
}

func (d *acceptDelay) Delay() {
	d.Up()
	d.do()
}

func (d *acceptDelay) Reset() {
	d.duration = 0
}

func (d *acceptDelay) Up() {
	if d.duration == 0 {
		d.duration = 5 * time.Millisecond
		return
	}
	d.duration = 2 * d.duration
	if d.duration > maxDelay {
		d.duration = maxDelay
	}
}

func (d *acceptDelay) do() {
	if d.duration > 0 {
		time.Sleep(d.duration)
	}
}

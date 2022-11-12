package reqcap

import (
	"os/exec"
	"time"
)

type ReqCap struct {
	timer    *time.Timer
	callback func(o *ReqCap)
	duration time.Duration
	state    bool
	Cmd      *exec.Cmd
	Name     string
	Pid      string
}

func New_reqCap(call func(o *ReqCap), dur time.Duration) *ReqCap {
	aux := ReqCap{
		timer:    time.NewTimer(time.Second * 20),
		callback: call,
		duration: dur,
		state:    true,
	}
	aux.timer.Stop()

	go func(r *ReqCap) {
		for r.state {
			<-r.timer.C
			r.callback(r)
		}
	}(&aux)

	return &aux
}

func (o *ReqCap) Close_reqCap() {
	o.timer.Stop()
	o.state = false
}

func (o *ReqCap) Capture(name string) bool {
	o.Name = name
	return o.timer.Reset(o.duration)
}

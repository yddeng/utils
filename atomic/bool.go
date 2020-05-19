package atomic

import "sync/atomic"

type Bool struct {
	b int32
}

func (a *Bool) Set(b bool) {
	var newV int32
	if b {
		newV = 1
	}
	atomic.SwapInt32(&a.b, newV)
}

func (a *Bool) Get() bool {
	return atomic.LoadInt32(&a.b) == 1
}

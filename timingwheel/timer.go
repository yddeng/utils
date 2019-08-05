package timingwheel

import (
	"time"
)

type TimerI interface {
	//到期时间与现在时间的时间间隔
	SubDuration(time.Time) time.Duration
	//到期后行为
	Do()
}

type Timer struct {
	expiredTime time.Time
	delayCall   func()
}

func newTimer(expiredTime time.Time, cb func()) *Timer {
	return &Timer{
		expiredTime: expiredTime,
		delayCall:   cb,
	}
}

func (t *Timer) SubDuration(now time.Time) time.Duration {
	return t.expiredTime.Sub(now)
}

func (t *Timer) Do() {
	//fmt.Println("do", t.expiredTime)
	t.delayCall()
}

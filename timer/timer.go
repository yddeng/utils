package timer

import (
	"github.com/yddeng/dutil"
	"log"
	"time"
)

type Timer interface {
	// 重置到期时间
	// 已执行或未执行都会重置 timer 的到期时间。重复定时器会一同改变周期
	// 当且仅当 定时器已经停止，返回 false。
	Reset(duration time.Duration) bool

	// 停止定时器
	Stop() bool
}

type runtimeTimer struct {
	when     int64 // 到期时间
	ctx      interface{}
	fn       func(ctx interface{})
	repeated bool
	period   int64 // 周期
}

func goFunc(fn func(ctx interface{}), ctx interface{}) {
	go func() {
		defer func() {
			err := dutil.Recover()
			if err != nil {
				log.Printf("timer: goFunc Recover %s", err)
			}
		}()
		fn(ctx)
	}()
}

func sendSignal(ch chan struct{}) {
	select {
	case ch <- struct{}{}:
	default:
	}
}

// when is a helper function for setting the 'when' field of a runtimeTimer.
// It returns what the time will be, in nanoseconds, Duration d in the future.
// If d is negative, it is ignored. If the returned value would be less than
// zero because of an overflow, MaxInt64 is returned.
func when(d time.Duration) int64 {
	if d <= 0 {
		return time.Now().UnixNano()
	}
	t := time.Now().UnixNano() + int64(d)
	if t < 0 {
		t = 1<<63 - 1 // math.MaxInt64
	}
	return t
}

type TimerMgr interface {
	// 一次性定时器
	OnceTimer(d time.Duration, ctx interface{}, f func(ctx interface{})) Timer

	// 重复定时器
	RepeatTimer(d time.Duration, ctx interface{}, f func(ctx interface{})) Timer
}

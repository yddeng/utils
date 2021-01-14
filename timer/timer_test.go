package timer

import (
	"fmt"
	"testing"
	"time"
)

func TestNewHeapTimerMgr(t *testing.T) {
	mgr := NewHeapTimerMgr()

	fmt.Println(time.Now().String())
	timer1 := mgr.OnceTimer(time.Second, "once1", func(ctx interface{}) {
		fmt.Println(ctx, time.Now().String())
	})

	timer2 := mgr.RepeatTimer(time.Second*2, "repeat1", func(ctx interface{}) {
		fmt.Println(ctx, time.Now().String())
	})

	time.Sleep(time.Second * 5)
	timer1.Reset(time.Second * 3)
	timer2.Reset(time.Second)

	time.Sleep(time.Second * 5)
	timer1.Stop()
	timer2.Stop()

	fmt.Println(timer1.Reset(time.Second))
	select {}
}

func TestNewTimeWheelMgr(t *testing.T) {
	mgr := NewTimeWheelMgr(time.Millisecond*200, 10)

	fmt.Println(time.Now().String())
	timer1 := mgr.OnceTimer(time.Second, "once1", func(ctx interface{}) {
		fmt.Println(ctx, time.Now().String())
	})

	timer2 := mgr.RepeatTimer(time.Second*3, "repeat1", func(ctx interface{}) {
		fmt.Println(ctx, time.Now().String())
	})

	time.Sleep(time.Second * 5)
	timer1.Reset(time.Second * 3)
	timer2.Reset(time.Second)

	time.Sleep(time.Second * 5)
	timer1.Stop()
	timer2.Stop()

	fmt.Println(timer1.Reset(time.Second))
	select {}
}

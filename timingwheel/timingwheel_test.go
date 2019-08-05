package timingwheel

import (
	"fmt"
	"testing"
	"time"
)

func TestNewTimingWheel(t *testing.T) {
	tw, _ := NewTimingWheel(1*time.Second, 20)
	fmt.Println(time.Now())

	tw.DelayFunc(5*time.Second, func() {
		fmt.Println("5 delayTime", time.Now())
	})

	tw.DelayFunc(8*time.Second, func() {
		fmt.Println("8 delayTime", time.Now())
	})
	tw.DelayFunc(10*time.Second, func() {
		fmt.Println("10 delayTime", time.Now())
	})
	tw.DelayFunc(16*time.Second, func() {
		fmt.Println("16 delayTime", time.Now())
	})

	tw.DelayFunc(3*time.Second, func() {
		fmt.Println("3 delayTime", time.Now())
	})

	select {}
}

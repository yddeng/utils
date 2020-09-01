package ticker

import (
	"fmt"
	"testing"
	"time"
)

func TestNewTicker(t *testing.T) {

	RegisterOnceTicker(time.Millisecond*500, func(t time.Time) {
		fmt.Println("ticker1 once", t)
	})
	RegisterOnceTicker(time.Millisecond*1000, func(t time.Time) {
		fmt.Println("ticker2 once", t)
	})
	RegisterRepeatTicker(time.Millisecond*1500, func(t time.Time) {
		fmt.Println("ticker3 repeat", t)
	})
	ticker4 := RegisterRepeatTicker(time.Millisecond*1000, func(t time.Time) {
		fmt.Println("ticker4 repeat", t)
	})

	time.Sleep(time.Second * 3)
	ticker4.Stop()
	select {}
}

package ticker_test

import (
	"fmt"
	"github.com/tagDong/dutil/ticker"
	"testing"
	"time"
)

func TestNewTicker(t *testing.T) {

	ticker.RegisterOnceTicker(time.Millisecond*500, func(t time.Time) {
		fmt.Println("ticker1 once", t)
	})
	ticker.RegisterOnceTicker(time.Millisecond*1000, func(t time.Time) {
		fmt.Println("ticker2 once", t)
	})
	ticker.RegisterRepeatTicker(time.Millisecond*1500, func(t time.Time) {
		fmt.Println("ticker3 repeat", t)
	})
	ticker2 := ticker.RegisterRepeatTicker(time.Millisecond*1000, func(t time.Time) {
		fmt.Println("ticker4 repeat", t)
	})

	time.Sleep(time.Second * 3)
	ticker2.Stop()
	select {}
}

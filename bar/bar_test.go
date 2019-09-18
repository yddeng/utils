/*
 *
 */
package bar_test

import (
	"fmt"
	"github.com/tagDong/dutil/bar"
	"math/rand"
	"testing"
	"time"
)

func TestNewBar(t *testing.T) {

	c := make(chan struct{})
	fmt.Println("bar bar")
	total := 10000
	b := bar.NewBar("bar", total)
	go func() {
		for i := 0; i < total; {
			count := rand.Int() % 10
			b.Add(count)
			i += count
			time.Sleep(time.Millisecond)
		}
		//fmt.Println("bar bar end")
		c <- struct{}{}
	}()

	<-c
	b1 := bar.NewBar("bar1", total)
	go func() {
		for i := 0; i < total; {
			count := rand.Int() % 10
			b1.Add(count)
			i += count
			time.Sleep(time.Millisecond)
		}
	}()
	select {}
}

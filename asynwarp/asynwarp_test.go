package asynwarp

import (
	"sync"
	"testing"
)

func TestWrapFunc(t *testing.T) {
	wg := sync.WaitGroup{}
	wg.Add(2)
	fn1 := func(i int, s string) {
		t.Log(i, s)
		wg.Done()
	}

	fn2 := func(s string) string {
		return s + "..."
	}

	WrapFunc(fn1)(nil, 1, "sss")
	WrapFunc(fn2)(func(ret string) {
		t.Log("fn2 result", ret)
		wg.Done()
	}, "sdsd")

	wg.Wait()
}

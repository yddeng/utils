package transaction

import (
	"fmt"
	"testing"
)

func TestNewTransaction(t *testing.T) {
	tt1, tt2 := 3, 4
	ok := false
	trans := NewTransaction()
	trans.Push(func() bool {
		tt1 += 3
		return true
	}, func() bool {
		tt1 -= 3
		return true
	})

	trans.Push(func() bool {
		tt2 += 3
		if ok {
			return true
		} else {
			return false
		}
	}, func() bool {
		tt2 -= 3
		return true
	})

	trans.Do(func(susses bool) {
		fmt.Println(susses, tt1, tt2, ok)
	})
}

package pool

import (
	"fmt"
	"testing"
)

func TestNewTreadPool(t *testing.T) {
	tp := NewTreadPool(2, 256)

	for i := 0; i < 500; i++ {
		id := i
		tp.AddTask(func() {
			fmt.Println("---", id)
		})
	}
}

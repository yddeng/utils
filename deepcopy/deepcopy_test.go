package deepcopy

import (
	"fmt"
	"testing"
)

func TestDeepCopy(t *testing.T) {
	src := []int{1, 2, 3}

	var dst []int
	DeepCopy(&dst, &src)

	for _, v := range dst {
		fmt.Println(v)
	}
}

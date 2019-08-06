package heap_test

import (
	"fmt"
	"github.com/tagDong/dutil/heap"
	"testing"
)

type HElement struct {
	name  string
	value int
}

func (e *HElement) Less(h heap.Element) bool {
	return e.value < h.(*HElement).value
}

func TestNewMinHeap(t *testing.T) {
	h := heap.NewHeap()
	h.Push(&HElement{name: "e1", value: 5})
	h.Push(&HElement{name: "e2", value: 12})
	h.Push(&HElement{name: "e3", value: 8})
	e4 := &HElement{name: "e4", value: 7}
	h.Push(e4)
	fmt.Println(h.Peek())

	e4.value = 10
	h.Fix(e4)
	fmt.Println(h.Peek())

	h.Remove(e4)
	fmt.Println(h.Peek())

	fmt.Println(h.Pop())
	fmt.Println(h.Pop())
	fmt.Println(h.Pop())
	fmt.Println(h.Pop())
}

package heap

import (
	"fmt"
	"testing"
)

type HElement struct {
	name  string
	value int
}

func (e *HElement) Less(h interface{}) bool {
	return e.value < h.(*HElement).value
}

func TestNewMinHeap(t *testing.T) {
	h := New()
	h.Push(&HElement{name: "e1", value: 5})
	h.Push(&HElement{name: "e2", value: 12})
	h.Push(&HElement{name: "e3", value: 8})
	e4 := &HElement{name: "e4", value: 7}
	h.Push(e4)
	fmt.Println(h.Top())

	e4.value = 2
	h.Fix(e4)
	fmt.Println(h.Top(), h.IsExist(e4))

	h.Remove(e4)
	fmt.Println(h.Top(), h.IsExist(e4))

	fmt.Println(h.Pop())
	fmt.Println(h.Pop())
	fmt.Println(h.Pop())
	fmt.Println(h.Pop())

}

func TestHeap_Fix(t *testing.T) {
	h := New()
	e1 := &HElement{name: "e1", value: 5}
	e2 := &HElement{name: "e2", value: 12}
	e3 := &HElement{name: "e3", value: 8}
	e4 := &HElement{name: "e4", value: 7}

	h.PushList(e1, e2, e3, e4)
	fmt.Println(h.Pop(), h.Pop(), h.Pop(), h.Pop())

	h.PushList(e1, e2, e3, e4)
	e1.value = 20
	h.Fix(e1)
	fmt.Println(h.Pop(), h.Pop(), h.Pop(), h.Pop())
}

func TestHeap_Reset(t *testing.T) {
	h := New()

	h.Push(&HElement{name: "e1", value: 5})
	h.Push(&HElement{name: "e1", value: 12})

	fmt.Println(h.Len())
	h.Reset()
	fmt.Println(h.Len())
}

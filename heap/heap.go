package heap

import "container/heap"

type Element interface {
	Less(element Element) bool
}

type myHeap struct {
	elements []Element
	elemIdx  map[Element]int
}

type Heap struct {
	myHeap *myHeap
}

func NewHeap() *Heap {
	h := &Heap{
		myHeap: &myHeap{
			elements: []Element{},
			elemIdx:  map[Element]int{},
		},
	}
	return h
}

func (h *myHeap) Less(i, j int) bool {
	return h.elements[i].Less(h.elements[j])
}

func (h *myHeap) Swap(i, j int) {
	h.elemIdx[h.elements[i]] = j
	h.elemIdx[h.elements[j]] = i
	h.elements[i], h.elements[j] = h.elements[j], h.elements[i]
}

func (h *myHeap) Len() int {
	return len(h.elements)
}

func (h *myHeap) Pop() (v interface{}) {
	h.elements, v = h.elements[:h.Len()-1], h.elements[h.Len()-1]
	item := v.(Element)
	delete(h.elemIdx, item)
	return item
}

func (h *myHeap) Push(v interface{}) {
	h.elemIdx[v.(Element)] = h.Len()
	h.elements = append(h.elements, v.(Element))
}

func (h *Heap) Len() int {
	return h.myHeap.Len()
}

func (h *Heap) Push(item Element) {
	heap.Push(h.myHeap, item)
}

func (h *Heap) Pop() Element {
	if h.Len() > 0 {
		return heap.Pop(h.myHeap).(Element)
	}
	return nil
}

func (h *Heap) Peek() Element {
	if h.Len() > 0 {
		return h.myHeap.elements[0]
	}
	return nil
}

func (h *Heap) In(ele Element) bool {
	_, ok := h.myHeap.elemIdx[ele]
	return ok
}

func (h *Heap) Remove(ele Element) {
	i, ok := h.myHeap.elemIdx[ele]
	if ok {
		heap.Remove(h.myHeap, i)
	}
}

func (h *Heap) Fix(ele Element) {
	i, ok := h.myHeap.elemIdx[ele]
	if ok {
		heap.Fix(h.myHeap, i)
	}
}

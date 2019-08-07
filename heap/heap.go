package heap

type Element interface {
	//实现大根堆、小根堆
	Less(element Element) bool
}

type Heap struct {
	data  []Element
	idMap map[Element]int
}

func NewHeap() *Heap {
	h := &Heap{
		data:  []Element{},
		idMap: map[Element]int{},
	}
	return h
}

func (h *Heap) Len() int {
	return len(h.data)
}

func (h *Heap) Push(item Element) {
	n := h.Len()
	h.idMap[item] = n
	h.data = append(h.data, item)

	h.up(h.Len() - 1)
}

func (h *Heap) Pop() Element {
	if h.Len() > 0 {
		n := h.Len() - 1
		h.swap(0, n)
		h.down(0, n)

		item := h.data[n]
		delete(h.idMap, item)
		h.data = h.data[0:n]

		return item
	}
	return nil
}

//返回根节点
func (h *Heap) Peek() Element {
	if h.Len() > 0 {
		return h.data[0]
	}
	return nil
}

func (h *Heap) Remove(ele Element) {
	i, ok := h.idMap[ele]
	if ok {
		n := h.Len() - 1
		if n != i {
			h.swap(i, n)
			if !h.down(i, n) {
				h.up(i)
			}
		}

		item := h.data[n]
		delete(h.idMap, item)
		h.data = h.data[0:n]
	}
}

//改变值,需要从新排序
//比 先remove再push 更少开销
func (h *Heap) Fix(ele Element) {
	i, ok := h.idMap[ele]
	if ok {
		if !h.down(i, h.Len()) {
			h.up(i)
		}
	}
}

//
func (h *Heap) up(j int) {
	for {
		i := (j - 1) / 2 // parent
		if i == j || !h.less(j, i) {
			break
		}
		h.swap(i, j)
		j = i
	}
}

//[i0,n)间做交换
func (h *Heap) down(i0, n int) bool {
	i := i0
	for {
		j1 := 2*i + 1
		if j1 >= n || j1 < 0 { // j1 < 0 after int overflow
			break
		}
		j := j1 // left child
		if j2 := j1 + 1; j2 < n && h.less(j2, j1) {
			j = j2 // = 2*i + 2  // right child
		}
		if !h.less(j, i) {
			break
		}
		h.swap(i, j)
		i = j
	}
	return i > i0
}

func (h *Heap) less(i, j int) bool {
	return h.data[i].Less(h.data[j])
}

func (h Heap) swap(i, j int) {
	h.data[i], h.data[j] = h.data[j], h.data[i]
	h.idMap[h.data[i]] = i
	h.idMap[h.data[j]] = j
}

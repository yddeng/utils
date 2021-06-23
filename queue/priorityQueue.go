package queue

import (
	"sort"
	"sync"
)

// 优先级队列

type priorityQueue struct {
	priority uint32
	data     []interface{}
	f, r     int
	size     int
}

func newPQ(size int, priority uint32) *priorityQueue {
	return &priorityQueue{
		priority: priority,
		data:     make([]interface{}, size),
		size:     size,
	}
}

func (q *priorityQueue) full() bool {
	return (q.r+1)%q.size == q.f
}

func (q *priorityQueue) empty() bool {
	return q.r == q.f
}

func (q *priorityQueue) add(v interface{}) {
	q.data[q.r] = v
	q.r = (q.r + 1) % q.size
}

func (q *priorityQueue) get() (v interface{}) {
	v = q.data[q.f]
	q.f = (q.f + 1) % q.size
	return
}

func (q *priorityQueue) len() int {
	rear := q.r
	if rear < q.f {
		rear += q.size
	}
	return rear - q.f
}

type PriorityQueue struct {
	size   int
	count  int
	pqList []*priorityQueue

	lock   *sync.Mutex
	fullC  *sync.Cond
	emptyC *sync.Cond

	closeOnce sync.Once
	closeCh   chan struct{}
}

func (q *PriorityQueue) sort() {
	sort.Slice(q.pqList, func(i, j int) bool {
		return q.pqList[i].priority > q.pqList[j].priority
	})
}

func (q *PriorityQueue) push(b bool, v interface{}, priority uint32) error {
	q.lock.Lock()
	for q.count >= q.size {
		select {
		case <-q.closeCh:
			q.lock.Unlock()
			return ErrClosed
		default:
			if !b {
				q.lock.Unlock()
				return ErrQueueFull
			}
			q.fullC.Wait()
		}
	}

	if q.Closed() {
		q.lock.Unlock()
		return ErrClosed
	}

	var pq *priorityQueue
	for _, v := range q.pqList {
		if v.priority == priority {
			pq = v
			break
		}
		if v.priority < priority {
			break
		}
	}

	if pq == nil {
		pq = newPQ(q.size, priority)
		q.pqList = append(q.pqList, pq)
		q.sort()
	}

	notify := q.count == 0
	pq.add(v)
	q.count++
	q.lock.Unlock()

	if notify {
		q.emptyC.Signal()
	}
	return nil
}

func (q *PriorityQueue) Push(v interface{}, priority ...uint32) error {
	var pro uint32
	if len(priority) > 0 {
		pro = priority[0]
	}

	select {
	case <-q.closeCh:
		return ErrClosed
	default:
		return q.push(true, v, pro)
	}
}

func (q *PriorityQueue) PushN(v interface{}, priority ...uint32) error {
	var pro uint32
	if len(priority) > 0 {
		pro = priority[0]
	}

	select {
	case <-q.closeCh:
		return ErrClosed
	default:
		return q.push(false, v, pro)
	}
}

func (q *PriorityQueue) Pop() (interface{}, bool) {
	q.lock.Lock()
	for q.count == 0 {
		select {
		case <-q.closeCh:
			q.lock.Unlock()
			return nil, false
		default:
			q.emptyC.Wait()
		}
	}

	notify := q.count == q.size

	v := q.pqList[0].get()
	if q.pqList[0].len() == 0 && len(q.pqList) > 1 {
		q.pqList = q.pqList[1:]
	}
	q.count--
	q.lock.Unlock()

	if notify {
		q.fullC.Signal()
	}
	return v, true
}

func (cq *PriorityQueue) Closed() bool {
	select {
	case <-cq.closeCh:
		return true
	default:
		return false
	}
}

func (q *PriorityQueue) Close() {
	q.closeOnce.Do(func() {
		close(q.closeCh)
		q.emptyC.Broadcast()
		q.fullC.Broadcast()
	})
}

func (q *PriorityQueue) Len() int {
	q.lock.Lock()
	defer q.lock.Unlock()
	return q.count
}

func NewPriorityQueue(size int) *PriorityQueue {
	lock := new(sync.Mutex)
	pq := &PriorityQueue{
		size:      size,
		lock:      lock,
		pqList:    make([]*priorityQueue, 0, 4),
		fullC:     sync.NewCond(lock),
		emptyC:    sync.NewCond(lock),
		closeOnce: sync.Once{},
		closeCh:   make(chan struct{}),
	}
	pq.pqList = append(pq.pqList, newPQ(size, 0))
	return pq
}

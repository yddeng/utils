package queue

import (
	"fmt"
	"sync"
)

var (
	ErrClosed = fmt.Errorf("queue is closed")
)

type BlockQueue struct {
	cap    int
	data   []interface{}
	mutex  *sync.Mutex
	emptyC *sync.Cond
	fullC  *sync.Cond
	front  int
	rear   int
	closed bool
}

func (q *BlockQueue) Push(v interface{}) error {
	q.mutex.Lock()
	if q.closed {
		q.mutex.Unlock()
		return ErrClosed
	}

	for !q.closed && q.full() {
		q.fullC.Wait()
		if q.closed {
			q.mutex.Unlock()
			return ErrClosed
		}
	}
	needSignal := q.empty()
	q.data[q.rear] = v
	q.rear = (q.rear + 1) % q.cap
	q.mutex.Unlock()

	if needSignal {
		q.emptyC.Broadcast()
	}
	return nil
}

func (q *BlockQueue) Peek() interface{} {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	if q.empty() {
		return nil
	}
	return q.data[q.front]
}

func (q *BlockQueue) Pop() (ret interface{}, closed bool) {
	q.mutex.Lock()
	for !q.closed && q.empty() {
		q.emptyC.Wait()
	}
	needSignal := q.full()
	if q.len() > 0 {
		ret = q.data[q.front]
		q.front = (q.front + 1) % q.cap
	}
	closed = q.closed
	q.mutex.Unlock()
	if needSignal {
		q.fullC.Signal()
	}
	return
}

func (q *BlockQueue) GetAll() (ret []interface{}, closed bool) {
	q.mutex.Lock()
	for !q.closed && q.empty() {
		q.emptyC.Wait()
	}
	needSignal := q.full()
	if q.len() > 0 {
		ret = []interface{}{}
		for i := 0; i < q.len(); i++ {
			idx := (q.front + i) % q.cap
			ret = append(ret, q.data[idx])
		}
		q.front = 0
		q.rear = 0
	}
	closed = q.closed
	q.mutex.Unlock()
	if needSignal {
		q.fullC.Signal()
	}
	return
}

func (q *BlockQueue) Len() int {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	return q.len()
}

func (q *BlockQueue) Close() {
	q.mutex.Lock()
	if q.closed {
		return
	}
	q.closed = true
	q.mutex.Unlock()

	q.emptyC.Broadcast()
	q.fullC.Broadcast()
}

func (q *BlockQueue) Closed() bool {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	return q.closed
}

func (q *BlockQueue) full() bool {
	return (q.rear+1)%q.cap == q.front
}

func (q *BlockQueue) empty() bool {
	return q.front == q.rear
}

func (q *BlockQueue) len() int {
	rear := q.rear
	if rear < q.front {
		rear += q.cap
	}
	return rear - q.front
}

func NewBlockQueue(cap int) *BlockQueue {
	mutex := &sync.Mutex{}
	return &BlockQueue{
		cap:    cap,
		data:   make([]interface{}, cap),
		mutex:  mutex,
		emptyC: sync.NewCond(mutex),
		fullC:  sync.NewCond(mutex),
	}
}

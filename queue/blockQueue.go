package queue

import (
	"fmt"
	"sync"
)

var (
	errClosed = fmt.Errorf("queue is closed")
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
		return errClosed
	}

	for !q.closed && q.full() {
		q.fullC.Wait()
		if q.closed {
			q.mutex.Unlock()
			return errClosed
		}
	}
	q.data[q.rear] = v
	q.rear = (q.rear + 1) % q.cap
	q.mutex.Unlock()
	q.emptyC.Broadcast()
	return nil
}

/*
func (q *BlockQueue) ForcePush(v interface{}) (interface{}, error) {
	var poped interface{}
	q.mutex.Lock()
	if q.closed {
		q.mutex.Unlock()
		return nil,ErrClosed
	}
	if q.full() {
		poped = q.data[q.front]
		q.front = (q.front + 1) % q.cap
	}
	q.data[q.rear] = v
	q.rear = (q.rear + 1) % q.cap
	q.emptyC.Broadcast()

	return poped,nil
}
*/

func (q *BlockQueue) Pop() (ret interface{}, closed bool) {
	q.mutex.Lock()
	for !q.closed && q.empty() {
		q.emptyC.Wait()
	}
	if q.len() > 0 {
		ret = q.data[q.front]
		q.front = (q.front + 1) % q.cap
	}
	closed = q.closed
	q.mutex.Unlock()
	q.fullC.Broadcast()
	return
}

func (q *BlockQueue) GetAll() (ret []interface{}, closed bool) {
	q.mutex.Lock()
	for !q.closed && q.empty() {
		q.emptyC.Wait()
	}
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
	q.fullC.Broadcast()
	return
}

func (q *BlockQueue) Len() int {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	return q.len()
}

//
func (q *BlockQueue) Do(f func(v interface{})) {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	for i := 0; i < q.len(); i++ {
		idx := (q.front + i) % q.cap
		f(q.data[idx])
	}
}

func (q *BlockQueue) Close() {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	if q.closed {
		return
	}

	q.closed = true
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
	data := make([]interface{}, cap)
	return &BlockQueue{cap: cap, data: data,
		mutex: mutex, emptyC: sync.NewCond(mutex),
		fullC: sync.NewCond(mutex),
	}
}

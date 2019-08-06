package queue

import (
	"fmt"
	"sync/atomic"
)

var ErrClosed = fmt.Errorf("queue is closed")

//close flag
const queClosed = 1

type BlockQueue struct {
	buf    chan interface{}
	closed int32
}

func NewBlockQueue(size int) *BlockQueue {
	return &BlockQueue{
		buf: make(chan interface{}, size),
	}
}

func (q *BlockQueue) Len() int {
	return len(q.buf)
}

func (q *BlockQueue) Push(elem interface{}) error {
	if atomic.LoadInt32(&q.closed) == queClosed {
		return ErrClosed
	}

	q.buf <- elem
	return nil
}

func (q *BlockQueue) Pop() (elem interface{}, closed bool) {
	if atomic.LoadInt32(&q.closed) == queClosed {
		closed = true
		if q.Len() > 0 {
			elem = <-q.buf
		}
	} else {
		elem = <-q.buf
	}
	return
}

func (q *BlockQueue) Close() (elems []interface{}) {
	if !atomic.CompareAndSwapInt32(&q.closed, 0, queClosed) {
		return
	}

	elems = []interface{}{}
	for q.Len() > 0 {
		e := <-q.buf
		elems = append(elems, e)
	}

	return
}

func (q *BlockQueue) Closed() bool {
	if atomic.LoadInt32(&q.closed) == queClosed {
		return true
	}
	return false
}

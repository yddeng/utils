package queue

import (
	"fmt"
	"sync"
)

var ErrQueueFull = fmt.Errorf("queue is full")

type ChannelQueue struct {
	queue     chan interface{}
	fullSize  int
	closeOnce sync.Once
	closeCh   chan struct{}
}

func (cq *ChannelQueue) Closed() bool {
	select {
	case <-cq.closeCh:
		return true
	default:
		return false
	}
}

func (cq *ChannelQueue) Close() {
	cq.closeOnce.Do(func() {
		close(cq.closeCh)
		close(cq.queue)
	})
}

// 阻塞投递
func (cq *ChannelQueue) PushB(e interface{}) error {
	select {
	case <-cq.closeCh:
		return ErrClosed
	case cq.queue <- e:
		return nil
	}
}

// 非阻塞投递
func (cq *ChannelQueue) PushN(e interface{}) error {
	select {
	case <-cq.closeCh:
		return ErrClosed
	case cq.queue <- e:
		return nil
	default:
		return ErrQueueFull
	}
}

// 堵塞出
func (cq *ChannelQueue) Pop() (interface{}, bool) {
	elem, open := <-cq.queue
	return elem, open
}

func (cq *ChannelQueue) Len() int {
	return len(cq.queue)
}

func NewChannelQueue(fullSize int) *ChannelQueue {
	return &ChannelQueue{
		queue:    make(chan interface{}, fullSize),
		fullSize: fullSize,
		closeCh:  make(chan struct{}),
	}
}

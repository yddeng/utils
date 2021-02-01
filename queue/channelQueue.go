package queue

import (
	"fmt"
	"sync/atomic"
)

var ErrQueueFull = fmt.Errorf("queue is full")

type ChannelQueue struct {
	channel  chan interface{}
	fullSize int
	opened   int32
}

func (cq *ChannelQueue) open() bool {
	return atomic.LoadInt32(&cq.opened) == 1
}

// 阻塞投递
func (cq *ChannelQueue) PushB(e interface{}) error {
	if !cq.open() {
		return ErrClosed
	}
	cq.channel <- e
	return nil
}

// 非阻塞投递
func (cq *ChannelQueue) PushN(e interface{}) error {
	if !cq.open() {
		return ErrClosed
	}

	if len(cq.channel) == cq.fullSize {
		return ErrQueueFull
	}

	cq.channel <- e
	return nil
}

// 堵塞出
func (cq *ChannelQueue) Pop() (interface{}, bool) {
	elem, open := <-cq.channel
	return elem, open
}

func (cq *ChannelQueue) Close() {
	if atomic.CompareAndSwapInt32(&cq.opened, 1, 0) {
		close(cq.channel)
	}
}

func (cq *ChannelQueue) IsOpen() bool {
	return cq.open()
}

func NewChannelQueue(fullSize int) *ChannelQueue {
	return &ChannelQueue{
		channel:  make(chan interface{}, fullSize),
		fullSize: fullSize,
		opened:   1,
	}
}

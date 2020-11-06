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
func (cq *ChannelQueue) PopB() (interface{}, bool) {
	elem := <-cq.channel
	if cq.open() {
		return elem, true
	} else {
		return nil, false
	}
}

// 非堵塞出，当队列为空时，返回nil
func (cq *ChannelQueue) PopN() (interface{}, bool) {
	if len(cq.channel) > 0 {
		elem := <-cq.channel
		return elem, cq.open()
	} else {
		return nil, cq.open()
	}
}

func (cq *ChannelQueue) Close() {
	if atomic.CompareAndSwapInt32(&cq.opened, 1, 0) {
		close(cq.channel)
	}
}

func NewChannelQueue(fullSize int) *ChannelQueue {
	return &ChannelQueue{
		channel:  make(chan interface{}, fullSize),
		fullSize: fullSize,
		opened:   1,
	}
}

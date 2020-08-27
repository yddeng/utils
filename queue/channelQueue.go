package queue

import (
	"fmt"
	"sync/atomic"
)

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
		return fmt.Errorf("channel queue closed")
	}
	cq.channel <- e
	return nil
}

// 非阻塞投递
func (cq *ChannelQueue) PushN(e interface{}) error {
	if !cq.open() {
		return fmt.Errorf("channel queue closed")
	}

	if len(cq.channel) == cq.fullSize {
		return fmt.Errorf("channel queue fullSize")
	}

	cq.channel <- e
	return nil
}

// 堵塞出
func (cq *ChannelQueue) PopB() (interface{}, bool) {
	elem, b := <-cq.channel
	return elem, b
}

// 非堵塞出，当队列为空时，返回nil
func (cq *ChannelQueue) PopN() (interface{}, bool) {
	if len(cq.channel) > 0 {
		elem, b := <-cq.channel
		return elem, b
	} else {
		return nil, cq.open()
	}
}

func (cq *ChannelQueue) Close() {
	atomic.StoreInt32(&cq.opened, 0)
	close(cq.channel)
}

func NewChannelQueue(fullSize ...int) *ChannelQueue {
	size := 1024
	if len(fullSize) > 0 {
		size = fullSize[0]
	}

	return &ChannelQueue{
		channel:  make(chan interface{}, size),
		fullSize: size,
		opened:   1,
	}
}

package queue

import (
	"fmt"
	"sync"
)

type ChannelQueue struct {
	channel  chan interface{}
	fullSize int
	opened   bool
	mtx      sync.Mutex
}

// 阻塞投递
func (cq *ChannelQueue) PushB(e interface{}) error {
	cq.mtx.Lock()
	if !cq.opened {
		cq.mtx.Unlock()
		return fmt.Errorf("channel queue closed")
	}
	cq.mtx.Unlock()

	cq.channel <- e
	return nil
}

// 非阻塞投递
func (cq *ChannelQueue) PushN(e interface{}) error {
	cq.mtx.Lock()
	if !cq.opened {
		cq.mtx.Unlock()
		return fmt.Errorf("channel queue closed")
	}

	if len(cq.channel) == cq.fullSize {
		cq.mtx.Unlock()
		return fmt.Errorf("channel queue fullSize")
	}
	cq.mtx.Unlock()

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
		cq.mtx.Lock()
		opened := cq.opened
		cq.mtx.Unlock()
		return nil, opened
	}
}

func (cq *ChannelQueue) Close() {
	cq.mtx.Lock()
	defer cq.mtx.Unlock()
	cq.opened = false
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
		mtx:      sync.Mutex{},
		opened:   true,
	}
}

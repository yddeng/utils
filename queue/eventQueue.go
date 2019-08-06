package queue

import "sync/atomic"

type EventQueue struct {
	queue   *BlockQueue
	onEvent func(interface{})
	started int32
}

func NewEventQueue(size int, onEvent func(interface{})) *EventQueue {
	e := &EventQueue{
		queue:   NewBlockQueue(size),
		onEvent: onEvent,
	}

	return e
}

func (e *EventQueue) Push(i interface{}) error {
	return e.queue.Push(i)
}

func (e *EventQueue) Close() {
	_ = e.queue.Close()
}

func (e *EventQueue) Run() {
	if !atomic.CompareAndSwapInt32(&e.started, 0, 1) {
		return
	}

	for {
		ele, closed := e.queue.Pop()
		if closed {
			return
		} else {
			e.onEvent(ele)
		}
	}
}

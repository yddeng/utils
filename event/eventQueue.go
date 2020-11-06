package event

import (
	"github.com/yddeng/dutil/queue"
	"sync/atomic"
)

type EventQueue struct {
	fullSize int
	inQueue  *queue.ChannelQueue
	state    int32
}

func NewEventQueue(fullSize int) *EventQueue {
	e := &EventQueue{
		fullSize: fullSize,
		inQueue:  queue.NewChannelQueue(fullSize),
	}
	return e
}

func (e *EventQueue) Push(fn interface{}, args ...interface{}) error {
	if atomic.LoadInt32(&e.state) != 1 {
		panic("eventQueue is't started")
	}

	event_, err := NewEvent(fn, args...)
	if err != nil {
		return err
	}

	return e.inQueue.PushB(event_)
}

func (e *EventQueue) PushEvent(event_ EventI) error {
	if atomic.LoadInt32(&e.state) != 1 {
		panic("eventQueue is't started")
	}

	return e.inQueue.PushB(event_)
}

func (e *EventQueue) Stop() {
	if atomic.CompareAndSwapInt32(&e.state, 1, 0) {
		e.inQueue.Close()
	}
}

func (e *EventQueue) Run() {
	if !atomic.CompareAndSwapInt32(&e.state, 0, 1) {
		return
	}

	go func() {
		for {
			ele, opened := e.inQueue.PopB()
			if !opened {
				return
			}

			event_ := ele.(EventI)
			event_.Call()
		}
	}()
}

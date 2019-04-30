package dutil

import (
	"sync"
)

const defSize = 1 << 4 //16

type EventQueue struct {
	queue chan interface{}
}

func NewEventQueue() *EventQueue {
	e := &EventQueue{
		queue: make(chan interface{}, defSize),
	}
	e.run()

	return e
}

func (e *EventQueue) Add(cb func()) {
	e.queue <- cb
}

func (e *EventQueue) run() {
	go func() {
		for {
			//队列有数据就唤醒，没有数据就堵塞
			select {
			case v := <-e.queue:
				v.(func())()
			}
		}
	}()

}

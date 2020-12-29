package task

import (
	"github.com/yddeng/dutil/queue"
	"sync/atomic"
)

type TaskQueue struct {
	fullSize int
	inQueue  *queue.ChannelQueue
	state    int32
}

func NewTaskQueue(fullSize int) *TaskQueue {
	e := &TaskQueue{
		fullSize: fullSize,
		inQueue:  queue.NewChannelQueue(fullSize),
	}
	return e
}

func (e *TaskQueue) Push(fn interface{}, args ...interface{}) error {
	if atomic.LoadInt32(&e.state) != 1 {
		panic("eventQueue is't started")
	}
	return e.inQueue.PushB(NewFuncTask(fn, args...))
}

func (e *TaskQueue) PushTask(task Task) error {
	if atomic.LoadInt32(&e.state) != 1 {
		panic("eventQueue is't started")
	}

	return e.inQueue.PushB(task)
}

func (e *TaskQueue) Stop() {
	if atomic.CompareAndSwapInt32(&e.state, 1, 0) {
		e.inQueue.Close()
	}
}

func (e *TaskQueue) Run() {
	if !atomic.CompareAndSwapInt32(&e.state, 0, 1) {
		return
	}

	go func() {
		for {
			ele, opened := e.inQueue.PopB()
			if !opened {
				return
			}
			task := ele.(Task)
			_, _ = task.Do()
		}
	}()
}

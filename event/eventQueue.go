package event

import (
	"github.com/yddeng/dutil/queue"
	"sync/atomic"
)

type EventQueue struct {
	fullSize   int
	inQueue    *queue.ChannelQueue
	routineCnt int32 //执行队列的携程数量，等于0时表示没有启动
}

func NewEventQueue(size int) *EventQueue {
	e := &EventQueue{
		fullSize: size,
		inQueue:  queue.NewChannelQueue(size),
	}
	return e
}

func (e *EventQueue) Push(fn interface{}, args ...interface{}) error {
	event_, err := NewEvent(fn, args...)
	if err != nil {
		return err
	}
	return e.inQueue.PushB(event_)
}

func (e *EventQueue) PushEvent(event_ EventI) error {
	return e.inQueue.PushB(event_)
}

func (e *EventQueue) Stop() {
	if atomic.LoadInt32(&e.routineCnt) == 0 {
		return
	}

	e.inQueue.Close()
	atomic.StoreInt32(&e.routineCnt, 0)
}

//创建一定数目的协程来处理
//count = 1 if routineCnt <= 0
//当线程数大于1时，不能保证完成顺序与投递顺序一致
func (e *EventQueue) Run(routineCnt int) {
	if atomic.LoadInt32(&e.routineCnt) != 0 {
		return
	}

	count := routineCnt
	if routineCnt <= 0 {
		count = 1
	}

	atomic.StoreInt32(&e.routineCnt, int32(count))

	for i := 0; i < count; i++ {
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

}

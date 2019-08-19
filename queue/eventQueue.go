package queue

import (
	"sync/atomic"
)

type EventQueue struct {
	queue      *BlockQueue
	onEvent    func(interface{})
	routineCnt int32 //执行队列的携程数量，等于0时表示没有启动
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

func (e *EventQueue) Stop() {
	if atomic.LoadInt32(&e.routineCnt) == 0 {
		return
	}

	e.queue.Close()
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
				ele, closed := e.queue.Pop()
				if closed && ele == nil {
					return
				} else {
					e.onEvent(ele)
				}
			}
		}()
	}

}

/*
 如何保证同一个客户端的消息的顺序处理？
 把不同的client hash到不同的逻辑线程上，今后来自该客户端的所有的消息都由该逻辑线程来处理，这就保证了顺序
*/

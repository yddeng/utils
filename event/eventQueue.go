package event

import (
	"github.com/yddeng/dutil/queue"
	"sync/atomic"
)

type EventQueue struct {
	fullSize   int
	inQueue    *queue.ChannelQueue
	outQueue   *queue.BlockQueue
	routineCnt int32         //执行队列的携程数量，等于0时表示没有启动
	signal     chan struct{} // 通知事件完成
}

func NewEventQueue(fullSize int) *EventQueue {
	e := &EventQueue{
		fullSize: fullSize,
		inQueue:  queue.NewChannelQueue(fullSize),
	}
	return e
}

type eventEntity struct {
	event_     EventI
	isComplete int32 // 0 for false,1 for true.
}

func (e *EventQueue) sendSignal() {
	select {
	case <-e.signal:
	default:
	}
	e.signal <- struct{}{}
}

func (e *EventQueue) Push(fn interface{}, args ...interface{}) error {
	cnt := atomic.LoadInt32(&e.routineCnt)
	if cnt == 0 {
		panic("eventQueue is't started")
	}

	event_, err := NewEvent(fn, args...)
	if err != nil {
		return err
	}

	if cnt == 1 {
		return e.inQueue.PushB(event_)
	} else {
		entity := &eventEntity{
			event_:     event_,
			isComplete: 0,
		}
		if err := e.outQueue.Push(entity); err != nil {
			return err
		}
		return e.inQueue.PushB(entity)
	}
}

func (e *EventQueue) PushEvent(event_ EventI) error {
	cnt := atomic.LoadInt32(&e.routineCnt)
	if cnt == 0 {
		panic("eventQueue is't started")
	}

	if cnt == 1 {
		return e.inQueue.PushB(event_)
	} else {
		entity := &eventEntity{
			event_:     event_,
			isComplete: 0,
		}
		if err := e.outQueue.Push(entity); err != nil {
			return err
		}
		return e.inQueue.PushB(entity)
	}
}

func (e *EventQueue) Stop() {
	if atomic.LoadInt32(&e.routineCnt) == 0 {
		return
	}
	e.inQueue.Close()
	atomic.StoreInt32(&e.routineCnt, 0)
}

/* 创建一定数目的协程来处理
 * 按照进入顺序，执行调用
 * 多携程暂时无意义
 */
func (e *EventQueue) Run(routineCnt int) {
	if atomic.LoadInt32(&e.routineCnt) != 0 {
		return
	}

	count := routineCnt
	if routineCnt <= 0 {
		count = 1
	}

	atomic.StoreInt32(&e.routineCnt, int32(count))

	if count == 1 {
		e.runOne()
	} else {
		e.signal = make(chan struct{}, 1)
		e.outQueue = queue.NewBlockQueue(e.fullSize)
		e.runMultiple(count)
	}

}

func (e *EventQueue) runOne() {
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

func (e *EventQueue) runMultiple(count int) {
	for i := 0; i < count; i++ {
		go func() {
			for {
				ele, opened := e.inQueue.PopB()
				if !opened {
					return
				}

				entity := ele.(*eventEntity)
				atomic.StoreInt32(&entity.isComplete, 1)
				e.sendSignal()
			}
		}()
	}

	go func() {
		breakSignal := true
		for {
			select {
			case <-e.signal:
			}

			breakSignal = true
			for {
				ele := e.outQueue.Peek()
				if entity, ok := ele.(*eventEntity); ok {
					if atomic.LoadInt32(&entity.isComplete) == 1 {
						entity.event_.Call()
						e.outQueue.Pop()
						breakSignal = false
					}
				}
				if breakSignal {
					break
				}
			}
		}
	}()
}

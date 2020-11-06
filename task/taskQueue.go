package task

import (
	"fmt"
	"github.com/yddeng/dutil/queue"
	"sync/atomic"
)

type TaskQueue struct {
	fullSize   int
	inQueue    *queue.ChannelQueue
	outQueue   *queue.BlockQueue
	state      int32
	routineCnt int32         //执行队列的携程数量
	signal     chan struct{} // 通知事件完成
}

func NewTaskQueue(fullSize int) *TaskQueue {
	e := &TaskQueue{
		fullSize: fullSize,
		inQueue:  queue.NewChannelQueue(fullSize),
	}
	return e
}

type taskEntity struct {
	task       Task
	isComplete int32 // 0 for false,1 for true.
}

func (e *TaskQueue) sendSignal() {
	select {
	case <-e.signal:
	default:
	}
	e.signal <- struct{}{}
}

func (e *TaskQueue) AddTask(task Task) error {
	cnt := atomic.LoadInt32(&e.routineCnt)
	if cnt == 0 {
		panic("eventQueue is't started")
	}

	if cnt == 1 {
		return e.inQueue.PushB(task)
	} else {
		entity := &taskEntity{
			task:       task,
			isComplete: 0,
		}
		if err := e.outQueue.Push(entity); err != nil {
			return err
		}
		return e.inQueue.PushB(entity)
	}
}

/*
 停止时，已经执行的继续执行，保证单个任务整个流程执行完。
	未开始执行的不在继续执行。
*/
func (e *TaskQueue) Stop() {
	if atomic.CompareAndSwapInt32(&e.state, 1, 0) {
		e.inQueue.Close()
	}
}

/* 创建一定数目的协程来处理
 * begin多线程运行。执行end调用按照进入顺序。
 */
func (e *TaskQueue) Run(routineCnt int) {
	if !atomic.CompareAndSwapInt32(&e.state, 0, 1) {
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

func (e *TaskQueue) runOne() {
	go func() {
		for {
			ele, opened := e.inQueue.PopB()
			if !opened {
				return
			}

			task := ele.(Task)
			task.Begin()
			task.End()
		}
	}()
}

func (e *TaskQueue) runMultiple(count int) {

	for i := 0; i < count; i++ {
		go func() {
			defer atomic.AddInt32(&e.routineCnt, -1)
			for {
				ele, opened := e.inQueue.PopB()
				fmt.Println(ele, opened)
				if !opened {
					return
				}

				entity := ele.(*taskEntity)
				entity.task.Begin()
				atomic.StoreInt32(&entity.isComplete, 1)
				e.sendSignal()
			}
		}()
	}

	go func() {
		breakSignal := true
		for {
			<-e.signal
			for {
				breakSignal = true
				ele := e.outQueue.Peek()
				if entity, ok := ele.(*taskEntity); ok {
					if atomic.LoadInt32(&entity.isComplete) == 1 {
						entity.task.End()
						e.outQueue.Pop()
						breakSignal = false
					}
				}
				if breakSignal {
					break
				}
			}

			if atomic.LoadInt32(&e.routineCnt) == 0 {
				return
			}
		}
	}()
}

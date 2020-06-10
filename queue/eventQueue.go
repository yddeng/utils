package queue

import (
	"fmt"
	"reflect"
	"sync/atomic"
)

type event struct {
	args []interface{}
	fn   interface{}
}

type EventQueue struct {
	queue      *ChannelQueue
	routineCnt int32 //执行队列的携程数量，等于0时表示没有启动
}

func NewEventQueue(size int) *EventQueue {
	e := &EventQueue{
		queue: NewChannelQueue(size),
	}

	return e
}

func preparePost(fn interface{}, args ...interface{}) (*event, error) {
	e := &event{fn: fn}
	switch fn.(type) {
	case func():
	case func([]interface{}), func(...interface{}):
		e.args = args
	default:
		return nil, fmt.Errorf("invaild callback type %s", reflect.TypeOf(fn).String())
	}
	return e, nil
}

func (e *EventQueue) Push(fn interface{}, args ...interface{}) error {
	event, err := preparePost(fn, args...)
	if err != nil {
		return err
	}
	return e.queue.PushB(event)
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
				ele, opened := e.queue.PopB()
				if !opened {
					return
				}

				event_ := ele.(*event)
				pcall1(event_.fn, event_.args)
			}
		}()
	}

}

func pcall1(fn interface{}, args []interface{}) {
	switch fn.(type) {
	case func():
		fn.(func())()
	case func([]interface{}):
		fn.(func([]interface{}))(args)
	case func(...interface{}):
		fn.(func(...interface{}))(args...)
	default:
	}
}

/*
 如何保证同一个客户端的消息的顺序处理？
 把不同的client hash到不同的逻辑线程上，今后来自该客户端的所有的消息都由该逻辑线程来处理，这就保证了顺序
*/

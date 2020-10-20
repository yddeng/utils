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
	fullSize   int
	inQueue    *ChannelQueue
	routineCnt int32 //执行队列的携程数量，等于0时表示没有启动
}

func NewEventQueue(size int) *EventQueue {
	e := &EventQueue{
		fullSize: size,
		inQueue:  NewChannelQueue(size),
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
	return e.inQueue.PushB(event)
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

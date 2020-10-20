package pool

import (
	"fmt"
	"github.com/yddeng/dutil/queue"
	"reflect"
	"sync"
)

/*
   线程池。
*/

type ThreadPool struct {
	maxCount     int32
	currentCount int32
	taskQueue    *queue.ChannelQueue
	taskCount    int
	mtx          sync.Mutex
}

func NewTreadPool(threadMaxCount, channelSize int) *ThreadPool {
	return &ThreadPool{
		currentCount: 0,
		maxCount:     int32(threadMaxCount),
		taskQueue:    queue.NewChannelQueue(channelSize),
	}
}

func (p *ThreadPool) AddTask(fn func(), args ...interface{}) error {
	event_, err := preparePost(fn, args...)
	if err != nil {
		return err
	}
	_ = p.taskQueue.PushB(event_)
	p.mtx.Lock()
	defer p.mtx.Unlock()
	p.taskCount++
	if p.currentCount < p.maxCount {
		p.currentCount++
		go p.newTread()
	}
	return nil
}

func (p *ThreadPool) newTread() {
	for {
		p.mtx.Lock()
		if p.taskCount == 0 {
			p.mtx.Unlock()
			break
		}
		p.mtx.Unlock()

		ele, opened := p.taskQueue.PopB()
		if !opened {
			break
		}

		p.mtx.Lock()
		p.taskCount--
		p.mtx.Unlock()

		event_ := ele.(*event)
		pcall(event_.fn, event_.args)

	}
	p.mtx.Lock()
	p.taskCount--
	p.mtx.Unlock()
}

type event struct {
	args []interface{}
	fn   interface{}
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

func pcall(fn interface{}, args []interface{}) {
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

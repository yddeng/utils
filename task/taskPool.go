package task

import (
	"errors"
	"github.com/sniperHW/kendynet/util"
	"sync"
	"sync/atomic"
)

const minChannelSize = 1024

type task struct {
	fn   interface{}
	args []interface{}
}

type TaskPool struct {
	maxTreadCnt int32
	crtTreadCnt int32
	taskChan    chan *task
	maxTaskCnt  int32
	crtTaskCnt  int32
	stopped     int32
	mtx         sync.Mutex
}

func NewTaskPool(threadMaxCount, channelSize int) *TaskPool {
	if channelSize < minChannelSize {
		channelSize = minChannelSize
	}
	return &TaskPool{
		crtTreadCnt: 0,
		maxTreadCnt: int32(threadMaxCount),
		taskChan:    make(chan *task, channelSize),
		maxTaskCnt:  int32(channelSize),
	}
}

func (p *TaskPool) AddTask(fn interface{}, args ...interface{}) error {
	if atomic.LoadInt32(&p.stopped) == 1 {
		return errors.New("taskPool : AddTask failed, pool is stopped")
	}

	p.mtx.Lock()
	defer p.mtx.Unlock()

	if len(p.taskChan) == int(p.maxTaskCnt) {
		return errors.New("taskPool : AddTask failed, task is full")
	}

	p.taskChan <- &task{fn: fn, args: args}

	p.crtTaskCnt++
	if p.crtTreadCnt < p.maxTreadCnt {
		p.crtTreadCnt++
		go p.newTread()
	}
	return nil
}

func (p *TaskPool) newTread() {
	for {
		p.mtx.Lock()
		if p.crtTaskCnt == 0 {
			p.mtx.Unlock()
			break
		}
		p.mtx.Unlock()

		task, opened := <-p.taskChan
		if !opened {
			break
		}

		p.mtx.Lock()
		p.crtTaskCnt--
		p.mtx.Unlock()

		_, _ = util.ProtectCall(task.fn, task.args)

	}
	p.mtx.Lock()
	p.crtTreadCnt--
	p.mtx.Unlock()
}

func (p *TaskPool) Stop() {
	if atomic.CompareAndSwapInt32(&p.stopped, 0, 1) {
		close(p.taskChan)
	}
}

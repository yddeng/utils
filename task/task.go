package task

import (
	"github.com/yddeng/utils"
	"runtime"
	"sync"
	"time"
)

type Task interface {
	Do() (result []interface{}, err error)
}

type funcTask struct {
	fn   interface{}
	args []interface{}
}

func FuncTask(f interface{}, args ...interface{}) *funcTask {
	return &funcTask{
		fn:   f,
		args: args,
	}
}

func (this *funcTask) Do() (result []interface{}, err error) {
	return utils.CallFunc(this.fn, this.args...)
}

type taskMgr interface {
	dataCh() chan Task
	idleWorker(*taskWorker)
}

type taskWorker struct {
	mgr  taskMgr
	idle time.Duration
}

func (this *taskWorker) run() {
	go func() {
		defer this.mgr.idleWorker(this)
		timer := time.NewTimer(this.idle)
		for {
			select {
			case <-timer.C:
				return
			case task := <-this.mgr.dataCh():
				timer.Reset(this.idle)
				task.Do()
			}
		}
	}()
}

var (
	defaultTaskPool *TaskPool
	createOnce      sync.Once
)

func Go(fn interface{}, args ...interface{}) error {
	createOnce.Do(func() {
		defaultTaskPool = NewTaskPool(runtime.NumCPU(), defaultTaskSize, defaultIdleTime)
	})
	return defaultTaskPool.Submit(fn, args...)
}

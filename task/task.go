package task

import (
	"github.com/yddeng/utils"
	"runtime"
	"sync"
)

type Task interface {
	Do()
}

type funcTask struct {
	fn   interface{}
	args []interface{}
}

func (this *funcTask) Do() {
	_, _ = utils.CallFunc(this.fn, this.args...)
}

type taskMgr interface {
	freeWorker(*taskWorker)
}

type taskWorker struct {
	mgr taskMgr
}

func (this *taskWorker) run(taskC chan Task) {
	defer this.mgr.freeWorker(this)
	for {
		select {
		case task := <-taskC:
			task.Do()
		default:
			return
		}
	}
}

var (
	defaultTaskPool *TaskPool
	createOnce      sync.Once
)

func Submit(fn interface{}, args ...interface{}) error {
	createOnce.Do(func() {
		defaultTaskPool = NewTaskPool(runtime.NumCPU()*2, defaultTaskSize)
	})
	return defaultTaskPool.Submit(fn, args...)
}

func SubmitTask(task Task) error {
	createOnce.Do(func() {
		defaultTaskPool = NewTaskPool(runtime.NumCPU()*2, defaultTaskSize)
	})
	return defaultTaskPool.SubmitTask(task)
}

package task

import (
	"errors"
	"sync"
)

const (
	defaultTaskSize = 1024
)

type TaskPool struct {
	workers    int
	workerSize int
	workerLock sync.Mutex
	workerPool sync.Pool

	taskChan chan Task

	die     chan struct{}
	dieOnce sync.Once
}

func (this *TaskPool) freeWorker(worker *taskWorker) {
	this.workerLock.Lock()
	this.workers--
	this.workerPool.Put(worker)
	this.workerLock.Unlock()
}

func (this *TaskPool) Running() int {
	this.workerLock.Lock()
	defer this.workerLock.Unlock()
	return this.workers
}

func (this *TaskPool) submit(task Task) error {
	select {
	case <-this.die:
		return errors.New("taskPool:Submit pool is stopped")
	default:
	}

	taskChan := this.taskChan
	if this.workerSize == 0 {
		taskChan = make(chan Task, 1)
	}

	select {
	case taskChan <- task:
	default:
		return errors.New("taskPool:Submit task channel is full")
	}

	this.workerLock.Lock()
	defer this.workerLock.Unlock()

	if this.workerSize == 0 || this.workers < this.workerSize {
		this.workers++
		w := this.workerPool.Get().(*taskWorker)
		go w.run(taskChan)
	}
	return nil
}

func (this *TaskPool) Submit(fn interface{}, args ...interface{}) error {
	return this.submit(&funcTask{fn: fn, args: args})
}

func (this *TaskPool) SubmitTask(task Task) error {
	return this.submit(task)
}

func (this *TaskPool) Stop() {
	this.dieOnce.Do(func() {
		close(this.die)
	})
}

// NewTaskPool
// workerSize > 0 , 限制goroutine的数量; workerSize = 0 , 不限制
func NewTaskPool(workerSize, taskSize int) *TaskPool {
	if taskSize < defaultTaskSize {
		taskSize = defaultTaskSize
	}
	if workerSize < 0 {
		workerSize = 0
	}

	pool := new(TaskPool)
	pool.die = make(chan struct{})
	pool.workerSize = workerSize
	pool.taskChan = make(chan Task, taskSize)
	pool.workerPool.New = func() interface{} {
		return &taskWorker{mgr: pool}
	}

	return pool
}

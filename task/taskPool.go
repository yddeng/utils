package task

import (
	"errors"
	"sync"
	"time"
)

const (
	defaultTaskSize = 1024
	defaultIdleTime = time.Second * 10
)

type TaskPool struct {
	workers    []*taskWorker
	workerSize int
	lock       sync.Mutex

	taskChan   chan Task
	workerPool sync.Pool

	die     chan struct{}
	dieOnce sync.Once
}

func (this *TaskPool) dataCh() chan Task {
	return this.taskChan
}

func (this *TaskPool) idleWorker(worker *taskWorker) {
	this.lock.Lock()
	for i, w := range this.workers {
		if w == worker {
			this.workers = append(this.workers[:i], this.workers[i+1:]...)
			break
		}
	}
	this.workerPool.Put(worker)
	this.lock.Unlock()
}

func (this *TaskPool) Running() int {
	this.lock.Lock()
	defer this.lock.Unlock()
	return len(this.workers)
}

func (this *TaskPool) Submit(fn interface{}, args ...interface{}) error {
	select {
	case <-this.die:
		return errors.New("taskPool:Submit pool is stopped")
	default:
	}

	task := FuncTask(fn, args...)

	select {
	case this.taskChan <- task:
	default:
		return errors.New("taskPool:Submit task channel is full")
	}

	this.lock.Lock()
	defer this.lock.Unlock()

	if len(this.workers) < this.workerSize {
		w := this.workerPool.Get().(*taskWorker)
		this.workers = append(this.workers, w)
		w.run()
	}
	return nil
}

func (this *TaskPool) Stop() {
	this.dieOnce.Do(func() {
		close(this.die)
	})
}

func NewTaskPool(workerSize, taskSize int, idleTimes ...time.Duration) *TaskPool {
	pool := new(TaskPool)
	pool.die = make(chan struct{})

	var workerIdle time.Duration
	if len(idleTimes) == 0 || idleTimes[0] < time.Millisecond {
		workerIdle = defaultIdleTime
	} else {
		workerIdle = idleTimes[0]
	}

	if taskSize < defaultTaskSize {
		pool.taskChan = make(chan Task, defaultTaskSize)
	} else {
		pool.taskChan = make(chan Task, taskSize)
	}

	if workerSize <= 0 {
		workerSize = 1
	}
	pool.workerSize = workerSize
	pool.workers = make([]*taskWorker, 0, workerSize)

	pool.workerPool.New = func() interface{} {
		return &taskWorker{
			mgr:  pool,
			idle: workerIdle,
		}
	}

	return pool
}

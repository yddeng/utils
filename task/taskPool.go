package task

import (
	"github.com/yddeng/dutil/queue"
	"sync"
)

/*
   任务线程池。多个线程并列消费
*/

type TaskPool struct {
	maxCount     int32
	currentCount int32
	taskQueue    *queue.ChannelQueue
	taskCount    int
	mtx          sync.Mutex
}

func NewTaskPool(threadMaxCount, channelSize int) *TaskPool {
	return &TaskPool{
		currentCount: 0,
		maxCount:     int32(threadMaxCount),
		taskQueue:    queue.NewChannelQueue(channelSize),
	}
}

func (p *TaskPool) AddTask(fn interface{}, args ...interface{}) error {
	_ = p.taskQueue.PushB(NewFuncTask(fn, args...))

	p.mtx.Lock()
	defer p.mtx.Unlock()
	p.taskCount++
	if p.currentCount < p.maxCount {
		p.currentCount++
		go p.newTread()
	}
	return nil
}

func (p *TaskPool) newTread() {
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

		task := ele.(Task)
		_, _ = task.Do()

	}
	p.mtx.Lock()
	p.currentCount--
	p.mtx.Unlock()
}

package event

import (
	"github.com/yddeng/dutil/queue"
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

func (p *ThreadPool) AddTask(fn interface{}, args ...interface{}) error {
	event_, err := NewEvent(fn, args...)
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

		event_ := ele.(*Event)
		event_.Call()

	}
	p.mtx.Lock()
	p.currentCount--
	p.mtx.Unlock()
}

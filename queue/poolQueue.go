package queue

import (
	"sync"
)

/*
 提供一个队列容器，里面有多个队列，每个队列都可固定存放一定数量的消息。
 网络IO线程要给逻辑线程投递消息时，会从队列容器中取一个空队列来使用, 直到将该队列填满后再放回容器中换另一个空队列。
 而逻辑线程取消息时是从队列容器中取一个有消息的队列来读取，处理完后清空队列再放回到容器中。
 这样便使得只有在对队列容器进行操作时才需要加锁，而IO线程和逻辑线程在操作自己当前使用的队列时都不需要加锁，所以锁竞争的机会大大减少了。
 这样有时也会出现IO线程未写满一个队列，而逻辑线程又没有数据可处理的情况，特别是当数据量很少时可能会很容易出现。
 这个可以通过设置超时来处理, 如果当前时间-向队列放入第一个包的时间> 50 ms, 就将其放回到容器中换另一个队列。
*/

type PoolQueue struct {
	input    *BlockQueue
	out      *BlockQueue
	queues   []*BlockQueue
	mutex    sync.Mutex
	queueCnt int
	queueCap int
}

func NewPoolQueue(queueCnt, queueCap int) *PoolQueue {
	pq := &PoolQueue{
		queues: []*BlockQueue{},
	}

	for i := 0; i < queueCnt; i++ {
		pq.queues = append(pq.queues, NewBlockQueue(queueCap))
	}

	return pq
}

func (pq *PoolQueue) Input(i interface{}) error {

	if pq.input.Len() >= pq.queueCap-1 {
		pq.inputQueue()
	}

	return pq.input.Push(i)
}

func (pq *PoolQueue) inputQueue() {
	pq.mutex.Lock()
	defer pq.mutex.Unlock()

	for _, bq := range pq.queues {
		if bq.Len() == 0 {
			pq.input = bq
			return
		}
	}

	pq.queueCnt += 1
	newQue := NewBlockQueue(pq.queueCap)
	pq.queues = append(pq.queues, newQue)

	pq.input = newQue

}

func (pq *PoolQueue) Get() interface{} {

	if pq.out.Len() == 0 {
		pq.outQueue()
	}

	if pq.out.Len() == 0 {
		return nil
	} else {
		ret, _ := pq.out.Pop()
		return ret
	}
}

func (pq *PoolQueue) outQueue() {
	pq.mutex.Lock()
	defer pq.mutex.Unlock()

	for _, bq := range pq.queues {
		if bq.Len() > 0 {
			pq.out = bq
			return
		}
	}

	for _, bq := range pq.queues {
		if bq.Len() == 0 {
			pq.out = bq
			return
		}
	}

}

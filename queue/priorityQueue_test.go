package queue

import (
	"sync"
	"testing"
	"time"
)

func TestNewPriorityQueue(t *testing.T) {
	p := NewPriorityQueue(10)

	p.Push(1, 1)
	p.Push(2, 3)
	p.Push(3, 2)
	p.Push(4, 1)
	p.Push(5, 1)
	p.Push(6, 2)
	p.Push(7, 0)
	p.Push(8, 0)

	for i := 0; i < 8; i++ {
		v, opend := p.Pop()
		t.Log(v, opend)
	}

	// 2 3 6 1 4 5 7 8
}

func TestPriorityQueue_Close(t *testing.T) {
	pq := NewPriorityQueue(10)
	wg := sync.WaitGroup{}
	for i := 0; i < 1; i++ {
		wg.Add(1)
		go func(id int) {
			t.Logf("id(%d) pop start", id)
			defer wg.Done()
			defer t.Logf("id(%d) pop end", id)
			for {
				v, b := pq.Pop()
				t.Logf("id(%d) %v %v", id, v, b)
				if !b {
					return
				}
			}
		}(i)
	}

	for i := 0; i < 20; i++ {
		pq.Push(i)
	}
	time.Sleep(time.Second)
	pq.Close()
	wg.Wait()

}

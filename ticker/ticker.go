package ticker

/*
 高精度 ticker
*/

import (
	"github.com/yddeng/dutil/heap"
	"sync"
	"sync/atomic"
	"time"
)

type Ticker struct {
	nextDoTime time.Time
	duration   time.Duration
	doFunc     func(time.Time)
	repeat     bool
	stopped    int32
	mgr        *tickerMgr
}

// 小根堆
func (t *Ticker) Less(elem heap.Element) bool {
	return t.nextDoTime.Before(elem.(*Ticker).nextDoTime)
}

func (t *Ticker) Stop() {
	atomic.StoreInt32(&t.stopped, 1)
}

func (t *Ticker) do(now time.Time) {
	if atomic.LoadInt32(&t.stopped) == 1 {
		return
	}

	t.doFunc(now)
	//repeat
	if t.repeat {
		if atomic.LoadInt32(&t.stopped) == 1 {
			return
		}
		t.mgr.setTicker(t)
	}

}

var (
	runOnce = sync.Once{}
	tMgr    = newMgr()
)

type tickerMgr struct {
	minHeap    *heap.Heap
	notionChan chan struct{}
	sync.Mutex
}

func newMgr() *tickerMgr {
	return &tickerMgr{
		minHeap:    heap.NewHeap(),
		notionChan: make(chan struct{}, 1),
	}
}

func (mgr *tickerMgr) setTicker(t *Ticker) {
	t.nextDoTime = time.Now().Add(t.duration)
	needNotion := false

	mgr.Lock()
	if mgr.minHeap.Len() == 0 { //min.(*Ticker).nextDoTime.Before(t.nextDoTime) {
		needNotion = true
	}
	mgr.minHeap.Push(t)
	mgr.Unlock()

	if needNotion {
		mgr.notionChan <- struct{}{}
	}
}

func (mgr *tickerMgr) run() {

	defaultSleepTime := 10 * time.Second
	var tt *time.Timer
	var min heap.Element
	for {
		now := time.Now()
		for {
			mgr.Lock()
			min = mgr.minHeap.Peek()
			if nil != min && now.After(min.(*Ticker).nextDoTime) {
				t := min.(*Ticker)
				mgr.minHeap.Pop()
				mgr.Unlock()
				t.do(now)
			} else {
				mgr.Unlock()
				break
			}
		}

		sleepTime := defaultSleepTime
		if nil != min {
			sleepTime = min.(*Ticker).nextDoTime.Sub(now)
		}
		if nil != tt {
			tt.Reset(sleepTime)
		} else {
			tt = time.AfterFunc(sleepTime, func() {
				mgr.notionChan <- struct{}{}
			})
		}

		<-mgr.notionChan
		tt.Stop()
	}
}

func newTicker(d time.Duration, cb func(time.Time), repeat bool) *Ticker {
	runOnce.Do(func() {
		go tMgr.run()
	})

	if d <= 0 {
		panic("newTicker error: duration value is failed")
	}

	t := &Ticker{
		duration: d,
		doFunc:   cb,
		repeat:   repeat,
		mgr:      tMgr,
	}

	tMgr.setTicker(t)
	return t
}

func RegisterOnceTicker(duration time.Duration, callback func(time.Time)) *Ticker {
	return newTicker(duration, callback, false)
}

func RegisterRepeatTicker(duration time.Duration, callback func(time.Time)) *Ticker {
	return newTicker(duration, callback, true)
}

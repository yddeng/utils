package ticker

import (
	"sync"
	"time"
)

//低精度ticker
type Ticker struct {
	deadline time.Time
	duration time.Duration
	doFunc   func(time.Time)
	repeat   bool
}

func (t *Ticker) Stop() {
	tMgr.mapLock.Lock()
	defer tMgr.mapLock.Unlock()

	if _, ok := tMgr.tickers[t]; ok {
		delete(tMgr.tickers, t)
	}
}

func (t *Ticker) do(now time.Time) {
	if now.After(t.deadline) {
		t.doFunc(now)
		//repeat
		if t.repeat {
			t.deadline = now.Add(t.duration)
		} else {
			delete(tMgr.tickers, t)
		}
	}
}

var (
	defDuration = 50 * time.Millisecond
	runOnce     = sync.Once{}
	tMgr        = newMgr()
)

type tickerMgr struct {
	ticker  *time.Ticker
	tickers map[*Ticker]struct{}
	mapLock sync.Mutex
}

func newMgr() *tickerMgr {
	return &tickerMgr{
		tickers: map[*Ticker]struct{}{},
	}
}

func (mgr *tickerMgr) run() {
	mgr.ticker = time.NewTicker(defDuration)
	for {
		tt := <-mgr.ticker.C
		mgr.mapLock.Lock()
		for t := range mgr.tickers {
			t.do(tt)
		}
		mgr.mapLock.Unlock()
	}
}

func newTicker(d time.Duration, cb func(time.Time), repeat bool) *Ticker {
	defer runOnce.Do(func() {
		go tMgr.run()
	})

	if d <= 0 {
		panic("newTicker error: interval value is err")
	}

	t := &Ticker{
		deadline: time.Now().Add(d),
		duration: d,
		doFunc:   cb,
		repeat:   repeat,
	}

	tMgr.mapLock.Lock()
	tMgr.tickers[t] = struct{}{}
	tMgr.mapLock.Unlock()

	return t
}

func RegisterOnceTicker(duration time.Duration, callback func(time.Time)) *Ticker {
	return newTicker(duration, callback, false)
}

func RegisterRepeatTicker(duration time.Duration, callback func(time.Time)) *Ticker {
	return newTicker(duration, callback, true)
}

package timer

import (
	"sync"
	"sync/atomic"
	"time"
)

/*
 低精度 timer, 最低精度为毫秒
 如果时间轮精度为10ms， 那么他的误差在 （0，10）ms之间。如果一个任务延迟 500ms，那它的执行时间在490～500ms之间。
 按平均来讲，出错的概率均等的情况下，那么这个出错可能会延迟或提前最小刻度的一半，在这里就是10ms/2=5ms.
 故，时间轮的 tick 单位 在总延迟时间上，应该不足以影响延迟执行函数处理的事务。
*/

const (
	defaultBucketNum = 10
	defaultInterval  = time.Millisecond
)

type WheelTimer struct {
	key     int64
	circle  int // 需要转动多少圈
	rt      *runtimeTimer
	mgr     *TimeWheelMgr
	stopped int32
}

func newWheelTimer(d time.Duration, repeated bool, ctx interface{}, f func(ctx interface{})) *WheelTimer {
	return &WheelTimer{
		key: 0,
		rt: &runtimeTimer{
			when:     when(d),
			ctx:      ctx,
			fn:       f,
			repeated: repeated,
			period:   int64(d),
		},
		mgr:     nil,
		stopped: 0,
	}
}

func (t *WheelTimer) Reset(d time.Duration) bool {
	if t.mgr == nil {
		panic("timer: Reset called on uninitialized WheelTimer")
	}

	if atomic.LoadInt32(&t.stopped) == 1 {
		return false
	}

	t.mgr.removeTimer(t)
	t.rt.when = when(d)
	t.rt.period = int64(d)
	t.mgr.addTimer(t)
	return true
}

func (t *WheelTimer) Stop() bool {
	if t.mgr == nil {
		panic("timer: Stop called on uninitialized WheelTimer")
	}
	atomic.StoreInt32(&t.stopped, 1)
	t.mgr.removeTimer(t)
	return true
}

func (t *WheelTimer) do() {
	if atomic.LoadInt32(&t.stopped) == 1 {
		return
	}

	goFunc(t.rt.fn, t.rt.ctx)
	//repeat
	if t.rt.repeated {
		if atomic.LoadInt32(&t.stopped) == 1 {
			return
		}
		t.rt.when = when(time.Duration(t.rt.period))
		t.mgr.addTimer(t)
	}
}

type TimeWheelMgr struct {
	interval     time.Duration
	bucketNum    int
	buckets      []*bucket
	timer2bucket map[int64]int // 定时器所在的槽, 主要用于删除定时器
	accumulator  int64         // 计数器
	currentIdx   int           // 当前游标
	mtx          sync.Mutex
}

type bucket struct {
	timers map[int64]*WheelTimer
}

func NewTimeWheelMgr(interval time.Duration, bucketNum int) *TimeWheelMgr {
	if bucketNum < defaultBucketNum {
		bucketNum = defaultBucketNum
	}

	if interval < defaultInterval {
		interval = defaultInterval
	}

	mgr := &TimeWheelMgr{
		interval:     interval,
		bucketNum:    bucketNum,
		buckets:      make([]*bucket, bucketNum),
		timer2bucket: map[int64]int{},
		accumulator:  0,
	}

	for i := 0; i < bucketNum; i++ {
		mgr.buckets[i] = &bucket{timers: map[int64]*WheelTimer{}}
	}

	go mgr.run()
	return mgr
}

func (mgr *TimeWheelMgr) addTimer(t *WheelTimer) {
	if t.key == 0 {
		key := atomic.AddInt64(&mgr.accumulator, 1)
		t.key = key
		t.mgr = mgr
	}

	// 如果过期时间小于单位的一半直接完成
	delay := t.rt.when - time.Now().UnixNano()
	if delay < int64(mgr.interval) {
		t.do()
	} else {
		mgr.mtx.Lock()
		defer mgr.mtx.Unlock()

		circle := int(delay / int64(mgr.interval) / int64(mgr.bucketNum))
		idx := (mgr.currentIdx + int(delay/int64(mgr.interval))) % mgr.bucketNum
		t.circle = circle

		mgr.timer2bucket[t.key] = idx
		mgr.buckets[idx].timers[t.key] = t
	}
}

func (mgr *TimeWheelMgr) removeTimer(t *WheelTimer) {
	mgr.mtx.Lock()
	defer mgr.mtx.Unlock()
	if idx, ok := mgr.timer2bucket[t.key]; ok {
		delete(mgr.timer2bucket, t.key)
		b := mgr.buckets[idx]
		delete(b.timers, t.key)
	}
}

func (mgr *TimeWheelMgr) run() {
	ticker := time.NewTicker(mgr.interval)
	for {
		<-ticker.C

		mgr.mtx.Lock()
		b := mgr.buckets[mgr.currentIdx]
		mgr.currentIdx = (mgr.currentIdx + 1) % mgr.bucketNum
		mgr.mtx.Unlock()
		for _, t := range b.timers {
			if t.circle > 0 {
				t.circle--
			} else {
				mgr.removeTimer(t)
				t.do()
			}
		}
	}
}

func (mgr *TimeWheelMgr) OnceTimer(d time.Duration, ctx interface{}, f func(ctx interface{})) Timer {
	t := newWheelTimer(d, false, ctx, f)
	mgr.addTimer(t)
	return t
}

func (mgr *TimeWheelMgr) RepeatTimer(d time.Duration, ctx interface{}, f func(ctx interface{})) Timer {
	t := newWheelTimer(d, true, ctx, f)
	mgr.addTimer(t)
	return t
}

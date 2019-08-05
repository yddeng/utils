package timingwheel

import (
	"fmt"
	"sync"
	"time"
)

/*
 分级时间轮
 延迟时间溢出最小时间轮的总量，自动创建更大的时间轮
*/

const (
	//默认队列缓冲大小
	DoCHANMAX = 1024
)

type TimingWheel struct {
	lock      sync.Mutex
	tickDur   time.Duration //tick时间间隔
	wheelSize int           //时间轮存储数量

	ticker      *time.Ticker
	crtWheel    int                   //当前时间轮索引
	wheelTimers []map[TimerI]struct{} //时间轮片

	levUpWheel   *TimingWheel //上一级轮子
	levDownWheel *TimingWheel //下一级轮子

	doChan chan TimerI //超时事务channer
}

/*
 时间间隔，时间轮长度
 当前时间轮的总量：间隔*长度
*/
func NewTimingWheel(tickDuration time.Duration, wheelSize int) (*TimingWheel, error) {
	if wheelSize <= 0 {
		return nil, fmt.Errorf("wheelSize value:%d error", wheelSize)
	}
	return newTimingWheel(tickDuration, wheelSize, nil), nil
}

func newTimingWheel(tickDuration time.Duration, wheelSize int, levDownWheel *TimingWheel) *TimingWheel {
	tw := &TimingWheel{
		tickDur:      tickDuration,
		wheelSize:    wheelSize,
		crtWheel:     0,
		wheelTimers:  make([]map[TimerI]struct{}, wheelSize),
		levDownWheel: levDownWheel,
	}

	for i := 0; i < wheelSize; i++ {
		tw.wheelTimers[i] = map[TimerI]struct{}{}
	}

	if levDownWheel == nil {
		tw.doChan = make(chan TimerI, DoCHANMAX)
		go tw.do()
	} else {
		tw.doChan = levDownWheel.doChan
	}

	go tw.run()

	return tw
}

func (tw *TimingWheel) addTimer(delayTime time.Duration, t TimerI) {
	//延迟时间小于最小时间轮的单位长度，默认时间到达直接处理
	if int(delayTime/tw.tickDur) <= 0 {
		tw.doChan <- t
		return
	}

	index := (int(delayTime/tw.tickDur) + tw.crtWheel) % tw.wheelSize
	tw.wheelTimers[index][t] = struct{}{}
	//fmt.Println(index, delayTime, tw)
}

func (tw *TimingWheel) DelayFunc(delayTime time.Duration, callback func()) *Timer {

	expiredTime := time.Now().Add(delayTime)
	timer := newTimer(expiredTime, callback)

	tw.lock.Lock()
	defer tw.lock.Unlock()

	//超过当前时间轮的总量，创建更大的时间轮
	if int(delayTime/tw.tickDur) > tw.wheelSize {
		levUpWheel := tw.levUpWheel
		if levUpWheel == nil {
			newTW := newTimingWheel(tw.tickDur*time.Duration(tw.wheelSize), tw.wheelSize, tw)
			tw.levUpWheel = newTW
			levUpWheel = tw.levUpWheel
		}
		levUpWheel.addTimer(delayTime, timer)

	} else {
		tw.addTimer(delayTime, timer)

	}

	return timer
}

func (tw *TimingWheel) RemoveTimer(timer *Timer) bool {

	tw.lock.Lock()
	defer tw.lock.Unlock()
	for _, timers := range tw.wheelTimers {
		if _, ok := timers[timer]; ok {
			delete(timers, timer)
			return true
		}
	}

	if tw.levUpWheel != nil {
		return tw.levUpWheel.RemoveTimer(timer)
	}

	return false
}

func (tw *TimingWheel) run() {

	tw.ticker = time.NewTicker(tw.tickDur)
	for {
		select {
		case now := <-tw.ticker.C:
			tw.lock.Lock()
			tw.crtWheel = (tw.crtWheel + 1) % tw.wheelSize
			//fmt.Println(&tw, tw.crtWheel, tw.wheelTimers[tw.crtWheel])

			//将大时间轮上的事件交给更小的时间轮处理
			if tw.levDownWheel != nil {
				for t := range tw.wheelTimers[tw.crtWheel] {
					delayTime := t.SubDuration(now)
					//fmt.Println("---", delayTime, now)
					tw.levDownWheel.addTimer(delayTime, t)
				}
			} else {
				for t := range tw.wheelTimers[tw.crtWheel] {
					tw.doChan <- t
				}
			}
			tw.wheelTimers[tw.crtWheel] = map[TimerI]struct{}{}

			tw.lock.Unlock()
		}
	}
}

//超时事务处理
func (tw *TimingWheel) do() {
	for {
		select {
		case t := <-tw.doChan:
			t.Do()
		}
	}
}

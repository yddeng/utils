package pipeline

import (
	"sync"
	"sync/atomic"
)

/* *** engine *** */

type call struct {
	f     func(out interface{}, err *Error)
	value interface{}
	doneC chan *call
	en    *Engine
}

func (this *call) done(err *Error) {

	this.f(this.value, err)
	this.en.addTask(-1)

	select {
	case this.doneC <- this:
		// ok
	default:
		// We don't want to block here. It is the caller's responsibility to make
		// sure the channel has enough buffer space. See comment in Go().
	}
}

type Error struct {
	err  error
	step int
}

func (e *Error) Error() string {
	return e.err.Error()
}

func (e *Error) Step() int {
	return e.step
}

type Engine struct {
	runner *stepRunner
	taskC  chan func()
	steps  []StepFunc

	maxWorker int32
	crtWorker int32
	maxTask   int32
	crtTask   int32
	mtx       sync.Mutex
}

func (this *Engine) Exec(value interface{}) (out interface{}, err *Error) {
	c := this.invoke(value, func(out_ interface{}, err_ *Error) {
		out = out_
		err = err_
	})
	<-c.doneC
	return
}

func (this *Engine) AsyncExec(value interface{}, callback func(out interface{}, err *Error)) {
	this.invoke(value, callback)
}

func (this *Engine) invoke(value interface{}, callback func(out interface{}, err *Error)) *call {

	c := &call{f: callback, value: value, doneC: make(chan *call, 1), en: this}

	if this.runner == nil {
		c.done(nil)
		return c
	}

	this.mtx.Lock()
	defer this.mtx.Unlock()

	if atomic.LoadInt32(&this.crtTask) == this.maxTask {
		c.done(&Error{err: ErrChannelFull})
		return c
	}

	this.taskC <- this.runner.wrapper(c)

	this.addTask(1)
	if atomic.LoadInt32(&this.crtWorker) < this.maxWorker {
		this.addWorker(1)
		go this.newWorker()
	}

	return c
}

func (this *Engine) addTask(delta int32) {
	atomic.AddInt32(&this.crtTask, delta)
}

func (this *Engine) isEmpty() bool {
	return atomic.LoadInt32(&this.crtTask) == 0
}

func (this *Engine) addWorker(delta int32) {
	atomic.AddInt32(&this.crtWorker, delta)
}

func (this *Engine) newWorker() {
	for {
		if this.isEmpty() {
			break
		}

		f := <-this.taskC
		pcall(f)
	}
	this.addWorker(-1)
}

func (this *Engine) wrapper(f StepFunc, call_ *call) func() {
	return func() {
		var err error
		if call_.value, err = f(call_.value); err == nil {
			if r.next != nil {
				r.next.wrapper(call_)
			} else {
				call_.done(nil)
			}
		} else {
			call_.done(&Error{err: err, step: r.step})
		}
	}
}

func pcall(f func()) {
	f()
}

type stepRunner struct {
	step int
	f    StepFunc
	next *stepRunner
}

func (r *stepRunner) wrapper(call_ *call) func() {
	return func() {
		var err error
		if call_.value, err = r.f(call_.value); err == nil {
			if r.next != nil {
				r.next.wrapper(call_)
			} else {
				call_.done(nil)
			}
		} else {
			call_.done(&Error{err: err, step: r.step})
		}
	}
}

// default
func NewEngine(worker, size int32, steps ...StepFunc) *Engine {
	if len(steps) == 0 {
		return nil
	}

	if worker <= 0 {
		worker = 1
	}
	if size <= 0 {
		size = 256
	}

	engine := &Engine{
		maxWorker: worker,
		maxTask:   size,
		taskC:     make(chan func(), size),
	}

	var runner *stepRunner
	for i, f := range steps {
		r := &stepRunner{
			step: i,
			f:    f,
		}

		if runner == nil {
			engine.runner = r
		} else {
			runner.next = r
		}
		runner = r
	}

	return engine
}

type Option struct {
	MaxTask   int
	MaxWorker int
}

func NewEngineWithOption(options ...Option) {

}

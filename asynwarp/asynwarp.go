package asynwarp

import (
	"reflect"
	"sync"
)

func WrapFunc(oriFunc interface{}) func(callback func([]interface{}), args ...interface{}) {
	oriF := reflect.ValueOf(oriFunc)

	if oriF.Kind() != reflect.Func {
		return nil
	}

	return func(callback func([]interface{}), args ...interface{}) {
		f := func() {
			in := []reflect.Value{}
			for _, v := range args {
				in = append(in, reflect.ValueOf(v))
			}
			out := oriF.Call(in)

			if len(out) > 0 {

				ret := make([]interface{}, 0, len(out))
				for _, v := range out {
					ret = append(ret, v.Interface())
				}
				if nil != callback {
					callback(ret)
				}
			} else {
				if nil != callback {
					callback(nil)
				}
			}
		}

		if threadPool_ == nil {
			go f()
		} else {
			threadPool_.addTask(f)
		}
	}
}

var threadPool_ *threadPool

type threadPool struct {
	threadCount int32
	threadMax   int32
	taskCount   int
	taskCh      chan func()
	mu          sync.Mutex
}

func newTreadPool(threadMax int) *threadPool {
	return &threadPool{
		threadCount: 0,
		threadMax:   int32(threadMax),
		taskCh:      make(chan func(), 1024),
	}
}

func (p *threadPool) addTask(fn func()) {
	p.taskCh <- fn
	p.mu.Lock()
	p.taskCount++
	if p.threadCount < p.threadMax {
		p.threadCount++
		p.mu.Unlock()
		go p.newTread()
		return
	}
	p.mu.Unlock()
}

func (p *threadPool) newTread() {

	for {
		p.mu.Lock()
		if p.taskCount == 0 {
			p.mu.Unlock()
			break
		}
		p.mu.Unlock()

		fn := <-p.taskCh

		p.mu.Lock()
		p.taskCount--
		p.mu.Unlock()

		fn()

	}
	p.mu.Lock()
	p.threadCount--
	p.mu.Unlock()
}

package pipeline

import (
	"errors"
	"sync/atomic"
)

type Pipeline struct {
	StepFuncs []StepFunc
}

type StepFunc func(in interface{}) (out interface{}, err error)

func NewPipeline() *Pipeline {
	return &Pipeline{StepFuncs: []StepFunc{}}
}

func (this *Pipeline) AddStep(steps ...StepFunc) {
	this.StepFuncs = append(this.StepFuncs, steps...)
}

func (this *Pipeline) Run(in interface{}) (out interface{}, err error) {
	out = in
	for _, fn := range this.StepFuncs {
		if out, err = fn(out); err != nil {
			return
		}
	}
	return
}

func (this *Pipeline) RunStep(in interface{}, step int) (out interface{}, err error) {
	var fn StepFunc
	if step >= 0 && step < len(this.StepFuncs) {
		fn = this.StepFuncs[step]
	}
	if fn != nil {
		return fn(in)
	}
	return in, errors.New("step out of range")
}

func (this *Pipeline) GoStep(in chan interface{}, stepWorkers ...int) (out chan *Element) {
	chans := make([]chan *Element, len(this.StepFuncs)+1)
	for i := 0; i < cap(chans); i++ {
		chans[i] = make(chan *Element, cap(in))
	}
	out = chans[len(chans)-1]

	go func() {
		input := chans[0]
		defer close(input)
		for v := range in {
			val := v
			input <- &Element{Value: val, Step: -1}
		}
	}()

	for idx, fn := range this.StepFuncs {
		engine := &stepEngine{
			inChan:  chans[idx],
			outChan: chans[idx+1],
			errChan: out,
			step:    idx,
			stepFn:  fn,
		}

		if idx >= len(stepWorkers) {
			engine.maxWorker = 1
		} else {
			engine.maxWorker = stepWorkers[idx]
			if engine.maxWorker < 1 {
				engine.maxWorker = 1
			}
		}

		engine.run()
	}
	return
}

type Element struct {
	Value interface{}
	Step  int // 执行到哪一步
	Err   error
}

type stepEngine struct {
	step   int
	stepFn StepFunc

	inChan  chan *Element
	outChan chan *Element
	errChan chan *Element

	maxWorker int
	runWorker int32
}

func (this *stepEngine) run() {
	newWorker := func() {
		atomic.AddInt32(&this.runWorker, 1)
		defer func() {
			if atomic.LoadInt32(&this.runWorker) == 0 {
				close(this.outChan)
			}
		}()

		for e := range this.inChan {
			elem := e
			elem.Step = this.step
			elem.Value, elem.Err = this.stepFn(elem.Value)
			if elem.Err != nil {
				this.errChan <- elem
			} else {
				this.outChan <- elem
			}
		}

		atomic.AddInt32(&this.runWorker, -1)
	}

	for i := 0; i < this.maxWorker; i++ {
		go newWorker()
	}
}

package pipeline

import (
	"errors"
)

// Errors that are used throughout the pipeline API.
var (
	ErrChannelFull  = errors.New("the engine channel is full")
	ErrStepNotExist = errors.New("the step func is not exist")
)

type StepFunc func(in interface{}) (out interface{}, err error)

type Pipeline struct {
	steps []StepFunc
}

func New(steps ...StepFunc) *Pipeline {
	return &Pipeline{steps: steps}
}

func (this *Pipeline) Run(in interface{}) (out interface{}, err error) {
	out = in
	for _, fn := range this.steps {
		if out, err = fn(out); err != nil {
			return
		}
	}
	return
}

func (this *Pipeline) RunStep(in interface{}, step int) (out interface{}, err error) {
	var fn StepFunc
	if step >= 0 && step < len(this.steps) {
		fn = this.steps[step]
	}
	if fn != nil {
		return fn(in)
	}
	return in, ErrStepNotExist
}

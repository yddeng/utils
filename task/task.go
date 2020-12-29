package task

import (
	"github.com/yddeng/dutil/callFunc"
)

type Task interface {
	Do() (result []interface{}, err error)
}

type FuncTask struct {
	fn   interface{}
	args []interface{}
}

func NewFuncTask(f interface{}, args ...interface{}) *FuncTask {
	return &FuncTask{
		fn:   f,
		args: args,
	}
}

func (this *FuncTask) Do() (result []interface{}, err error) {
	return callFunc.CallFunc(this.fn, this.args...)
}

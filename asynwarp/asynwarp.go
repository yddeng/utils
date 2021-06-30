package asynwarp

import (
	"github.com/yddeng/utils"
	"github.com/yddeng/utils/task"
	"reflect"
)

type wrapFunc func(callback interface{}, args ...interface{})

func WrapFunc(oriFunc interface{}) wrapFunc {
	oriF := reflect.ValueOf(oriFunc)

	if oriF.Kind() != reflect.Func {
		panic("asynwarp: WrapFunc oriFunc is not a func")
	}

	return func(callback interface{}, args ...interface{}) {
		f := func() {
			out, err := utils.CallFunc(oriFunc, args...)
			if err != nil {
				panic(err)
			}

			if len(out) > 0 {
				if nil != callback {
					utils.CallFunc(callback, out...)
				}
			} else {
				if nil != callback {
					utils.CallFunc(callback)
				}
			}
		}

		if taskPool != nil {
			taskPool.Submit(f)
		} else {
			go f()
		}
	}
}

var taskPool *task.TaskPool

func SetTaskPool(p *task.TaskPool) {
	taskPool = p
}

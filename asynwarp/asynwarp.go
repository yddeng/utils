package asynwarp

import (
	"github.com/yddeng/dutil/callFunc"
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
			out, err := callFunc.CallFunc(oriFunc, args...)
			if err != nil {
				panic(err)
			}

			if len(out) > 0 {
				if nil != callback {
					callFunc.CallFunc(oriFunc, out...)
				}
			} else {
				if nil != callback {
					callFunc.CallFunc(oriFunc)
				}
			}
		}

		go f()
	}
}

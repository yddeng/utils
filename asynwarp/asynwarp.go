package asynwarp

import (
	"reflect"
)

func WrapFunc(oriFunc interface{}) func(callback func([]interface{}), args ...interface{}) {
	oriF := reflect.ValueOf(oriFunc)

	if oriF.Kind() != reflect.Func {
		return nil
	}
	fnType := reflect.TypeOf(oriF)

	return func(callback func([]interface{}), args ...interface{}) {
		f := func() {
			var in []reflect.Value
			numIn := fnType.NumIn()
			if numIn > 0 {
				in = make([]reflect.Value, numIn)
				for i := 0; i < numIn; i++ {
					if i >= len(args) || args[i] == nil {
						in[i] = reflect.Zero(fnType.In(i))
					} else {
						in[i] = reflect.ValueOf(args[i])
					}
				}
			}

			out := oriF.Call(in)

			if len(out) > 0 {
				ret := make([]interface{}, len(out))
				for i, v := range out {
					ret[i] = v.Interface()
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

		go f()
	}
}

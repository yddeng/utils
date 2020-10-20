package asynwarp

import (
	"reflect"
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
		go f()
	}
}

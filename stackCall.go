package dutil

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"
)

func FormatFileLine(format string, v ...interface{}) string {
	_, file, line, ok := runtime.Caller(1)
	if ok {
		s := fmt.Sprintf("[%s:%d]", file, line)
		return strings.Join([]string{s, fmt.Sprintf(format, v...)}, "")
	} else {
		return fmt.Sprintf(format, v...)
	}
}

func CallStack(maxStack int) string {
	var str string
	i := 1
	for {
		pc, file, line, ok := runtime.Caller(i)
		if !ok || i > maxStack {
			break
		}
		str += fmt.Sprintf("    stack: %d %v [file: %s] [func: %s] [line: %d]\n", i-1, ok, file, runtime.FuncForPC(pc).Name(), line)
		i++
	}
	return str
}

func ProtectCall(fn interface{}, args ...interface{}) (ret []interface{}, err error) {
	defer func() {
		if r := recover(); r != nil {
			buf := make([]byte, 65535)
			l := runtime.Stack(buf, false)
			err = fmt.Errorf(fmt.Sprintf("%v: %s", r, buf[:l]))

		}
	}()

	oriF := reflect.ValueOf(fn)

	if oriF.Kind() != reflect.Func {
		err = fmt.Errorf("not func")
		return
	}

	fnType := reflect.TypeOf(fn)

	in := make([]reflect.Value, len(args))
	for i, v := range args {
		if v == nil {
			in[i] = reflect.Zero(fnType.In(i))
		} else {
			in[i] = reflect.ValueOf(v)
		}
	}

	out := oriF.Call(in)
	if len(out) > 0 {
		ret = make([]interface{}, len(out))
		for i, v := range out {
			ret[i] = v.Interface()
		}
	}
	return
}

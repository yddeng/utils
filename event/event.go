package event

import (
	"fmt"
	"reflect"
)

type EventI interface {
	Call()
}

type Event struct {
	args []interface{}
	fn   interface{}
}

func NewEvent(fn interface{}, args ...interface{}) (EventI, error) {
	e := &Event{fn: fn}
	switch fn.(type) {
	case func():
	case func([]interface{}), func(...interface{}):
		e.args = args
	default:
		return nil, fmt.Errorf("invaild callback type %s", reflect.TypeOf(fn).String())
	}
	return e, nil
}

func (this *Event) Call() {
	fn := this.fn
	args := this.args
	switch fn.(type) {
	case func():
		fn.(func())()
	case func([]interface{}):
		fn.(func([]interface{}))(args)
	case func(...interface{}):
		fn.(func(...interface{}))(args...)
	default:
	}
}

type Task interface {
}

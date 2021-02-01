package queue

import (
	"fmt"
	"reflect"
	"testing"
	"time"
)

type Node struct {
	name string
}

func TestNewQueue(t *testing.T) {
	q := NewBlockQueue(3)
	fmt.Println("push 1:", q.Push(1))
	fmt.Println("push 2:", q.Push(2))
	e, ok := q.Pop()
	fmt.Println("pop", e, ok)

	fmt.Println("push 3:", q.Push(3))

	q.Close()
	fmt.Println(q.Closed())

	e, ok = q.Pop()
	fmt.Println("pop after closed", e, ok)
	e, ok = q.GetAll()
	fmt.Println("GetAll after closed", e, ok)

	fmt.Println("push 4:", q.Push(4))

}

func TestNewChannelQueue(t *testing.T) {
	cq := NewChannelQueue(1)

	fmt.Println(cq.PushB(1))
	fmt.Println(cq.PushN(2))

	go func() {
		for {
			time.Sleep(time.Millisecond * 100)
			elem, b := cq.Pop()
			if !b {
				break
			}
			fmt.Println(elem, b)
		}
	}()

	fmt.Println(cq.PushB(2))

	time.Sleep(time.Second)
	cq.Close()
	select {}
}

func TestNewEventHandler(t *testing.T) {
	fn1 := func() {}
	fn2 := func(i int) {}
	fn3 := func(i int, b bool) {}
	fn4 := func(i int) error { return nil }
	fn5 := func(i int) (bool, error) { return false, nil }
	fn6 := func(args ...interface{}) {}

	fmt.Println(reflect.TypeOf(fn1).String())
	fmt.Println(reflect.TypeOf(fn2).String())
	fmt.Println(reflect.TypeOf(fn3).String())
	fmt.Println(reflect.TypeOf(fn4).String())
	fmt.Println(reflect.TypeOf(fn5).String())
	fmt.Println(reflect.TypeOf(fn6).String())

	fn := func(f interface{}) {
		switch f.(type) {
		case func():
			fmt.Println(reflect.TypeOf(f).String(), "func()")
		case func([]interface{}), func(...interface{}):
			fmt.Println(reflect.TypeOf(f).String(), "func([]interface{})")
		//case func(...interface{}):
		//	fmt.Println(reflect.TypeOf(f).String(), "func(...interface{})")
		default:
			fmt.Println(reflect.TypeOf(f).String())
		}
	}
	fn(fn1)
	fn(fn2)
	fn(fn3)
	fn(fn4)
	fn(fn5)
	fn(fn6)
}

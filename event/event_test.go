package event

import (
	"fmt"
	"testing"
	"time"
)

func TestNewEvent(t *testing.T) {
	e, _ := NewEvent(func() {
		fmt.Println("func()")
	})
	e.Call()

	e, _ = NewEvent(func(args ...interface{}) {
		fmt.Println("func(args ...interface)", args)
	}, "args")
	e.Call()

	e, _ = NewEvent(func(args []interface{}) {
		fmt.Println("func(args []interface)", args)
	}, "args")
	e.Call()

	e, err := NewEvent(func(s string) {
		fmt.Println("func(string)", s)
	})
	fmt.Println(err)
	if err == nil {
		e.Call()
	}
}

func TestNewEventQueueOne(t *testing.T) {
	eq := NewEventQueue(2)

	eq.Run()

	fmt.Println("push 1:", eq.Push(func() {
		fmt.Println("1")
	}))
	fmt.Println("push 2:", eq.Push(func() {
		fmt.Println("2")
	}))
	fmt.Println("push 3:", eq.Push(func() {
		fmt.Println("3")
	}))

	eq.Stop()
	time.Sleep(time.Second)
	fmt.Println("push 4:", eq.Push(4))

}

func TestNewEventQueueMultiple(t *testing.T) {
	eq := NewEventQueue(100)

	eq.Run()

	for i := 1; i <= 100; i++ {
		k := i
		//fmt.Println("push", k, eq.Push(func() {
		//	fmt.Println(k)
		//}))

		eq.Push(func() {
			fmt.Println(k)
		})
	}

	eq.Stop()
	time.Sleep(time.Second)

}

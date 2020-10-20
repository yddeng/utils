package event

import (
	"fmt"
	"testing"
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

func TestNewEventQueue(t *testing.T) {
	eq := NewEventQueue(2)

	eq.Run(1)

	fmt.Println("push 1:", eq.Push(1))
	fmt.Println("push 2:", eq.Push(2))
	fmt.Println("push 3:", eq.Push(3))

	eq.Stop()
	fmt.Println("push 4:", eq.Push(4))

}

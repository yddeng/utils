package events

import (
	"fmt"
	"testing"
)

type Num struct {
	num int
}

func TestNewEventGroup(t *testing.T) {
	h := NewEventHandler()

	t1 := h.Listen(&Num{}, func(event interface{}) {
		a := event.(*Num)
		fmt.Println("trigger1", a)

	}, false)

	h.Listen(&Num{}, func(event interface{}) {
		a := event.(*Num)
		fmt.Println("trigger2", a)
		t1.Release()

		h.Listen(&Num{}, func(event interface{}) {
			a := event.(*Num)
			fmt.Println("trigger4", a)
		}, false)

	}, true)

	h.Listen(&Num{}, func(event interface{}) {
		a := event.(*Num)
		fmt.Println("trigger3", a)
	}, false)

	h.Trigger(&Num{num: 1})
	fmt.Println()
	h.Trigger(&Num{num: 2})

}

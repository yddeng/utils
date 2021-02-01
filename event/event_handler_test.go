package event

import (
	"fmt"
	"testing"
)

type AddNum struct {
	num int
}

func TestNewEventGroup(t *testing.T) {
	g := NewEventGroup(&AddNum{})

	num := 0
	t1 := g.Handle(func(event interface{}) {
		a := event.(*AddNum)
		num += a.num
		fmt.Println("trigger1", a, num)

	})

	num2 := 0
	cc := 0
	g.Handle(func(event interface{}) {
		a := event.(*AddNum)
		num2 += a.num
		fmt.Println("trigger2", a, num2)
		cc++
		if cc == 1 {
			g.Remove(t1)
		} else if cc == 2 {
			g.Handle(func(event interface{}) {
				a := event.(*AddNum)
				fmt.Println("trigger4", a)
			})
		}

	})

	c := 0
	num3 := 0
	g.Handle(func(event interface{}) {
		a := event.(*AddNum)
		num3 += a.num
		fmt.Println("trigger3", a, num3)
		c++
		if c < 3 {
			g.Execute(&AddNum{num: num3 + 1})
		}
	})

	fmt.Println(g.Execute(&AddNum{num: 2}))
	g.Execute(&AddNum{num: 10})
	//h.Trigger(&AddNum{num: 4})

}

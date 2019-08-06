package queue_test

import (
	"fmt"
	"github.com/tagDong/dutil/queue"
	"testing"
)

type Node struct {
	name string
}

func TestNewLinkList(t *testing.T) {
	llist := queue.NewLinkList()
	n1 := llist.Push(&Node{name: "no1"})
	n2 := llist.Push(&Node{name: "no2"})
	fmt.Println(n1, n2)

	llist.Remove(n1)
	fmt.Println(llist.Pop())
	fmt.Println(llist.Pop())
}

func TestNewQueue(t *testing.T) {
	q := queue.NewBlockQueue(2)
	fmt.Println("push 1:", q.Push(1))
	fmt.Println("push 2:", q.Push(2))
	e, ok := q.Pop()
	fmt.Println("pop", e, ok)

	fmt.Println("push 3:", q.Push(3))

	es := q.Close()
	fmt.Println("close", es)
	fmt.Println(q.Closed())

	e, ok = q.Pop()
	fmt.Println("pop after closed", e, ok)

	fmt.Println("push 4:", q.Push(4))

}

func TestNewEventQueue(t *testing.T) {
	eq := queue.NewEventQueue(2, func(i interface{}) {
		fmt.Println("onEvent", i)
	})

	go eq.Run()

	fmt.Println("push 1:", eq.Push(1))
	fmt.Println("push 2:", eq.Push(2))
	fmt.Println("push 3:", eq.Push(3))

	eq.Close()
	fmt.Println("push 4:", eq.Push(4))

}

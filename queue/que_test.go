package queue

import (
	"fmt"
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

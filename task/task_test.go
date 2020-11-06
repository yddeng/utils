package task

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

type myTask struct {
	id  int
	num int
}

func (this *myTask) Begin() {
	time.Sleep(time.Millisecond * time.Duration(rand.Intn(200)+100))
	this.num = rand.Intn(10000)
}

func (this *myTask) End() {
	fmt.Println(this.id, this.num)
}

func TestNewTaskQueue(t *testing.T) {
	tq := NewTaskQueue(1024)
	tq.Run(1)

	for i := 1; i <= 20; i++ {
		tq.AddTask(&myTask{id: i})
	}

	time.Sleep(time.Second)
	fmt.Println("stop")
	tq.Stop()
	time.Sleep(time.Second)

}

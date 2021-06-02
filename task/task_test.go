package task

import (
	"fmt"
	"testing"
	"time"
)

func TestNewFuncTask(t *testing.T) {
	t1 := NewFuncTask(f1)
	if _, err := t1.Do(); err != nil {
		fmt.Println(err)
	}

	t2 := NewFuncTask(f2, 4)
	if _, err := t2.Do(); err != nil {
		fmt.Println(err)
	}

	t2_1 := NewFuncTask(f2)
	if _, err := t2_1.Do(); err != nil {
		fmt.Println("t2_1", err)
	}

	t2_2 := NewFuncTask(f2, 3, 4)
	if _, err := t2_2.Do(); err != nil {
		fmt.Println("t2_2", err)
	}

	t3 := NewFuncTask(f3, "t3", 4, 5)
	if _, err := t3.Do(); err != nil {
		fmt.Println(err)
	}

	t3_1 := NewFuncTask(f3, "t3")
	if _, err := t3_1.Do(); err != nil {
		fmt.Println(err)
	}

	t3_2 := NewFuncTask(f3)
	if _, err := t3_2.Do(); err != nil {
		fmt.Println("t3_2", err)
	}

	t4 := NewFuncTask(f4, "t4")
	if _, err := t4.Do(); err != nil {
		fmt.Println(err)
	}

	t5 := NewFuncTask(f5, 5)
	if _, err := t5.Do(); err != nil {
		fmt.Println(err)
	}

	t6 := NewFuncTask(f6) //, fmt.Errorf("f6 error"))
	if _, err := t6.Do(); err != nil {
		fmt.Println(err)
	}

	t6_1 := NewFuncTask(f6, nil, fmt.Errorf("f6 error"), nil)
	if _, err := t6_1.Do(); err != nil {
		fmt.Println(err)
	}

	t7 := NewFuncTask(f7, nil)
	if _, err := t7.Do(); err != nil {
		fmt.Println(err)
	}

	t8 := NewFuncTask(f8, 3)
	if _, err := t8.Do(); err != nil {
		fmt.Println(err)
	}
	f := func(args ...interface{}) {
		tf := NewFuncTask(ff, args...)
		if _, err := tf.Do(); err != nil {
			fmt.Println(err)
		}
	}
	f(2)
}

func f1() {
	fmt.Println("f1")
}

func f2(i int) {
	fmt.Println("f2", i)
}

func f3(n string, i ...int) {
	fmt.Println("f3", n, i)
}

func f4(k ...string) {
	fmt.Println("f4", k)
}

func f5(m interface{}) {
	fmt.Println("f5", m)
}

func f6(m ...error) {
	fmt.Println("f6", m)
}

func f7(e error) {
	fmt.Println("f7", e)
}

func f8(args ...interface{}) {
	fmt.Println("f8", args)
}

func ff(args ...interface{}) {
	f := func(args1 ...interface{}) {
		fmt.Println(2, args1)
		t := NewFuncTask(f2, args1...)
		if _, err := t.Do(); err != nil {
			fmt.Println(err)
		}
	}
	fmt.Println(1, args)
	f(args...)
}

func TestTaskQueue_Push(t *testing.T) {
	taskQueue := NewTaskQueue(2)

	taskQueue.Run()

	fmt.Println("push 1:", taskQueue.Push(func() {
		fmt.Println("1")
	}))
	fmt.Println("push 2:", taskQueue.Push(func() {
		fmt.Println("2")
	}))
	fmt.Println("push 3:", taskQueue.Push(func() {
		fmt.Println("3")
	}))

	taskQueue.Stop()
	time.Sleep(time.Second)
	fmt.Println("push 4:", taskQueue.Push(4))
}

func TestTaskQueue_WaitPush(t *testing.T) {
	taskQueue := NewTaskQueue(2)

	taskQueue.Run()

	fmt.Println("push 1:", taskQueue.WaitPush(func() {
		fmt.Println("1")
	}))
	fmt.Println("push 2:", taskQueue.WaitPush(func() {
		fmt.Println("2")
	}))
	fmt.Println("push 3:", taskQueue.WaitPush(func() {
		fmt.Println("3")
	}))

}

package task

import (
	"fmt"
	"testing"
	"time"
)

func TestFuncTask(t *testing.T) {
	t1 := FuncTask(f1)
	if _, err := t1.Do(); err != nil {
		fmt.Println(err)
	}

	t2 := FuncTask(f2, 4)
	if _, err := t2.Do(); err != nil {
		fmt.Println(err)
	}

	t2_1 := FuncTask(f2)
	if _, err := t2_1.Do(); err != nil {
		fmt.Println("t2_1", err)
	}

	t2_2 := FuncTask(f2, 3, 4)
	if _, err := t2_2.Do(); err != nil {
		fmt.Println("t2_2", err)
	}

	t3 := FuncTask(f3, "t3", 4, 5)
	if _, err := t3.Do(); err != nil {
		fmt.Println(err)
	}

	t3_1 := FuncTask(f3, "t3")
	if _, err := t3_1.Do(); err != nil {
		fmt.Println(err)
	}

	t3_2 := FuncTask(f3)
	if _, err := t3_2.Do(); err != nil {
		fmt.Println("t3_2", err)
	}

	t4 := FuncTask(f4, "t4")
	if _, err := t4.Do(); err != nil {
		fmt.Println(err)
	}

	t5 := FuncTask(f5, 5)
	if _, err := t5.Do(); err != nil {
		fmt.Println(err)
	}

	t6 := FuncTask(f6) //, fmt.Errorf("f6 error"))
	if _, err := t6.Do(); err != nil {
		fmt.Println(err)
	}

	t6_1 := FuncTask(f6, nil, fmt.Errorf("f6 error"), nil)
	if _, err := t6_1.Do(); err != nil {
		fmt.Println(err)
	}

	t7 := FuncTask(f7, nil)
	if _, err := t7.Do(); err != nil {
		fmt.Println(err)
	}

	t8 := FuncTask(f8, 3)
	if _, err := t8.Do(); err != nil {
		fmt.Println(err)
	}
	f := func(args ...interface{}) {
		tf := FuncTask(ff, args...)
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
		t := FuncTask(f2, args1...)
		if _, err := t.Do(); err != nil {
			fmt.Println(err)
		}
	}
	fmt.Println(1, args)
	f(args...)
}

func TestNewTaskPool(t *testing.T) {
	p := NewTaskPool(2, 100, time.Second)
	p.Submit(func() {
		t.Log("f1")
	})
	p.Submit(func() {
		t.Log("f2")
	})
	time.Sleep(time.Millisecond)
	t.Log(p.Running())

	time.Sleep(time.Second * 2)
	t.Log(p.Running())

	p.Stop()
	t.Log(p.Submit(func() {
		t.Log("f3")
	}))

}

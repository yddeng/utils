package pipeline

import (
	"fmt"
	"math/rand"
	"runtime"
	"testing"
	"time"
)

func init() {
	rand.Seed(time.Now().Unix())
}

func NewStep(in interface{}) (out interface{}, err error) {
	if rand.Int()%10 == 0 {
		return nil, fmt.Errorf("new err")
	}
	return in.(int), nil
}

func AddNum(in interface{}) (out interface{}, err error) {
	time.Sleep(time.Duration(rand.Intn(200)) * time.Millisecond)
	num := in.(int) + 5
	if rand.Int()%10 == 0 {
		return nil, fmt.Errorf("add err")
	}

	return num, nil
}

func SubNum(in interface{}) (out interface{}, err error) {
	return in.(int) - 2, nil
}

func TestPipeline_RunAll(t *testing.T) {
	pipe := New(NewStep, AddNum, SubNum)
	pipe.Run(2)
}

func TestPipeline_RunStep(t *testing.T) {
	pipe := New(NewStep, AddNum, SubNum)
	pipe.RunStep(2, 2)
}

func TestPipeline_RunNull(t *testing.T) {
	pipe := New()
	out, err := pipe.Run(1)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(out)
}

func TestNewEngine(t *testing.T) {
	e := NewEngine(0, 0, NewStep, AddNum, SubNum)

	f := func(out interface{}, err *Error) {
		fmt.Println(out, err)
		if err != nil {
			//fmt.Println(err.Step(), err.Error())
		}
	}

	for i := 1; i < 100; i++ {
		num := i
		e.AsyncExec(num, f)
	}

	fmt.Println("NumGoroutine", runtime.NumGoroutine())
	time.Sleep(time.Second)
	fmt.Println("NumGoroutine", runtime.NumGoroutine())
	time.Sleep(time.Second)
	fmt.Println("NumGoroutine", runtime.NumGoroutine())
}

package pipeline

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func init() {
	rand.Seed(time.Now().Unix())
}

type step1 struct {
	num int
}

type step2 struct {
	old int
	num int
}

func NewStep(in interface{}) (out interface{}, err error) {
	if rand.Int()%10 == 0 {
		return nil, fmt.Errorf("new err")
	}
	return &step1{num: in.(int)}, nil
}

func AddNum(in interface{}) (out interface{}, err error) {
	time.Sleep(time.Duration(rand.Intn(200)) * time.Millisecond)
	p := in.(*step1)
	num := p.num + 5
	if rand.Int()%10 == 0 {
		return nil, fmt.Errorf("add err")
	}

	return &step2{old: p.num, num: num}, nil
}

func OutNum(in interface{}) (out interface{}, err error) {
	p := in.(*step2)
	fmt.Println(p.old, p.num)
	return nil, nil
}

func TestPipeline_RunAll(t *testing.T) {
	pipe := NewPipeline()
	pipe.AddStep(NewStep, AddNum, OutNum)
	pipe.Run(2)
}

func TestPipeline_RunStep(t *testing.T) {
	pipe := NewPipeline()
	pipe.AddStep(NewStep, AddNum, OutNum)
	pipe.RunStep(2, 2)
}

func TestPipeline_RunNull(t *testing.T) {
	pipe := NewPipeline()
	out, err := pipe.Run(1)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(out)
}

func TestPipeline_GoStep(t *testing.T) {
	pipe := NewPipeline()
	pipe.AddStep(NewStep, AddNum, OutNum)

	in := make(chan interface{}, 5)
	out := pipe.GoStep(in, 1, 2, 1)

	go func() {
		for i := 1; i < 10; i++ {
			in <- i
		}

		time.Sleep(time.Second)
		close(in)
	}()

	for v := range out {
		if v.Err != nil {
			fmt.Println(v.Step, v.Err)
		}
	}
	fmt.Println("end")
}

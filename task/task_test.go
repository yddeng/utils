package task

import (
	"testing"
	"time"
)

func TestNewTaskPool(t *testing.T) {
	p := NewTaskPool(0, 100)
	p.Submit(func() { t.Log("f1") })
	p.Submit(func() { t.Log("f2") })
	p.Submit(func() { t.Log("f3") })
	p.Submit(func() { t.Log("f4") })
	p.Submit(func() { t.Log("f5") })

	t.Log(p.NumWorker())
	time.Sleep(time.Millisecond)
	t.Log(p.NumWorker())

	p.Submit(func() { t.Log("f6") })
	p.Submit(func() { t.Log("f7") })
	p.Submit(func() { t.Log("f8") })

	t.Log(p.NumWorker())
	time.Sleep(time.Millisecond)
	t.Log(p.NumWorker())

	p.Stop()
	t.Log(p.Submit(func() { t.Log("f10") }))
}
func TestNewTaskPool2(t *testing.T) {
	p := NewTaskPool(2, 100)
	p.Submit(func() { t.Log("f1") })
	p.Submit(func() { t.Log("f2") })
	p.Submit(func() { t.Log("f3") })
	p.Submit(func() { t.Log("f4") })
	p.Submit(func() { t.Log("f5") })
	t.Log(p.NumWorker())

	time.Sleep(time.Millisecond)
	t.Log(p.NumWorker())

	p.Stop()
	t.Log(p.Submit(func() { t.Log("f10") }))
}

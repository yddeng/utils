package cache

import (
	"fmt"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	m := NewMemoryCache()

	m.Put("1", 1, 0)
	m.Put("2", 2, time.Second)
	m.Put("3", 3, time.Second*5)

	fmt.Println(m.Get("1"), m.Get("2"), m.Get("3"))
	fmt.Println(m.GetMulti([]string{"1", "2", "3"}))
	time.Sleep(time.Second * 2)
	fmt.Println(m.Get("1"), m.Get("2"), m.Get("3"))
	time.Sleep(time.Second * 4)
	fmt.Println(m.Get("1"), m.Get("2"), m.Get("3"))
}

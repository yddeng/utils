package lset

import (
	"testing"
)

func TestNew(t *testing.T) {
	ls := New()

	var i int = 4
	ls.Store(1, "1")
	ls.Store(2, "2")
	ls.Store(3, "3")
	ls.Store(4, &i)

	ls.Range(func(key, value interface{}) bool {
		t.Log(key, value)
		return true
	})

	t.Log(ls.Load(2))

	//ls.Delete(3)
	ls.Store(4, "4")
	ls.Range(func(key, value interface{}) bool {
		t.Log(key, value)
		return true
	})

	// sort
	ls.Sort(func(key1, value1, key2, value2 interface{}) bool {
		return key1.(int) > key2.(int)
	})
	ls.Range(func(key, value interface{}) bool {
		t.Log(key, value)
		return true
	})

	t.Log(ls.Load(2))
	ls.Delete(2)
	t.Log(ls.Len())
	ls.Range(func(key, value interface{}) bool {
		t.Log(key, value)
		return true
	})

}

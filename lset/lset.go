package lset

import (
	"container/list"
	"sort"
)

type element struct {
	key   interface{}
	value interface{}
}

type LSet struct {
	l *list.List
	m map[interface{}]*list.Element
}

func New() *LSet {
	return &LSet{
		l: list.New(),
		m: map[interface{}]*list.Element{},
	}
}

func (this *LSet) Store(key, value interface{}) {
	elem := &element{key: key, value: value}
	if e, ok := this.m[key]; ok {
		e.Value = elem
	} else {
		e = this.l.PushBack(elem)
		this.m[key] = e
	}
}

func (this *LSet) Load(key interface{}) (value interface{}, ok bool) {
	if e, ok := this.m[key]; ok {
		return e.Value.(*element).value, true
	} else {
		return nil, false
	}
}

func (this *LSet) Delete(key interface{}) {
	if e, ok := this.m[key]; ok {
		this.l.Remove(e)
		delete(this.m, key)
	}
}

func (this *LSet) Len() int {
	return this.l.Len()
}

func (this *LSet) Range(f func(key, value interface{}) bool) {
	for e := this.l.Front(); e != nil; e = e.Next() {
		elem := e.Value.(*element)
		if !f(elem.key, elem.value) {
			break
		}
	}
}

func (this *LSet) Sort(less func(key1, value1, key2, value2 interface{}) bool) {
	elements := make([]*element, 0, this.Len())
	for e := this.l.Front(); e != nil; e = e.Next() {
		elements = append(elements, e.Value.(*element))
	}

	sort.Slice(elements, func(i, j int) bool {
		return less(elements[i].key, elements[i].value, elements[j].key, elements[j].value)
	})

	this.l.Init()
	this.m = map[interface{}]*list.Element{}
	for _, elem := range elements {
		e := this.l.PushBack(elem)
		this.m[elem.key] = e
	}
}

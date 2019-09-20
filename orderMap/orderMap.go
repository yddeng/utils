package orderMap

import (
	"fmt"
	"math/rand"
	"reflect"
	"sync"
)

type OrderMap struct {
	keyType    reflect.Type
	valType    reflect.Type
	idxToKey   []interface{}
	keyToIdx   map[interface{}]int
	keyToValue map[interface{}]interface{}
	mu         sync.Mutex
}

func New() *OrderMap {
	o := &OrderMap{
		idxToKey:   []interface{}{},
		keyToIdx:   map[interface{}]int{},
		keyToValue: map[interface{}]interface{}{},
	}
	return o
}

func (o *OrderMap) SetType(keyType, valType reflect.Type) {
	o.mu.Lock()
	defer o.mu.Unlock()
	o.keyType = keyType
	o.valType = valType
}

func (o *OrderMap) Store(key, val interface{}) error {

	if o.keyType != nil && o.keyType != reflect.TypeOf(key) {
		return fmt.Errorf("invaild key type(%v), need type(%v)", reflect.TypeOf(key).String(), o.keyType.String())
	}
	if o.valType != nil && o.valType != reflect.TypeOf(val) {
		return fmt.Errorf("invaild val type(%v), need type(%v)", reflect.TypeOf(val).String(), o.valType.String())
	}

	o.mu.Lock()
	o.keyToValue[key] = val

	o.keyToIdx[key] = len(o.idxToKey)
	o.idxToKey = append(o.idxToKey, key)
	o.mu.Unlock()
	return nil
}

func (o *OrderMap) Del(key interface{}) {
	o.mu.Lock()
	defer o.mu.Unlock()

	idx, ok := o.keyToIdx[key]
	if ok {
		o.idxToKey = append(o.idxToKey[:idx], o.idxToKey[idx+1:]...)
		delete(o.keyToValue, key)

		o.keyToIdx = map[interface{}]int{}
		for i, k := range o.idxToKey {
			o.keyToIdx[k] = i
		}
	}
}

func (o *OrderMap) Len() int {
	o.mu.Lock()
	defer o.mu.Unlock()
	return len(o.keyToValue)
}

func (o *OrderMap) LoadByIdx(idx int) (interface{}, bool) {
	o.mu.Lock()
	defer o.mu.Unlock()
	if idx < 0 || idx >= len(o.idxToKey) {
		return nil, false
	}

	v := o.keyToValue[o.idxToKey[idx]]
	return v, true
}

func (o *OrderMap) LoadByKey(key interface{}) (interface{}, bool) {
	o.mu.Lock()
	defer o.mu.Unlock()

	val, ok := o.keyToValue[key]
	if !ok {
		return nil, false
	} else {
		return val, true
	}

}

func (o *OrderMap) Range(f func(key, value interface{}) bool) {
	o.mu.Lock()
	defer o.mu.Unlock()

	for _, k := range o.idxToKey {
		v, ok := o.keyToValue[k]
		if !ok {
			continue
		}
		if !f(k, v) {
			break
		}
	}
}

func (o *OrderMap) RandOneKey() interface{} {
	o.mu.Lock()
	defer o.mu.Unlock()

	size := len(o.idxToKey)
	if size > 0 {
		i := rand.Int() % size
		return o.idxToKey[i]
	} else {
		return nil
	}
}

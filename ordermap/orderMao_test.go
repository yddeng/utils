/*
 *
 */
package ordermap

import (
	"fmt"
	"reflect"
	"testing"
)

type TT struct {
	ID int32
	S  string
}

func TestNew(t *testing.T) {
	oMap := New()
	fmt.Println(oMap.Store(int32(1), &TT{ID: 1, S: "11"}))
	fmt.Println(oMap.Store(int32(2), &TT{ID: 2, S: "22"}))
	fmt.Println(oMap.Store(int32(3), &TT{ID: 3, S: "33"}))

	oMap.Range(func(key, value interface{}) bool {
		fmt.Println(key, value)
		return true
	})
	val, ok := oMap.LoadByIdx(0)
	fmt.Println("get index 0 :", val, ok)
	val, ok = oMap.LoadByKey(int32(1))
	fmt.Println("get key 1 :", val, ok)

	oMap.Del(int32(1))
	fmt.Println("------------------")

	oMap.SetType(reflect.TypeOf(int32(1)), reflect.TypeOf(&TT{}))
	fmt.Println(oMap.Store("key", "value"))

	var er error
	er.Error()
	val, ok = oMap.LoadByIdx(0)
	fmt.Println("get index 0 :", val, ok)
	val, ok = oMap.LoadByKey(int32(1))
	fmt.Println("get key 1 :", val, ok)
}

package protocol

import (
	"fmt"
	"reflect"
)

//协议
type Protoc interface {
	//反序列化
	Unmarshaler(data []byte, o interface{}) (err error)
	//序列化
	Marshaler(data interface{}) ([]byte, error)
}

type Protocol struct {
	id2Type map[uint16]reflect.Type
	type2Id map[reflect.Type]uint16
	protoc  Protoc
}

var dProto *Protocol

//注册协议的类型（json,protobuf）
func InitProtocol(protoType Protoc) {
	dProto = &Protocol{
		id2Type: map[uint16]reflect.Type{},
		type2Id: map[reflect.Type]uint16{},
		protoc:  protoType,
	}
}

//注册协议，ID <-> 协议结构。编解码时使用
func Register(id uint16, msg interface{}) error {
	if dProto == nil {
		return fmt.Errorf("protocol is nil,need init")
	}

	tt := reflect.TypeOf(msg)

	if _, ok := dProto.id2Type[id]; ok {
		return fmt.Errorf("%d already register to type:%s", id, tt)
	}

	dProto.id2Type[id] = tt
	dProto.type2Id[tt] = id

	return nil
}

//序列化，根据反射类型，获取协议ID、序列化后的二进制数据
func Marshal(data interface{}) (uint16, []byte, error) {
	if dProto == nil {
		return 0, nil, fmt.Errorf("protocol is nil,need init")
	}

	id, ok := dProto.type2Id[reflect.TypeOf(data)]
	if !ok {
		return 0, nil, fmt.Errorf("type: %s undefined", reflect.TypeOf(data))
	}

	ret, err := dProto.protoc.Marshaler(data)
	if err != nil {
		return 0, nil, err
	}

	return id, ret, nil
}

//反序列化
func Unmarshal(id uint16, data []byte) (msg interface{}, err error) {
	if dProto == nil {
		return nil, fmt.Errorf("protocol is nil,need init")
	}

	tt, ok := dProto.id2Type[id]
	if !ok {
		err = fmt.Errorf("ID: %d undefined", id)
		return
	}

	//反序列化的结构
	msg = reflect.New(tt.Elem()).Interface()
	err = dProto.protoc.Unmarshaler(data, msg)
	if err != nil {
		return nil, err
	}

	return msg, nil
}

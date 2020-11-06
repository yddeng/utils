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

func NewProtoc(protoc Protoc) *Protocol {
	return &Protocol{
		id2Type: map[uint16]reflect.Type{},
		type2Id: map[reflect.Type]uint16{},
		protoc:  protoc,
	}
}

func (this *Protocol) Register(id uint16, msg interface{}) {
	tt := reflect.TypeOf(msg)

	if _, ok := this.id2Type[id]; ok {
		panic(fmt.Sprintf("%d already register to type:%s\n", id, tt))
	}

	this.id2Type[id] = tt
	this.type2Id[tt] = id
}

func (this *Protocol) Marshal(data interface{}) (uint16, []byte, error) {
	id, ok := this.type2Id[reflect.TypeOf(data)]
	if !ok {
		return 0, nil, fmt.Errorf("marshal type: %s undefined", reflect.TypeOf(data))
	}

	ret, err := this.protoc.Marshaler(data)
	if err != nil {
		return 0, nil, err
	}

	return id, ret, nil
}

func (this *Protocol) Unmarshal(msgID uint16, data []byte) (msg interface{}, err error) {
	tt, ok := this.id2Type[msgID]
	if !ok {
		err = fmt.Errorf("unmarshal msgID: %d undefined", msgID)
		return
	}

	//反序列化的结构
	msg = reflect.New(tt.Elem()).Interface()
	err = this.protoc.Unmarshaler(data, msg)
	if err != nil {
		return nil, err
	}

	return msg, nil
}

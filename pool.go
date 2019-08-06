package dutil

import "sync"

var (
	uint16BytePool = sync.Pool{
		New: func() interface{} {
			data := make([]byte, 2)
			return &data
		},
	}

	uint32BytePool = sync.Pool{
		New: func() interface{} {
			data := make([]byte, 4)
			return &data
		},
	}
)

func GetUint16Byte() []byte {
	return uint16BytePool.Get().([]byte)
}

func PutUint16Byte(v []byte) {
	uint16BytePool.Put(v)
}

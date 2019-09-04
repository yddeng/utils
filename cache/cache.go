/*
 * cache
 */
package cache

import (
	"github.com/tagDong/dutil/heap"
	"reflect"
	"sync"
	"time"
)

type DataBase interface {
	Save(key string, value interface{}) error
	Load(key string, ret interface{}) error
}

type TT int16

const (
	LFU TT = iota //（Least Frequently Used）算法根据数据的历史访问频率来淘汰数据，其核心思想是“如果数据过去被访问多次，那么将来被访问的频率也更高”。
	LRU           //（Least recently used，最近最少使用）算法根据数据的历史访问记录来进行淘汰数据，其核心思想是“如果数据最近被访问过，那么将来被访问的几率也更高”。
)

// 淘汰算法
var cacheTT = LFU

// 字符串类型的缓存
type Cache struct {
	dbs     DataBase
	size    int //存储量
	mu      sync.Mutex
	valType reflect.Type
	data    map[string]*val
	minHeap *heap.Heap
}

type val struct {
	key          string
	value        interface{}
	lastUsedTime time.Time //最近使用时间，LRU算法
	usedTimes    int       //使用次数，LFU算法
	dbDirty      bool      //脏数据，存储
}

func (v *val) Less(e heap.Element) bool {
	switch cacheTT {
	case LFU:
		return v.usedTimes < (e.(*val).usedTimes)
	case LRU:
		return v.lastUsedTime.Before(e.(*val).lastUsedTime)
	default:
		return v.lastUsedTime.Before(e.(*val).lastUsedTime)
	}
}

func New(dbs DataBase, size int, vT interface{}, tt TT) *Cache {
	cacheTT = tt
	return &Cache{
		dbs:     dbs,
		size:    size,
		valType: reflect.TypeOf(vT),
		data:    map[string]*val{},
		minHeap: heap.NewHeap(),
	}
}

func (c *Cache) Set(key string, value interface{}) {
	//类型判断
	if reflect.TypeOf(value) != c.valType {
		return
	}

	c.mu.Lock()
	defer c.mu.Unlock()
	_val, ok := c.data[key]
	if ok {
		_val.value = value
	} else {
		c.checkAndRemove()
		_val = &val{
			key:          key,
			value:        value,
			lastUsedTime: time.Now(),
			usedTimes:    1,
		}
		c.data[key] = _val
		c.minHeap.Push(_val)
	}
	_val.dbDirty = true

}

func (c *Cache) Get(key string) (interface{}, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	_val, ok := c.data[key]
	if ok {
		_val.usedTimes++
		_val.lastUsedTime = time.Now()
		c.minHeap.Fix(_val)
		return _val.value, true
	} else {

		if c.dbs != nil {
			// 从本地读取数据,加载到缓存
			ret := reflect.New(c.valType.Elem()).Interface()
			err := c.dbs.Load(key, ret)
			if err != nil {
				return nil, false
			}

			c.checkAndRemove()
			_val = &val{
				key:          key,
				value:        ret,
				lastUsedTime: time.Now(),
				usedTimes:    1,
			}
			c.data[key] = _val
			c.minHeap.Push(_val)

			return ret, true
		}

		return nil, false
	}
}

func (c *Cache) GetAll() map[string]interface{} {
	c.mu.Lock()
	defer c.mu.Unlock()
	ret := map[string]interface{}{}
	for _, v := range c.data {
		ret[v.key] = v.value
	}
	return ret
}

func (c *Cache) Size() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return len(c.data)
}

func (c *Cache) checkAndRemove() {
	if len(c.data) >= c.size {
		min := c.minHeap.Pop()
		v := min.(*val)
		delete(c.data, v.key)

		if c.dbs != nil {
			// 保存到数据库
			c.dbs.Save(v.key, v.value)
		}
	}
}

//将脏标记数据保存本地
func (c *Cache) SaveDirty() {
	if c.dbs != nil {
		c.mu.Lock()
		defer c.mu.Unlock()

		for _, v := range c.data {
			if v.dbDirty {
				if err := c.dbs.Save(v.key, v.value); err == nil {
					v.dbDirty = false
				}
			}
		}
	}
}

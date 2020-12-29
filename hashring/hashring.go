package hashring

import (
	"errors"
	"hash/crc32"
	"hash/fnv"
	"sort"
	"strconv"
	"sync"
)

type uints []uint32

func (x uints) Len() int {
	return len(x)
}

func (x uints) Less(i, j int) bool {
	return x[i] < x[j]
}

func (x uints) Swap(i, j int) {
	x[i], x[j] = x[j], x[i]
}

var ErrEmptyCircle = errors.New("empty circle")

type HashRing struct {
	circle       map[uint32]string
	nodes        map[string]struct{}
	sortedHashes uints
	count        int64
	NumOfVirtual int  // 虚拟节点数量,默认20
	UseFnv       bool // go的fnv非加密算法，默认使用crc32
	sync.RWMutex
}

func New() *HashRing {
	c := new(HashRing)
	c.NumOfVirtual = 20
	c.circle = make(map[uint32]string)
	c.nodes = make(map[string]struct{})
	return c
}

func (c *HashRing) eltKey(elt string, idx int) string {
	return strconv.Itoa(idx) + elt
}

func (c *HashRing) Add(elt string) {
	c.Lock()
	defer c.Unlock()
	c.add(elt)
}

func (c *HashRing) add(elt string) {
	for i := 0; i < c.NumOfVirtual; i++ {
		c.circle[c.hashKey(c.eltKey(elt, i))] = elt
	}
	c.nodes[elt] = struct{}{}
	c.updateSortedHashes()
	c.count++
}

func (c *HashRing) Remove(elt string) {
	c.Lock()
	defer c.Unlock()
	c.remove(elt)
}

func (c *HashRing) remove(elt string) {
	for i := 0; i < c.NumOfVirtual; i++ {
		delete(c.circle, c.hashKey(c.eltKey(elt, i)))
	}
	delete(c.nodes, elt)
	c.updateSortedHashes()
	c.count--
}

// 替换
func (c *HashRing) Set(elts []string) {
	c.Lock()
	defer c.Unlock()
	for k := range c.nodes {
		found := false
		for _, v := range elts {
			if k == v {
				found = true
				break
			}
		}
		if !found {
			c.remove(k)
		}
	}
	for _, v := range elts {
		_, exists := c.nodes[v]
		if exists {
			continue
		}
		c.add(v)
	}
}

func (c *HashRing) Nodes() []string {
	c.RLock()
	defer c.RUnlock()
	var m []string
	for k := range c.nodes {
		m = append(m, k)
	}
	return m
}

// 获取键 hash的节点
func (c *HashRing) Get(name string) (string, error) {
	c.RLock()
	defer c.RUnlock()
	if len(c.circle) == 0 {
		return "", ErrEmptyCircle
	}
	key := c.hashKey(name)
	i := c.search(key)
	return c.circle[c.sortedHashes[i]], nil
}

// 查找key的下一节点，不包含key位置节点。 节点范围：[上一节点位置:本节点位置)
func (c *HashRing) search(key uint32) (i int) {
	i = sort.Search(len(c.sortedHashes), func(x int) bool {
		return c.sortedHashes[x] > key
	})
	if i >= len(c.sortedHashes) {
		i = 0
	}
	return
}

// 查找多个不同节点
func (c *HashRing) GetN(name string, n int) ([]string, error) {
	c.RLock()
	defer c.RUnlock()

	if len(c.circle) == 0 {
		return nil, ErrEmptyCircle
	}

	if c.count < int64(n) {
		n = int(c.count)
	}

	var (
		key   = c.hashKey(name)
		i     = c.search(key)
		start = i
		res   = make([]string, 0, n)
		elem  = c.circle[c.sortedHashes[i]]
	)

	res = append(res, elem)

	if len(res) == n {
		return res, nil
	}

	for i = start + 1; i != start; i++ {
		if i >= len(c.sortedHashes) {
			i = 0
		}
		elem = c.circle[c.sortedHashes[i]]
		if !sliceContainsNodes(res, elem) {
			res = append(res, elem)
		}
		if len(res) == n {
			break
		}
	}

	return res, nil
}

func (c *HashRing) hashKey(key string) uint32 {
	if c.UseFnv {
		return c.hashKeyFnv(key)
	}
	return c.hashKeyCRC32(key)
}

func (c *HashRing) hashKeyCRC32(key string) uint32 {
	if len(key) < 64 {
		var scratch [64]byte
		copy(scratch[:], key)
		return crc32.ChecksumIEEE(scratch[:len(key)])
	}
	return crc32.ChecksumIEEE([]byte(key))
}

func (c *HashRing) hashKeyFnv(key string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(key))
	return h.Sum32()
}

func (c *HashRing) updateSortedHashes() {
	hashes := c.sortedHashes[:0]
	if cap(c.sortedHashes)/(c.NumOfVirtual*4) > len(c.circle) {
		hashes = nil
	}
	for k := range c.circle {
		hashes = append(hashes, k)
	}
	sort.Sort(hashes)
	c.sortedHashes = hashes
}

func sliceContainsNodes(set []string, node string) bool {
	for _, m := range set {
		if m == node {
			return true
		}
	}
	return false
}

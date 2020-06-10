package skiplist

import (
	"bytes"
	"math/rand"
)

// Propability is the propability to create a new skiplist level.
const Propability = 0x3FFF

var (
	// DefaultMaxLevel is the default max level of a skiplist.
	DefaultMaxLevel = 32
	defaultSource   = defaultRandSource{}
)

// All built-in comparasion functions.
var (
	Byte GreaterThanFunc = func(lhs, rhs interface{}) bool {
		return lhs.(byte) > rhs.(byte)
	}
	ByteAscending               = Byte
	ByteAsc                     = Byte
	ByteDescending LessThanFunc = func(lhs, rhs interface{}) bool {
		return lhs.(byte) < rhs.(byte)
	}
	ByteDesc = ByteDescending

	Float32 GreaterThanFunc = func(lhs, rhs interface{}) bool {
		return lhs.(float32) > rhs.(float32)
	}
	Float32Ascending               = Float32
	Float32Asc                     = Float32
	Float32Descending LessThanFunc = func(lhs, rhs interface{}) bool {
		return lhs.(float32) < rhs.(float32)
	}
	Float32Desc = Float32Descending

	Float64 GreaterThanFunc = func(lhs, rhs interface{}) bool {
		return lhs.(float64) > rhs.(float64)
	}
	Float64Ascending               = Float64
	Float64Asc                     = Float64
	Float64Descending LessThanFunc = func(lhs, rhs interface{}) bool {
		return lhs.(float64) < rhs.(float64)
	}
	Float64Desc = Float64Descending

	Int GreaterThanFunc = func(lhs, rhs interface{}) bool {
		return lhs.(int) > rhs.(int)
	}
	IntAscending               = Int
	IntAsc                     = Int
	IntDescending LessThanFunc = func(lhs, rhs interface{}) bool {
		return lhs.(int) < rhs.(int)
	}
	IntDesc = IntDescending

	Int16 GreaterThanFunc = func(lhs, rhs interface{}) bool {
		return lhs.(int16) > rhs.(int16)
	}
	Int16Ascending               = Int16
	Int16Asc                     = Int16
	Int16Descending LessThanFunc = func(lhs, rhs interface{}) bool {
		return lhs.(int16) < rhs.(int16)
	}
	Int16Desc = Int16Descending

	Int32 GreaterThanFunc = func(lhs, rhs interface{}) bool {
		return lhs.(int32) > rhs.(int32)
	}
	Int32Ascending               = Int32
	Int32Asc                     = Int32
	Int32Descending LessThanFunc = func(lhs, rhs interface{}) bool {
		return lhs.(int32) < rhs.(int32)
	}
	Int32Desc = Int32Descending

	Int64 GreaterThanFunc = func(lhs, rhs interface{}) bool {
		return lhs.(int64) > rhs.(int64)
	}
	Int64Ascending               = Int64
	Int64Asc                     = Int64
	Int64Descending LessThanFunc = func(lhs, rhs interface{}) bool {
		return lhs.(int64) < rhs.(int64)
	}
	Int64Desc = Int64Descending

	Int8 GreaterThanFunc = func(lhs, rhs interface{}) bool {
		return lhs.(int8) > rhs.(int8)
	}
	Int8Ascending               = Int8
	Int8Asc                     = Int8
	Int8Descending LessThanFunc = func(lhs, rhs interface{}) bool {
		return lhs.(int8) < rhs.(int8)
	}
	Int8Desc = Int8Descending

	Rune GreaterThanFunc = func(lhs, rhs interface{}) bool {
		return lhs.(rune) > rhs.(rune)
	}
	RuneAscending               = Rune
	RuneAsc                     = Rune
	RuneDescending LessThanFunc = func(lhs, rhs interface{}) bool {
		return lhs.(rune) < rhs.(rune)
	}
	RuneDesc = RuneDescending

	String GreaterThanFunc = func(lhs, rhs interface{}) bool {
		return lhs.(string) > rhs.(string)
	}
	StringAscending               = String
	StringAsc                     = String
	StringDescending LessThanFunc = func(lhs, rhs interface{}) bool {
		return lhs.(string) < rhs.(string)
	}
	StringDesc = StringDescending

	Uint GreaterThanFunc = func(lhs, rhs interface{}) bool {
		return lhs.(uint) > rhs.(uint)
	}
	UintAscending               = Uint
	UintAsc                     = Uint
	UintDescending LessThanFunc = func(lhs, rhs interface{}) bool {
		return lhs.(uint) < rhs.(uint)
	}
	UintDesc = UintDescending

	Uint16 GreaterThanFunc = func(lhs, rhs interface{}) bool {
		return lhs.(uint16) > rhs.(uint16)
	}
	Uint16Ascending               = Uint16
	Uint16Asc                     = Uint16
	Uint16Descending LessThanFunc = func(lhs, rhs interface{}) bool {
		return lhs.(uint16) < rhs.(uint16)
	}
	Uint16Desc = Uint16Descending

	Uint32 GreaterThanFunc = func(lhs, rhs interface{}) bool {
		return lhs.(uint32) > rhs.(uint32)
	}
	Uint32Ascending               = Uint32
	Uint32Asc                     = Uint32
	Uint32Descending LessThanFunc = func(lhs, rhs interface{}) bool {
		return lhs.(uint32) < rhs.(uint32)
	}
	Uint32Desc = Uint32Descending

	Uint64 GreaterThanFunc = func(lhs, rhs interface{}) bool {
		return lhs.(uint64) > rhs.(uint64)
	}
	Uint64Ascending               = Uint64
	Uint64Asc                     = Uint64
	Uint64Descending LessThanFunc = func(lhs, rhs interface{}) bool {
		return lhs.(uint64) < rhs.(uint64)
	}
	Uint64Desc = Uint64Descending

	Uint8 GreaterThanFunc = func(lhs, rhs interface{}) bool {
		return lhs.(uint8) > rhs.(uint8)
	}
	Uint8Ascending               = Uint8
	Uint8Asc                     = Uint8
	Uint8Descending LessThanFunc = func(lhs, rhs interface{}) bool {
		return lhs.(uint8) < rhs.(uint8)
	}
	Uint8Desc = Uint8Descending

	Uintptr GreaterThanFunc = func(lhs, rhs interface{}) bool {
		return lhs.(uintptr) > rhs.(uintptr)
	}
	UintptrAscending               = Uintptr
	UintptrAsc                     = Uintptr
	UintptrDescending LessThanFunc = func(lhs, rhs interface{}) bool {
		return lhs.(uintptr) < rhs.(uintptr)
	}
	UintptrDesc = UintptrDescending

	// the type []byte.
	Bytes GreaterThanFunc = func(lhs, rhs interface{}) bool {
		return bytes.Compare(lhs.([]byte), rhs.([]byte)) > 0
	}
	BytesAscending = Bytes
	BytesAsc       = Bytes
	// the type []byte. reversed order.
	BytesDescending LessThanFunc = func(lhs, rhs interface{}) bool {
		return bytes.Compare(lhs.([]byte), rhs.([]byte)) < 0
	}
	BytesDesc = BytesDescending
)

// GreaterThanFunc returns true if lhs greater than rhs
type GreaterThanFunc func(lhs, rhs interface{}) bool

// LessThanFunc returns true if lhs less than rhs
type LessThanFunc GreaterThanFunc

type defaultRandSource struct{}

// Comparable defines a comparable element.
type Comparable interface {
	Descending() bool
	Compare(lhs, rhs interface{}) bool
}

type elementNode struct {
	next []*Element
}

// Element is an element in skiplist.
type Element struct {
	elementNode
	key, Value interface{}
	score      float64
}

// SkipList represents a skiplist header node.
type SkipList struct {
	elementNode
	level      int
	length     int
	keyFunc    Comparable
	randSource rand.Source
	reversed   bool

	prevNodesCache []*elementNode // a cache for Set/Remove
}

// Scorable is used by skip list using customized key comparing function.
// For built-in functions, there is no need to care of this interface.
//
// Every skip list element with customized key must have a score value
// to indicate its sequence.
// For any two elements with key "k1" and "k2":
// - If Compare(k1, k2) is true, k1.Score() >= k2.Score() must be true.
// - If Compare(k1, k2) is false and k1 doesn't equal to k2, k1.Score() < k2.Score() must be true.
type Scorable interface {
	Score() float64
}

func (r defaultRandSource) Int63() int64 {
	return rand.Int63()
}

func (r defaultRandSource) Seed(seed int64) {
	rand.Seed(seed)
}

// Descending always returns false to sort list in ascending order.
func (f GreaterThanFunc) Descending() bool {
	return false
}

// Compare compares lhs and rhs using f.
func (f GreaterThanFunc) Compare(lhs, rhs interface{}) bool {
	return f(lhs, rhs)
}

// Descending always returns true to sort list in descending order.
func (f LessThanFunc) Descending() bool {
	return true
}

// Compare compares lhs and rhs using f.
func (f LessThanFunc) Compare(lhs, rhs interface{}) bool {
	return f(lhs, rhs)
}

// Next returns the ajancent next element.
func (element *Element) Next() *Element {
	return element.next[0]
}

// NextLevel returns next element at a specific level.
func (element *Element) NextLevel(level int) *Element {
	if level >= len(element.next) || level < 0 {
		panic("invalid argument to NextLevel")
	}

	return element.next[level]
}

// Key returns the key of element.
func (element *Element) Key() interface{} {
	return element.key
}

func New(keyFunc Comparable) *SkipList {
	if DefaultMaxLevel <= 0 {
		panic("skiplist default level must not be zero or negative")
	}

	return &SkipList{
		elementNode:    elementNode{next: make([]*Element, DefaultMaxLevel)},
		prevNodesCache: make([]*elementNode, DefaultMaxLevel),
		level:          DefaultMaxLevel,
		keyFunc:        keyFunc,
		randSource:     defaultSource,
		reversed:       keyFunc.Descending(),
	}
}

// Init resets a skiplist and discards all exists elements.
func (list *SkipList) Init() *SkipList {
	list.next = make([]*Element, list.level)
	list.length = 0
	return list
}

// SetRandSource sets a new rand source.
//
// Skiplist uses global rand defined in math/rand by default.
// The default rand acquires a global mutex before generating any number.
// It's not necessary if the skiplist is well protected by caller.
func (list *SkipList) SetRandSource(source rand.Source) {
	list.randSource = source
}

// Front returns the first element.
func (list *SkipList) Front() *Element {
	return list.next[0]
}

// Len returns list length.
func (list *SkipList) Len() int {
	return list.length
}

// Set sets a value in the list with key.
// If the key exists, change element value to the new one.
// Returns new element pointer.
func (list *SkipList) Set(key, value interface{}) *Element {
	var element *Element

	score := getScore(key, list.reversed)
	prevs := list.getPrevElementNodes(key, score)

	// found an element with the same key, replace its value
	if element = prevs[0].next[0]; element != nil && !list.keyFunc.Compare(element.key, key) {
		element.Value = value
		return element
	}

	element = &Element{
		elementNode: elementNode{
			next: make([]*Element, list.randLevel()),
		},
		key:   key,
		score: score,
		Value: value,
	}

	for i := range element.next {
		element.next[i] = prevs[i].next[i]
		prevs[i].next[i] = element
	}

	list.length++
	return element
}

// Get returns an element.
// Returns element pointer if found, nil if not found.
func (list *SkipList) Get(key interface{}) *Element {
	prev := &list.elementNode
	var next *Element
	score := getScore(key, list.reversed)

	for i := list.level - 1; i >= 0; i-- {
		next = prev.next[i]

		for next != nil &&
			(score > next.score || (score == next.score && list.keyFunc.Compare(key, next.key))) {
			prev = &next.elementNode
			next = next.next[i]
		}
	}

	if next != nil && score == next.score && !list.keyFunc.Compare(next.key, key) {
		return next
	}

	return nil
}

// GetValue returns a value. It's a short hand for Get().Value.
// Returns value and its existence status.
func (list *SkipList) GetValue(key interface{}) (interface{}, bool) {
	element := list.Get(key)

	if element == nil {
		return nil, false
	}

	return element.Value, true
}

// MustGetValue returns a value. It will panic if key doesn't exist.
// Returns value.
func (list *SkipList) MustGetValue(key interface{}) interface{} {
	element := list.Get(key)

	if element == nil {
		panic("cannot find key in skiplist")
	}

	return element.Value
}

// Remove removes an element.
// Returns removed element pointer if found, nil if not found.
func (list *SkipList) Remove(key interface{}) *Element {
	score := getScore(key, list.reversed)
	prevs := list.getPrevElementNodes(key, score)

	// found the element, remove it
	if element := prevs[0].next[0]; element != nil && !list.keyFunc.Compare(element.key, key) {
		for k, v := range element.next {
			prevs[k].next[k] = v
		}

		list.length--
		return element
	}

	return nil
}

func (list *SkipList) getPrevElementNodes(key interface{}, score float64) []*elementNode {
	prev := &list.elementNode
	var next *Element

	prevs := list.prevNodesCache

	for i := list.level - 1; i >= 0; i-- {
		next = prev.next[i]

		for next != nil &&
			(score > next.score || (score == next.score && list.keyFunc.Compare(key, next.key))) {
			prev = &next.elementNode
			next = next.next[i]
		}

		prevs[i] = prev
	}

	return prevs
}

// MaxLevel returns current max level value.
func (list *SkipList) MaxLevel() int {
	return list.level
}

// SetMaxLevel changes skip list max level.
// If level is not greater than 0, just panic.
func (list *SkipList) SetMaxLevel(level int) (old int) {
	if level <= 0 {
		panic("invalid argument to SetLevel")
	}

	old, list.level = list.level, level

	if old == level {
		return
	}

	if old > level {
		list.next = list.next[:level]
		list.prevNodesCache = list.prevNodesCache[:level]
		return
	}

	next := make([]*Element, level)
	copy(next, list.next)
	list.next = next
	list.prevNodesCache = make([]*elementNode, level)

	return
}

func (list *SkipList) randLevel() int {
	l := 1

	for ((list.randSource.Int63() >> 32) & 0xFFFF) < Propability {
		l++
	}

	if l > list.level {
		l = list.level
	}

	return l
}

func getScore(key interface{}, reversed bool) (score float64) {
	switch t := key.(type) {
	case []byte:
		var result uint64
		data := []byte(t)
		l := len(data)

		// only use first 8 bytes
		if l > 8 {
			l = 8
		}

		for i := 0; i < l; i++ {
			result |= uint64(data[i]) << uint(8*(7-i))
		}

		score = float64(result)

	case float32:
		score = float64(t)

	case float64:
		score = t

	case int:
		score = float64(t)

	case int16:
		score = float64(t)

	case int32:
		score = float64(t)

	case int64:
		score = float64(t)

	case int8:
		score = float64(t)

	case string:
		var result uint64
		data := string(t)
		l := len(data)

		// only use first 8 bytes
		if l > 8 {
			l = 8
		}

		for i := 0; i < l; i++ {
			result |= uint64(data[i]) << uint(8*(7-i))
		}

		score = float64(result)

	case uint:
		score = float64(t)

	case uint16:
		score = float64(t)

	case uint32:
		score = float64(t)

	case uint64:
		score = float64(t)

	case uint8:
		score = float64(t)

	case uintptr:
		score = float64(t)

	case Scorable:
		score = t.Score()
	}

	if reversed {
		score = -score
	}

	return
}

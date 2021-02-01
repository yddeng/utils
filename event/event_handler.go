package event

import (
	"container/list"
	"fmt"
	"reflect"
	"sync"
)

type Handler func(e interface{})
type Element *list.Element

type EventGroup struct {
	mtx         sync.Mutex
	eventType   reflect.Type
	handlerList *list.List
	eventChan   chan *event
}

type event struct {
	fn  func(event interface{})
	arg interface{}
}

func NewEventGroup(e interface{}) *EventGroup {
	return &EventGroup{
		eventType:   reflect.TypeOf(e),
		handlerList: list.New(),
		eventChan:   make(chan *event, 128),
	}
}

func (g *EventGroup) Handle(h Handler) Element {
	g.mtx.Lock()
	defer g.mtx.Unlock()
	return g.handlerList.PushBack(h)
}

func (g *EventGroup) Remove(e Element) {
	g.mtx.Lock()
	defer g.mtx.Unlock()
	g.handlerList.Remove(e)
}

func (g *EventGroup) Clear() {
	g.mtx.Lock()
	defer g.mtx.Unlock()
	g.handlerList.Init()
}

func (g *EventGroup) Execute(arg interface{}) error {
	if g.eventType != reflect.TypeOf(arg) {
		return fmt.Errorf("type error")
	}

	g.mtx.Lock()
	isExecute := len(g.eventChan) != 0
	for el := g.handlerList.Front(); el != nil; el = el.Next() {
		h := el.Value.(Handler)
		g.eventChan <- &event{fn: h, arg: arg}
	}
	g.mtx.Unlock()

	if !isExecute {
		for len(g.eventChan) > 0 {
			e := <-g.eventChan
			e.fn(e.arg)
		}
	}
	return nil
}

type EventHandler struct {
	groups map[reflect.Type]*EventGroup
}

func NewEventHandler() *EventHandler {
	return &EventHandler{groups: map[reflect.Type]*EventGroup{}}
}

func (h *EventHandler) Handle(event interface{}, handler Handler) Element {
	eventType := reflect.TypeOf(event)
	g, ok := h.groups[eventType]
	if !ok {
		g = NewEventGroup(event)
		h.groups[eventType] = g
	}
	return g.Handle(handler)
}

func (h *EventHandler) Execute(event interface{}) {
	eventType := reflect.TypeOf(event)
	if g, ok := h.groups[eventType]; ok {
		g.Execute(event)
	}
}

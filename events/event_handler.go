package events

import (
	"container/list"
	"reflect"
)

type eventGroup struct {
	handlers *list.List
}

type Handle struct {
	group *eventGroup
	elem  *list.Element
	fn    func(event interface{})
	once  bool
}

func (h *Handle) Release() {
	if h.group != nil {
		h.group.remove(h)
		h.elem = nil
		h.group = nil
	}
}

func newGroup() *eventGroup {
	return &eventGroup{
		handlers: list.New(),
	}
}

func (g *eventGroup) listen(f func(event interface{}), once bool) *Handle {
	h := &Handle{
		group: g,
		fn:    f,
		once:  once,
	}
	h.elem = g.handlers.PushBack(h)
	return h
}

func (g *eventGroup) trigger(arg interface{}) {
	count := g.handlers.Len()
	for i := 0; i < count; i++ {
		e := g.handlers.Front()
		h := e.Value.(*Handle)
		g.execute(arg, h.fn)
		if !h.once {
			g.handlers.MoveToBack(e)
		} else {
			g.handlers.Remove(e)
		}
	}
}

func (g *eventGroup) execute(event interface{}, h func(event interface{})) {
	h(event)
}

func (g *eventGroup) remove(h *Handle) {
	g.handlers.Remove(h.elem)
}

type EventHandler struct {
	groups map[reflect.Type]*eventGroup
}

func NewEventHandler() *EventHandler {
	return &EventHandler{groups: map[reflect.Type]*eventGroup{}}
}

func (eh *EventHandler) Listen(event interface{}, handler func(event interface{}), once bool) *Handle {
	t := reflect.TypeOf(event)
	g, ok := eh.groups[t]
	if !ok {
		g = newGroup()
		eh.groups[t] = g
	}
	return g.listen(handler, once)
}

func (eh *EventHandler) Trigger(event interface{}) {
	t := reflect.TypeOf(event)
	if g, ok := eh.groups[t]; ok {
		g.trigger(event)
	}
}

func (eh *EventHandler) Clear(event interface{}) {
	delete(eh.groups, reflect.TypeOf(event))
}

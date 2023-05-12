package forddb

import (
	"reflect"
	"sync"
)

type Listener interface {
	OnResourceChanged(id BasicResourceID, previous, current BasicResource)
}

type ListenerFunc func(id BasicResourceID, previous, current BasicResource)

func (f ListenerFunc) OnResourceChanged(id BasicResourceID, previous, current BasicResource) {
	f(id, previous, current)
}

type TypedListenerFunc[ID ResourceID[T], T Resource[ID]] func(id ID, previous, current T)

func (t TypedListenerFunc[ID, T]) OnResourceChanged(id BasicResourceID, previous, current BasicResource) {
	myType := reflect.TypeOf((*T)(nil)).Elem()
	myTypeId := TypeSystem().LookupByResourceType(myType).GetResourceID()

	if previous != nil && previous.GetResourceTypeID() != myTypeId {
		return
	}

	if current != nil && current.GetResourceTypeID() != myTypeId {
		return
	}

	var previousT, currentT T

	if previous != nil {
		previousT = previous.(T)
	}

	if current != nil {
		currentT = current.(T)
	}

	t(id.(ID), previousT, currentT)
}

type HasListeners interface {
	Subscribe(listener Listener) func()
}

type HasListenersBase struct {
	listeners ListenerSet
}

func (h *HasListenersBase) Subscribe(listener Listener) func() {
	return h.listeners.Subscribe(listener)
}

func FireListeners(hlb *HasListenersBase, id BasicResourceID, previous, current BasicResource) {
	hlb.listeners.OnResourceChanged(id, previous, current)
}

type ListenerSet struct {
	m         sync.RWMutex
	listeners map[uint64]Listener
	counter   uint64
}

func (l *ListenerSet) Subscribe(listener Listener) func() {
	l.m.Lock()
	defer l.m.Unlock()

	l.counter++

	id := l.counter

	if l.listeners == nil {
		l.listeners = map[uint64]Listener{}
	}

	l.listeners[id] = listener

	return func() {
		l.m.Lock()
		defer l.m.Unlock()

		delete(l.listeners, id)
	}
}

func (l *ListenerSet) OnResourceChanged(id BasicResourceID, previous, current BasicResource) {
	l.m.RLock()
	defer l.m.RUnlock()

	if l.listeners == nil {
		return
	}

	for _, l := range l.listeners {
		l.OnResourceChanged(id, previous, current)
	}
}

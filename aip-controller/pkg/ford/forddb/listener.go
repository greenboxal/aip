package forddb

import (
	"reflect"
	"sync"

	"golang.org/x/exp/slices"
)

type Listener interface {
	OnResourceChanged(id BasicResourceID, previous, current BasicResource)
}

type ListenerFunc func(resource BasicResource)

func (f ListenerFunc) OnResourceChanged(resource BasicResource) {
	f(resource)
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
	AddListener(listener Listener)
	RemoveListener(listener Listener)
}

type HasListenersBase struct {
	listeners ListenerSet
}

func (h *HasListenersBase) AddListener(listener Listener)    { h.listeners.AddListener(listener) }
func (h *HasListenersBase) RemoveListener(listener Listener) { h.listeners.RemoveListener(listener) }

func FireListeners(hlb *HasListenersBase, id BasicResourceID, previous, current BasicResource) {
	hlb.listeners.OnResourceChanged(id, previous, current)
}

type ListenerSet struct {
	m         sync.RWMutex
	listeners []Listener
}

func (l *ListenerSet) AddListener(listener Listener) {
	l.m.Lock()
	defer l.m.Unlock()

	index := slices.Index(l.listeners, listener)

	if index != -1 {
		return
	}

	l.listeners = append(l.listeners, listener)
}

func (l *ListenerSet) RemoveListener(listener Listener) {
	l.m.Lock()
	defer l.m.Unlock()

	index := slices.Index(l.listeners, listener)

	if index == -1 {
		return
	}

	l.listeners = slices.Delete(l.listeners, index, index+1)
}

func (l *ListenerSet) OnResourceChanged(id BasicResourceID, previous, current BasicResource) {
	l.m.RLock()
	defer l.m.RUnlock()

	for _, l := range l.listeners {
		l.OnResourceChanged(id, previous, current)
	}
}

package forddb

import (
	"reflect"
	"sync"

	"golang.org/x/exp/slices"
)

type Listener interface {
	OnResourceChanged(resource BasicResource)
}

type ListenerFunc func(resource BasicResource)

func (f ListenerFunc) OnResourceChanged(resource BasicResource) {
	f(resource)
}

type TypedListenerFunc[ID ResourceID[T], T Resource[ID]] func(resource T)

func (t TypedListenerFunc[ID, T]) OnResourceChanged(resource BasicResource) {
	if resource == nil {
		return
	}

	myType := reflect.TypeOf((*T)(nil)).Elem()
	myTypeId := typeSystem.LookupByResourceType(myType).ID()

	if resource.GetType() != myTypeId {
		return
	}

	t(resource.(T))
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

func (l *ListenerSet) OnResourceChanged(resource BasicResource) {
	l.m.RLock()
	defer l.m.RUnlock()

	for _, l := range l.listeners {
		l.OnResourceChanged(resource)
	}
}

type HasListeners interface {
	AddListener(listener Listener)
	RemoveListener(listener Listener)
}

type HasListenersBase struct {
	listeners ListenerSet
}

func FireListeners(hlb *HasListenersBase, resource BasicResource) {
	hlb.listeners.OnResourceChanged(resource)
}

func (h *HasListenersBase) AddListener(listener Listener) {
	h.listeners.AddListener(listener)
}

func (h *HasListenersBase) RemoveListener(listener Listener) {
	h.listeners.RemoveListener(listener)
}

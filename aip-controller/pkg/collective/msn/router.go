package msn

import (
	"sync"
)

type Router struct {
	m                    sync.RWMutex
	channelSubscriptions map[ChannelID]*subscriptionTable
}

func NewRouter() *Router {
	return &Router{
		channelSubscriptions: map[ChannelID]*subscriptionTable{},
	}
}

func (r *Router) Dispatch(ev Event) {
	table := r.getTable(ev.MessageEvent.Message.Channel, false)

	if table != nil {
		table.Dispatch(ev)
	}
}

func (r *Router) getTable(endpoint ChannelID, create bool) *subscriptionTable {
	r.m.Lock()
	defer r.m.Unlock()

	if st, ok := r.channelSubscriptions[endpoint]; ok {
		return st
	}

	if !create {
		return nil
	}

	st := newSubscriptionTable()

	r.channelSubscriptions[endpoint] = st

	return st
}

type subscriptionTable struct {
	m                   sync.RWMutex
	subscriptions       map[uint64]chan<- Event
	subscriptionCounter uint64
}

func newSubscriptionTable() *subscriptionTable {
	return &subscriptionTable{
		subscriptions: map[uint64]chan<- Event{},
	}
}

func (st *subscriptionTable) Dispatch(ev Event) {
	st.m.RLock()
	defer st.m.RUnlock()

	for _, ch := range st.subscriptions {
		ch <- ev
	}
}

func (st *subscriptionTable) Subscribe(ch chan<- Event) func() {
	st.m.Lock()
	defer st.m.Unlock()

	st.subscriptionCounter++

	subscriptionId := st.subscriptionCounter

	st.subscriptions[subscriptionId] = ch

	return func() {
		st.m.Lock()
		defer st.m.Unlock()

		delete(st.subscriptions, subscriptionId)
	}
}

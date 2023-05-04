package comms

import (
	"sync"

	"github.com/zyedidia/generic/mapset"
)

type Room struct {
	m sync.RWMutex

	name    string
	members mapset.Set[string]
}

func NewRoom(name string) *Room {
	return &Room{
		name:    name,
		members: mapset.New[string](),
	}
}

func (r *Room) Join(port string) {
	r.m.Lock()
	defer r.m.Unlock()

	r.members.Put(port)
}

func (r *Room) Leave(port string) {
	r.m.Lock()
	defer r.m.Unlock()

	r.members.Remove(port)
}

func (r *Room) Name() string {
	return r.name
}

func (r *Room) Members() []string {
	result := make([]string, 0, r.members.Size())

	r.members.Each(func(key string) {
		result = append(result, key)
	})

	return result
}

package comms

import (
	"sync"

	"github.com/greenboxal/aip/pkg/daemon"
)

type Manager struct {
	m     sync.RWMutex
	rooms map[string]*Room

	routing *daemon.Routing
}

func NewManager(routing *daemon.Routing) *Manager {
	return &Manager{
		rooms: map[string]*Room{},

		routing: routing,
	}
}

func (m *Manager) CreateRoom(name string) *Room {
	m.m.Lock()
	defer m.m.Unlock()

	if exising := m.rooms[name]; exising != nil {
		return exising
	}

	room := NewRoom(name)

	m.rooms[name] = room

	return room
}

func (m *Manager) GetRoom(name string) *Room {
	m.m.RLock()
	defer m.m.RUnlock()

	return m.rooms[name]
}

func (m *Manager) JoinRoom(room, port string) *Room {
	r := m.CreateRoom(room)

	r.Join(port)

	return r
}

func (m *Manager) LeaveRoom(room, port string) *Room {
	r := m.CreateRoom(room)

	r.Leave(port)

	return r
}

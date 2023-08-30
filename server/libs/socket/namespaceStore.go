package socket

import "sync"

type NamespaceStore interface {
	Add(name string, roomName string, connect Socket)
	Get(name string) RoomStore
	Remove(nsp string, room string)
}

type namespaceStore struct {
	store map[string]RoomStore
	lock  sync.RWMutex
}

func (s *namespaceStore) Add(name string, roomName string, connect Socket) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	nsp, ok := s.store[name]
	if ok {
		nsp.Add(roomName, connect)
	} else {
		s.store[name] = newRoomStore()
		s.store[name].Add(roomName, connect)
	}
}

func (s *namespaceStore) Get(name string) RoomStore {
	c, ok := s.store[name]

	if ok {
		return c
	}
	return nil
}

func (s *namespaceStore) Remove(nsp string, room string) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	rs, ok := s.store[nsp]
	if ok {
		rs.Remove(room)
	}
}

func NewNamespaceStore() NamespaceStore {
	store := make(map[string]RoomStore)
	return &namespaceStore{
		store: store,
	}
}

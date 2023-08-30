package socket

import (
	"github.com/google/uuid"
	"sync"
)

var DefaultRoom = ""

type RoomStore interface {
	Add(name string, socket Socket)
	Get(name string) SocketStore
	Remove(name string)
	LeaveRoom(sk Socket, room string)
	GetSocket(socketId uuid.UUID) Socket
}

type roomStore struct {
	store map[string]SocketStore
	lock  sync.RWMutex
}

func (s *roomStore) Add(name string, socket Socket) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	defaultRoomStore, ok := s.store[DefaultRoom]
	if ok {
		defaultRoomStore.Add(socket)
	}

	room, ok := s.store[name]
	if ok {
		room.Add(socket)
	} else {
		s.store[name] = NewSocketStore()
		s.store[name].Add(socket)
	}
}

func (s *roomStore) Get(name string) SocketStore {
	skStore, ok := s.store[name]
	if ok {
		return skStore
	}
	return nil
}

func (s *roomStore) Remove(name string) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	delete(s.store, name)
}

func (s *roomStore) LeaveRoom(sk Socket, room string) {
	rs, ok := s.store[room]
	if ok {
		rs.Remove(sk)
	}
}

func (s *roomStore) GetSocket(socketId uuid.UUID) Socket {
	defaultRoomStore, ok := s.store[DefaultRoom]
	if !ok {
		return nil
	}

	sk, err := defaultRoomStore.Get(socketId)

	if err != nil {
		return nil
	}

	return sk
}

func newRoomStore() RoomStore {
	store := make(map[string]SocketStore)
	store[DefaultRoom] = NewSocketStore()
	return &roomStore{
		store: store,
	}
}

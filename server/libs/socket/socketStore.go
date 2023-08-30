package socket

import (
	"errors"
	"github.com/google/uuid"
	"sync"
)

type SocketMap map[uuid.UUID]Socket

type SocketStore interface {
	Add(socket Socket)
	Get(socketId uuid.UUID) (Socket, error)
	RemoveID(socketId uuid.UUID)
	Remove(socket Socket)
	GetAll() SocketMap
}

type socketStore struct {
	sockets SocketMap
	lock    sync.RWMutex
}

func (s *socketStore) Add(socket Socket) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	s.sockets[socket.GetId()] = socket
}

func (s *socketStore) Get(socketId uuid.UUID) (Socket, error) {
	socket, ok := s.sockets[socketId]
	if !ok {
		return nil, errors.New("not_found")
	}

	return socket, nil
}

func (s *socketStore) RemoveID(socketId uuid.UUID) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	delete(s.sockets, socketId)
}

func (s *socketStore) Remove(socket Socket) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	delete(s.sockets, socket.GetId())
}

func (s *socketStore) GetAll() SocketMap {
	return s.sockets
}

func NewSocketStore() SocketStore {
	sockets := make(SocketMap)
	return &socketStore{
		sockets: sockets,
	}
}

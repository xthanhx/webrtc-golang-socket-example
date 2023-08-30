package socket

import (
	"github.com/google/uuid"
)

var DefaultNamespace = "/"

type SocketSet map[uuid.UUID]struct{}

type NamespaceBroadcast interface {
	Of(path string) NamespaceBroadcast
	ToRoom(room string) NamespaceBroadcast
	ToRooms(room []string) NamespaceBroadcast
	Emit(event string, message any)
	RemoveRoom(room string)
	GetSocketSelected() SocketMap
	To(socketID uuid.UUID) NamespaceBroadcast
}

func NewNamespaceBroadcast(store NamespaceStore) NamespaceBroadcast {
	roomSet := make(RoomSet)
	socketSet := make(SocketSet)
	return &namespaceBroadcast{
		namespaceStore: store,
		namespace:      DefaultNamespace,
		roomSet:        roomSet,
		socketSet:      socketSet,
	}
}

type namespaceBroadcast struct {
	namespaceStore NamespaceStore ``
	namespace      string
	roomSet        RoomSet
	socketSet      SocketSet
}

func (n *namespaceBroadcast) RemoveRoom(room string) {
	n.namespaceStore.Get(n.namespace).Remove(room)
}

func (n *namespaceBroadcast) Of(path string) NamespaceBroadcast {
	roomSet := make(RoomSet)
	socketSet := make(SocketSet)
	roomSet[DefaultRoom] = empty

	return &namespaceBroadcast{
		namespaceStore: n.namespaceStore,
		namespace:      path,
		roomSet:        roomSet,
		socketSet:      socketSet,
	}
}

func (n *namespaceBroadcast) ToRoom(room string) NamespaceBroadcast {
	roomSet := copyRoomSet(n.roomSet)
	roomSet[room] = empty

	return &namespaceBroadcast{
		namespaceStore: n.namespaceStore,
		namespace:      n.namespace,
		roomSet:        roomSet,
		socketSet:      n.socketSet,
	}
}

func (n *namespaceBroadcast) ToRooms(rooms []string) NamespaceBroadcast {
	roomSet := copyRoomSet(n.roomSet)

	for _, room := range rooms {
		roomSet[room] = empty
	}

	return &namespaceBroadcast{
		namespaceStore: n.namespaceStore,
		namespace:      n.namespace,
		roomSet:        roomSet,
		socketSet:      n.socketSet,
	}
}

func (n *namespaceBroadcast) Emit(event string, message any) {
	skMap := n.GetSocketSelected()

	if skMap == nil {
		return
	}

	for _, sk := range skMap {
		sk.Emit(event, message)
	}
}

func (n *namespaceBroadcast) GetSocketSelected() SocketMap {
	rStore := n.namespaceStore.Get(n.namespace)
	if rStore == nil {
		return nil
	}

	skStore := NewSocketStore()
	for room, _ := range n.roomSet {
		sockets := rStore.Get(room)
		if sockets == nil {
			return nil
		}

		for _, sk := range sockets.GetAll() {
			skStore.Add(sk)
		}
	}

	for skId, _ := range n.socketSet {
		sk := rStore.GetSocket(skId)
		if sk == nil {
			continue
		}

		skStore.Add(sk)
	}

	return skStore.GetAll()
}

func (n *namespaceBroadcast) To(socketID uuid.UUID) NamespaceBroadcast {
	socketSet := copySocketSet(n.socketSet)
	socketSet[socketID] = empty

	return &namespaceBroadcast{
		namespaceStore: n.namespaceStore,
		namespace:      n.namespace,
		roomSet:        n.roomSet,
		socketSet:      socketSet,
	}
}

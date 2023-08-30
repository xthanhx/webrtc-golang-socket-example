package socket

import (
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"net/http"
)

type MessageContent []byte

type Connect interface {
	GetConnect() *websocket.Conn
	GetId() uuid.UUID
	ReadMessage() (messageType int, messageContent []byte, err error)
	WriteMessage(messageType int, messageContent MessageContent) error
	Close()
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type connectHandle struct {
	conn *websocket.Conn
	id   uuid.UUID
}

func (c *connectHandle) GetConnect() *websocket.Conn {
	return c.conn
}

func (c *connectHandle) GetId() uuid.UUID {
	return c.id
}

func (c *connectHandle) Close() {
	c.conn.Close()
}

func (c *connectHandle) ReadMessage() (messageType int, messageContent []byte, err error) {
	return c.conn.ReadMessage()
}

func (c *connectHandle) WriteMessage(messageType int, content MessageContent) error {
	return c.conn.WriteMessage(messageType, content)
}

func NewSocketConnect(w http.ResponseWriter, r *http.Request) (Connect, error) {
	connect, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return nil, err
	}

	return &connectHandle{
		conn: connect,
		id:   uuid.New(),
	}, nil
}

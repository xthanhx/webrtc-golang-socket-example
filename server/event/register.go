package event

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	SK "gitlab.com/circa/services/signal/libs/socket"
	"log"
)

func RegisterEvent(socketServer SK.Server) {
	logger := log.Default()

	socketServer.Register("/signal", func(socket SK.Socket) {
		socket.On("connected", func(_ any) {
			socket.Emit("send_connect_id", socket.GetId())
		})

		socket.On("connect_signal", func(payload any) {
			sks := socketServer.Broadcast().Of("/test").GetSocketSelected()
			for _, sk := range sks {
				if sk.GetId() != socket.GetId() {
					sk.Emit("signal", payload)
				}
			}
		})
	})
}

package helios

import (
	"fmt"

	"github.com/googollee/go-socket.io"
	log "gopkg.in/inconshreveable/log15.v2"
)

const socketRoom = "helios"

func initSocket() *socketio.Server {
	server, err := socketio.NewServer(nil)
	if err != nil {
		fmt.Println("Failed to start socket server")
		return nil
	}

	server.On("connection", func(so socketio.Socket) {
		fmt.Printf("New socket.io connection: %s", so.Id())
		so.Join(socketRoom)
		so.On("disconnection", func() {
			// no op
		})
	})

	server.On("error", func(so socketio.Socket, err error) {
		log.Error("Error on socket.io server", "error", err.Error())
	})

	return server
}

func (h *Engine) NewBroadcastChannel(message string) chan interface{} {
	chReceiver := make(chan interface{})
	go func() {
		for {
			msg := <-chReceiver
			h.Info("Got message to broadcast", "socket", message)
			h.Socket.BroadcastTo(socketRoom, message, msg)
		}
	}()
	return chReceiver
}

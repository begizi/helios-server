package main

import (
	"fmt"
	"log"

	"github.com/google/go-github/github"
	"github.com/googollee/go-socket.io"
)

func setupSocketIO() (*socketio.Server, error) {
	server, err := socketio.NewServer(nil)
	if err != nil {
		return nil, err
	}

	server.On("connection", func(so socketio.Socket) {
		fmt.Printf("New socket.io connection: %s", so.Id())
		so.Join("helios")
		so.On("disconnection", func() {
			// no op
		})
	})

	server.On("error", func(so socketio.Socket, err error) {
		log.Fatalf("Error on socket.io server", err.Error())
	})

	return server, nil
}

func startSocketPusher(s *socketio.Server, c <-chan []github.Event) error {
	go func() {
		for {
			events := <-c
			s.BroadcastTo("feedbag", "activity", events)
		}
	}()

	return nil
}

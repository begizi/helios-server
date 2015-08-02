package helios

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/googollee/go-socket.io"
)

type Engine struct {
	HTTPEngine *gin.Engine
	Socket     *socketio.Server
	SocketChan chan interface{}
	services   []ServiceHandler
}

type ServiceHandler func(*Engine) error

func New() *Engine {
	// package instance of the helios type
	var server = &Engine{
		HTTPEngine: gin.Default(),
		Socket:     initSocket(),
		SocketChan: make(chan interface{}),
	}

	return server
}

func (h *Engine) Use(mw ServiceHandler) {
	h.services = append(h.services, mw)
}

func (h *Engine) Run(port string) {
	// Start services services
	h.startServices()

	// Start engine now that all services have loaded
	h.HTTPEngine.Run(":" + port)
}

func (h *Engine) startServices() {
	for _, mw := range h.services {
		err := mw(h)
		if err != nil {
			fmt.Println("Failed to start service: ", err)
		}
	}
}

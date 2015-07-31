package helios

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/googollee/go-socket.io"
	"github.com/tommy351/gin-cors"
)

type Engine struct {
	HTTPEngine *gin.Engine
	Socket     *socketio.Server
	SocketChan chan interface{}
	middleware []MiddlewareFunc
}

type MiddlewareFunc func(*Engine) error

func New() *Engine {
	// package instance of the helios type
	var server = &Engine{
		Socket:     initSocket(),
		SocketChan: make(chan interface{}),
	}

	return server
}

func (h *Engine) Use(mw MiddlewareFunc) {
	h.middleware = append(h.middleware, mw)
}

func (h *Engine) Run(port string) {
	// Setup HTTP Engine
	h.HTTPEngine = initHTTPEngine(port)

	// Start middleware services
	h.startMiddleware()

	// Start engine now that all middelware have loaded
	h.HTTPEngine.Run(":" + port)
}

func (h *Engine) startMiddleware() {
	for _, mw := range h.middleware {
		err := mw(h)
		if err != nil {
			fmt.Println("Failed to start middleware: ", err)
		}
	}
}

func initHTTPEngine(port string) *gin.Engine {
	// Create Engine Instance
	r := gin.Default()

	// Server settings
	r.RedirectTrailingSlash = true
	r.RedirectFixedPath = true
	r.HandleMethodNotAllowed = true

	// Middleware
	r.Use(cors.Middleware(cors.Options{AllowCredentials: true}))

	return r
}

package helios

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/googollee/go-socket.io"
	log "gopkg.in/inconshreveable/log15.v2"
)

type Engine struct {
	HTTPEngine *gin.Engine
	Socket     *socketio.Server
	services   []ServiceHandler
	log.Logger
}

type ServiceHandler func(*Engine) error

func New() *Engine {
	// package instance of the helios type
	engine := &Engine{
		HTTPEngine: gin.Default(),
		Socket:     initSocket(),
		Logger:     log.New(),
	}

	fileHandler, _ := log.FileHandler("./log.debug", log.LogfmtFormat())
	engine.SetHandler(log.MultiHandler(log.LvlFilterHandler(log.LvlWarn, fileHandler), log.StreamHandler(os.Stdout, log.TerminalFormat())))

	return engine
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
			h.Warn("Failed to start service: ", err)
		}
	}
}

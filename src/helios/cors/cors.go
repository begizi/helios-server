package cors

import (
	"helios/helios"

	"github.com/tommy351/gin-cors"
)

func Plugin() helios.MiddlewareFunc {
	return func(h *helios.Engine) error {

		h.HTTPEngine.Use(cors.Middleware(cors.Options{AllowCredentials: true}))

		return nil
	}
}

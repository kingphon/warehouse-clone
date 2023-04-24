package route

import (
	"github.com/labstack/echo/v4"

	"git.selly.red/Selly-Server/warehouse/pkg/app/handler"
)

// common ...
func common(e *echo.Group) {
	var (
		g = e.Group("")
		h = handler.Common{}
	)

	// Ping
	g.GET("/ping", h.Ping)
}

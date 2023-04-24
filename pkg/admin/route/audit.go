package route

import (
	"github.com/labstack/echo/v4"

	"git.selly.red/Selly-Server/warehouse/pkg/admin/handler"
	routeauth "git.selly.red/Selly-Server/warehouse/pkg/admin/route/auth"
	routevalidation "git.selly.red/Selly-Server/warehouse/pkg/admin/route/validation"
)

type Audit struct{}

// audit ...
func audit(e *echo.Group) {
	var (
		g = e.Group("/audits", routeauth.RequiredLogin)
		h = handler.Audit{}
		v = routevalidation.Audit{}
	)

	// All
	g.GET("", h.All, v.GetAll)

}

package route

import (
	"github.com/labstack/echo/v4"

	routevalidation "git.selly.red/Selly-Server/warehouse/pkg/app/route/validation"

	"git.selly.red/Selly-Server/warehouse/pkg/app/handler"
)

type Warehouse struct{}

// warehouse ...
func warehouse(e *echo.Group) {
	var (
		g = e.Group("/warehouses")
		h = handler.Warehouse{}
		v = routevalidation.Warehouse{}
	)

	// Detail
	g.GET("/:id", h.Detail, v.Detail)

}

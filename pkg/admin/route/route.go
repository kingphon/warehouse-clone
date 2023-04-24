package route

import (
	"github.com/labstack/echo/v4"

	"git.selly.red/Selly-Server/warehouse/external/utils/routemiddleware"
	routeauth "git.selly.red/Selly-Server/warehouse/pkg/admin/route/auth"
)

// Init ...
func Init(e *echo.Echo) {

	// Middlewares ...
	e.Use(routeauth.Jwt())

	e.Use(routemiddleware.CORSConfig())
	e.Use(routemiddleware.Locale)

	r := e.Group("/admin/warehouse")

	// Components

	common(r)
	audit(r)
	warehouse(r)
	warehouseConfiguration(r)
	migrate(r)
	outboundRequest(r)
	outboundRequestHistory(r)
	supplierHoliday(r)
}

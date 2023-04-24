package route

import (
	"github.com/labstack/echo/v4"

	"git.selly.red/Selly-Server/warehouse/external/auth/permission"
	"git.selly.red/Selly-Server/warehouse/pkg/admin/handler"
	routeauth "git.selly.red/Selly-Server/warehouse/pkg/admin/route/auth"
	routevalidation "git.selly.red/Selly-Server/warehouse/pkg/admin/route/validation"
)

func outboundRequest(e *echo.Group) {
	var (
		g = e.Group("/outbound-requests", routeauth.RequiredLogin)
		h = handler.OutboundRequest{}
		v = routevalidation.OutboundRequest{}
	)

	// Permission
	view := routeauth.CheckPermission(permission.Warehouse.OutboundRequest.View)

	g.GET("", h.GetList, view, v.List)
}

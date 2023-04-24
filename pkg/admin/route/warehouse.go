package route

import (
	"github.com/labstack/echo/v4"

	"git.selly.red/Selly-Server/warehouse/external/auth/permission"
	"git.selly.red/Selly-Server/warehouse/pkg/admin/handler"
	routeauth "git.selly.red/Selly-Server/warehouse/pkg/admin/route/auth"
	routevalidation "git.selly.red/Selly-Server/warehouse/pkg/admin/route/validation"
)

type Warehouse struct{}

// warehouse ...
func warehouse(e *echo.Group) {
	var (
		g = e.Group("/warehouses", routeauth.RequiredLogin)
		h = handler.Warehouse{}
		v = routevalidation.Warehouse{}
	)

	// Permission
	edit := routeauth.CheckPermission(permission.Warehouse.Self.Edit)
	view := routeauth.CheckPermission(permission.Warehouse.Self.View)

	// Create
	g.POST("", h.Create, edit, v.Create)

	// All
	g.GET("", h.All, view, v.All)

	// Detail
	g.GET("/:id", h.Detail, view, v.Detail)

	// Update
	g.PUT("/:id", h.Update, edit, v.Update)

	// UpdateStatus
	g.PATCH("/:id/status", h.UpdateStatus, edit, v.UpdateStatus)

	// RunJobSetIsClosed ...
	g.GET("/run-job-set-is-closed", h.RunJobSetIsClosed)
}

package route

import (
	"git.selly.red/Selly-Server/warehouse/external/auth/permission"
	"git.selly.red/Selly-Server/warehouse/pkg/admin/handler"
	routeauth "git.selly.red/Selly-Server/warehouse/pkg/admin/route/auth"
	routevalidation "git.selly.red/Selly-Server/warehouse/pkg/admin/route/validation"
	"github.com/labstack/echo/v4"
)

// SupplierHoliday ...
type SupplierHoliday struct{}

func supplierHoliday(e *echo.Group) {
	var (
		g = e.Group("/supplier-holidays", routeauth.RequiredLogin)
		h = handler.SupplierHoliday{}
		v = routevalidation.SupplierHoliday{}
	)

	edit := routeauth.CheckPermission(permission.Warehouse.SupplierHoliday.Edit)
	view := routeauth.CheckPermission(permission.Warehouse.SupplierHoliday.View)

	// Create
	g.POST("", h.Create, edit, v.Create)

	// All
	g.GET("", h.All, view, v.All)

	// Detail
	g.GET("/:id", h.Detail, view, v.Detail)

	// Update
	g.PUT("/:id", h.Update, edit, v.Update, v.Detail)

	// Change status
	g.PATCH("/:id/status", h.ChangeStatus, edit, v.ChangeStatus, v.Detail)

	// JobUpdateHolidayWarehouse
	g.GET("/run-job-update-holiday-warehouses", h.RunJobUpdateHolidayWarehouses, edit)

	// RunJobUpdateHolidayStatusForSupplier
	g.GET("/run-job-update-holiday-status-for-supplier", h.RunJobUpdateHolidayStatusForSupplier, edit)
}

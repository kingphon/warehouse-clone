package route

import (
	"git.selly.red/Selly-Server/warehouse/external/auth/permission"
	"git.selly.red/Selly-Server/warehouse/pkg/app/handler"
	routeauth "git.selly.red/Selly-Server/warehouse/pkg/app/route/auth"
	routevalidation "git.selly.red/Selly-Server/warehouse/pkg/app/route/validation"
	"github.com/labstack/echo/v4"
)

// SupplierHoliday ...
type SupplierHoliday struct{}

func supplierHoliday(e *echo.Group) {
	var (
		g = e.Group("/supplier-holidays")
		h = handler.SupplierHoliday{}
		v = routevalidation.SupplierHoliday{}
	)

	edit := routeauth.CheckTokenSupplier(permission.Warehouse.SupplierHolidayMode.Edit)
	view := routeauth.CheckTokenSupplier(permission.Warehouse.SupplierHolidayMode.View)

	// Create
	g.POST("", h.Create, edit, v.Create)

	// Detail
	g.GET("/supplier", h.Detail, view)

	// Update
	g.PUT("", h.Update, edit, v.Update)

	// Change status
	g.PATCH("/change-status", h.ChangeStatus, edit, v.ChangeStatus)

}

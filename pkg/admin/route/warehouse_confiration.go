package route

import (
	"github.com/labstack/echo/v4"

	"git.selly.red/Selly-Server/warehouse/external/auth/permission"
	"git.selly.red/Selly-Server/warehouse/pkg/admin/handler"
	routeauth "git.selly.red/Selly-Server/warehouse/pkg/admin/route/auth"
	routevalidation "git.selly.red/Selly-Server/warehouse/pkg/admin/route/validation"
)

type WarehouseConfiguration struct{}

// warehouse ...
func warehouseConfiguration(e *echo.Group) {
	var (
		g = e.Group("/configurations", routeauth.RequiredLogin)
		h = handler.WarehouseConfiguration{}
		v = routevalidation.WarehouseConfiguration{}
	)

	// Permission
	view := routeauth.CheckPermission(permission.Warehouse.Configuration.View)
	edit := routeauth.CheckPermission(permission.Warehouse.Configuration.Edit)

	// Detail
	g.GET("/:id", h.Detail, view, v.Detail)

	// Update supplier
	g.PUT("/:id/supplier", h.UpdateSupplier, edit, v.UpdateSupplier)

	// Update food
	g.PUT("/:id/food", h.UpdateFood, edit, v.UpdateFood)

	// Update partner
	g.PUT("/:id/partner", h.UpdatePartner, edit, v.UpdatePartner)

	// Update order
	g.PUT("/:id/order", h.UpdateOrder, edit, v.UpdateOrder)

	// Update delivery
	g.PUT("/:id/delivery", h.UpdateDelivery, edit, v.UpdateDelivery)

	// Update other
	g.PUT("/:id/other", h.UpdateOther, edit, v.UpdateOther)

	// Update order confirm
	g.PUT("/:id/order-confirm", h.UpdateOrderConfirm, edit, v.UpdateOrderConfirm)
}

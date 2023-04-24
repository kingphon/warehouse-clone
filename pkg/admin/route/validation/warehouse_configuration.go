package routevalidation

import (
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"git.selly.red/Selly-Server/warehouse/external/response"
	"git.selly.red/Selly-Server/warehouse/external/utils/echocontext"
	"git.selly.red/Selly-Server/warehouse/pkg/admin/errorcode"
	requestmodel "git.selly.red/Selly-Server/warehouse/pkg/admin/model/request"
)

// WarehouseConfiguration ...
type WarehouseConfiguration struct {
}

// Detail ...
func (WarehouseConfiguration) Detail(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		var (
			id = c.Param("id")
		)

		if !primitive.IsValidObjectID(id) {
			return response.R404(c, nil, errorcode.WarehouseCfgInvalidID)
		}

		objID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return response.R400(c, nil, "")
		}

		echocontext.SetParam(c, "id", objID)
		return next(c)
	}
}

// UpdateFood ...
func (WarehouseConfiguration) UpdateFood(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		var (
			id      = c.Param("id")
			payload requestmodel.WarehouseCfgFoodUpdate
		)

		if !primitive.IsValidObjectID(id) {
			return response.R404(c, nil, errorcode.WarehouseCfgInvalidID)
		}

		objID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return response.R400(c, nil, "")
		}

		if err := c.Bind(&payload); err != nil {
			return response.R400(c, nil, "")
		}

		if err := payload.Validate(); err != nil {
			return response.RouteValidation(c, err)
		}

		echocontext.SetPayload(c, payload)
		echocontext.SetParam(c, "id", objID)
		return next(c)
	}
}

// UpdateSupplier ...
func (WarehouseConfiguration) UpdateSupplier(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		var (
			id      = c.Param("id")
			payload requestmodel.WarehouseCfgSupplierUpdate
		)

		if !primitive.IsValidObjectID(id) {
			return response.R404(c, nil, errorcode.WarehouseCfgInvalidID)
		}

		objID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return response.R400(c, nil, "")
		}

		if err := c.Bind(&payload); err != nil {
			return response.R400(c, nil, "")
		}

		if err := payload.Validate(); err != nil {
			return response.RouteValidation(c, err)
		}

		echocontext.SetPayload(c, payload)
		echocontext.SetParam(c, "id", objID)
		return next(c)
	}
}

// UpdatePartner ...
func (WarehouseConfiguration) UpdatePartner(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		var (
			id      = c.Param("id")
			payload requestmodel.WarehouseCfgPartnerUpdate
		)

		if !primitive.IsValidObjectID(id) {
			return response.R404(c, nil, errorcode.WarehouseCfgInvalidID)
		}

		objID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return response.R400(c, nil, "")
		}

		if err := c.Bind(&payload); err != nil {
			return response.R400(c, nil, "")
		}

		if err := payload.Validate(); err != nil {
			return response.RouteValidation(c, err)
		}

		echocontext.SetPayload(c, payload)
		echocontext.SetParam(c, "id", objID)
		return next(c)
	}
}

// UpdateOrder ...
func (WarehouseConfiguration) UpdateOrder(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		var (
			id      = c.Param("id")
			payload requestmodel.WarehouseCfgOrderUpdate
		)

		if !primitive.IsValidObjectID(id) {
			return response.R404(c, nil, errorcode.WarehouseCfgInvalidID)
		}

		objID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return response.R400(c, nil, "")
		}

		if err := c.Bind(&payload); err != nil {
			return response.R400(c, nil, "")
		}

		if err := payload.Validate(); err != nil {
			return response.RouteValidation(c, err)
		}

		echocontext.SetPayload(c, payload)
		echocontext.SetParam(c, "id", objID)
		return next(c)
	}
}

// UpdateDelivery ...
func (WarehouseConfiguration) UpdateDelivery(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		var (
			id      = c.Param("id")
			payload requestmodel.WarehouseCfgDeliveryUpdate
		)

		if !primitive.IsValidObjectID(id) {
			return response.R404(c, nil, errorcode.WarehouseCfgInvalidID)
		}

		objID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return response.R400(c, nil, "")
		}

		if err := c.Bind(&payload); err != nil {
			return response.R400(c, nil, "")
		}

		if err := payload.Validate(); err != nil {
			return response.RouteValidation(c, err)
		}

		echocontext.SetPayload(c, payload)
		echocontext.SetParam(c, "id", objID)
		return next(c)
	}
}

// UpdateOrderConfirm ...
func (WarehouseConfiguration) UpdateOrderConfirm(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		var (
			id      = c.Param("id")
			payload requestmodel.WarehouseCfgOrderConfirm
		)

		if !primitive.IsValidObjectID(id) {
			return response.R404(c, nil, errorcode.WarehouseCfgInvalidID)
		}

		objID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return response.R400(c, nil, "")
		}

		if err := c.Bind(&payload); err != nil {
			return response.R400(c, nil, "")
		}

		if err := payload.Validate(); err != nil {
			return response.RouteValidation(c, err)
		}

		echocontext.SetPayload(c, payload)
		echocontext.SetParam(c, "id", objID)
		return next(c)
	}
}

// UpdateOther ...
func (WarehouseConfiguration) UpdateOther(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		var (
			id      = c.Param("id")
			payload requestmodel.WarehouseCfgOtherUpdate
		)

		if !primitive.IsValidObjectID(id) {
			return response.R404(c, nil, errorcode.WarehouseCfgInvalidID)
		}

		objID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return response.R400(c, nil, "")
		}

		if err := c.Bind(&payload); err != nil {
			return response.R400(c, nil, "")
		}

		if err := payload.Validate(); err != nil {
			return response.RouteValidation(c, err)
		}

		echocontext.SetPayload(c, payload)
		echocontext.SetParam(c, "id", objID)
		return next(c)
	}
}

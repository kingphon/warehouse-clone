package routevalidation

import (
	"fmt"

	"git.selly.red/Selly-Server/warehouse/external/response"
	"git.selly.red/Selly-Server/warehouse/external/utils/echocontext"
	"git.selly.red/Selly-Server/warehouse/pkg/admin/errorcode"
	requestmodel "git.selly.red/Selly-Server/warehouse/pkg/admin/model/request"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Warehouse ...
type Warehouse struct{}

// Create ...
func (Warehouse) Create(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		var (
			payload requestmodel.WarehouseCreate
		)

		if err := c.Bind(&payload); err != nil {
			fmt.Println("err Bind payload: ", err)
			return response.R400(c, nil, "")
		}

		if err := payload.Validate(); err != nil {
			return response.RouteValidation(c, err)
		}

		// assign courierID
		supplierID, _ := primitive.ObjectIDFromHex(payload.Supplier)
		payload.SupplierID = supplierID

		fmt.Println("payload: ", payload)
		echocontext.SetPayload(c, payload)
		return next(c)
	}
}

// All ...
func (Warehouse) All(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		var (
			query requestmodel.WarehouseAll
		)

		if err := c.Bind(&query); err != nil {
			return response.R400(c, nil, "")
		}

		if err := query.Validate(); err != nil {
			return response.RouteValidation(c, err)
		}

		echocontext.SetQuery(c, query)
		return next(c)
	}
}

// Detail ...
func (Warehouse) Detail(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		var (
			id = c.Param("id")
		)

		if !primitive.IsValidObjectID(id) {
			return response.R404(c, nil, errorcode.WarehouseInvalidID)
		}

		objID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return response.R400(c, nil, "")
		}

		echocontext.SetParam(c, "id", objID)
		return next(c)
	}
}

// Update ...
func (Warehouse) Update(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		var (
			payload requestmodel.WarehouseUpdate
			id      = c.Param("id")
		)

		if !primitive.IsValidObjectID(id) {
			return response.R404(c, nil, errorcode.WarehouseInvalidID)
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

// UpdateStatus ...
func (Warehouse) UpdateStatus(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		var (
			id      = c.Param("id")
			payload requestmodel.WarehouseUpdateStatus
		)

		if !primitive.IsValidObjectID(id) {
			return response.R404(c, nil, errorcode.WarehouseInvalidID)
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

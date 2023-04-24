package routevalidation

import (
	"git.selly.red/Selly-Modules/mongodb"
	"git.selly.red/Selly-Server/warehouse/external/response"
	"git.selly.red/Selly-Server/warehouse/external/utils/echocontext"
	"git.selly.red/Selly-Server/warehouse/pkg/app/errorcode"
	requestmodel "git.selly.red/Selly-Server/warehouse/pkg/app/model/request"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// SupplierHoliday ...
type SupplierHoliday struct{}

// Create ...
func (SupplierHoliday) Create(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		var payload requestmodel.SupplierHolidayCreate

		if err := c.Bind(&payload); err != nil {
			return response.R400(c, nil, "")
		}

		if err := payload.Validate(); err != nil {
			return response.RouteValidation(c, err)
		}

		echocontext.SetPayload(c, payload)
		return next(c)
	}
}

// Detail ...
func (SupplierHoliday) Detail(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		var id = c.Param("id")

		if !primitive.IsValidObjectID(id) {
			return response.R400(c, nil, errorcode.SupplierHolidayInvalidID)
		}

		objID, ok := mongodb.NewIDFromString(id)
		if !ok {
			return response.R400(c, nil, errorcode.SupplierHolidayInvalidID)
		}

		echocontext.SetParam(c, "id", objID)
		return next(c)
	}
}

// Update ...
func (SupplierHoliday) Update(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		var payload requestmodel.SupplierHolidayUpdate

		if err := c.Bind(&payload); err != nil {
			return response.R400(c, nil, "")
		}

		if err := payload.Validate(); err != nil {
			return response.RouteValidation(c, err)
		}

		echocontext.SetPayload(c, payload)
		return next(c)
	}
}

// ChangeStatus ...
func (SupplierHoliday) ChangeStatus(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		var payload requestmodel.SupplierHolidayChangeStatus

		if err := c.Bind(&payload); err != nil {
			return response.R400(c, nil, "")
		}

		if err := payload.Validate(); err != nil {
			return response.RouteValidation(c, err)
		}

		echocontext.SetPayload(c, payload)
		return next(c)
	}
}

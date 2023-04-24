package handler

import (
	"git.selly.red/Selly-Server/warehouse/external/response"
	"git.selly.red/Selly-Server/warehouse/external/utils/echocontext"
	requestmodel "git.selly.red/Selly-Server/warehouse/pkg/app/model/request"
	responsemodel "git.selly.red/Selly-Server/warehouse/pkg/app/model/response"
	"git.selly.red/Selly-Server/warehouse/pkg/app/service"
	"github.com/labstack/echo/v4"
)

// SupplierHoliday ...
type SupplierHoliday struct{}

// Create godoc
// @tags SupplierHoliday
// @summary Create
// @id app-supplier-holiday-create
// @security ApiKeyAuth
// @accept json
// @produce json
// @param device-id header string false "Device id"
// @param payload body requestmodel.SupplierHolidayCreate true "Payload"
// @success 200 {object} responsemodel.ResponseCreate
// @router /supplier-holidays [post]
func (SupplierHoliday) Create(c echo.Context) error {
	var (
		ctx     = echocontext.GetContext(c)
		payload = echocontext.GetPayload(c).(requestmodel.SupplierHolidayCreate)
		user    = c.Get("user").(*responsemodel.ResponseUserInfo)
		s       = service.SupplierHoliday(user)
	)

	doc, err := s.CreateWithClientData(ctx, payload)
	if err != nil {
		return response.R404(c, echo.Map{}, err.Error())
	}

	return response.R200(c, responsemodel.ResponseCreate{ID: doc.ID.Hex()}, "")
}

// Detail godoc
// @tags SupplierHoliday
// @summary Detail
// @id app-supplier-holiday-detail
// @security ApiKeyAuth
// @accept json
// @produce json
// @param device-id header string false "Device id"
// @success 200 {object} responsemodel.ResponseSupplierHolidayDetail
// @router /supplier-holidays/supplier [get]
func (SupplierHoliday) Detail(c echo.Context) error {
	var (
		ctx  = echocontext.GetContext(c)
		user = c.Get("user").(*responsemodel.ResponseUserInfo)
		s    = service.SupplierHoliday(user)
	)

	result := s.Detail(ctx)

	return response.R200(c, result, "")
}

// Update godoc
// @tags SupplierHoliday
// @summary Update
// @id app-supplier-holiday-update
// @security ApiKeyAuth
// @accept json
// @produce json
// @param device-id header string false "Device id"
// @param payload body requestmodel.SupplierHolidayUpdate true "Payload"
// @success 200 {object} responsemodel.ResponseChangeStatus
// @router /supplier-holidays [put]
func (SupplierHoliday) Update(c echo.Context) error {
	var (
		ctx     = echocontext.GetContext(c)
		payload = echocontext.GetPayload(c).(requestmodel.SupplierHolidayUpdate)
		user    = c.Get("user").(*responsemodel.ResponseUserInfo)
		s       = service.SupplierHoliday(user)
	)

	result, err := s.UpdateWithClientData(ctx, payload)
	if err != nil {
		return response.R404(c, echo.Map{}, err.Error())
	}

	return response.R200(c, result, "")
}

// ChangeStatus godoc
// @tags SupplierHoliday
// @summary Change status
// @id app-supplier-holiday-change-status
// @security ApiKeyAuth
// @accept json
// @produce json
// @param device-id header string false "Device id"
// @param payload body requestmodel.SupplierHolidayChangeStatus true "Payload"
// @success 200 {object} responsemodel.ResponseUpdate
// @router /supplier-holidays/change-status [patch]
func (SupplierHoliday) ChangeStatus(c echo.Context) error {
	var (
		ctx     = echocontext.GetContext(c)
		payload = echocontext.GetPayload(c).(requestmodel.SupplierHolidayChangeStatus)
		user    = c.Get("user").(*responsemodel.ResponseUserInfo)
		s       = service.SupplierHoliday(user)
	)

	result, err := s.ChangeStatus(ctx, payload)
	if err != nil {
		return response.R404(c, echo.Map{}, err.Error())
	}

	return response.R200(c, result, "")
}

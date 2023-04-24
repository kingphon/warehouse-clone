package handler

import (
	externalauth "git.selly.red/Selly-Server/warehouse/external/auth"
	"git.selly.red/Selly-Server/warehouse/external/response"
	"git.selly.red/Selly-Server/warehouse/external/utils/echocontext"
	"git.selly.red/Selly-Server/warehouse/external/utils/mgquery"
	"git.selly.red/Selly-Server/warehouse/external/utils/ptime"
	"git.selly.red/Selly-Server/warehouse/pkg/admin/config"
	requestmodel "git.selly.red/Selly-Server/warehouse/pkg/admin/model/request"
	responsemodel "git.selly.red/Selly-Server/warehouse/pkg/admin/model/response"
	"git.selly.red/Selly-Server/warehouse/pkg/admin/service"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// SupplierHoliday ...
type SupplierHoliday struct{}

// Create godoc
// @tags SupplierHoliday
// @summary Create
// @id admin-supplier-holiday-create
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
		staff   = echocontext.GetStaff(c).(externalauth.User)
		s       = service.SupplierHoliday(staff)
	)

	doc, err := s.CreateWithClientData(ctx, payload)
	if err != nil {
		return response.R404(c, echo.Map{}, err.Error())
	}

	return response.R200(c, responsemodel.ResponseCreate{ID: doc.ID.Hex()}, "")
}

// All godoc
// @tags SupplierHoliday
// @summary All
// @id admin-supplier-holiday-all
// @security ApiKeyAuth
// @accept json
// @produce json
// @param device-id header string false "Device id"
// @param payload query requestmodel.SupplierHolidayAll true "Query"
// @success 200 {object} responsemodel.ResponseSupplierHolidayAll
// @router /supplier-holidays [get]
func (SupplierHoliday) All(c echo.Context) error {

	var (
		ctx     = echocontext.GetContext(c)
		qParams = echocontext.GetQuery(c).(requestmodel.SupplierHolidayAll)
		s       = service.SupplierHoliday(externalauth.User{})
		q       = mgquery.AppQuery{
			Page:          qParams.Page,
			Limit:         qParams.Limit,
			SortInterface: bson.M{"createdAt": -1},
			Status:        qParams.Status,
			Warehouse: mgquery.Warehouse{
				Warehouse: qParams.Warehouse,
				Keyword:   qParams.Keyword,
				Supplier:  qParams.Supplier,
				FromAt:    ptime.TimeParseISODate(qParams.FromAt),
				ToAt:      ptime.TimeParseISODate(qParams.ToAt),
			},
		}
	)

	result := s.All(ctx, q)
	return response.R200(c, result, "")
}

// Detail godoc
// @tags SupplierHoliday
// @summary Detail
// @id admin-supplier-holiday-detail
// @security ApiKeyAuth
// @accept json
// @produce json
// @param device-id header string false "Device id"
// @Param  id path string true "Supplier holiday id"
// @success 200 {object} responsemodel.ResponseSupplierHolidayDetail
// @router /supplier-holidays/{id} [get]
func (SupplierHoliday) Detail(c echo.Context) error {
	var (
		ctx = echocontext.GetContext(c)
		s   = service.SupplierHoliday(externalauth.User{})
		id  = echocontext.GetParam(c, "id").(primitive.ObjectID)
	)

	result, err := s.Detail(ctx, id)
	if err != nil {
		return response.R404(c, echo.Map{}, err.Error())
	}
	return response.R200(c, result, "")
}

// Update godoc
// @tags SupplierHoliday
// @summary Update
// @id admin-supplier-holiday-update
// @security ApiKeyAuth
// @accept json
// @produce json
// @param device-id header string false "Device id"
// @Param  id path string true "Supplier holiday id"
// @param payload body requestmodel.SupplierHolidayUpdate true "Payload"
// @success 200 {object} responsemodel.ResponseChangeStatus
// @router /supplier-holidays/{id} [put]
func (SupplierHoliday) Update(c echo.Context) error {
	var (
		ctx     = echocontext.GetContext(c)
		payload = echocontext.GetPayload(c).(requestmodel.SupplierHolidayUpdate)
		id      = echocontext.GetParam(c, "id").(primitive.ObjectID)
		staff   = echocontext.GetStaff(c).(externalauth.User)
		s       = service.SupplierHoliday(staff)
	)

	result, err := s.UpdateWithClientData(ctx, id, payload)
	if err != nil {
		return response.R404(c, echo.Map{}, err.Error())
	}

	return response.R200(c, result, "")
}

// ChangeStatus godoc
// @tags SupplierHoliday
// @summary Change status
// @id admin-supplier-holiday-change-status
// @security ApiKeyAuth
// @accept json
// @produce json
// @param device-id header string false "Device id"
// @Param  id path string true "Supplier holiday id"
// @param payload body requestmodel.SupplierHolidayChangeStatus true "Payload"
// @success 200 {object} responsemodel.ResponseUpdate
// @router /supplier-holidays/{id}/status [patch]
func (SupplierHoliday) ChangeStatus(c echo.Context) error {
	var (
		ctx     = echocontext.GetContext(c)
		id      = echocontext.GetParam(c, "id").(primitive.ObjectID)
		payload = echocontext.GetPayload(c).(requestmodel.SupplierHolidayChangeStatus)
		staff   = echocontext.GetStaff(c).(externalauth.User)
		s       = service.SupplierHoliday(staff)
	)

	result, err := s.ChangeStatus(ctx, id, payload)
	if err != nil {
		return response.R404(c, echo.Map{}, err.Error())
	}

	return response.R200(c, result, "")
}

// RunJobUpdateHolidayWarehouses godoc
// @tags SupplierHoliday
// @summary RunJobUpdateHolidayWarehouses
// @id run-job-update-holiday-warehouses
// @security ApiKeyAuth
// @accept json
// @produce json
// @param device-id header string false "Device id"
// @success 200 {object} nil
// @router /supplier-holidays/run-job-update-holiday-warehouses [get]
func (SupplierHoliday) RunJobUpdateHolidayWarehouses(c echo.Context) error {
	var (
		s = service.Schedule{}
	)
	if config.IsEnvDevelop() {
		go s.RunJobUpdateHolidayWarehouses()
	}

	return response.R200(c, nil, "")
}

// RunJobUpdateHolidayStatusForSupplier godoc
// @tags SupplierHoliday
// @summary RunJobUpdateHolidayStatusForSupplier
// @id run-job-update-holiday-status-for-supplier
// @security ApiKeyAuth
// @accept json
// @produce json
// @param device-id header string false "Device id"
// @success 200 {object} nil
// @router /supplier-holidays/run-job-update-holiday-status-for-supplier [get]
func (SupplierHoliday) RunJobUpdateHolidayStatusForSupplier(c echo.Context) error {
	var (
		s = service.Schedule{}
	)

	if config.IsEnvDevelop() {
		go s.RunJobUpdateHolidayStatusForSupplier()
	}

	return response.R200(c, nil, "")
}

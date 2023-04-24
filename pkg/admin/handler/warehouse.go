package handler

import (
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	externalauth "git.selly.red/Selly-Server/warehouse/external/auth"
	"git.selly.red/Selly-Server/warehouse/pkg/admin/config"

	"git.selly.red/Selly-Server/warehouse/external/response"
	"git.selly.red/Selly-Server/warehouse/external/utils/echocontext"
	"git.selly.red/Selly-Server/warehouse/external/utils/mgquery"
	requestmodel "git.selly.red/Selly-Server/warehouse/pkg/admin/model/request"
	responsemodel "git.selly.red/Selly-Server/warehouse/pkg/admin/model/response"
	"git.selly.red/Selly-Server/warehouse/pkg/admin/service"
)

// Warehouse ...
type Warehouse struct{}

// RunJobSetIsClosed godoc
// @tags Warehouse
// @summary RunJobSetIsClosed
// @id run-job-set-is-close-warehouse-all
// @security ApiKeyAuth
// @accept json
// @produce json
// @param device-id header string true "Device id"
// @success 200 {object} nil
// @router /warehouses/run-job-set-is-closed [get]
func (Warehouse) RunJobSetIsClosed(c echo.Context) error {
	var (
		s = service.Schedule{}
	)
	if config.IsEnvDevelop() {
		go s.RunJobUpdateIsClosed()
	}

	return response.R200(c, nil, "")
}

// Create godoc
// @tags Warehouse
// @summary Create
// @id admin-warehouse-create
// @security ApiKeyAuth
// @accept json
// @produce json
// @param device-id header string true "Device id"
// @param payload body requestmodel.WarehouseCreate true "Payload"
// @success 200 {object} responsemodel.Upsert
// @router /warehouses [post]
func (Warehouse) Create(c echo.Context) error {
	var (
		ctx     = echocontext.GetContext(c)
		payload = echocontext.GetPayload(c).(requestmodel.WarehouseCreate)
		cs      = echocontext.GetStaff(c).(externalauth.User)
		s       = service.Warehouse(cs)
	)

	doc, err := s.CreateWithClientData(ctx, payload)
	if err != nil {
		return response.R404(c, echo.Map{}, err.Error())
	}

	return response.R200(c, responsemodel.Upsert{ID: doc.ID.Hex()}, "")
}

// All godoc
// @tags Warehouse
// @summary All
// @id admin-warehouse-all
// @security ApiKeyAuth
// @accept json
// @produce json
// @param device-id header string true "Device id"
// @param payload query requestmodel.WarehouseAll true "Query"
// @success 200 {object} responsemodel.ResponseWarehouseAll
// @router /warehouses [get]
func (Warehouse) All(c echo.Context) error {
	var (
		ctx     = echocontext.GetContext(c)
		qParams = echocontext.GetQuery(c).(requestmodel.WarehouseAll)
		s       = service.Warehouse(externalauth.User{})
		q       = mgquery.AppQuery{
			Page:          qParams.Page,
			Limit:         qParams.Limit,
			SortInterface: bson.M{"createdAt": -1},
			Keyword:       qParams.Keyword,
			BusinessType:  qParams.BusinessType,
			Warehouse: mgquery.Warehouse{
				Status:   qParams.Status,
				Supplier: qParams.Supplier,
				Partner:  qParams.Partner,
			},
		}
	)

	result := s.All(ctx, q)
	return response.R200(c, result, "")
}

// Detail godoc
// @tags Warehouse
// @summary Detail
// @id admin-warehouse-detail
// @security ApiKeyAuth
// @accept json
// @produce json
// @param device-id header string true "Device id"
// @Param  id path string true "Warehouse id"
// @success 200 {object} responsemodel.ResponseWarehouseDetail
// @router /warehouses/{id} [get]
func (Warehouse) Detail(c echo.Context) error {
	var (
		ctx = echocontext.GetContext(c)
		s   = service.Warehouse(externalauth.User{})
		id  = echocontext.GetParam(c, "id").(primitive.ObjectID)
	)

	result, err := s.Detail(ctx, id)
	if err != nil {
		return response.R404(c, echo.Map{}, err.Error())
	}
	return response.R200(c, result, "")
}

// Update godoc
// @tags Warehouse
// @summary Update
// @id admin-warehouse-update
// @security ApiKeyAuth
// @accept json
// @produce json
// @param device-id header string true "Device id"
// @Param  id path string true "routeId"
// @param payload body requestmodel.WarehouseUpdate true "Payload"
// @success 200 {object} responsemodel.Upsert
// @router /warehouses/{id} [put]
func (Warehouse) Update(c echo.Context) error {
	var (
		ctx     = echocontext.GetContext(c)
		payload = echocontext.GetPayload(c).(requestmodel.WarehouseUpdate)
		id      = echocontext.GetParam(c, "id").(primitive.ObjectID)
		cs      = echocontext.GetStaff(c).(externalauth.User)
		s       = service.Warehouse(cs)
	)

	if err := s.UpdateWithClientData(ctx, id, payload); err != nil {
		return response.R404(c, echo.Map{}, err.Error())
	}

	return response.R200(c, nil, "")
}

// UpdateStatus godoc
// @tags Warehouse
// @summary Update
// @id admin-warehouse-updateStatus
// @security ApiKeyAuth
// @accept json
// @produce json
// @param device-id header string true "Device id"
// @Param  id path string true "Warehouse id"
// @param payload body requestmodel.WarehouseUpdateStatus true "Payload"
// @success 200 {object} responsemodel.Upsert
// @router /warehouses/{id}/status [patch]
func (Warehouse) UpdateStatus(c echo.Context) error {
	var (
		ctx     = echocontext.GetContext(c)
		id      = echocontext.GetParam(c, "id").(primitive.ObjectID)
		payload = echocontext.GetPayload(c).(requestmodel.WarehouseUpdateStatus)
		cs      = echocontext.GetStaff(c).(externalauth.User)
		s       = service.Warehouse(cs)
	)

	if err := s.UpdateStatus(ctx, id, payload); err != nil {
		return response.R404(c, echo.Map{}, err.Error())
	}

	return response.R200(c, nil, "")
}

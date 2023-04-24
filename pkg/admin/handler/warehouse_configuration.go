package handler

import (
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"

	externalauth "git.selly.red/Selly-Server/warehouse/external/auth"
	"git.selly.red/Selly-Server/warehouse/external/response"
	"git.selly.red/Selly-Server/warehouse/external/utils/echocontext"
	requestmodel "git.selly.red/Selly-Server/warehouse/pkg/admin/model/request"
	"git.selly.red/Selly-Server/warehouse/pkg/admin/service"
)

type WarehouseConfiguration struct{}

// Detail godoc
// @tags WarehouseConfiguration
// @summary Detail
// @id admin-warehouseConfiguration-detail
// @security ApiKeyAuth
// @accept json
// @produce json
// @param device-id header string true "Device id"
// @Param  id path string true "Warehouse id"
// @success 200 {object} responsemodel.ResponseWarehouseConfigurationDetail
// @router /configurations/{id} [get]
func (WarehouseConfiguration) Detail(c echo.Context) error {
	var (
		ctx = echocontext.GetContext(c)
		s   = service.WarehouseConfiguration(externalauth.User{})
		id  = echocontext.GetParam(c, "id").(primitive.ObjectID)
	)

	result, err := s.DetailByWarehouseID(ctx, id)
	if err != nil {
		return response.R404(c, echo.Map{}, err.Error())
	}
	return response.R200(c, result, "")
}

// UpdateFood godoc
// @tags WarehouseConfiguration
// @summary Update Food
// @id admin-warehouseConfiguration-updateFood
// @security ApiKeyAuth
// @accept json
// @produce json
// @param device-id header string true "Device id"
// @Param  id path string true "Warehouse config id"
// @param payload body requestmodel.WarehouseCfgFoodUpdate true "Payload"
// @success 200 {object} nil
// @router /configurations/{id}/food [PUT]
func (WarehouseConfiguration) UpdateFood(c echo.Context) error {
	var (
		ctx     = echocontext.GetContext(c)
		id      = echocontext.GetParam(c, "id").(primitive.ObjectID)
		payload = echocontext.GetPayload(c).(requestmodel.WarehouseCfgFoodUpdate)
		cs      = echocontext.GetStaff(c).(externalauth.User)
		s       = service.WarehouseConfiguration(cs)
	)

	if err := s.UpdateFood(ctx, id, payload); err != nil {
		return response.R404(c, echo.Map{}, err.Error())
	}

	return response.R200(c, nil, "")
}

// UpdateSupplier godoc
// @tags WarehouseConfiguration
// @summary Update
// @id admin-warehouseConfiguration-updateSupplier
// @security ApiKeyAuth
// @accept json
// @produce json
// @param device-id header string true "Device id"
// @Param  id path string true "Warehouse config id"
// @param payload body requestmodel.WarehouseCfgSupplierUpdate true "Payload"
// @success 200 {object} responsemodel.Upsert
// @router /configurations/{id}/supplier [PUT]
func (WarehouseConfiguration) UpdateSupplier(c echo.Context) error {
	var (
		ctx     = echocontext.GetContext(c)
		id      = echocontext.GetParam(c, "id").(primitive.ObjectID)
		payload = echocontext.GetPayload(c).(requestmodel.WarehouseCfgSupplierUpdate)
		cs      = echocontext.GetStaff(c).(externalauth.User)
		s       = service.WarehouseConfiguration(cs)
	)

	if err := s.UpdateSupplier(ctx, id, payload); err != nil {
		return response.R404(c, echo.Map{}, err.Error())
	}

	return response.R200(c, nil, "")
}

// UpdatePartner godoc
// @tags WarehouseConfiguration
// @summary Update
// @id admin-warehouseConfiguration-updatePartner
// @security ApiKeyAuth
// @accept json
// @produce json
// @param device-id header string true "Device id"
// @Param  id path string true "Warehouse config id"
// @param payload body requestmodel.WarehouseCfgPartnerUpdate true "Payload"
// @success 200 {object} responsemodel.Upsert
// @router /configurations/{id}/partner [PUT]
func (WarehouseConfiguration) UpdatePartner(c echo.Context) error {
	var (
		ctx     = echocontext.GetContext(c)
		id      = echocontext.GetParam(c, "id").(primitive.ObjectID)
		payload = echocontext.GetPayload(c).(requestmodel.WarehouseCfgPartnerUpdate)
		cs      = echocontext.GetStaff(c).(externalauth.User)
		s       = service.WarehouseConfiguration(cs)
	)

	if err := s.UpdatePartner(ctx, id, payload); err != nil {
		return response.R404(c, echo.Map{}, err.Error())
	}

	return response.R200(c, nil, "")
}

// UpdateOrder godoc
// @tags WarehouseConfiguration
// @summary Update
// @id admin-warehouseConfiguration-updateOrder
// @security ApiKeyAuth
// @accept json
// @produce json
// @param device-id header string true "Device id"
// @Param  id path string true "Warehouse config id"
// @param payload body requestmodel.WarehouseCfgOrderUpdate true "Payload"
// @success 200 {object} responsemodel.Upsert
// @router /configurations/{id}/order [PUT]
func (WarehouseConfiguration) UpdateOrder(c echo.Context) error {
	var (
		ctx     = echocontext.GetContext(c)
		id      = echocontext.GetParam(c, "id").(primitive.ObjectID)
		payload = echocontext.GetPayload(c).(requestmodel.WarehouseCfgOrderUpdate)
		cs      = echocontext.GetStaff(c).(externalauth.User)
		s       = service.WarehouseConfiguration(cs)
	)

	if err := s.UpdateOrder(ctx, id, payload); err != nil {
		return response.R404(c, echo.Map{}, err.Error())
	}

	return response.R200(c, nil, "")
}

// UpdateDelivery godoc
// @tags WarehouseConfiguration
// @summary Update
// @id admin-warehouseConfiguration-updateDelivery
// @security ApiKeyAuth
// @accept json
// @produce json
// @param device-id header string true "Device id"
// @Param  id path string true "Warehouse config id"
// @param payload body requestmodel.WarehouseCfgDeliveryUpdate true "Payload"
// @success 200 {object} responsemodel.Upsert
// @router /configurations/{id}/delivery [PUT]
func (WarehouseConfiguration) UpdateDelivery(c echo.Context) error {
	var (
		ctx     = echocontext.GetContext(c)
		id      = echocontext.GetParam(c, "id").(primitive.ObjectID)
		payload = echocontext.GetPayload(c).(requestmodel.WarehouseCfgDeliveryUpdate)
		cs      = echocontext.GetStaff(c).(externalauth.User)
		s       = service.WarehouseConfiguration(cs)
	)

	if err := s.UpdateDelivery(ctx, id, payload); err != nil {
		return response.R404(c, echo.Map{}, err.Error())
	}

	return response.R200(c, nil, "")
}

// UpdateOther godoc
// @tags WarehouseConfiguration
// @summary Update
// @id admin-warehouseConfiguration-updateOther
// @security ApiKeyAuth
// @accept json
// @produce json
// @param device-id header string true "Device id"
// @Param  id path string true "Warehouse config id"
// @param payload body requestmodel.WarehouseCfgOtherUpdate true "Payload"
// @success 200 {object} responsemodel.Upsert
// @router /configurations/{id}/other [PUT]
func (WarehouseConfiguration) UpdateOther(c echo.Context) error {
	var (
		ctx     = echocontext.GetContext(c)
		id      = echocontext.GetParam(c, "id").(primitive.ObjectID)
		payload = echocontext.GetPayload(c).(requestmodel.WarehouseCfgOtherUpdate)
		cs      = echocontext.GetStaff(c).(externalauth.User)
		s       = service.WarehouseConfiguration(cs)
	)

	if err := s.UpdateOther(ctx, id, payload); err != nil {
		return response.R404(c, echo.Map{}, err.Error())
	}

	return response.R200(c, nil, "")
}

// UpdateOrderConfirm godoc
// @tags WarehouseConfiguration
// @summary UpdateOrderConfirm
// @id admin-warehouseConfiguration-order-confirm
// @security ApiKeyAuth
// @accept json
// @produce json
// @param device-id header string true "Device id"
// @Param  id path string true "Warehouse config id"
// @param payload body requestmodel.WarehouseCfgOrderConfirm true "Payload"
// @success 200 {object} nil
// @router /configurations/{id}/order-confirm [PUT]
func (WarehouseConfiguration) UpdateOrderConfirm(c echo.Context) error {
	var (
		ctx     = echocontext.GetContext(c)
		id      = echocontext.GetParam(c, "id").(primitive.ObjectID)
		payload = echocontext.GetPayload(c).(requestmodel.WarehouseCfgOrderConfirm)
		cs      = echocontext.GetStaff(c).(externalauth.User)
		s       = service.WarehouseConfiguration(cs)
	)

	if err := s.UpdateOrderConfirm(ctx, id, payload); err != nil {
		return response.R404(c, echo.Map{}, err.Error())
	}

	return response.R200(c, nil, "")
}

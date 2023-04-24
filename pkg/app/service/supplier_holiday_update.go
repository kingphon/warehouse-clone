package service

import (
	"context"
	"git.selly.red/Selly-Modules/mongodb"
	"git.selly.red/Selly-Server/warehouse/external/utils/ptime"

	requestmodel "git.selly.red/Selly-Server/warehouse/pkg/app/model/request"
	responsemodel "git.selly.red/Selly-Server/warehouse/pkg/app/model/response"

	"git.selly.red/Selly-Server/warehouse/external/constant"
	"git.selly.red/Selly-Server/warehouse/pkg/app/dao"
	"git.selly.red/Selly-Server/warehouse/pkg/app/errorcode"
	"github.com/friendsofgo/errors"
	"go.mongodb.org/mongo-driver/bson"
)

// UpdateWithClientData ...
func (s supplierHolidayImplement) UpdateWithClientData(ctx context.Context, payload requestmodel.SupplierHolidayUpdate) (result responsemodel.ResponseUpdate, err error) {
	var (
		warehouseSvc = warehouseImplement{}
		d            = dao.SupplierHoliday()
		supplierId   = mongodb.ConvertStringToObjectID(s.CurrentUser.SupplierID)
	)

	// Check existed
	supplierHoliday := d.FindOneByCondition(ctx, bson.M{"supplier": supplierId})
	if supplierHoliday.ID.IsZero() {
		err = errors.New(errorcode.SupplierHolidayNotFound)
		return
	}

	if len(payload.Warehouses) == 0 {
		err = errors.New(errorcode.SupplierHolidayMusHaveAtLeastOneWarehouse)
		return
	}

	// Check warehouses invalid
	if err = s.checkInvalidWarehouses(ctx, requestmodel.SupplierHolidayPayload{
		Update: &payload,
	}); err != nil {
		return
	}

	var (
		cond          = bson.M{"supplier": supplierId}
		data          = payload.ConvertToBSON()
		payloadUpdate = bson.M{
			"from":       data.From,
			"to":         data.To,
			"updatedAt":  data.UpdatedAt,
			"isApplyAll": data.IsApplyAll,
			"warehouses": data.Warehouses,
			"reason":     data.Reason,
			"status":     data.Status,
		}
	)

	// Update supplier-holiday
	if err = d.UpdateOneByCondition(ctx, cond, bson.M{"$set": payloadUpdate}); err != nil {
		err = errors.New(errorcode.SupplierHolidayErrorWhenUpdate)
		return
	}

	//  Update holiday warehouse
	go func() {
		var (
			ctxBg = context.Background()
		)

		warehousesByStatuses := warehouseSvc.FindByCondition(ctxBg, bson.M{
			"supplier": supplierId,
			"status": bson.M{"$in": []string{
				constant.WarehouseStatusHoliday,
				constant.StatusActive,
			}},
		})

		warehouseSvc.UpdateWarehousesStatusByIDs(warehousesByStatuses)
	}()

	// Audit
	auditSvc := Audit(s.CurrentUser)
	go auditSvc.Create(
		constant.AuditTargetWarehouseSupplierHoliday,
		supplierHoliday.ID.Hex(),
		constant.MsgEditSupplierHoliday,
		constant.AuditActionCreate,
		payloadUpdate,
	)

	// Response
	result.ID = supplierHoliday.ID.Hex()
	return
}

// ChangeStatus ...
func (s supplierHolidayImplement) ChangeStatus(ctx context.Context, payload requestmodel.SupplierHolidayChangeStatus) (result responsemodel.ResponseChangeStatus, err error) {
	var (
		warehouseSvc  = warehouseImplement{}
		d             = dao.SupplierHoliday()
		supplierId    = mongodb.ConvertStringToObjectID(s.CurrentUser.SupplierID)
		payloadUpdate = bson.M{
			"status":    payload.Status,
			"updatedAt": ptime.Now(),
		}
	)

	// Check isExist supplier-holiday
	supplierHoliday := d.FindOneByCondition(ctx, bson.M{"supplier": supplierId})
	if supplierHoliday.ID.IsZero() {
		err = errors.New(errorcode.SupplierHolidayNotFound)
		return
	}

	// Update supplier-holiday
	if err = d.UpdateOneByCondition(ctx, bson.M{"supplier": supplierId}, bson.M{"$set": payloadUpdate}); err != nil {
		err = errors.New(errorcode.SupplierHolidayErrorWhenUpdate)
		return
	}

	// Update warehouse
	go func() {
		var (
			ctxBg = context.Background()
		)
		warehousesByStatuses := warehouseSvc.FindByCondition(ctxBg, bson.M{
			"supplier": supplierHoliday.Supplier,
			"status": bson.M{"$in": []string{
				constant.WarehouseStatusHoliday,
				constant.StatusActive,
			}},
		})

		warehouseSvc.UpdateWarehousesStatusByIDs(warehousesByStatuses)
	}()

	// Create audit
	auditSvc := Audit(s.CurrentUser)
	go auditSvc.Create(
		constant.AuditTargetWarehouseSupplierHoliday,
		supplierId.Hex(),
		constant.MsgEditStatusSupplierHoliday,
		constant.AuditActionCreate,
		payloadUpdate,
	)

	// Response
	result.ID = supplierId.Hex()
	return
}

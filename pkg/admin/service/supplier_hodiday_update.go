package service

import (
	"context"
	"git.selly.red/Selly-Modules/logger"
	"go.mongodb.org/mongo-driver/mongo"

	"git.selly.red/Selly-Server/warehouse/external/constant"
	"git.selly.red/Selly-Server/warehouse/external/utils/ptime"
	"git.selly.red/Selly-Server/warehouse/pkg/admin/dao"
	"git.selly.red/Selly-Server/warehouse/pkg/admin/errorcode"
	requestmodel "git.selly.red/Selly-Server/warehouse/pkg/admin/model/request"
	responsemodel "git.selly.red/Selly-Server/warehouse/pkg/admin/model/response"
	"github.com/friendsofgo/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//
// PUBLIC METHODS
//

// UpdateWithClientData ...
func (s supplierHolidayImplement) UpdateWithClientData(ctx context.Context, id primitive.ObjectID, payload requestmodel.SupplierHolidayUpdate) (result responsemodel.ResponseUpdate, err error) {
	var (
		warehouseSvc = warehouseImplement{s.CurrentStaff}
		d            = dao.SupplierHoliday()
	)

	// 1. Check existed
	supplierHoliday := d.FindOneByCondition(ctx, bson.M{"_id": id})
	if supplierHoliday.ID.IsZero() {
		err = errors.New(errorcode.SupplierHolidayNotFound)
		return
	}

	// Check warehouses invalid
	if err = s.checkInvalidWarehouses(ctx, requestmodel.SupplierHolidayPayload{
		Update: &payload,
	}); err != nil {
		return
	}

	var (
		cond          = bson.M{"_id": id}
		data          = payload.ConvertToBSON()
		payloadUpdate = bson.M{
			"title":        data.Title,
			"from":         data.From,
			"to":           data.To,
			"reason":       data.Reason,
			"searchString": data.SearchString,
			"updatedAt":    data.UpdatedAt,
			"isApplyAll":   data.IsApplyAll,
			"warehouses":   data.Warehouses,
			"supplier":     data.Supplier,
		}
	)

	// Update supplier-holiday
	if err = d.UpdateOneByCondition(ctx, cond, bson.M{"$set": payloadUpdate}); err != nil {
		err = errors.New(errorcode.SupplierHolidayErrorWhenUpdate)
		return
	}

	// Update holiday warehouse ...
	go func() {
		var ctxBg = context.Background()
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
	auditSvc := Audit(s.CurrentStaff)
	go auditSvc.Create(
		constant.AuditTargetWarehouseSupplierHoliday,
		id.Hex(),
		constant.MsgEditSupplierHoliday,
		constant.AuditActionCreate,
		payloadUpdate,
	)

	// Response
	result.ID = id.Hex()

	return
}

// ChangeStatus ...
func (s supplierHolidayImplement) ChangeStatus(ctx context.Context, id primitive.ObjectID, payload requestmodel.SupplierHolidayChangeStatus) (result responsemodel.ResponseChangeStatus, err error) {
	var (
		warehouseSvc = warehouseImplement{CurrentStaff: s.CurrentStaff}
		d            = dao.SupplierHoliday()
	)

	// Check isExist supplier-holiday
	supplierHoliday := d.FindOneByCondition(ctx, bson.M{"_id": id})
	if supplierHoliday.ID.IsZero() {
		err = errors.New(errorcode.SupplierHolidayNotFound)
		return
	}

	var payloadUpdate = bson.M{
		"status":    payload.Status,
		"updatedAt": ptime.Now(),
	}

	// Update supplier-holiday
	if err = d.UpdateOneByCondition(ctx, bson.M{"_id": id}, bson.M{"$set": payloadUpdate}); err != nil {
		err = errors.New(errorcode.SupplierHolidayErrorWhenUpdate)
		return
	}

	// Update holiday warehouse ...
	go func() {
		var ctxBg = context.Background()
		warehousesByStatuses := warehouseSvc.FindByCondition(ctxBg, bson.M{
			"supplier": supplierHoliday.Supplier,
			"status": bson.M{"$in": []string{
				constant.WarehouseStatusHoliday,
				constant.StatusActive,
			}},
		})

		warehouseSvc.UpdateWarehousesStatusByIDs(warehousesByStatuses)
	}()

	// 6. Create audit
	auditSvc := Audit(s.CurrentStaff)
	go auditSvc.Create(
		constant.AuditTargetWarehouseSupplierHoliday,
		id.Hex(),
		constant.MsgEditStatusSupplierHoliday,
		constant.AuditActionCreate,
		payloadUpdate,
	)

	// Response
	result.ID = id.Hex()

	return
}

// UpdateHolidayStatusForSupplier ...
func (s supplierHolidayImplement) UpdateHolidayStatusForSupplier() {
	var (
		wModel = make([]mongo.WriteModel, 0)
		d      = dao.SupplierHoliday()
		ctx    = context.Background()
	)

	supplierHolidays := s.FindByCondition(ctx, bson.M{"source": constant.WarehouseSupplierHolidaySourceSupplier})
	if len(supplierHolidays) < 0 {
		logger.Error("Error len(holiday)<0 in holiday-update-service : ", logger.LogData{})
	}

	for _, h := range supplierHolidays {
		update := bson.M{
			"updatedAt": ptime.Now(),
		}

		if h.To.Before(ptime.Now()) {
			update["status"] = constant.StatusInactive
		}

		if h.From.After(ptime.Now()) && h.To.After(ptime.Now()) {
			update["status"] = constant.StatusActive
		}

		wModel = append(wModel, mongo.NewUpdateOneModel().SetFilter(bson.M{
			"_id": h.ID,
		}).SetUpdate(bson.M{
			"$set": update,
		}))
	}

	if len(wModel) > 0 {
		if err := d.BulkWrite(ctx, wModel); err != nil {
			logger.Error("Error dao.SupplierHoliday().BulkWrite in supplier-holiday service app : ", logger.LogData{
				"error": err.Error(),
			})
		}
	}

}

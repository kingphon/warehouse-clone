package service

import (
	"context"
	"git.selly.red/Selly-Modules/logger"
	"git.selly.red/Selly-Modules/natsio/client"
	"git.selly.red/Selly-Modules/natsio/model"
	"git.selly.red/Selly-Server/warehouse/external/constant"
	mgwarehouse "git.selly.red/Selly-Server/warehouse/external/model/mg/warehouse"
	"git.selly.red/Selly-Server/warehouse/external/utils/ptime"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"git.selly.red/Selly-Server/warehouse/pkg/app/dao"
)

//
// PUBLIC METHOD
//

// UpdateManyByCondition ...
func (warehouseImplement) UpdateManyByCondition(ctx context.Context, cond interface{}, payload interface{}) error {
	return dao.Warehouse().UpdateManyByCondition(ctx, cond, payload)
}

// BulkWriteWhenChangeStatusSupplierHoliday ...
func (warehouseImplement) BulkWriteWhenChangeStatusSupplierHoliday(ctx context.Context, wModel []mongo.WriteModel) error {
	return dao.Warehouse().BulkWrite(ctx, wModel)
}

// UpdateWarehousesStatusByIDs ...
func (s warehouseImplement) UpdateWarehousesStatusByIDs(warehouses []mgwarehouse.Warehouse) {
	var (
		wModel = make([]mongo.WriteModel, 0)
		ctx    = context.Background()
	)

	listWarehouseUpdatePendingInActive := make([]model.UpdateStatusWarehousePendingInactive, 0)

	for _, w := range warehouses {
		// Check isActive holiday
		newStatus, holiday := s.checkWarehouseActiveHoliday(ctx, w)

		update := bson.M{
			"updatedAt": ptime.Now(),
		}

		if newStatus != w.Status {
			update["status"] = newStatus
			update["reasonPendingInactive"] = ""
			if newStatus == constant.WarehouseStatusHoliday {
				update["statusBeforeHoliday"] = w.Status
				update["reasonPendingInactive"] = holiday.Reason
			}

			listWarehouseUpdatePendingInActive = append(listWarehouseUpdatePendingInActive, model.UpdateStatusWarehousePendingInactive{
				WarehouseID:     w.ID.Hex(),
				PendingInactive: newStatus == constant.WarehouseStatusHoliday,
			})
		}

		if newStatus == constant.WarehouseStatusHoliday {
			update["reasonPendingInactive"] = holiday.Reason
		}

		wModel = append(wModel, mongo.NewUpdateOneModel().SetFilter(bson.M{
			"_id": w.ID,
		}).SetUpdate(bson.M{
			"$set": update,
		}))
	}

	if len(wModel) > 0 {
		var warehouseSvc = warehouseImplement{}
		if err := warehouseSvc.BulkWriteWhenChangeStatusSupplierHoliday(ctx, wModel); err != nil {
			logger.Error("Error dao.Warehouse().BulkWrite in supplier-holiday service : ", logger.LogData{
				"error": err.Error(),
			})
		}
	}

	// Update pendingInactive product
	var bodyUpdateProduct = model.UpdateStatusWarehousePendingInactiveRequest{Warehouses: listWarehouseUpdatePendingInActive}
	if err := client.GetWarehouse().UpdateStatusWarehousePendingInactive(bodyUpdateProduct); err != nil {
		logger.Error("Error Update pending inActive product by warehouseIDs : ", logger.LogData{
			"error": err.Error(),
		})
	}

}

//
// PRIVATE METHOD
//

// checkWarehouseActiveHoliday ...
func (warehouseImplement) checkWarehouseActiveHoliday(ctx context.Context, warehouse mgwarehouse.Warehouse) (string, mgwarehouse.SupplierHoliday) {
	var (
		supplierHolidaySvc = supplierHolidayImplement{nil}
		cond               = bson.M{
			"supplier": warehouse.Supplier,
			"status":   "active",
			"from":     bson.M{"$lte": ptime.Now()}, // from < now < to
			"to":       bson.M{"$gte": ptime.Now()},
			"$or": []bson.M{
				{
					"isApplyAll": true,
				},
				{
					"isApplyAll": false,
					"warehouses": warehouse.ID,
				},
			},
		}
	)

	holidays := supplierHolidaySvc.FindByCondition(ctx, cond)

	if len(holidays) == 0 {
		if warehouse.StatusBeforeHoliday == "" {
			return warehouse.Status, mgwarehouse.SupplierHoliday{}
		}

		return warehouse.StatusBeforeHoliday, mgwarehouse.SupplierHoliday{}
	}

	return constant.WarehouseStatusHoliday, holidays[0]
}

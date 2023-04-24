package service

import (
	"context"

	"git.selly.red/Selly-Server/warehouse/external/constant"
	"git.selly.red/Selly-Server/warehouse/external/utils/ptime"

	"git.selly.red/Selly-Modules/mongodb"
	"git.selly.red/Selly-Server/warehouse/pkg/admin/errorcode"
	requestmodel "git.selly.red/Selly-Server/warehouse/pkg/admin/model/request"
	"github.com/friendsofgo/errors"
	"go.mongodb.org/mongo-driver/bson"
)

// checkWarehouseInvalid...
func (s supplierHolidayImplement) checkInvalidWarehouses(ctx context.Context, inputPayload requestmodel.SupplierHolidayPayload) (err error) {
	var (
		isApplyALl bool
		supplier   string
		warehouses []string
	)

	if inputPayload.Create != nil {
		isApplyALl = inputPayload.Create.IsApplyAll
		supplier = inputPayload.Create.Supplier
		warehouses = inputPayload.Create.Warehouses
	}

	if inputPayload.Update != nil {
		isApplyALl = inputPayload.Update.IsApplyAll
		supplier = inputPayload.Update.Supplier
		warehouses = inputPayload.Update.Warehouses
	}

	var warehouseSvc = warehouseImplement{s.CurrentStaff}

	if !isApplyALl {
		warehousesByIDs := warehouseSvc.FindByCondition(ctx, bson.M{
			"supplier": mongodb.ConvertStringToObjectID(supplier),
			"_id": bson.M{
				"$in": mongodb.ConvertStringsToObjectIDs(warehouses),
			},
		})
		if len(warehouses) != len(warehousesByIDs) {
			err = errors.New(errorcode.WarehouseNotBelongSupplier)
			return
		}
	}

	return
}

// updateWarehousesOnHoliday ...
func (s supplierHolidayImplement) updateWarehousesOnHoliday(ctx context.Context, inputPayload requestmodel.SupplierHolidayPayload) (err error) {
	var (
		from, to   string
		warehouses []string
	)

	if inputPayload.Create != nil {
		from = inputPayload.Create.From
		to = inputPayload.Create.To
		warehouses = inputPayload.Create.Warehouses
	}

	if inputPayload.Update != nil {
		from = inputPayload.Update.From
		to = inputPayload.Update.To
		warehouses = inputPayload.Update.Warehouses
	}

	var warehouseSvc = warehouseImplement{s.CurrentStaff}
	if ptime.Now().After(ptime.TimeParseISODate(from)) && ptime.Now().Before(ptime.TimeParseISODate(to)) {
		var (
			warehouseObjIDs = mongodb.ConvertStringsToObjectIDs(warehouses)
			cond            = bson.M{"_id": bson.M{"$in": warehouseObjIDs}}
			whUpdate        = bson.M{"status": constant.WarehouseStatusHoliday}
		)
		if err = warehouseSvc.UpdateManyByCondition(ctx, cond, bson.M{"$set": whUpdate}); err != nil {
			err = errors.New(errorcode.WarehouseErrorWhenUpdate)
			return
		}
	}

	return
}

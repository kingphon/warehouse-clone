package service

import (
	"context"
	"git.selly.red/Selly-Modules/mongodb"
	"git.selly.red/Selly-Server/warehouse/pkg/app/errorcode"
	requestmodel "git.selly.red/Selly-Server/warehouse/pkg/app/model/request"
	"github.com/friendsofgo/errors"
	"go.mongodb.org/mongo-driver/bson"
)

// checkWarehouseInvalid...
func (s supplierHolidayImplement) checkInvalidWarehouses(ctx context.Context, inputPayload requestmodel.SupplierHolidayPayload) (err error) {
	var (
		isApplyALl bool
		warehouses []string
	)

	if inputPayload.Create != nil {
		isApplyALl = inputPayload.Create.IsApplyAll
		warehouses = inputPayload.Create.Warehouses
	}

	if inputPayload.Update != nil {
		isApplyALl = inputPayload.Update.IsApplyAll
		warehouses = inputPayload.Update.Warehouses
	}

	var warehouseSvc = warehouseImplement{}

	if !isApplyALl {
		warehousesByIDs := warehouseSvc.FindByCondition(ctx, bson.M{
			"supplier": mongodb.ConvertStringToObjectID(s.CurrentUser.SupplierID),
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

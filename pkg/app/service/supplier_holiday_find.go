package service

import (
	"context"
	"git.selly.red/Selly-Modules/mongodb"
	"git.selly.red/Selly-Server/warehouse/external/constant"
	"git.selly.red/Selly-Server/warehouse/external/utils/parray"
	"git.selly.red/Selly-Server/warehouse/external/utils/ptime"

	mgwarehouse "git.selly.red/Selly-Server/warehouse/external/model/mg/warehouse"
	"git.selly.red/Selly-Server/warehouse/pkg/app/dao"
	responsemodel "git.selly.red/Selly-Server/warehouse/pkg/app/model/response"
	"go.mongodb.org/mongo-driver/bson"
)

//
// PUBLIC METHODS
//

// FindByCondition ...
func (s supplierHolidayImplement) FindByCondition(ctx context.Context, cond interface{}) []mgwarehouse.SupplierHoliday {
	return dao.SupplierHoliday().FindByCondition(ctx, cond)
}

// Detail ...
func (s supplierHolidayImplement) Detail(ctx context.Context) (result *responsemodel.ResponseSupplierHolidayDetail) {
	var (
		d    = dao.SupplierHoliday()
		user = s.CurrentUser
		cond = bson.M{"supplier": mongodb.ConvertStringToObjectID(user.SupplierID)}
	)

	supplierHoliday := d.FindOneByCondition(ctx, cond)
	if supplierHoliday.ID.IsZero() {
		return
	}

	warehouses := s.getWarehouseBySupplierHolidays(ctx, supplierHoliday)

	result = s.getSupplierHolidayDetailInfo(ctx, supplierHoliday, warehouses)
	return
}

//
// PRIVATE METHODS
//

// getWarehousesBySupplierHolidays ...
func (s supplierHolidayImplement) getWarehouseBySupplierHolidays(ctx context.Context, doc mgwarehouse.SupplierHoliday) (result []responsemodel.ResponseWarehouseInfo) {
	result = make([]responsemodel.ResponseWarehouseInfo, 0)
	var warehouseSvc = warehouseImplement{}

	if doc.IsApplyAll == true {
		getWarehousesBySuppliers := warehouseSvc.FindByCondition(ctx, bson.M{"supplier": doc.Supplier})
		for _, w := range getWarehousesBySuppliers {
			result = append(result, responsemodel.ResponseWarehouseInfo{
				ID:     w.ID.Hex(),
				Name:   w.Name,
				Status: w.Status,
			})
		}
		return
	}

	listWarehouse := warehouseSvc.FindByCondition(ctx, bson.M{"_id": bson.M{"$in": mongodb.UniqObjectIds(doc.Warehouses)}})
	for _, w := range listWarehouse {
		result = append(result, responsemodel.ResponseWarehouseInfo{
			ID:     w.ID.Hex(),
			Name:   w.Name,
			Status: w.Status,
		})
	}
	return
}

// getSupplierHolidayDetailInfo ...
func (s supplierHolidayImplement) getSupplierHolidayDetailInfo(ctx context.Context, doc mgwarehouse.SupplierHoliday, warehouses []responsemodel.ResponseWarehouseInfo) *responsemodel.ResponseSupplierHolidayDetail {
	if doc.IsApplyAll == true {
		return s.detail(ctx, doc, warehouses)
	}

	var warehouseList []responsemodel.ResponseWarehouseInfo
	for _, id := range doc.Warehouses {
		if found := parray.Find(warehouses, func(item responsemodel.ResponseWarehouseInfo) bool {
			return item.ID == id.Hex()
		}); found != nil {
			warehouseList = append(warehouseList, found.(responsemodel.ResponseWarehouseInfo))
		}
	}

	return s.detail(ctx, doc, warehouseList)
}

// detail ...
func (s supplierHolidayImplement) detail(ctx context.Context, doc mgwarehouse.SupplierHoliday, warehouses []responsemodel.ResponseWarehouseInfo) *responsemodel.ResponseSupplierHolidayDetail {
	result := &responsemodel.ResponseSupplierHolidayDetail{
		ID:         doc.ID.Hex(),
		Title:      doc.Title,
		From:       ptime.TimeResponseInit(doc.From),
		To:         ptime.TimeResponseInit(doc.To),
		Reason:     doc.Reason,
		Source:     doc.Source,
		Status:     doc.Status,
		Warehouses: make([]responsemodel.ResponseWarehouseInfo, 0),
		IsApplyAll: doc.IsApplyAll,
	}

	if len(warehouses) > 0 {
		result.Warehouses = warehouses
	}

	if doc.To.Before(ptime.Now()) {
		result.Status = constant.StatusInactive
	}

	return result
}

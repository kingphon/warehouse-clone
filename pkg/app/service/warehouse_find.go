package service

import (
	"context"

	"git.selly.red/Selly-Modules/natsio/client"
	"git.selly.red/Selly-Modules/natsio/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"git.selly.red/Selly-Server/warehouse/pkg/app/dao"
	responsemodel "git.selly.red/Selly-Server/warehouse/pkg/app/model/response"

	mgwarehouse "git.selly.red/Selly-Server/warehouse/external/model/mg/warehouse"
)

//
// PUBLIC METHODS
//

// Detail ...
func (s warehouseImplement) Detail(ctx context.Context, id primitive.ObjectID) (*responsemodel.ResponseWarehouseDetail, error) {
	var (
		d    = dao.Warehouse()
		cond = bson.M{"_id": id}
		err  error
	)

	// get warehouse
	warehouse := d.FindOneByCondition(ctx, cond)

	wConfig := dao.WarehouseConfiguration().FindByWarehouseID(ctx, warehouse.ID)
	return s.detail(ctx, warehouse, wConfig), err
}

//
// PRIVATE METHODS
//

// detail ...
func (warehouseImplement) detail(ctx context.Context, d mgwarehouse.Warehouse, config mgwarehouse.Configuration) *responsemodel.ResponseWarehouseDetail {
	response, _ := client.GetNews().GetProductNoticesByInventory(model.GetProductNoticesByInventoryRequest{InventoryIds: []string{d.ID.Hex()}})
	res := &responsemodel.ResponseWarehouseDetail{
		ID:                              d.ID.Hex(),
		Name:                            d.Name,
		Code:                            0,
		CanIssueInvoice:                 config.Supplier.InvoiceDeliveryMethod != "none" && config.Supplier.InvoiceDeliveryMethod != "",
		InvoiceDeliveryMethod:           config.Supplier.InvoiceDeliveryMethod,
		DoesSupportSellyExpress:         config.Other.DoesSupportSellyExpress,
		LimitedNumberOfProductsPerOrder: int(config.Order.LimitNumberOfPurchases),
	}
	if response != nil {
		res.Notices = response.Notices
	}
	return res
}

func (warehouseImplement) FindByCondition(ctx context.Context, cond interface{}) []mgwarehouse.Warehouse {
	return dao.Warehouse().FindByCondition(ctx, cond)
}

//
// NATS METHODS
//

// GetLocationByCode ...
func (warehouseImplement) GetLocationByCode(ctx context.Context, w mgwarehouse.Warehouse) (*model.ResponseLocationAddress, error) {
	body := model.LocationRequestPayload{
		Province: w.Location.Province,
		District: w.Location.District,
		Ward:     w.Location.Ward,
	}
	return client.GetLocation().GetLocationByCode(body)
}

// GetSupplierByIDs ...
func (warehouseImplement) GetSupplierByIDs(ctx context.Context, listID []primitive.ObjectID) ([]*model.ResponseSupplierInfo, error) {
	body := model.GetSupplierRequest{ListID: listID}
	return client.GetSupplier().GetListSupplierInfo(body)
}

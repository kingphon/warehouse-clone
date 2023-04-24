package service

import (
	"context"
	mgwarehouse "git.selly.red/Selly-Server/warehouse/external/model/mg/warehouse"

	"go.mongodb.org/mongo-driver/bson/primitive"

	responsemodel "git.selly.red/Selly-Server/warehouse/pkg/app/model/response"
)

// WarehouseInterface ...
type WarehouseInterface interface {
	// Detail get detail warehouse from web app ...
	Detail(ctx context.Context, id primitive.ObjectID) (*responsemodel.ResponseWarehouseDetail, error)

	// UpdateManyByCondition ...
	UpdateManyByCondition(ctx context.Context, cond interface{}, payload interface{}) error

	// UpdateWarehousesStatusByIDs ...
	UpdateWarehousesStatusByIDs(warehouses []mgwarehouse.Warehouse)
}

// WarehouseImplement ...
type warehouseImplement struct {
}

// Warehouse return warehouse service
func Warehouse() WarehouseInterface {
	return warehouseImplement{}
}

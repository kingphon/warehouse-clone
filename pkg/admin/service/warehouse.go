package service

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"

	"go.mongodb.org/mongo-driver/bson/primitive"

	externalauth "git.selly.red/Selly-Server/warehouse/external/auth"
	mgwarehouse "git.selly.red/Selly-Server/warehouse/external/model/mg/warehouse"
	"git.selly.red/Selly-Server/warehouse/external/utils/mgquery"
	requestmodel "git.selly.red/Selly-Server/warehouse/pkg/admin/model/request"
	responsemodel "git.selly.red/Selly-Server/warehouse/pkg/admin/model/response"
)

// WarehouseInterface ...
type WarehouseInterface interface {
	// CreateWithClientData create new courier from web admin
	CreateWithClientData(ctx context.Context, payload requestmodel.WarehouseCreate) (mgwarehouse.Warehouse, error)

	// All return warehouse with condition ...
	All(ctx context.Context, q mgquery.AppQuery) (result responsemodel.ResponseWarehouseAll)

	// Detail get detail warehouse from web admin ...
	Detail(ctx context.Context, id primitive.ObjectID) (*responsemodel.ResponseWarehouseDetail, error)

	// UpdateWithClientData Update get update warehouse from web admin ...
	UpdateWithClientData(ctx context.Context, id primitive.ObjectID, payload requestmodel.WarehouseUpdate) (err error)

	// UpdateStatus Update get update status warehouse from web admin ...
	UpdateStatus(ctx context.Context, id primitive.ObjectID, payload requestmodel.WarehouseUpdateStatus) (err error)

	// FindByCondition ...
	FindByCondition(ctx context.Context, cond interface{}) (result []mgwarehouse.Warehouse)

	// FindOneByCondition ...
	FindOneByCondition(ctx context.Context, cond interface{}) (result mgwarehouse.Warehouse)

	// BulkWriteWhenChangeStatusSupplierHoliday ...
	BulkWriteWhenChangeStatusSupplierHoliday(ctx context.Context, wModel []mongo.WriteModel) error

	// UpdateManyByCondition ...
	UpdateManyByCondition(ctx context.Context, cond interface{}, payload interface{}) error

	// UpdateWarehousesStatusByIDs ...
	UpdateWarehousesStatusByIDs(warehouses []mgwarehouse.Warehouse)
}

// WarehouseImplement ...
type warehouseImplement struct {
	CurrentStaff externalauth.User
}

// Warehouse return warehouse service
func Warehouse(cs externalauth.User) WarehouseInterface {
	return warehouseImplement{
		CurrentStaff: cs,
	}
}

package dao

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"

	mgwarehouse "git.selly.red/Selly-Server/warehouse/external/model/mg/warehouse"
	"go.mongodb.org/mongo-driver/mongo/options"
)

/* Entity
- Warehouse
*/

// WarehouseInterface ...
type WarehouseInterface interface {
	// FindOneByCondition ...
	FindOneByCondition(ctx context.Context, cond interface{}, opts ...*options.FindOneOptions) (doc mgwarehouse.Warehouse)

	// FindByCondition ...
	FindByCondition(ctx context.Context, cond interface{}, opts ...*options.FindOptions) (docs []mgwarehouse.Warehouse)

	// UpdateManyByCondition ...
	UpdateManyByCondition(ctx context.Context, cond interface{}, payload interface{}, opts ...*options.UpdateOptions) (err error)

	// BulkWrite ...
	BulkWrite(ctx context.Context, models []mongo.WriteModel, opts ...*options.BulkWriteOptions) error
}

// WarehouseImplement ...
type warehouseImplement struct{}

// Warehouse return warehouse dao
func Warehouse() WarehouseInterface {
	return warehouseImplement{}
}

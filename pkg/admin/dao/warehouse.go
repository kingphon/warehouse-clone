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
	// InsertOne ...
	InsertOne(ctx context.Context, payload mgwarehouse.Warehouse) error

	// FindByCondition ...
	FindByCondition(ctx context.Context, cond interface{}, opts ...*options.FindOptions) (docs []mgwarehouse.Warehouse)

	// CountByCondition ...
	CountByCondition(ctx context.Context, cond interface{}) int64

	// DistinctWithField ...
	DistinctWithField(ctx context.Context, cond interface{}, field string) ([]interface{}, error)

	// FindOneByCondition ...
	FindOneByCondition(ctx context.Context, cond interface{}, opts ...*options.FindOneOptions) (doc mgwarehouse.Warehouse)

	// UpdateOneByCondition ...
	UpdateOneByCondition(ctx context.Context, cond interface{}, payload interface{}, opts ...*options.UpdateOptions) (err error)

	// UpdateStatus ...
	UpdateStatus(ctx context.Context, cond interface{}, status string) (err error)

	// BulkWrite ...
	BulkWrite(ctx context.Context, models []mongo.WriteModel, opts ...*options.BulkWriteOptions) error

	UpdateManyByCondition(ctx context.Context, cond interface{}, payload interface{}, opts ...*options.UpdateOptions) (err error)
}

// WarehouseImplement ...
type warehouseImplement struct{}

// Warehouse return warehouse dao
func Warehouse() WarehouseInterface {
	return warehouseImplement{}
}

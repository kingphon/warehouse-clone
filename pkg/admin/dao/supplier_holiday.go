package dao

import (
	"context"
	mgwarehouse "git.selly.red/Selly-Server/warehouse/external/model/mg/warehouse"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// SupplierHolidayInterface ...
type SupplierHolidayInterface interface {
	// InsertOne ...
	InsertOne(ctx context.Context, payload mgwarehouse.SupplierHoliday) error

	// FindByCondition ...
	FindByCondition(ctx context.Context, cond interface{}, opts ...*options.FindOptions) (docs []mgwarehouse.SupplierHoliday)

	// FindOneByCondition ...
	FindOneByCondition(ctx context.Context, cond interface{}, opts ...*options.FindOneOptions) (doc mgwarehouse.SupplierHoliday)

	// UpdateOneByCondition ...
	UpdateOneByCondition(ctx context.Context, cond interface{}, payload interface{}, opts ...*options.UpdateOptions) (err error)

	// CountByCondition ...
	CountByCondition(ctx context.Context, cond interface{}) int64

	// BulkWrite ...
	BulkWrite(ctx context.Context, models []mongo.WriteModel, opts ...*options.BulkWriteOptions) error
}

// supplierHolidayImplement ..
type supplierHolidayImplement struct{}

// SupplierHoliday ...
func SupplierHoliday() SupplierHolidayInterface {
	return &supplierHolidayImplement{}
}

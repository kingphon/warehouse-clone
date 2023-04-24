package dao

import (
	"context"
	mgwarehouse "git.selly.red/Selly-Server/warehouse/external/model/mg/warehouse"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// SupplierHolidayInterface ...
type SupplierHolidayInterface interface {
	// FindOneByCondition ...
	FindOneByCondition(ctx context.Context, cond interface{}, opts ...*options.FindOneOptions) (doc mgwarehouse.SupplierHoliday)

	// InsertOne ...
	InsertOne(ctx context.Context, payload mgwarehouse.SupplierHoliday) error

	// UpdateOneByCondition ...
	UpdateOneByCondition(ctx context.Context, cond interface{}, payload interface{}, opts ...*options.UpdateOptions) (err error)

	// FindByCondition ...
	FindByCondition(ctx context.Context, cond interface{}, opts ...*options.FindOptions) (docs []mgwarehouse.SupplierHoliday)
}

// supplierHolidayImplement ..
type supplierHolidayImplement struct{}

// SupplierHoliday ...
func SupplierHoliday() SupplierHolidayInterface {
	return &supplierHolidayImplement{}
}

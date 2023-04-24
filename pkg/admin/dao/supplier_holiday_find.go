package dao

import (
	"context"
	"git.selly.red/Selly-Modules/logger"
	mgwarehouse "git.selly.red/Selly-Server/warehouse/external/model/mg/warehouse"
	"git.selly.red/Selly-Server/warehouse/pkg/admin/database"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// FindByCondition ...
func (supplierHolidayImplement) FindByCondition(ctx context.Context, cond interface{}, opts ...*options.FindOptions) (docs []mgwarehouse.SupplierHoliday) {
	var (
		col = database.SupplierHolidayCol()
	)

	cursor, err := col.Find(ctx, cond, opts...)
	if err != nil {
		logger.Error("dao.SupplierHoliday - FindByCondition cursor", logger.LogData{
			"cond":  cond,
			"opts":  opts,
			"error": err.Error(),
		})
		return
	}

	// Close cursor
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &docs); err != nil {
		logger.Error("dao.SupplierHoliday - FindByCondition decode", logger.LogData{
			"cond":  cond,
			"opts":  opts,
			"error": err.Error(),
		})
	}
	return
}

// FindOneByCondition ...
func (supplierHolidayImplement) FindOneByCondition(ctx context.Context, cond interface{}, opts ...*options.FindOneOptions) (doc mgwarehouse.SupplierHoliday) {
	var (
		col = database.SupplierHolidayCol()
	)

	if err := col.FindOne(ctx, cond, opts...).Decode(&doc); err != nil {
		logger.Error("dao.SupplierHoliday - FindOneByCondition err", logger.LogData{
			"cond":  cond,
			"opts":  opts,
			"error": err.Error(),
		})
	}
	return
}

// CountByCondition ...
func (supplierHolidayImplement) CountByCondition(ctx context.Context, cond interface{}) int64 {
	var (
		col = database.SupplierHolidayCol()
	)

	total, _ := col.CountDocuments(ctx, cond)
	return total
}

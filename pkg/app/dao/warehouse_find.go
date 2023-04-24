package dao

import (
	"context"

	"git.selly.red/Selly-Modules/logger"
	"go.mongodb.org/mongo-driver/mongo/options"

	mgwarehouse "git.selly.red/Selly-Server/warehouse/external/model/mg/warehouse"
	"git.selly.red/Selly-Server/warehouse/pkg/app/database"
)

// CountByCondition ...
func (warehouseImplement) CountByCondition(ctx context.Context, cond interface{}) int64 {
	var (
		col = database.WarehouseCol()
	)

	total, _ := col.CountDocuments(ctx, cond)
	return total
}

// FindOneByCondition ...
func (warehouseImplement) FindOneByCondition(ctx context.Context, cond interface{}, opts ...*options.FindOneOptions) (doc mgwarehouse.Warehouse) {
	var (
		col = database.WarehouseCol()
	)

	if err := col.FindOne(ctx, cond, opts...).Decode(&doc); err != nil {
		logger.Error("dao.Warehouse - FindOneByCondition err", logger.LogData{
			"cond":  cond,
			"opts":  opts,
			"error": err.Error(),
		})
	}
	return
}

// FindByCondition ...
func (warehouseImplement) FindByCondition(ctx context.Context, cond interface{}, opts ...*options.FindOptions) (docs []mgwarehouse.Warehouse) {
	var (
		col = database.WarehouseCol()
	)

	cursor, err := col.Find(ctx, cond, opts...)
	if err != nil {
		logger.Error("dao.Warehouse - FindByCondition cursor", logger.LogData{
			"cond":  cond,
			"opts":  opts,
			"error": err.Error(),
		})
		return
	}

	// Close cursor
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &docs); err != nil {
		logger.Error("dao.Warehouse - FindByCondition decode", logger.LogData{
			"cond":  cond,
			"opts":  opts,
			"error": err.Error(),
		})
	}
	return
}

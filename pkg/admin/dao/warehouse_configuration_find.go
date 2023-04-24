package dao

import (
	"context"

	"git.selly.red/Selly-Modules/logger"
	"go.mongodb.org/mongo-driver/mongo/options"

	mgwarehouse "git.selly.red/Selly-Server/warehouse/external/model/mg/warehouse"
	"git.selly.red/Selly-Server/warehouse/pkg/admin/database"
)

// FindByCondition ...
func (d warehouseConfigurationImplement) FindByCondition(ctx context.Context, cond interface{}, opts ...*options.FindOptions) (docs []mgwarehouse.Configuration) {
	var (
		col = database.WarehouseConfigurationCol()
	)

	cursor, err := col.Find(ctx, cond, opts...)
	if err != nil {
		logger.Error("dao.Warehouse Configuration - FindByCondition cursor", logger.LogData{
			"cond":  cond,
			"opts":  opts,
			"error": err.Error(),
		})
		return
	}

	// Close cursor
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &docs); err != nil {
		logger.Error("dao.Warehouse Configuration - FindByCondition decode", logger.LogData{
			"cond":  cond,
			"opts":  opts,
			"error": err.Error(),
		})
	}
	return
}

// FindOneByCondition ...
func (warehouseConfigurationImplement) FindOneByCondition(ctx context.Context, cond interface{}, opts ...*options.FindOneOptions) (doc mgwarehouse.Configuration) {
	var (
		col = database.WarehouseConfigurationCol()
	)

	if err := col.FindOne(ctx, cond, opts...).Decode(&doc); err != nil {
		logger.Error("dao.WarehouseConfigurations - FindOneByCondition err", logger.LogData{
			"cond":  cond,
			"opts":  opts,
			"error": err.Error(),
		})
	}
	return
}

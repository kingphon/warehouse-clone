package dao

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"

	"git.selly.red/Selly-Modules/logger"
	"go.mongodb.org/mongo-driver/mongo/options"

	"git.selly.red/Selly-Server/warehouse/pkg/admin/database"
)

// UpdateOne ...
func (warehouseConfigurationImplement) UpdateOne(ctx context.Context, cond interface{}, payload interface{}, opts ...*options.UpdateOptions) (err error) {
	var (
		col = database.WarehouseConfigurationCol()
	)

	if _, err := col.UpdateOne(ctx, cond, payload, opts...); err != nil {
		logger.Error("dao.WarehouseConfiguration - UpdateSupplier", logger.LogData{
			"payload": payload,
			"error":   err.Error(),
		})
	}

	return
}

// BulkWrite ...
func (warehouseConfigurationImplement) BulkWrite(ctx context.Context, models []mongo.WriteModel, opts ...*options.BulkWriteOptions) error {
	_, err := database.WarehouseConfigurationCol().BulkWrite(ctx, models, opts...)
	if err != nil {
		logger.Error("dao.WarehouseConfiguration - BulkWrite", logger.LogData{
			"error": err.Error(),
		})
	}
	return err
}

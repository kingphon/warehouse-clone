package dao

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"

	"git.selly.red/Selly-Modules/logger"
	"git.selly.red/Selly-Server/warehouse/pkg/app/database"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// UpdateManyByCondition ...
func (d warehouseImplement) UpdateManyByCondition(ctx context.Context, cond interface{}, payload interface{}, opts ...*options.UpdateOptions) (err error) {
	var (
		col = database.WarehouseCol()
	)

	if _, err := col.UpdateMany(ctx, cond, payload, opts...); err != nil {
		logger.Error("dao.Warehouse - UpdateOne", logger.LogData{
			"payload": payload,
			"error":   err.Error(),
		})
	}
	return
}

// BulkWrite ...
func (d warehouseImplement) BulkWrite(ctx context.Context, models []mongo.WriteModel, opts ...*options.BulkWriteOptions) error {
	_, err := database.WarehouseCol().BulkWrite(ctx, models, opts...)
	if err != nil {
		logger.Error("dao.Warehouse - BulkWrite", logger.LogData{
			"error": err.Error(),
		})
	}
	return err
}

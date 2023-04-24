package dao

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"time"

	"git.selly.red/Selly-Modules/logger"
	"git.selly.red/Selly-Server/warehouse/pkg/admin/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// UpdateOneByCondition ...
func (d warehouseImplement) UpdateOneByCondition(ctx context.Context, cond interface{}, payload interface{}, opts ...*options.UpdateOptions) (err error) {
	var (
		col = database.WarehouseCol()
	)

	if _, err := col.UpdateOne(ctx, cond, payload, opts...); err != nil {
		logger.Error("dao.Warehouse - UpdateOne", logger.LogData{
			"payload": payload,
			"error":   err.Error(),
		})
	}
	return
}

// UpdateStatus ...
func (d warehouseImplement) UpdateStatus(ctx context.Context, cond interface{}, status string) (err error) {
	var (
		col = database.WarehouseCol()
	)
	if _, err := col.UpdateOne(ctx, cond, bson.D{{"$set", bson.D{
		{"status", status},
		{"updatedAt", time.Now()},
	}}}); err != nil {
		logger.Error("dao.Warehouse - UpdateStatus", logger.LogData{
			"status": status,
			"error":  err.Error(),
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

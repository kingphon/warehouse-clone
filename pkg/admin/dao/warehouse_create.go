package dao

import (
	"context"

	"git.selly.red/Selly-Modules/logger"

	mgwarehouse "git.selly.red/Selly-Server/warehouse/external/model/mg/warehouse"
	"git.selly.red/Selly-Server/warehouse/pkg/admin/database"
)

// InsertOne ...
func (d warehouseImplement) InsertOne(ctx context.Context, payload mgwarehouse.Warehouse) error {
	var (
		col = database.WarehouseCol()
	)

	_, err := col.InsertOne(ctx, payload)
	if err != nil {
		logger.Error("dao.Warehouse - InsertOne", logger.LogData{
			"payload": payload,
			"error":   err.Error(),
		})
	}
	return err
}

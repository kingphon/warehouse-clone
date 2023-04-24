package dao

import (
	"context"
	"git.selly.red/Selly-Modules/logger"
	mgwarehouse "git.selly.red/Selly-Server/warehouse/external/model/mg/warehouse"
	"git.selly.red/Selly-Server/warehouse/pkg/admin/database"
)

// InsertOne ...
func (d supplierHolidayImplement) InsertOne(ctx context.Context, payload mgwarehouse.SupplierHoliday) error {
	var (
		col = database.SupplierHolidayCol()
	)

	_, err := col.InsertOne(ctx, payload)
	if err != nil {
		logger.Error("dao.SupplierHoliday - InsertOne", logger.LogData{
			"payload": payload,
			"error":   err.Error(),
		})
	}
	return err
}

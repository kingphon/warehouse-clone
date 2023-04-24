package dao

import (
	"context"
	"git.selly.red/Selly-Modules/logger"
	"git.selly.red/Selly-Server/warehouse/pkg/app/database"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// UpdateOneByCondition ...
func (supplierHolidayImplement) UpdateOneByCondition(ctx context.Context, cond interface{}, payload interface{}, opts ...*options.UpdateOptions) (err error) {
	var (
		col = database.SupplierHolidayCol()
	)

	if _, err := col.UpdateOne(ctx, cond, payload, opts...); err != nil {
		logger.Error("dao.SupplierHoliday - UpdateOne", logger.LogData{
			"payload": payload,
			"error":   err.Error(),
		})
	}
	return
}

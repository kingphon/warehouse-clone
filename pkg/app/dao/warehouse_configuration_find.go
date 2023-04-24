package dao

import (
	"context"

	"git.selly.red/Selly-Modules/logger"
	mgwarehouse "git.selly.red/Selly-Server/warehouse/external/model/mg/warehouse"
	"git.selly.red/Selly-Server/warehouse/pkg/app/database"
	"go.mongodb.org/mongo-driver/bson"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// FindByWarehouseID ...
func (d warehouseConfigurationImplement) FindByWarehouseID(ctx context.Context, id primitive.ObjectID) (doc mgwarehouse.Configuration) {
	if err := database.WarehouseConfigurationCol().FindOne(ctx, bson.M{"warehouse": id}).Decode(&doc); err != nil {
		logger.Error("dao.warehouseConfigurationImplement.FindByWarehouseID", logger.LogData{"err": err.Error()})
	}
	return doc
}

package dao

import (
	"context"

	mgwarehouse "git.selly.red/Selly-Server/warehouse/external/model/mg/warehouse"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// WarehouseConfigurationInterface ...
type WarehouseConfigurationInterface interface {
	// FindByWarehouseID ...
	FindByWarehouseID(ctx context.Context, id primitive.ObjectID) mgwarehouse.Configuration
}

// WarehouseConfigurationImplement ...
type warehouseConfigurationImplement struct{}

// WarehouseConfiguration return courier dao
func WarehouseConfiguration() WarehouseConfigurationInterface {
	return warehouseConfigurationImplement{}
}

package dao

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"

	"git.selly.red/Selly-Modules/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"go.mongodb.org/mongo-driver/mongo/options"

	mgwarehouse "git.selly.red/Selly-Server/warehouse/external/model/mg/warehouse"
	"git.selly.red/Selly-Server/warehouse/pkg/admin/database"
)

/* Entity
- Warehouse
*/

// WarehouseConfigurationInterface ...
type WarehouseConfigurationInterface interface {
	// InsertOne ...
	InsertOne(ctx context.Context, payload mgwarehouse.Configuration) error

	// FindOneByCondition ...
	FindOneByCondition(ctx context.Context, cond interface{}, opts ...*options.FindOneOptions) (doc mgwarehouse.Configuration)

	// FindByWarehouseID ...
	FindByWarehouseID(ctx context.Context, id primitive.ObjectID) mgwarehouse.Configuration

	// FindByCondition ...
	FindByCondition(ctx context.Context, cond interface{}, opts ...*options.FindOptions) (docs []mgwarehouse.Configuration)

	// UpdateOne ...
	UpdateOne(ctx context.Context, cond interface{}, payload interface{}, opts ...*options.UpdateOptions) error

	// UpdateMany ...
	UpdateMany(ctx context.Context, cond interface{}, payload interface{}, opts ...*options.UpdateOptions) error

	// BulkWrite ...
	BulkWrite(ctx context.Context, models []mongo.WriteModel, opts ...*options.BulkWriteOptions) error
}

// WarehouseConfigurationImplement ...
type warehouseConfigurationImplement struct{}

func (d warehouseConfigurationImplement) UpdateMany(ctx context.Context, cond interface{}, payload interface{}, opts ...*options.UpdateOptions) error {
	_, err := database.WarehouseConfigurationCol().UpdateMany(ctx, cond, payload, opts...)
	if err != nil {
		logger.Error("dao.warehouseConfigurationImplement.UpdateMany", logger.LogData{"err": err.Error()})
	}
	return err
}

// FindByWarehouseID ...
func (d warehouseConfigurationImplement) FindByWarehouseID(ctx context.Context, id primitive.ObjectID) (doc mgwarehouse.Configuration) {
	if err := database.WarehouseConfigurationCol().FindOne(ctx, bson.M{"warehouse": id}).Decode(&doc); err != nil {
		logger.Error("dao.warehouseConfigurationImplement.FindByWarehouseID", logger.LogData{"err": err.Error()})
	}
	return doc
}

// WarehouseConfiguration return courier dao
func WarehouseConfiguration() WarehouseConfigurationInterface {
	return warehouseConfigurationImplement{}
}

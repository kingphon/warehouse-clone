package service

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"

	externalauth "git.selly.red/Selly-Server/warehouse/external/auth"
	requestmodel "git.selly.red/Selly-Server/warehouse/pkg/admin/model/request"
	responsemodel "git.selly.red/Selly-Server/warehouse/pkg/admin/model/response"
)

// WarehouseConfigurationInterface ...
type WarehouseConfigurationInterface interface {
	// CreateWithClientData create new warehouseConfiguration from web admin
	CreateWithClientData(ctx context.Context, payload requestmodel.WarehouseCfgCreate, id primitive.ObjectID) error

	// DetailByWarehouseID get detail by warehouse id
	DetailByWarehouseID(ctx context.Context, id primitive.ObjectID) (*responsemodel.ResponseWarehouseConfigurationDetail, error)

	// UpdateSupplier update supplier configuration from web admin
	UpdateSupplier(ctx context.Context, id primitive.ObjectID, payload requestmodel.WarehouseCfgSupplierUpdate) (err error)

	// UpdateFood update food configuration from web admin
	UpdateFood(ctx context.Context, id primitive.ObjectID, payload requestmodel.WarehouseCfgFoodUpdate) (err error)

	// UpdatePartner update partner configuration from web admin
	UpdatePartner(ctx context.Context, id primitive.ObjectID, payload requestmodel.WarehouseCfgPartnerUpdate) (err error)

	// UpdateOrder update order configuration from web admin
	UpdateOrder(ctx context.Context, id primitive.ObjectID, payload requestmodel.WarehouseCfgOrderUpdate) (err error)

	// UpdateDelivery update delivery configuration from web admin
	UpdateDelivery(ctx context.Context, id primitive.ObjectID, payload requestmodel.WarehouseCfgDeliveryUpdate) (err error)

	// UpdateOther update other configuration from web admin
	UpdateOther(ctx context.Context, id primitive.ObjectID, payload requestmodel.WarehouseCfgOtherUpdate) (err error)

	// UpdateOrderConfirm  ...
	UpdateOrderConfirm(ctx context.Context, id primitive.ObjectID, payload requestmodel.WarehouseCfgOrderConfirm) error
}

// WarehouseConfigurationImplement ...
type warehouseConfigurationImplement struct {
	CurrentStaff externalauth.User
}

// WarehouseConfiguration return warehouseConfiguration service
func WarehouseConfiguration(cs externalauth.User) WarehouseConfigurationInterface {
	return warehouseConfigurationImplement{
		CurrentStaff: cs,
	}
}

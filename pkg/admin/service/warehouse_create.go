package service

import (
	"context"
	"errors"

	"git.selly.red/Selly-Modules/natsio/client"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	responsemodel "git.selly.red/Selly-Server/warehouse/pkg/admin/model/response"

	"git.selly.red/Selly-Server/warehouse/external/constant"
	mgwarehouse "git.selly.red/Selly-Server/warehouse/external/model/mg/warehouse"
	"git.selly.red/Selly-Server/warehouse/pkg/admin/dao"
	"git.selly.red/Selly-Server/warehouse/pkg/admin/errorcode"
	requestmodel "git.selly.red/Selly-Server/warehouse/pkg/admin/model/request"
)

//
// PUBLIC METHODS
//

// CreateWithClientData ...
func (s warehouseImplement) CreateWithClientData(ctx context.Context, payload requestmodel.WarehouseCreate) (doc mgwarehouse.Warehouse, err error) {
	var (
		d               = dao.Warehouse()
		warehouseCfgSvc = WarehouseConfiguration(s.CurrentStaff)
		auditSvc        = Audit(s.CurrentStaff)
	)

	// Check supplier invalid and get info
	listSupplierID := []primitive.ObjectID{payload.SupplierID}
	suppliers, err := s.getSupplierByIDs(ctx, listSupplierID)
	if err != nil || len(suppliers) == 0 {
		err = errors.New(errorcode.WarehouseInvalidSupplier)
		return
	}
	supplier := suppliers[0]
	supplierID, _ := primitive.ObjectIDFromHex(supplier.ID)
	if supplierID.IsZero() {
		err = errors.New(errorcode.WarehouseInvalidSupplier)
		return
	}
	if supplier.BusinessType != payload.BusinessType {
		err = errors.New(errorcode.WarehouseCfgInvalidBusinessType)
		return
	}

	// Check warehouse existed by name
	if warehouse := s.isExistedByName(ctx, payload.Name); !warehouse.ID.IsZero() {
		err = errors.New(errorcode.WarehouseExistedName)
		return
	}
	// Create db document
	doc = payload.ConvertToBSON()
	if err = d.InsertOne(ctx, doc); err != nil {
		err = errors.New(errorcode.WarehouseErrorWhenInsert)
	}

	// Create warehouseConfiguration document
	if err := warehouseCfgSvc.CreateWithClientData(ctx, payload.Config, doc.ID); err != nil {
		err = errors.New(err.Error())
	}

	// publish after create
	go func() {
		config := dao.WarehouseConfiguration().FindByWarehouseID(context.Background(), doc.ID)
		response := responsemodel.GetWarehouseNatsResponse(doc, config, s.CurrentStaff.ID)
		client.GetWarehouse().AfterCreateWarehouse(response)
	}()

	// Create audit
	go auditSvc.Create(
		constant.AuditTargetWarehouse,
		doc.ID.Hex(),
		constant.MsgCreateWarehouse,
		constant.AuditActionCreate,
		payload,
	)
	return
}

//
// PRIVATE METHODS
//

// isExistedByCode ...
func (warehouseImplement) isExistedByName(ctx context.Context, name string) mgwarehouse.Warehouse {
	var (
		d    = dao.Warehouse()
		cond = bson.M{"name": name}
	)

	return d.FindOneByCondition(ctx, cond)
}

package service

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson/primitive"

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
func (s warehouseConfigurationImplement) CreateWithClientData(ctx context.Context, payload requestmodel.WarehouseCfgCreate, id primitive.ObjectID) (err error) {
	var (
		d        = dao.WarehouseConfiguration()
		auditSvc = Audit(s.CurrentStaff)
		doc      mgwarehouse.Configuration
	)

	// Create db document
	doc = payload.ConvertToBSON(id)
	if err = d.InsertOne(ctx, doc); err != nil {
		err = errors.New(errorcode.WarehouseErrorWhenInsert)
	}

	// Create audit
	go auditSvc.Create(
		constant.AuditTargetWarehouseConfiguration,
		doc.ID.Hex(),
		constant.MsgCreateWarehouseCfg,
		constant.AuditActionCreate,
		payload,
	)

	go s.updateIsClosedSupplier(doc)

	return
}

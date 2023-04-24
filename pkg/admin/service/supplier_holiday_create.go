package service

import (
	"context"

	"git.selly.red/Selly-Modules/mongodb"
	"go.mongodb.org/mongo-driver/bson"

	"git.selly.red/Selly-Server/warehouse/external/constant"
	mgwarehouse "git.selly.red/Selly-Server/warehouse/external/model/mg/warehouse"
	"git.selly.red/Selly-Server/warehouse/pkg/admin/dao"
	"git.selly.red/Selly-Server/warehouse/pkg/admin/errorcode"
	requestmodel "git.selly.red/Selly-Server/warehouse/pkg/admin/model/request"
	"github.com/friendsofgo/errors"
)

// CreateWithClientData ...
func (s supplierHolidayImplement) CreateWithClientData(ctx context.Context, payload requestmodel.SupplierHolidayCreate) (doc mgwarehouse.SupplierHoliday, err error) {
	var (
		d = dao.SupplierHoliday()
	)

	// Check is existed supplier-holiday
	supplierHoliday := d.FindOneByCondition(ctx, bson.M{"supplier": mongodb.ConvertStringToObjectID(payload.Supplier)})
	if !supplierHoliday.ID.IsZero() {
		err = errors.New(errorcode.SupplierHolidayIsExitedHoliday)
		return
	}
	// 1. Check warehouse invalid
	if err = s.checkInvalidWarehouses(ctx, requestmodel.SupplierHolidayPayload{
		Create: &payload,
	}); err != nil {
		return
	}

	// 2. Create supplier-holiday
	doc = payload.ConvertToBSON(s.CurrentStaff)
	if err = d.InsertOne(ctx, doc); err != nil {
		err = errors.New(errorcode.SupplierHolidayErrorWhenCreate)
	}

	// 4. Create audit
	auditSvc := Audit(s.CurrentStaff)
	go auditSvc.Create(
		constant.AuditTargetWarehouseSupplierHoliday,
		doc.ID.Hex(),
		constant.MsgCreateSupplierHoliday,
		constant.AuditActionCreate,
		doc,
	)

	return
}

package service

import (
	"context"
	"git.selly.red/Selly-Modules/natsio/client"
	natsmodel "git.selly.red/Selly-Modules/natsio/model"
	responsemodel "git.selly.red/Selly-Server/warehouse/pkg/app/model/response"

	"git.selly.red/Selly-Modules/mongodb"
	"go.mongodb.org/mongo-driver/bson"

	"git.selly.red/Selly-Server/warehouse/external/constant"
	mgwarehouse "git.selly.red/Selly-Server/warehouse/external/model/mg/warehouse"
	"git.selly.red/Selly-Server/warehouse/pkg/app/dao"
	"git.selly.red/Selly-Server/warehouse/pkg/app/errorcode"
	requestmodel "git.selly.red/Selly-Server/warehouse/pkg/app/model/request"
	"github.com/friendsofgo/errors"
)

// CreateWithClientData ...
func (s supplierHolidayImplement) CreateWithClientData(ctx context.Context, payload requestmodel.SupplierHolidayCreate) (doc mgwarehouse.SupplierHoliday, err error) {
	var (
		d          = dao.SupplierHoliday()
		supplierId = mongodb.ConvertStringToObjectID(s.CurrentUser.SupplierID)
	)

	// Check is existed supplier-holiday
	supplierHoliday := d.FindOneByCondition(ctx, bson.M{"supplier": supplierId})
	if !supplierHoliday.ID.IsZero() {
		err = errors.New(errorcode.SupplierHolidayIsExitedHoliday)
		return
	}

	// Check warehouse invalid
	if err = s.checkInvalidWarehouses(ctx, requestmodel.SupplierHolidayPayload{
		Create: &payload,
	}); err != nil {
		return
	}

	// Get Supplier
	supplier, err := s.getSupplier(ctx, s.CurrentUser.SupplierID)
	if err != nil {
		return
	}

	// Create supplier-holiday
	doc = payload.ConvertToBSON(supplier)
	if err = d.InsertOne(ctx, doc); err != nil {
		err = errors.New(errorcode.SupplierHolidayErrorWhenCreate)
	}

	//  Update holiday warehouse
	go func() {
		var (
			ctxBg        = context.Background()
			warehouseSvc = warehouseImplement{}
		)

		warehousesByStatuses := warehouseSvc.FindByCondition(ctxBg, bson.M{
			"supplier": supplierId,
			"status": bson.M{"$in": []string{
				constant.WarehouseStatusHoliday,
				constant.StatusActive,
			}},
		})

		warehouseSvc.UpdateWarehousesStatusByIDs(warehousesByStatuses)
	}()

	// Create audit
	auditSvc := Audit(s.CurrentUser)
	go auditSvc.Create(
		constant.AuditTargetWarehouseSupplierHoliday,
		doc.ID.Hex(),
		constant.MsgCreateSupplierHoliday,
		constant.AuditActionCreate,
		doc,
	)

	return
}

// getSupplier ...
func (s supplierHolidayImplement) getSupplier(ctx context.Context, supplierID string) (result responsemodel.ResponseSupplierInfo, err error) {
	suppliers, err := client.GetSupplier().GetDetailSupplierInfo(natsmodel.GetDetailSupplierRequest{Supplier: supplierID})
	if err != nil {
		return
	}

	result = responsemodel.ResponseSupplierInfo{
		ID:   suppliers.ID,
		Name: suppliers.Name,
	}

	return
}

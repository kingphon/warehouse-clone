package service

import (
	"context"
	"errors"
	"git.selly.red/Selly-Modules/natsio/model"
	"time"

	"git.selly.red/Selly-Modules/logger"
	"git.selly.red/Selly-Server/warehouse/external/utils/ptime"

	"go.mongodb.org/mongo-driver/mongo"

	"git.selly.red/Selly-Modules/natsio/client"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	responsemodel "git.selly.red/Selly-Server/warehouse/pkg/admin/model/response"

	mgwarehouse "git.selly.red/Selly-Server/warehouse/external/model/mg/warehouse"

	"git.selly.red/Selly-Server/warehouse/external/constant"
	"git.selly.red/Selly-Server/warehouse/pkg/admin/dao"
	"git.selly.red/Selly-Server/warehouse/pkg/admin/errorcode"
	requestmodel "git.selly.red/Selly-Server/warehouse/pkg/admin/model/request"
)

//
// PUBLIC METHODS
//

// UpdateWithClientData ...
func (s warehouseImplement) UpdateWithClientData(ctx context.Context, id primitive.ObjectID, payload requestmodel.WarehouseUpdate) (err error) {
	var (
		d        = dao.Warehouse()
		cond     = bson.M{"_id": id}
		auditSvc = Audit(s.CurrentStaff)
	)

	bodyConvert := payload.ConvertToBSON()
	update := bson.D{{"$set",
		bson.D{
			{"name", bodyConvert.Name},
			{"searchString", bodyConvert.SearchString},
			{"slug", bodyConvert.Slug},
			{"location", bodyConvert.Location},
			{"contact", bodyConvert.Contact},
			{"updatedAt", bodyConvert.UpdatedAt},
		},
	}}

	// Check route existed by ID
	doc := d.FindOneByCondition(ctx, cond)
	if doc.ID.IsZero() {
		err = errors.New(errorcode.WarehouseNotFound)
		return
	}

	doc.Name = bodyConvert.Name
	doc.SearchString = bodyConvert.SearchString
	doc.Slug = bodyConvert.Slug
	doc.Location = mgwarehouse.Location{
		Province: bodyConvert.Location.Province,
		District: bodyConvert.Location.District,
		Ward:     bodyConvert.Location.Ward,
		Address:  bodyConvert.Location.Address,
		LocationCoordinates: mgwarehouse.Coordinates{
			Type:        bodyConvert.Location.LocationCoordinates.Type,
			Coordinates: bodyConvert.Location.LocationCoordinates.Coordinates,
		},
	}
	doc.Contact = mgwarehouse.Contact{
		Name:    bodyConvert.Contact.Name,
		Phone:   bodyConvert.Contact.Phone,
		Address: bodyConvert.Contact.Address,
		Email:   bodyConvert.Contact.Email,
	}
	doc.UpdatedAt = bodyConvert.UpdatedAt

	go auditSvc.Create(
		constant.AuditTargetWarehouse,
		id.Hex(),
		constant.MsgEditWarehouse,
		constant.AuditActionEdit,
		payload,
	)

	err = d.UpdateOneByCondition(ctx, cond, update)
	if err != nil {
		return
	}
	go s.AfterUpdateWarehouse(doc, s.CurrentStaff.ID)
	return
}

// UpdateStatus ...
func (s warehouseImplement) UpdateStatus(ctx context.Context, id primitive.ObjectID, payload requestmodel.WarehouseUpdateStatus) (err error) {
	var (
		d        = dao.Warehouse()
		cond     = bson.M{"_id": id}
		auditSvc = Audit(s.CurrentStaff)
	)

	// Check route existed by ID
	doc := d.FindOneByCondition(ctx, cond)
	if doc.ID.IsZero() {
		err = errors.New(errorcode.WarehouseNotFound)
		return
	}

	// Get supplier-contract
	supplierContract, err := s.getSupplierContact(ctx, doc.Supplier)
	if err != nil {
		err = errors.New(errorcode.SupplierContractInvalidStatus)
		return
	}

	if doc.Status == constant.StatusInactive && (supplierContract.Status != constant.ContractStatus.Approved && supplierContract.Status != constant.ContractStatus.Completed) {
		return errors.New(errorcode.SupplierContractInvalidStatus)
	}

	err = d.UpdateStatus(ctx, cond, payload.Status)
	if err != nil {
		return err
	}

	go auditSvc.Create(
		constant.AuditTargetWarehouse,
		id.Hex(),
		constant.MsgEditStatusWarehouse,
		constant.AuditActionEdit,
		payload,
	)

	doc.Status = payload.Status
	doc.UpdatedAt = time.Now()
	go s.AfterUpdateWarehouse(doc, s.CurrentStaff.ID)
	return nil
}

// AfterUpdateWarehouse ...
func (warehouseImplement) AfterUpdateWarehouse(doc mgwarehouse.Warehouse, staff string) {
	ctxbg := context.Background()
	c := dao.WarehouseConfiguration().FindByWarehouseID(ctxbg, doc.ID)
	w := responsemodel.GetWarehouseNatsResponse(doc, c, staff)
	client.GetWarehouse().AfterUpdateWarehouse(w)
}

// BulkWriteWhenChangeStatusSupplierHoliday ...
func (warehouseImplement) BulkWriteWhenChangeStatusSupplierHoliday(ctx context.Context, wModel []mongo.WriteModel) error {
	return dao.Warehouse().BulkWrite(ctx, wModel)
}

// UpdateManyByCondition ...
func (warehouseImplement) UpdateManyByCondition(ctx context.Context, cond interface{}, payload interface{}) error {
	return dao.Warehouse().UpdateManyByCondition(ctx, cond, payload)
}

// UpdateWarehousesStatusByIDs ...
func (s warehouseImplement) UpdateWarehousesStatusByIDs(warehouses []mgwarehouse.Warehouse) {
	var (
		wModel = make([]mongo.WriteModel, 0)
		ctx    = context.Background()
	)

	listWarehouseUpdatePendingInActive := make([]model.UpdateStatusWarehousePendingInactive, 0)

	for _, w := range warehouses {
		// Check isActive holiday
		newStatus, holiday := s.checkWarehouseActiveHoliday(ctx, w)

		update := bson.M{
			"updatedAt": ptime.Now(),
		}

		if newStatus != w.Status {
			update["status"] = newStatus
			update["reasonPendingInactive"] = ""
			if newStatus == constant.WarehouseStatusHoliday {
				update["statusBeforeHoliday"] = w.Status
				update["reasonPendingInactive"] = holiday.Reason
			}

			listWarehouseUpdatePendingInActive = append(listWarehouseUpdatePendingInActive, model.UpdateStatusWarehousePendingInactive{
				WarehouseID:     w.ID.Hex(),
				PendingInactive: newStatus == constant.WarehouseStatusHoliday,
			})
		}

		if newStatus == constant.WarehouseStatusHoliday {
			update["reasonPendingInactive"] = holiday.Reason
		}

		wModel = append(wModel, mongo.NewUpdateOneModel().SetFilter(bson.M{
			"_id": w.ID,
		}).SetUpdate(bson.M{
			"$set": update,
		}))
	}

	if len(wModel) > 0 {
		var warehouseSvc = warehouseImplement{}
		if err := warehouseSvc.BulkWriteWhenChangeStatusSupplierHoliday(ctx, wModel); err != nil {
			logger.Error("Error dao.Warehouse().BulkWrite in supplier-holiday service : ", logger.LogData{
				"error": err.Error(),
			})
		}
	}

	// Update pendingInactive product
	var bodyUpdateProduct = model.UpdateStatusWarehousePendingInactiveRequest{Warehouses: listWarehouseUpdatePendingInActive}
	if err := client.GetWarehouse().UpdateStatusWarehousePendingInactive(bodyUpdateProduct); err != nil {
		logger.Error("Error Update pending inActive product by warehouseIDs : ", logger.LogData{
			"error": err.Error(),
		})
	}

}

//
// PRIVATE METHODS
//

// getSupplierContact ...
func (warehouseImplement) getSupplierContact(ctx context.Context, supplierID primitive.ObjectID) (*model.ResponseSupplierContract, error) {
	body := model.GetSupplierContractRequest{SupplierID: supplierID}
	return client.GetSupplier().GetSupplierContractBySupplierID(body)
}

// checkWarehouseActiveHoliday ...
func (S warehouseImplement) checkWarehouseActiveHoliday(ctx context.Context, warehouse mgwarehouse.Warehouse) (string, mgwarehouse.SupplierHoliday) {
	var (
		supplierHolidaySvc = supplierHolidayImplement{CurrentStaff: S.CurrentStaff}
		cond               = bson.M{
			"supplier": warehouse.Supplier,
			"status":   "active",
			"from":     bson.M{"$lte": ptime.Now()}, // from < now < to
			"to":       bson.M{"$gte": ptime.Now()},
			"$or": []bson.M{
				{
					"isApplyAll": true,
				},
				{
					"isApplyAll": false,
					"warehouses": warehouse.ID,
				},
			},
		}
	)

	holidays := supplierHolidaySvc.FindByCondition(ctx, cond)

	if len(holidays) == 0 {
		if warehouse.StatusBeforeHoliday == "" {
			return warehouse.Status, mgwarehouse.SupplierHoliday{}
		}
		return warehouse.StatusBeforeHoliday, mgwarehouse.SupplierHoliday{}
	}

	return constant.WarehouseStatusHoliday, holidays[0]
}

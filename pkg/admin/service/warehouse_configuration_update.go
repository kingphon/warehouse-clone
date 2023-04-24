package service

import (
	"context"
	"errors"

	"git.selly.red/Selly-Modules/natsio/client"
	"git.selly.red/Selly-Modules/natsio/model"
	responsemodel "git.selly.red/Selly-Server/warehouse/pkg/admin/model/response"

	mgwarehouse "git.selly.red/Selly-Server/warehouse/external/model/mg/warehouse"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"git.selly.red/Selly-Server/warehouse/external/constant"
	"git.selly.red/Selly-Server/warehouse/pkg/admin/dao"
	"git.selly.red/Selly-Server/warehouse/pkg/admin/errorcode"
	requestmodel "git.selly.red/Selly-Server/warehouse/pkg/admin/model/request"
)

// UpdateOrderConfirm ...
func (s warehouseConfigurationImplement) UpdateOrderConfirm(ctx context.Context, id primitive.ObjectID, payload requestmodel.WarehouseCfgOrderConfirm) (err error) {
	var (
		d        = dao.WarehouseConfiguration()
		cond     = bson.M{"_id": id}
		auditSvc = Audit(s.CurrentStaff)
	)

	update := bson.M{
		"$set": bson.M{
			"autoConfirmOrder": payload.OrderConfirm,
		},
	}

	// Check route existed by ID
	doc := d.FindOneByCondition(ctx, cond)
	if doc.ID.IsZero() {
		err = errors.New(errorcode.WarehouseCfgNotFound)
		return
	}

	err = d.UpdateOne(ctx, cond, update)
	if err != nil {
		return err
	}

	// after update config
	doc.AutoConfirmOrder = payload.OrderConfirm
	go s.afterUpdateWarehouseConfig(doc, s.CurrentStaff.ID)

	// Create audit
	go auditSvc.Create(
		constant.AuditTargetWarehouseConfiguration,
		id.Hex(),
		constant.MsgEditWarehouseCfgSupplier,
		constant.AuditActionEdit,
		payload,
	)
	return nil
}

// UpdateFood ...
func (s warehouseConfigurationImplement) UpdateFood(ctx context.Context, id primitive.ObjectID, payload requestmodel.WarehouseCfgFoodUpdate) (err error) {
	var (
		d        = dao.WarehouseConfiguration()
		cond     = bson.M{"_id": id}
		auditSvc = Audit(s.CurrentStaff)
	)

	bodyConvert := payload.ConvertTOBSON()
	update := bson.M{
		"$set": bson.M{
			"food": bodyConvert,
		},
	}

	// Check route existed by ID
	doc := d.FindOneByCondition(ctx, cond)
	if doc.ID.IsZero() {
		err = errors.New(errorcode.WarehouseCfgNotFound)
		return
	}

	err = d.UpdateOne(ctx, cond, update)
	if err != nil {
		return err
	}

	// after update config
	doc.Food = mgwarehouse.ConfigFood{
		ForceClosed: bodyConvert.ForceClosed,
		IsClosed:    bodyConvert.IsClosed || bodyConvert.ForceClosed,
		TimeRange:   make([]mgwarehouse.TimeRange, 0),
	}
	if len(bodyConvert.TimeRange) > 0 {
		for _, t := range bodyConvert.TimeRange {
			doc.Food.TimeRange = append(doc.Food.TimeRange, mgwarehouse.TimeRange{
				From: t.From,
				To:   t.To,
			})
		}
	}
	go s.afterUpdateWarehouseConfig(doc, s.CurrentStaff.ID)

	go s.updateIsClosedSupplier(doc)

	// Create audit
	go auditSvc.Create(
		constant.AuditTargetWarehouseConfiguration,
		id.Hex(),
		constant.MsgEditWarehouseCfgSupplier,
		constant.AuditActionEdit,
		payload,
	)
	return nil
}

// UpdateSupplier ...
func (s warehouseConfigurationImplement) UpdateSupplier(ctx context.Context, id primitive.ObjectID, payload requestmodel.WarehouseCfgSupplierUpdate) (err error) {
	var (
		d        = dao.WarehouseConfiguration()
		cond     = bson.M{"_id": id}
		auditSvc = Audit(s.CurrentStaff)
	)

	bodyConvert := payload.ConvertTOBSON()
	update := bson.D{{"$set",
		bson.D{
			{"supplier", bodyConvert},
		},
	}}

	// Check route existed by ID
	doc := d.FindOneByCondition(ctx, cond)
	if doc.ID.IsZero() {
		err = errors.New(errorcode.WarehouseCfgNotFound)
		return
	}
	// after update config
	doc.Supplier = mgwarehouse.ConfigSupplier{
		InvoiceDeliveryMethod: bodyConvert.InvoiceDeliveryMethod,
	}
	go s.afterUpdateWarehouseConfig(doc, s.CurrentStaff.ID)

	// Create audit
	go auditSvc.Create(
		constant.AuditTargetWarehouseConfiguration,
		id.Hex(),
		constant.MsgEditWarehouseCfgSupplier,
		constant.AuditActionEdit,
		payload,
	)

	return d.UpdateOne(ctx, cond, update)
}

// UpdatePartner ...
func (s warehouseConfigurationImplement) UpdatePartner(ctx context.Context, id primitive.ObjectID, payload requestmodel.WarehouseCfgPartnerUpdate) (err error) {
	var (
		d        = dao.WarehouseConfiguration()
		cond     = bson.M{"_id": id}
		auditSvc = Audit(s.CurrentStaff)
	)

	bodyConvert := payload.ConvertToBSON()
	update := bson.D{{"$set",
		bson.D{
			{"partner", bodyConvert},
		},
	}}

	// Check route existed by ID
	doc := d.FindOneByCondition(ctx, cond)
	if doc.ID.IsZero() {
		err = errors.New(errorcode.WarehouseCfgNotFound)
		return
	}
	// after update config
	doc.Partner = mgwarehouse.ConfigPartner{
		IdentityCode:   bodyConvert.IdentityCode,
		Code:           bodyConvert.Code,
		Enabled:        bodyConvert.Enabled,
		Authentication: bodyConvert.Authentication,
	}
	go s.afterUpdateWarehouseConfig(doc, s.CurrentStaff.ID)

	// Create audit
	go auditSvc.Create(
		constant.AuditTargetWarehouseConfiguration,
		id.Hex(),
		constant.MsgEditWarehouseCfgPartner,
		constant.AuditActionEdit,
		payload,
	)

	return d.UpdateOne(ctx, cond, update)
}

// UpdateOrder ...
func (s warehouseConfigurationImplement) UpdateOrder(ctx context.Context, id primitive.ObjectID, payload requestmodel.WarehouseCfgOrderUpdate) (err error) {
	var (
		d        = dao.WarehouseConfiguration()
		cond     = bson.M{"_id": id}
		auditSvc = Audit(s.CurrentStaff)
	)

	bodyConvert := payload.ConvertToBSON()
	update := bson.D{{"$set",
		bson.D{
			{"order", bodyConvert},
		},
	}}

	// Check route existed by ID
	doc := d.FindOneByCondition(ctx, cond)
	if doc.ID.IsZero() {
		err = errors.New(errorcode.WarehouseCfgNotFound)
		return
	}
	// after update config
	doc.Order = mgwarehouse.ConfigOrder{
		MinimumValue: bodyConvert.MinimumValue,
		PaymentMethod: mgwarehouse.ConfigPaymentMethod{
			Cod:          bodyConvert.PaymentMethod.Cod,
			BankTransfer: bodyConvert.PaymentMethod.BankTransfer,
		},
		IsLimitNumberOfPurchases: bodyConvert.IsLimitNumberOfPurchases,
		LimitNumberOfPurchases:   bodyConvert.LimitNumberOfPurchases,
	}
	go s.afterUpdateWarehouseConfig(doc, s.CurrentStaff.ID)

	// Create audit
	go auditSvc.Create(
		constant.AuditTargetWarehouseConfiguration,
		id.Hex(),
		constant.MsgEditWarehouseCfgOrder,
		constant.AuditActionEdit,
		payload,
	)

	return d.UpdateOne(ctx, cond, update)
}

// UpdateDelivery ...
func (s warehouseConfigurationImplement) UpdateDelivery(ctx context.Context, id primitive.ObjectID, payload requestmodel.WarehouseCfgDeliveryUpdate) (err error) {
	var (
		d        = dao.WarehouseConfiguration()
		cond     = bson.M{"_id": id}
		auditSvc = Audit(s.CurrentStaff)
	)

	bodyConvert := payload.ConvertToBSON()
	update := bson.D{{"$set",
		bson.D{
			{"delivery", bodyConvert},
		},
	}}

	// Check route existed by ID
	doc := d.FindOneByCondition(ctx, cond)
	if doc.ID.IsZero() {
		err = errors.New(errorcode.WarehouseCfgNotFound)
		return
	}
	// after update config
	doc.Delivery = mgwarehouse.ConfigDelivery{
		DeliveryMethods:      bodyConvert.DeliveryMethods,
		PriorityServiceCodes: bodyConvert.PriorityServiceCodes,
		EnabledSources:       bodyConvert.EnabledSources,
		Types:                bodyConvert.Types,
	}
	go s.afterUpdateWarehouseConfig(doc, s.CurrentStaff.ID)

	// Create audit
	go auditSvc.Create(
		constant.AuditTargetWarehouseConfiguration,
		id.Hex(),
		constant.MsgEditWarehouseCfgDelivery,
		constant.AuditActionEdit,
		payload,
	)

	return d.UpdateOne(ctx, cond, update)
}

// UpdateOther ...
func (s warehouseConfigurationImplement) UpdateOther(ctx context.Context, id primitive.ObjectID, payload requestmodel.WarehouseCfgOtherUpdate) (err error) {
	var (
		d        = dao.WarehouseConfiguration()
		cond     = bson.M{"_id": id}
		auditSvc = Audit(s.CurrentStaff)
	)

	bodyConvert := payload.ConvertToBSON()
	update := bson.D{{"$set",
		bson.D{
			{"other", bodyConvert},
		},
	}}

	// Check route existed by ID
	doc := d.FindOneByCondition(ctx, cond)
	if doc.ID.IsZero() {
		err = errors.New(errorcode.WarehouseCfgNotFound)
		return
	}
	// after update config
	doc.Other = mgwarehouse.ConfigOther{
		DoesSupportSellyExpress: bodyConvert.DoesSupportSellyExpress,
	}
	go s.afterUpdateWarehouseConfig(doc, s.CurrentStaff.ID)

	// Create audit
	go auditSvc.Create(
		constant.AuditTargetWarehouseConfiguration,
		id.Hex(),
		constant.MsgEditWarehouseCfgOther,
		constant.AuditActionEdit,
		payload,
	)

	return d.UpdateOne(ctx, cond, update)
}

func (s warehouseConfigurationImplement) updateIsClosedSupplier(doc mgwarehouse.Configuration) {
	var (
		ctx                = context.Background()
		listUpdateSupplier = make([]model.SupplierIsClosed, 0)
	)
	wh := dao.Warehouse().FindOneByCondition(ctx, bson.M{
		"_id": doc.Warehouse,
	})

	if wh.ID.IsZero() {
		return
	}

	warehouses := dao.Warehouse().FindByCondition(ctx, bson.M{
		"supplier": wh.Supplier,
	})
	listWarehouseID := make([]primitive.ObjectID, 0)
	for _, warehouse := range warehouses {
		listWarehouseID = append(listWarehouseID, warehouse.ID)
	}
	configurations := dao.WarehouseConfiguration().FindByCondition(ctx, bson.M{
		"warehouse": bson.M{
			"$in": listWarehouseID,
		},
	})
	if len(configurations) == 0 {
		listUpdateSupplier = append(listUpdateSupplier, model.SupplierIsClosed{
			Supplier: wh.Supplier.Hex(),
			IsClosed: true,
		})
		return
	}
	var isClosed = true
	for _, cf := range configurations {
		if !cf.Food.IsClosed {
			isClosed = false
			break
		}
	}
	listUpdateSupplier = append(listUpdateSupplier, model.SupplierIsClosed{
		Supplier: wh.Supplier.Hex(),
		IsClosed: isClosed,
	})
	client.GetWarehouse().UpdateIsClosedSupplier(model.UpdateSupplierIsClosedRequest{Suppliers: listUpdateSupplier})
}

func (s warehouseConfigurationImplement) afterUpdateWarehouseConfig(doc mgwarehouse.Configuration, staff string) {
	warehouse := dao.Warehouse().FindOneByCondition(context.Background(), bson.M{
		"_id": doc.Warehouse,
	})
	if !warehouse.ID.IsZero() {
		w := responsemodel.GetWarehouseNatsResponse(warehouse, doc, staff)
		client.GetWarehouse().AfterUpdateWarehouse(w)
	}
}

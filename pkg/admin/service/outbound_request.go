package service

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"git.selly.red/Selly-Modules/logger"
	"git.selly.red/Selly-Modules/mongodb"
	"git.selly.red/Selly-Modules/natsio/client"
	"git.selly.red/Selly-Modules/natsio/model"
	"github.com/panjf2000/ants/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"

	externalauth "git.selly.red/Selly-Server/warehouse/external/auth"
	"git.selly.red/Selly-Server/warehouse/external/constant"
	mgwarehouse "git.selly.red/Selly-Server/warehouse/external/model/mg/warehouse"
	"git.selly.red/Selly-Server/warehouse/external/response"
	"git.selly.red/Selly-Server/warehouse/external/utils/pstring"
	"git.selly.red/Selly-Server/warehouse/external/utils/ptime"
	"git.selly.red/Selly-Server/warehouse/pkg/admin/config"
	"git.selly.red/Selly-Server/warehouse/pkg/admin/dao"
	"git.selly.red/Selly-Server/warehouse/pkg/admin/errorcode"
	requestmodel "git.selly.red/Selly-Server/warehouse/pkg/admin/model/request"
	responsemodel "git.selly.red/Selly-Server/warehouse/pkg/admin/model/response"
)

type OutboundRequestInterface interface {
	GetList(ctx context.Context, q requestmodel.OutboundRequestQuery) ([]responsemodel.OutboundRequest, error)

	Create(req *model.OutboundRequestPayload) (*model.OutboundRequestResponse, error)
	UpdateLogisticInfo(req *model.UpdateOutboundRequestLogisticInfoPayload) error
	Cancel(req *model.CancelOutboundRequest) error
	UpdateStatusFromWebhook(b []byte, partnerIdentityCode string) error
	UpdateORDeliveryStatus(req *model.WarehouseORUpdateDeliveryStatus) error

	SyncStatus(req *model.SyncORStatusRequest) (*model.SyncORStatusResponse, error)
	SyncAllGlobalCare()
}

// OutboundRequest ...
func OutboundRequest() OutboundRequestInterface {
	return outboundRequestImplement{}
}

type outboundRequestImplement struct{}

// UpdateORDeliveryStatus ...
func (s outboundRequestImplement) UpdateORDeliveryStatus(req *model.WarehouseORUpdateDeliveryStatus) error {
	ctx := context.Background()
	doc, orSvc, err := s.getORInfoByCode(ctx, req.ORCode, req.OrderID)
	if err != nil {
		return err
	}
	// TODO: update delivery to OR if needed
	return orSvc.UpdateDeliveryStatus(*req, doc)
}

// GetList ...
func (s outboundRequestImplement) GetList(ctx context.Context, q requestmodel.OutboundRequestQuery) ([]responsemodel.OutboundRequest, error) {
	orderID, _ := mongodb.NewIDFromString(q.Order)
	cond := bson.M{"order": orderID}
	d := dao.OutboundRequest()
	docs := d.FindByCondition(ctx, cond)
	res := make([]responsemodel.OutboundRequest, len(docs))

	for i, doc := range docs {
		res[i] = s.getResponse(doc)
	}
	return res, nil
}

// Create ...
func (s outboundRequestImplement) Create(req *model.OutboundRequestPayload) (*model.OutboundRequestResponse, error) {
	ctx := context.Background()
	whID := mongodb.ConvertStringToObjectID(req.WarehouseID)
	whCfgService := WarehouseConfiguration(externalauth.User{})

	whCfg, err := whCfgService.DetailByWarehouseID(ctx, whID)
	if err != nil {
		return nil, err
	}
	partner := whCfg.Partner
	if !partner.Enabled {
		return nil, fmt.Errorf("warehouse_not_enable_partner_api")
	}
	var (
		now = ptime.Now()
		d   = dao.OutboundRequest()
	)
	// check order has OR success first
	orderID := mongodb.ConvertStringToObjectID(req.OrderID)
	or := s.getOutboundRequestSuccessByOrderID(ctx, orderID)
	if !or.ID.IsZero() {
		// skip if not change TPL code
		if or.TPLCode == req.TPLCode {
			return &model.OutboundRequestResponse{
				OrderCode:    or.OrderCode,
				TrackingCode: or.TrackingCode,
				ID:           or.ID.Hex(),
				ORCode:       or.Partner.Code,
				RequestID:    or.Partner.RequestID,
				Status:       or.Status,
				Reason:       or.Reason,
			}, nil
		}
		if or.Partner.IdentityCode == constant.WarehousePartnerCodeOnPoint {
			// update delivery info for onpoint
			err = s.UpdateLogisticInfo(&model.UpdateOutboundRequestLogisticInfoPayload{
				ShippingLabel: "",
				TrackingCode:  req.TrackingCode,
				ORCode:        or.Partner.Code,
				TPLCode:       req.TPLCode,
				OrderID:       req.OrderID,
			})
			if err != nil {
				return nil, err
			}
			// update to db
			update := bson.M{
				"trackingCode": req.TrackingCode,
				"tplCode":      req.TPLCode,
			}
			if _, err = dao.OutboundRequest().UpdateByID(ctx, or.ID, bson.M{"$set": update}); err != nil {
				return nil, err
			}
			return &model.OutboundRequestResponse{
				OrderCode:    or.OrderCode,
				TrackingCode: req.TrackingCode,
				ID:           or.ID.Hex(),
				ORCode:       or.Partner.Code,
				RequestID:    or.Partner.RequestID,
				Status:       or.Status,
				Reason:       or.Reason,
			}, nil
		}
		// cancel old OR and recreate below
		if err = s.Cancel(&model.CancelOutboundRequest{
			ORCode:  or.Partner.Code,
			OrderID: or.Order.Hex(),
			Note:    errorcode.ORMessageAutoCancelOR,
		}); err != nil {
			return nil, fmt.Errorf("warehouse: cancel_or err: %v", err)
		}
	}
	orSvc, err := NewOutboundRequestService(partner.IdentityCode, partner.Authentication)
	if err != nil {
		return nil, err
	}

	doc := mgwarehouse.OutboundRequest{
		ID:        primitive.NewObjectID(),
		Warehouse: whID,
		Supplier:  mongodb.ConvertStringToObjectID(req.SupplierID),
		Order:     orderID,
		Status:    constant.WarehouseORStatusWaitingToConfirm,
		Partner: mgwarehouse.PartnerOutboundRequest{
			IdentityCode: partner.IdentityCode,
		},
		CreatedAt:    now,
		UpdatedAt:    now,
		TrackingCode: req.TrackingCode,
		OrderCode:    req.OrderCode,
		TPLCode:      req.TPLCode,
	}
	if err = d.InsertOne(ctx, doc); err != nil {
		return nil, err
	}

	go func() { _, _ = OutboundRequestHistory().SaveByOR(context.Background(), doc) }()

	var res *model.OutboundRequestResponse
	defer func() {
		// update OR status
		update := bson.M{
			"updatedAt":         ptime.Now(),
			"status":            res.Status,
			"partner.code":      res.ORCode,
			"partner.requestId": res.RequestID,
		}
		if err != nil {
			res.Status = constant.WarehouseORStatusFailed
			res.Reason = err.Error()
			update["reason"] = res.Reason
			update["status"] = res.Status
		} else if res.ORCode == "" {
			// set flag for some partner not response OR code in creation OR time, update in webhook later
			update["partner.isWaitingCode"] = true
		}

		_, err = d.UpdateByID(ctx, doc.ID, bson.M{"$set": update})

		if res.Status != doc.Status && res.Status != "" {
			doc.Status = res.Status
			go func() { _, _ = OutboundRequestHistory().SaveByOR(context.Background(), doc) }()
		}
	}()
	s.preprocessPayload(req)
	res, err = orSvc.Create(*req, partner)
	if res == nil {
		res = &model.OutboundRequestResponse{}
	}
	if err != nil {
		if c := response.GetByKey(err.Error()); c.Code > 0 {
			err = errors.New(c.Message)
		}
	}
	res.ID = doc.ID.Hex()
	return res, err
}

func (s outboundRequestImplement) preprocessPayload(p *model.OutboundRequestPayload) {
	products := make([]model.OutboundRequestItem, 0)
	for _, item := range p.Items {
		// check exist sku first
		sku := pstring.TrimSpace(item.SupplierSKU)
		if sku == "" {
			continue
		}
		index := -1
		for i, product := range products {
			if product.SupplierSKU == sku {
				index = i
				break
			}
		}

		// sum quantity if exists
		if index > -1 {
			products[index].Quantity += item.Quantity
		} else {
			products = append(products, item)
		}
	}
	p.Items = products
}

func (s outboundRequestImplement) getOutboundRequestSuccessByOrderID(ctx context.Context, orderID primitive.ObjectID) mgwarehouse.OutboundRequest {
	cond := bson.D{
		{"order", orderID},
		{"status", bson.M{
			"$in": []string{
				constant.WarehouseORStatusWaitingToConfirm,
				constant.WarehouseORStatusConfirmed,
				constant.WarehouseORStatusPacked,
				constant.WarehouseORStatusPicked,
				constant.WarehouseORStatusReturned,
			},
		}},
	}
	d := dao.OutboundRequest()
	opts := options.FindOne().SetSort(bson.M{"_id": -1})
	or := d.FindOne(ctx, cond, opts)
	return or
}

// UpdateLogisticInfo ...
func (s outboundRequestImplement) UpdateLogisticInfo(req *model.UpdateOutboundRequestLogisticInfoPayload) error {
	ctx := context.Background()
	d := dao.OutboundRequest()
	doc, orSvc, err := s.getORInfoByCode(ctx, req.ORCode, req.OrderID)
	if err != nil {
		return err
	}
	if err = orSvc.UpdateLogisticInfo(*req, doc); err != nil {
		return err
	}
	// update info to OR
	update := bson.M{
		"trackingCode":  req.TrackingCode,
		"shippingLabel": req.ShippingLabel,
		"updatedAt":     ptime.Now(),
	}
	_, err = d.UpdateByID(ctx, doc.ID, bson.M{"$set": update})
	return err
}

// Cancel ...
func (s outboundRequestImplement) Cancel(req *model.CancelOutboundRequest) error {
	ctx := context.Background()
	d := dao.OutboundRequest()
	doc, orSvc, err := s.getORInfoByCode(ctx, req.ORCode, req.OrderID)
	if err != nil {
		return err
	}

	if err = orSvc.Cancel(doc, req.Note); err != nil {
		logger.Error("service.OutboundRequest.Cancel", logger.LogData{
			"err":    err.Error(),
			"orCode": req.ORCode,
		})
		return err
	}

	update := bson.M{
		"status":    constant.WarehouseORStatusCancelled,
		"updatedAt": ptime.Now(),
	}
	_, err = d.UpdateByID(ctx, doc.ID, bson.M{"$set": update})
	return err
}

// UpdateStatusFromWebhook ...
func (s outboundRequestImplement) UpdateStatusFromWebhook(b []byte, partnerIdentityCode string) error {
	orSvc, err := NewOutboundRequestService(partnerIdentityCode, "")
	if err != nil {
		return err
	}
	data, err := orSvc.GetWebhookData(b)
	if err != nil {
		return err
	}
	if data == nil || data.Status == "" {
		return nil
	}
	ctx := context.Background()
	// get OR
	d := dao.OutboundRequest()
	cond := bson.M{"orderCode": data.OrderCode}
	if data.ORRequestID != "" {
		cond["partner.requestId"] = data.ORRequestID
	}
	if data.ORCode != "" {
		cond["$or"] = []bson.M{
			{"partner.code": data.ORCode},
			{
				"partner.isWaitingCode": true,
				"status":                bson.M{"$ne": constant.WarehouseORStatusFailed},
			},
		}
	}
	or := d.FindOne(ctx, cond)
	if or.ID.IsZero() {
		return errors.New(errorcode.ORNotFound)
	}
	// update status
	update := bson.M{
		"status":    data.Status,
		"updatedAt": data.UpdatedAt,
		"reason":    data.Reason,
	}
	if data.ORCode != "" && or.Partner.IsWaitingCode {
		update["partner.code"] = data.ORCode
		update["partner.isWaitingCode"] = false
	}
	if data.DeliveryStatus != "" {
		update["deliveryStatus"] = data.DeliveryStatus
		or.DeliveryStatus = data.DeliveryStatus
	}
	if data.Link != "" {
		update["link"] = data.Link
		or.Link = data.Link
	}
	_, err = d.UpdateByID(ctx, or.ID, bson.M{"$set": update})
	if err != nil {
		return err
	}

	or.Status = data.Status
	or.UpdatedAt = data.UpdatedAt
	or.Reason = data.Reason

	// update data to order
	go func() {
		payload := model.OrderUpdateORStatus{
			ID:             or.ID.Hex(),
			OrderCode:      or.OrderCode,
			ORCode:         or.Partner.Code,
			Status:         or.Status,
			DeliveryStatus: or.DeliveryStatus,
			Reason:         or.Reason,
			Data:           model.OrderORData{Link: or.Link},
		}
		if err = client.GetOrder().UpdateORStatus(payload); err != nil {
			logger.Error("client.GetOrder().UpdateORStatus", logger.LogData{
				"err":     err.Error(),
				"payload": payload,
			})
		}
		// save history
		go func() { _, _ = OutboundRequestHistory().SaveByOR(context.Background(), or) }()
	}()
	return nil
}

// SyncStatus ...
func (s outboundRequestImplement) SyncStatus(req *model.SyncORStatusRequest) (*model.SyncORStatusResponse, error) {
	ctx := context.Background()
	or, p, err := s.getORInfoByCode(ctx, req.ORCode, req.OrderID)
	if err != nil {
		return nil, err
	}
	res, err := p.GetORStatus(or)
	if err != nil || res == nil {
		return nil, err
	}
	// update latest status
	if res.Status != "" && res.Status != or.Status {
		update := bson.M{
			"status":    res.Status,
			"updatedAt": ptime.Now(),
		}
		d := dao.OutboundRequest()
		if _, err = d.UpdateByID(ctx, or.ID, bson.M{"$set": update}); err != nil {
			return nil, err
		}
	}
	return res, nil
}

// SyncStatusByOR ...
func (s outboundRequestImplement) SyncStatusByOR(ctx context.Context, or mgwarehouse.OutboundRequest) (*responsemodel.OutboundRequestStatus, error) {
	p, err := s.getWarehousePartnerByDoc(ctx, or)
	if err != nil {
		return nil, err
	}
	res, err := p.GetORStatus(or)
	if err != nil {
		return nil, err
	}
	if res.Status == "" {
		res.Status = or.Status
	}
	if res.DeliveryStatus == "" {
		res.DeliveryStatus = or.DeliveryStatus
	}
	return &responsemodel.OutboundRequestStatus{
		ORCode:         res.ORCode,
		Status:         res.Status,
		DeliveryStatus: res.DeliveryStatus,
		UpdatedAt:      ptime.Now(),
		Link:           res.Data.Link,
	}, nil
}

// SyncAllGlobalCare ...
func (s outboundRequestImplement) SyncAllGlobalCare() {
	var (
		ctx  = context.Background()
		wg   = sync.WaitGroup{}
		cond = bson.D{
			{"partner.identityCode", constant.WarehousePartnerCodeGlobalCare},
			{"createdAt", bson.M{
				"$gte": ptime.Now().AddDate(0, 0, -7), // 7 days ago
			}},
			{"$or", []bson.M{
				{"status": constant.WarehouseORStatusWaitingToConfirm},
				{"deliveryStatus": constant.DeliveryPartnerStatusWaitingToPick},
			}},
		}
	)
	cursor, err := dao.OutboundRequest().FindCursor(ctx, cond)
	if err != nil {
		logger.Error("service.OutboundRequest.SyncAllGlobalCare - Find", logger.LogData{
			"cond": cond,
			"err":  err.Error(),
		})
		return
	}
	defer cursor.Close(ctx)

	notUpdateStatusORs := make(chan string)
	done := make(chan bool)
	var notUpdateStatusORCodes []string
	go func() {
		for orCode := range notUpdateStatusORs {
			notUpdateStatusORCodes = append(notUpdateStatusORCodes, orCode)
		}
		done <- true
	}()

	p, err := ants.NewPoolWithFunc(10, s.handleSyncORStatus(&wg, notUpdateStatusORs))
	if err != nil {
		return
	}
	defer p.Release()
	for cursor.Next(ctx) {
		var doc mgwarehouse.OutboundRequest
		if err = cursor.Decode(&doc); err != nil {
			logger.Error("service.OutboundRequest.SyncAllGlobalCare - Decode", logger.LogData{"err": err.Error()})
			continue
		}
		wg.Add(1)
		_ = p.Invoke(doc)
	}
	wg.Wait()
	close(notUpdateStatusORs)

	<-done
	if len(notUpdateStatusORCodes) == 0 {
		return
	}
	// notify OR Global Care not update status
	if err = client.GetOrder().ORNotUpdateStatus(model.OrderORsNotUpdateStatus{ORCodes: notUpdateStatusORCodes}); err != nil {
		logger.Error("service.OutboundRequest.SyncAllGlobalCare - ORNotUpdateStatus", logger.LogData{
			"err":   err.Error(),
			"codes": notUpdateStatusORCodes,
		})
	}
}

func (s outboundRequestImplement) handleSyncORStatus(wg *sync.WaitGroup, ch chan string) func(i interface{}) {
	return func(i interface{}) {
		ctx := context.Background()
		defer wg.Done()
		or, ok := i.(mgwarehouse.OutboundRequest)
		if !ok {
			return
		}
		res, err := s.SyncStatusByOR(ctx, or)
		if err != nil {
			return
		}
		if res.DeliveryStatus == or.DeliveryStatus && res.Status == or.Status {
			// not update status from GC
			now := ptime.Now()
			maxTime := config.GetENV().MaxTimeWaitForGCUpdateStatusInMinute
			isExceedUpdateTime := or.CreatedAt.Add(time.Duration(maxTime) * time.Minute).Before(now)
			if or.Status == constant.WarehouseORStatusWaitingToConfirm && isExceedUpdateTime {
				ch <- or.Partner.Code
			}
			return
		}
		// update OR status
		or.Status = res.Status
		or.DeliveryStatus = res.DeliveryStatus
		or.UpdatedAt = res.UpdatedAt
		or.Link = res.Link

		update := bson.M{
			"status":         res.Status,
			"deliveryStatus": res.DeliveryStatus,
			"updatedAt":      res.UpdatedAt,
			"link":           res.Link,
		}
		_, _ = dao.OutboundRequest().UpdateByID(ctx, or.ID, bson.M{"$set": update})
		// save history
		go func() { _, _ = OutboundRequestHistory().SaveByOR(context.Background(), or) }()
		// update status to order
		payload := model.OrderUpdateORStatus{
			ID:             or.ID.Hex(),
			OrderCode:      or.OrderCode,
			ORCode:         or.Partner.Code,
			Status:         res.Status,
			DeliveryStatus: res.DeliveryStatus,
			Data:           model.OrderORData{Link: res.Link},
		}
		if err = client.GetOrder().UpdateORStatus(payload); err != nil {
			logger.Error("service.outboundRequestImplement.handleSyncORStatus - client.GetOrder().UpdateORStatus", logger.LogData{
				"err":     err.Error(),
				"payload": payload,
			})
		}
	}
}

func (s outboundRequestImplement) getORInfoByCode(ctx context.Context, orCode, orderID string) (mgwarehouse.OutboundRequest, OutboundRequestPartner, error) {
	d := dao.OutboundRequest()
	id, valid := mongodb.NewIDFromString(orderID)
	if !valid {
		return mgwarehouse.OutboundRequest{}, nil, errors.New(errorcode.ORInvalidOrderID)
	}
	doc := d.FindOneByORCodeAndOrderID(ctx, orCode, id)
	if doc.ID.IsZero() {
		return mgwarehouse.OutboundRequest{}, nil, errors.New(errorcode.ORNotFound)
	}
	orSvc, err := s.getWarehousePartnerByDoc(ctx, doc)
	return doc, orSvc, err
}

func (s outboundRequestImplement) getWarehousePartnerByDoc(ctx context.Context, doc mgwarehouse.OutboundRequest) (OutboundRequestPartner, error) {
	whCfg := dao.WarehouseConfiguration().FindByWarehouseID(ctx, doc.Warehouse)
	return NewOutboundRequestService(doc.Partner.IdentityCode, whCfg.Partner.Authentication)
}

func (s outboundRequestImplement) getResponse(doc mgwarehouse.OutboundRequest) responsemodel.OutboundRequest {
	p := doc.Partner
	return responsemodel.OutboundRequest{
		ID:     doc.ID.Hex(),
		Status: doc.Status,
		Partner: responsemodel.OutboundRequestPartner{
			IdentityCode: p.IdentityCode,
			Code:         p.Code,
			RequestID:    p.RequestID,
		},
		TrackingCode: doc.TrackingCode,
		CreatedAt:    ptime.TimeResponseInit(doc.CreatedAt),
		UpdatedAt:    ptime.TimeResponseInit(doc.UpdatedAt),
	}
}

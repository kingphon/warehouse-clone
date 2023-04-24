package handler

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/nats-io/nats.go"

	"git.selly.red/Selly-Modules/logger"
	"git.selly.red/Selly-Modules/mongodb"
	"git.selly.red/Selly-Modules/natsio"
	"git.selly.red/Selly-Modules/natsio/model"

	externalauth "git.selly.red/Selly-Server/warehouse/external/auth"
	"git.selly.red/Selly-Server/warehouse/external/constant"
	"git.selly.red/Selly-Server/warehouse/external/response"
	"git.selly.red/Selly-Server/warehouse/external/utils/mgquery"
	"git.selly.red/Selly-Server/warehouse/pkg/admin/service"
)

// Nats ...
type Nats struct {
	EncodedConn *natsio.JSONEncoder
}

// CountWarehouseWithCondition ...
func (h *Nats) CountWarehouseWithCondition(msg *nats.Msg) {
	var (
		payload *model.FindWithCondition
	)

	if err := bson.Unmarshal(msg.Data, &payload); err != nil {
		logger.Error("handler.Nats - CountWarehouseWithCondition", logger.LogData{
			"err":  err.Error(),
			"data": string(msg.Data),
		})
		h.response(msg.Reply, nil, err)
		return
	}

	data, err := service.Nats{}.CountWarehouseWithCondition(context.Background(), payload)
	if err != nil {
		h.response(msg.Reply, nil, err)
		return
	}
	h.response(msg.Reply, data, nil)
}

// DistinctWarehouseWithField ...
func (h *Nats) DistinctWarehouseWithField(msg *nats.Msg) {
	var (
		payload *model.DistinctWithField
	)

	if err := bson.Unmarshal(msg.Data, &payload); err != nil {
		logger.Error("handler.Nats - DistinctWarehouseWithField", logger.LogData{
			"err":  err.Error(),
			"data": string(msg.Data),
		})
		h.response(msg.Reply, nil, err)
		return
	}

	data, err := service.Nats{}.DistinctWarehouseWithField(context.Background(), payload)
	if err != nil {
		h.response(msg.Reply, nil, err)
		return
	}
	h.response(msg.Reply, data, nil)
}

// GetWarehouseWithCondition ...
func (h *Nats) GetWarehouseWithCondition(msg *nats.Msg) {
	id := mongodb.NewObjectID().Hex()
	var (
		payload *model.FindWithCondition
		ctx     = context.Background()
	)

	if err := bson.Unmarshal(msg.Data, &payload); err != nil {
		logger.Error("Nats.GetWarehouseWithCondition - GetWarehouseWithCondition", logger.LogData{
			"err":  err.Error(),
			"data": string(msg.Data),
		})
		h.response(msg.Reply, nil, err)
		return
	}

	data, err := service.Nats{}.GetWarehouseWithCondition(ctx, payload, id)
	if err != nil {
		h.response(msg.Reply, nil, err)
		return
	}
	h.response(msg.Reply, data, nil)
}

// GetOneWarehouse ...
func (h *Nats) GetOneWarehouse(msg *nats.Msg) {
	var (
		payload *model.FindOneCondition
	)

	if err := bson.Unmarshal(msg.Data, &payload); err != nil {
		logger.Error("handler.Nats - GetOneWarehouse", logger.LogData{
			"err":  err.Error(),
			"data": string(msg.Data),
		})
		h.response(msg.Reply, nil, err)
		return
	}

	data, err := service.Nats{}.GetOneWarehouse(context.Background(), payload)
	if err != nil {
		h.response(msg.Reply, nil, err)
		return
	}
	h.response(msg.Reply, data, nil)
}

// CreateOutboundRequest ...
func (h *Nats) CreateOutboundRequest(subject, reply string, req *model.OutboundRequestPayload) {
	s := service.Nats{}
	res, err := s.CreateOutboundRequest(req)
	if err != nil {
		logger.Error("handler.Nats.CreateOutboundRequest", logger.LogData{
			"err": err.Error(),
			"req": req,
		})
	}
	h.response(reply, res, err)
}

// UpdateOutboundRequestLogistic ...
func (h *Nats) UpdateOutboundRequestLogistic(subject, reply string, req *model.UpdateOutboundRequestLogisticInfoPayload) {
	s := service.Nats{}
	err := s.UpdateLogisticInfo(req)
	if err != nil {
		logger.Error("handler.Nats.UpdateOutboundRequestLogistic", logger.LogData{
			"err":     err.Error(),
			"payload": req,
		})
	}
	h.response(reply, nil, err)
}

// UpdateORDeliveryStatus ...
func (h *Nats) UpdateORDeliveryStatus(subject, reply string, req *model.WarehouseORUpdateDeliveryStatus) {
	logger.Debug("Nats.UpdateORDeliveryStatus", logger.LogData{"req": req})
	s := service.OutboundRequest()
	err := s.UpdateORDeliveryStatus(req)
	if err != nil {
		logger.Error("handler.Nats.UpdateORDeliveryStatus", logger.LogData{
			"err":     err.Error(),
			"payload": req,
		})
	}
	h.response(reply, nil, err)
}

// CancelOutboundRequest ...
func (h *Nats) CancelOutboundRequest(subject, reply string, req *model.CancelOutboundRequest) {
	s := service.Nats{}
	err := s.CancelOutboundRequest(req)
	if err != nil {
		logger.Error("handler.Nats.CancelOutboundRequest", logger.LogData{
			"err":     err.Error(),
			"payload": req,
		})
	}
	h.response(reply, nil, err)
}

// GetWarehouseConfiguration ...
func (h *Nats) GetWarehouseConfiguration(subject, reply string, req *string) {
	s := service.WarehouseConfiguration(externalauth.User{})
	id, valid := mongodb.NewIDFromString(*req)
	if !valid {
		h.response(reply, nil, errors.New(response.CommonInvalidID))
		return
	}
	data, err := s.DetailByWarehouseID(context.Background(), id)
	if err != nil {
		logger.Error("handler.Nats.GetWarehouseConfiguration", logger.LogData{
			"err":     err.Error(),
			"payload": req,
		})
		h.response(reply, nil, err)
		return
	}
	res := model.WarehouseConfiguration{
		Warehouse:               data.Warehouse,
		DoesSupportSellyExpress: data.Other.DoesSupportSellyExpress,
		Supplier: model.WarehouseSupplier{
			CanAutoSendMail:       false,
			InvoiceDeliveryMethod: data.Supplier.InvoiceDeliveryMethod,
		},
		Order: model.WarehouseOrder{
			MinimumValue: data.Order.MinimumValue,
			PaymentMethod: model.WarehousePaymentMethod{
				Cod:          data.Order.PaymentMethod.Cod,
				BankTransfer: data.Order.PaymentMethod.BankTransfer,
			},
			IsLimitNumberOfPurchases: data.Order.IsLimitNumberOfPurchases,
			LimitNumberOfPurchases:   data.Order.LimitNumberOfPurchases,
			NotifyOnNewOrder: model.WarehouseConfigNotifyOnNewOrder{
				Enable:  data.Order.NotifyOnNewOrder.Enable,
				Channel: data.Order.NotifyOnNewOrder.Channel,
				RoomID:  data.Order.NotifyOnNewOrder.RoomID,
			},
		},
		Partner: model.WarehousePartner{
			IdentityCode:   data.Partner.IdentityCode,
			Code:           data.Partner.Code,
			Enabled:        data.Partner.Enabled,
			Authentication: data.Partner.Authentication,
		},
		Delivery: model.WarehouseDelivery{
			DeliveryMethods:      data.Delivery.DeliveryMethods,
			PriorityServiceCodes: data.Delivery.PriorityServiceCodes,
			EnabledSources:       data.Delivery.EnabledSources,
			Types:                data.Delivery.Types,
		},
		Food: model.WarehouseFood{
			ForceClosed: data.Food.ForceClosed,
			IsClosed:    data.Food.IsClosed,
			TimeRange:   make([]model.TimeRange, 0),
		},
		AutoConfirmOrder: model.WarehouseOrderConfirm{
			IsEnable:              data.AutoConfirmOrder.IsEnable,
			ConfirmDelayInSeconds: data.AutoConfirmOrder.ConfirmDelayInSeconds,
		},
	}

	if len(data.Food.TimeRange) > 0 {
		for _, t := range data.Food.TimeRange {
			res.Food.TimeRange = append(res.Food.TimeRange, model.TimeRange{
				From: t.From,
				To:   t.To,
			})
		}
	}
	h.response(reply, res, nil)
}

// SyncORStatus ...
func (h *Nats) SyncORStatus(subject, reply string, req *model.SyncORStatusRequest) {
	s := service.OutboundRequest()
	res, err := s.SyncStatus(req)
	h.response(reply, res, err)
}

// NewTNCWebhook ...
func (h *Nats) NewTNCWebhook(subject, reply string, req *string) {
	s := service.OutboundRequest()
	err := s.UpdateStatusFromWebhook([]byte(*req), constant.WarehousePartnerCodeTNC)
	h.response(reply, nil, err)
}

// NewGlobalCareWebhook ...
func (h *Nats) NewGlobalCareWebhook(subject, reply string, req *string) {
	s := service.OutboundRequest()
	err := s.UpdateStatusFromWebhook([]byte(*req), constant.WarehousePartnerCodeGlobalCare)
	h.response(reply, nil, err)
}

// GetWarehouses ...
func (h *Nats) GetWarehouses(subject, reply string, req *model.GetWarehousesRequest) {
	var (
		ctx = context.Background()
		q   = mgquery.AppQuery{
			Page:         req.Page,
			Limit:        req.Limit,
			Keyword:      req.Keyword,
			BusinessType: req.BusinessType,
			Warehouse: mgquery.Warehouse{
				Status:   req.Status,
				Supplier: req.Supplier,
			},
		}
	)
	s := service.Warehouse(externalauth.User{})
	data := s.All(ctx, q)
	res := model.GetWarehousesResponse{
		Total: data.Total,
		Limit: data.Limit,
		List:  make([]model.WarehouseInfo, len(data.List)),
	}
	for i, item := range data.List {
		res.List[i] = model.WarehouseInfo{
			ID:           item.ID,
			Name:         item.Name,
			BusinessType: item.BusinessType,
			Status:       item.Status,
			Slug:         item.Slug,
			Supplier: model.WarehouseSupplierInfo{
				ID:   item.Supplier.ID,
				Name: item.Supplier.Name,
			},
			Location: model.ResponseWarehouseLocation{
				Province: model.CommonLocation{
					ID:   item.Location.Province.ID,
					Name: item.Location.Province.Name,
					Code: item.Location.Province.Code,
				},
				District: model.CommonLocation{
					ID:   item.Location.District.ID,
					Name: item.Location.District.Name,
					Code: item.Location.District.Code,
				},
				Ward: model.CommonLocation{
					ID:   item.Location.Ward.ID,
					Name: item.Location.Ward.Name,
					Code: item.Location.Ward.Code,
				},
				Address: item.Location.Address,
				LocationCoordinates: model.ResponseLatLng{
					Latitude:  item.Location.LocationCoordinates.Latitude,
					Longitude: item.Location.LocationCoordinates.Longitude,
				},
			},
			Contact: model.ResponseWarehouseContact{
				Name:    item.Contact.Name,
				Phone:   item.Contact.Name,
				Address: item.Contact.Address,
				Email:   item.Contact.Email,
			},
			CreatedAt: item.CreatedAt.FormatISODate(),
			UpdatedAt: item.UpdatedAt.FormatISODate(),
		}
	}
	h.response(reply, res, nil)
}

// NewOnPointWebhook ...
func (h *Nats) NewOnPointWebhook(subject, reply string, req *string) {
	s := service.OutboundRequest()
	err := s.UpdateStatusFromWebhook([]byte(*req), constant.WarehousePartnerCodeOnPoint)
	h.response(reply, nil, err)
}

func (h *Nats) response(reply string, data interface{}, err error) {
	res := map[string]interface{}{"data": data}
	if err != nil {
		res["error"] = err.Error()
	}
	h.EncodedConn.Publish(reply, res)
}

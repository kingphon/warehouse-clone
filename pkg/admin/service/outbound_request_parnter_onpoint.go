package service

import (
	"fmt"

	"git.selly.red/Selly-Modules/3pl/partnerapi/onpoint"
	"git.selly.red/Selly-Modules/3pl/util/pjson"
	"git.selly.red/Selly-Modules/logger"
	"git.selly.red/Selly-Modules/natsio"
	"git.selly.red/Selly-Modules/natsio/model"

	"git.selly.red/Selly-Server/warehouse/external/constant"
	mgwarehouse "git.selly.red/Selly-Server/warehouse/external/model/mg/warehouse"
	"git.selly.red/Selly-Server/warehouse/external/utils/base64"
	"git.selly.red/Selly-Server/warehouse/external/utils/ptime"
	"git.selly.red/Selly-Server/warehouse/pkg/admin/config"
	"git.selly.red/Selly-Server/warehouse/pkg/admin/errorcode"
	responsemodel "git.selly.red/Selly-Server/warehouse/pkg/admin/model/response"
)

var (
	onpointStatusMapping = map[string]string{
		onpoint.OrderStatusNew:              constant.WarehouseORStatusConfirmed,
		onpoint.OrderStatusPendingWarehouse: constant.WarehouseORStatusConfirmed,
		onpoint.OrderStatusWhProcessing:     constant.WarehouseORStatusConfirmed,
		onpoint.OrderStatusWhCompleted:      constant.WarehouseORStatusPacked,
		onpoint.OrderStatusDlPending:        constant.WarehouseORStatusPicked,
		onpoint.OrderStatusDlIntransit:      constant.WarehouseORStatusPicked,
		onpoint.OrderStatusDLDelivered:      constant.WarehouseORStatusPicked,
		onpoint.OrderStatusDLReturning:      constant.WarehouseORStatusPicked,
		onpoint.OrderStatusReturned:         constant.WarehouseORStatusReturned,
		onpoint.OrderStatusPartialCancelled: constant.WarehouseORStatusCancelled,
		onpoint.OrderStatusCancelled:        constant.WarehouseORStatusCancelled,
		onpoint.OrderStatusCompleted:        constant.WarehouseORStatusPicked,
		onpoint.OrderStatusUnknown:          "",
	}
)

type outboundRequestOnPoint struct {
	auth string
}

func newOROnPointService(auth string) (OutboundRequestPartner, error) {
	return &outboundRequestOnPoint{auth: auth}, nil
}

// HasUpdateORCode ...
func (s *outboundRequestOnPoint) HasUpdateORCode() bool {
	return true
}

// UpdateDeliveryStatus ...
func (s *outboundRequestOnPoint) UpdateDeliveryStatus(req model.WarehouseORUpdateDeliveryStatus, or mgwarehouse.OutboundRequest) error {
	c, err := s.getClient()
	if err != nil {
		return err
	}
	_, err = c.UpdateDelivery(onpoint.UpdateOrderDeliveryRequest{
		OrderCode:      or.OrderCode,
		DeliveryStatus: req.DeliveryStatus,
	})
	return err
}

// Create ...
func (s *outboundRequestOnPoint) Create(p model.OutboundRequestPayload,
	whPartner responsemodel.ResponseWarehousePartnerConfig) (*model.OutboundRequestResponse, error) {
	c, err := s.getClient()
	if err != nil {
		return nil, err
	}
	r := &model.OutboundRequestResponse{
		OrderCode:    p.OrderCode,
		TrackingCode: p.TrackingCode,
		Status:       constant.WarehouseORStatusWaitingToConfirm,
	}
	payload := s.getCreateOrderPayload(p, whPartner)
	res, err := c.CreateOrder(payload)
	if err != nil {
		r.Status = constant.WarehouseORStatusFailed
		r.Reason = err.Error()
	} else {
		r.ORCode = res.OnpointOrderCode
	}
	return r, err
}

func (s *outboundRequestOnPoint) getCreateOrderPayload(p model.OutboundRequestPayload,
	partner responsemodel.ResponseWarehousePartnerConfig) onpoint.CreateOrderRequest {
	items := make([]onpoint.OrderItem, len(p.Items))
	for i, item := range p.Items {
		items[i] = onpoint.OrderItem{
			SellingPrice: int(item.Price),
			Quantity:     int(item.Quantity),
			Uom:          item.UnitCode,
			Amount:       int(item.Price) * int(item.Quantity),
			Name:         item.Name,
			PartnerSku:   item.SupplierSKU,
		}
	}
	m := map[string]string{
		constant.PaymentMethod.COD:          onpoint.PaymentMethodCOD,
		constant.PaymentMethod.BankTransfer: onpoint.PaymentMethodBankTransfer,
	}
	payload := onpoint.CreateOrderRequest{
		OrderCode:          p.OrderCode,
		OrderDate:          ptime.Now(),
		PickupLocationCode: partner.Code,
		Note:               p.Note,
		SubtotalPrice:      int(p.CODAmount),
		TotalDiscounts:     0,
		TotalPrice:         int(p.CODAmount),
		PaymentMethod:      m[p.PaymentMethod],
		Items:              items,
	}
	return payload
}

// UpdateLogisticInfo ...
func (s *outboundRequestOnPoint) UpdateLogisticInfo(p model.UpdateOutboundRequestLogisticInfoPayload,
	or mgwarehouse.OutboundRequest) error {
	c, err := s.getClient()
	if err != nil {
		return err
	}
	payload := onpoint.UpdateOrderDeliveryRequest{
		OrderCode:              or.OrderCode,
		DeliveryPlatform:       p.TPLCode,
		DeliveryTrackingNumber: p.TrackingCode,
		ShippingLabel:          p.ShippingLabel,
	}
	_, err = c.UpdateDelivery(payload)
	return err
}

// Cancel ...
func (s *outboundRequestOnPoint) Cancel(or mgwarehouse.OutboundRequest, note string) error {
	c, err := s.getClient()
	if err != nil {
		return err
	}
	_, err = c.CancelOrder(onpoint.CancelOrderRequest{OrderNo: or.Partner.Code})
	return err
}

// GetORStatus ...
func (s *outboundRequestOnPoint) GetORStatus(or mgwarehouse.OutboundRequest) (*model.SyncORStatusResponse, error) {
	return nil, nil
}

// GetWebhookData ...
func (s *outboundRequestOnPoint) GetWebhookData(b []byte) (*responsemodel.OutboundRequestWebhookData, error) {
	var p onpoint.WebhookPayload
	if err := pjson.Unmarshal(b, &p); err != nil {
		logger.Error("service.outboundRequestOnPoint.GetWebhookData - parse webhook", logger.LogData{
			"data": string(b),
			"err":  err.Error(),
		})
		return nil, err
	}
	data, ok := p.GetDataEventUpdateOrderStatus()
	// just ignore if webhook is not update order status, process for product inventory if we need later
	if !ok {
		return nil, nil
	}
	result := &responsemodel.OutboundRequestWebhookData{
		Status:         onpointStatusMapping[data.Status],
		DeliveryStatus: "",
		ORCode:         data.OnpointOrderCode,
		OrderCode:      data.OrderCode,
		UpdatedAt:      ptime.Now(),
	}
	t, err := ptime.ParseWithUTC(data.UpdatedAt, onpoint.TimeLayout)
	if err == nil {
		result.UpdatedAt = t
	}
	if result.Status == constant.WarehouseORStatusCancelled {
		result.Reason = errorcode.ORMessageCancelOR
	}

	return result, nil
}

// GetIdentityCode ...
func (s *outboundRequestOnPoint) GetIdentityCode() string {
	return constant.WarehousePartnerCodeOnPoint
}

func (s *outboundRequestOnPoint) getClient() (*onpoint.Client, error) {
	var info struct {
		APIKey    string `json:"apiKey"`
		SecretKey string `json:"secretKey"`
	}
	b := base64.Decode(s.auth)
	if err := pjson.Unmarshal(b, &info); err != nil {
		return nil, fmt.Errorf("getOnPointClient: parse_data: %v, %s", err, string(b))
	}
	env := onpoint.EnvStaging
	if config.IsEnvProduction() {
		env = onpoint.EnvProd
	}
	natsClient := natsio.GetServer()
	return onpoint.NewClient(env, info.APIKey, info.SecretKey, natsClient)
}

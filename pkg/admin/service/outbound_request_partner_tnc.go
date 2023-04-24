package service

import (
	"errors"
	"fmt"

	"git.selly.red/Selly-Modules/3pl/partnerapi/tnc"
	"git.selly.red/Selly-Modules/3pl/util/pjson"
	"git.selly.red/Selly-Modules/logger"
	"git.selly.red/Selly-Modules/natsio"
	"git.selly.red/Selly-Modules/natsio/model"

	"git.selly.red/Selly-Server/warehouse/external/constant"
	mgwarehouse "git.selly.red/Selly-Server/warehouse/external/model/mg/warehouse"
	"git.selly.red/Selly-Server/warehouse/external/utils/base64"
	"git.selly.red/Selly-Server/warehouse/external/utils/format"
	"git.selly.red/Selly-Server/warehouse/external/utils/pstring"
	"git.selly.red/Selly-Server/warehouse/external/utils/ptime"
	"git.selly.red/Selly-Server/warehouse/pkg/admin/config"
	responsemodel "git.selly.red/Selly-Server/warehouse/pkg/admin/model/response"
)

var (
	tncEventStatusMapping = map[string]string{
		"OR_CREATING_COMPLETED": constant.WarehouseORStatusConfirmed, // handle by error code
		"OR_PROCESSING_STARTED": constant.WarehouseORStatusConfirmed,
		"OR_CANCELLED":          constant.WarehouseORStatusCancelled,
		"OR_PACKED":             constant.WarehouseORStatusPacked,
		"OR_HANDED_ON":          constant.WarehouseORStatusPicked,
		"OR_RETURNED":           constant.WarehouseORStatusReturned,
	}
)

type outboundRequestTNC struct {
	auth string
	// client *tnc.Client
}

func (s outboundRequestTNC) UpdateDeliveryStatus(req model.WarehouseORUpdateDeliveryStatus, or mgwarehouse.OutboundRequest) error {
	return nil
}

func (s outboundRequestTNC) HasUpdateORCode() bool {
	return false
}

// GetORStatus ...
func (s outboundRequestTNC) GetORStatus(or mgwarehouse.OutboundRequest) (*model.SyncORStatusResponse, error) {
	// TODO implement me
	return nil, errors.New("tnc.GetORStatus - not implemented")
}

func newORTNCService(auth string) (OutboundRequestPartner, error) {
	s := outboundRequestTNC{auth: auth}
	// client, err := s.getTNCClient()
	// if err != nil {
	// 	return nil, err
	// }
	// s.client = client
	return s, nil
}

// GetIdentityCode ...
func (s outboundRequestTNC) GetIdentityCode() string {
	return constant.WarehousePartnerCodeTNC
}

// GetWebhookData ...
func (s outboundRequestTNC) GetWebhookData(b []byte) (*responsemodel.OutboundRequestWebhookData, error) {
	var data tnc.Webhook
	if err := pjson.Unmarshal(b, &data); err != nil {
		logger.Error("service.outboundRequestTNC.GetWebhookData - parse webhook", logger.LogData{
			"data": string(b),
			"err":  err.Error(),
		})
		return nil, err
	}
	result := &responsemodel.OutboundRequestWebhookData{
		Status:      tncEventStatusMapping[data.Event],
		ORCode:      data.OrCode,
		ORRequestID: format.ToString(data.OrId),
		UpdatedAt:   ptime.FromUnix(data.Timestamp),
		OrderCode:   data.PartnerORCode,
	}
	if s.hasError(data) {
		result.Status = constant.WarehouseORStatusFailed
		result.Reason = fmt.Sprintf("tnc: code %s - message %s", data.ErrorCode, data.ErrorMessage)
	}

	if result.Status == constant.WarehouseORStatusCancelled {
		result.Reason = data.Note
	}

	return result, nil
}

func (s outboundRequestTNC) hasError(t tnc.Webhook) bool {
	return t.ErrorCode != "" || (t.Event == "OR_CREATING_COMPLETED" && t.ErrorMessage != "")
}

// Cancel ...
func (s outboundRequestTNC) Cancel(or mgwarehouse.OutboundRequest, note string) error {
	client, err := s.getTNCClient()
	if err != nil {
		return err
	}
	requestID, err := format.StringToInt(or.Partner.RequestID)
	if err != nil {
		return fmt.Errorf("outboundRequestTNC.Cancel: invalid_request_id %s", or.Partner.RequestID)
	}
	return client.CancelOutboundRequest(requestID, note)
}

// UpdateLogisticInfo ...
func (s outboundRequestTNC) UpdateLogisticInfo(p model.UpdateOutboundRequestLogisticInfoPayload, or mgwarehouse.OutboundRequest) error {
	client, err := s.getTNCClient()
	if err != nil {
		return err
	}
	requestID, err := format.StringToInt(or.Partner.RequestID)
	if err != nil {
		return fmt.Errorf("outboundRequestTNC.UpdateLogisticInfo: invalid_request_id %s", or.Partner.RequestID)
	}
	payload := tnc.UpdateORLogisticInfoPayload{
		OrID:          requestID,
		TrackingCode:  s.getTrackingCode(p.TrackingCode, or.TPLCode),
		ShippingLabel: p.ShippingLabel,
		SlaShipDate:   "",
	}
	return client.UpdateOutboundRequestLogisticInfo(payload)
}

// Create ...
func (s outboundRequestTNC) Create(p model.OutboundRequestPayload, whp responsemodel.ResponseWarehousePartnerConfig) (
	*model.OutboundRequestResponse, error) {
	client, err := s.getTNCClient()
	if err != nil {
		return nil, err
	}

	r := &model.OutboundRequestResponse{
		OrderCode:    p.OrderCode,
		TrackingCode: p.TrackingCode,
		Status:       constant.WarehouseORStatusWaitingToConfirm,
	}

	customer := p.Customer
	payload := s.getCreateOrderPayload(p, whp, customer)
	res, err := client.CreateOutboundRequest(payload)
	if err != nil {
		r.Status = constant.WarehouseORStatusFailed
		r.Reason = err.Error()
	} else {
		r.ORCode = res.OrCode
		r.RequestID = pjson.ToJSONString(res.OrID)
	}
	return r, err
}

func (s outboundRequestTNC) getCreateOrderPayload(p model.OutboundRequestPayload, whp responsemodel.ResponseWarehousePartnerConfig,
	customer model.CustomerInfo) tnc.OutboundRequestPayload {
	return tnc.OutboundRequestPayload{
		WarehouseCode:       whp.Code,
		ShippingServiceCode: tnc.ShippingServiceCodeSTD,
		PartnerORCode:       p.OrderCode,
		PartnerRefId:        "",
		RefCode:             "",
		CodAmount:           p.CODAmount,
		PriorityType:        tnc.PriorityNormal,
		CustomerName:        customer.Name,
		CustomerPhoneNumber: customer.PhoneNumber,
		Type:                tnc.ORTypeOrder,
		ShippingType:        tnc.ShippingTypeSelfShip,
		VehicleNumber:       "",
		ContainerNumber:     "",
		PackType:            tnc.PackTypeNormal,
		PackingNote:         "",
		CustomLabel:         false,
		BizType:             tnc.BizTypeB2C,
		Note:                p.Note,
		ShippingAddress: tnc.Address{
			AddressNo: customer.Address.Address,
			// pass value when use TNC delivery service
			ProvinceCode: "",
			DistrictCode: "",
			WardCode:     "",
		},
		Products:            s.getProducts(p),
		PartnerCreationTime: ptime.Now().Format(tnc.TimeLayout),
		TPLCode:             p.TPLCode,
		TrackingCode:        s.getTrackingCode(p.TrackingCode, p.TPLCode),
	}
}

func (s outboundRequestTNC) getProducts(p model.OutboundRequestPayload) []tnc.Product {
	products := make([]tnc.Product, 0)
	for _, item := range p.Items {
		// check exist sku first
		sku := pstring.TrimSpace(item.SupplierSKU)
		index := -1
		for i, product := range products {
			if product.PartnerSKU == sku {
				index = i
				break
			}
		}
		if sku == "" {
			continue
		}
		// sum quantity if exists
		if index > -1 {
			products[index].Quantity += item.Quantity
		} else {
			products = append(products, tnc.Product{
				PartnerSKU:        sku,
				UnitCode:          item.UnitCode,
				ConditionTypeCode: tnc.ConditionTypeCodeNew,
				Quantity:          item.Quantity,
			})
		}
	}
	return products
}
func (s outboundRequestTNC) getTNCClient() (*tnc.Client, error) {
	var info struct {
		ClientID     string `json:"clientId"`
		ClientSecret string `json:"clientSecret"`
		Realm        string `json:"realm"`
	}
	b := base64.Decode(s.auth)
	if err := pjson.Unmarshal(b, &info); err != nil {
		return nil, fmt.Errorf("getTNCClient: parse_data: %v, %s", err, string(b))
	}
	env := tnc.EnvStaging
	if config.IsEnvProduction() {
		env = tnc.EnvProd
	}
	natsClient := natsio.GetServer()
	return tnc.NewClient(env, info.ClientID, info.ClientSecret, info.Realm, natsClient)
}

func (s outboundRequestTNC) getTrackingCode(trackingCode string, tplCode string) string {
	if tplCode == tnc.TPLCodeGHTK {
		if v := pstring.Split(trackingCode, "."); len(v) > 1 {
			trackingCode = v[len(v)-1]
		}
	}
	return trackingCode
}

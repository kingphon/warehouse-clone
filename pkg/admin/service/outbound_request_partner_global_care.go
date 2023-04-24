package service

import (
	"errors"
	"fmt"

	"git.selly.red/Selly-Modules/3pl/partnerapi/globalcare"
	"git.selly.red/Selly-Modules/3pl/util/pjson"
	"git.selly.red/Selly-Modules/logger"
	"git.selly.red/Selly-Modules/natsio"
	"git.selly.red/Selly-Modules/natsio/model"

	"git.selly.red/Selly-Server/warehouse/external/constant"
	mgwarehouse "git.selly.red/Selly-Server/warehouse/external/model/mg/warehouse"
	"git.selly.red/Selly-Server/warehouse/external/utils/base64"
	"git.selly.red/Selly-Server/warehouse/external/utils/format"
	"git.selly.red/Selly-Server/warehouse/pkg/admin/config"
	"git.selly.red/Selly-Server/warehouse/pkg/admin/errorcode"
	responsemodel "git.selly.red/Selly-Server/warehouse/pkg/admin/model/response"
)

var (
	globalCareStatusMapping = map[int]string{
		1: constant.WarehouseORStatusWaitingToConfirm,
		2: constant.WarehouseORStatusConfirmed,
		3: constant.WarehouseORStatusConfirmed,
	}
	globalCareDeliveryStatusMapping = map[int]string{
		1: constant.DeliveryPartnerStatusWaitingToPick,
		2: constant.DeliveryPartnerStatusWaitingToPick,
		3: constant.DeliveryPartnerStatusDelivered,
	}
)

type outboundRequestGlobalCare struct {
	auth   string
	client *globalcare.Client
}

func (s outboundRequestGlobalCare) UpdateDeliveryStatus(req model.WarehouseORUpdateDeliveryStatus, or mgwarehouse.OutboundRequest) error {
	return nil
}

func (s outboundRequestGlobalCare) HasUpdateORCode() bool {
	return false
}

func newORGlobalCareService(auth string) (OutboundRequestPartner, error) {
	s := outboundRequestGlobalCare{auth: auth}
	if auth != "" {
		client, err := s.getClient()
		if err != nil {
			return nil, err
		}
		s.client = client
	}
	return s, nil
}

// GetORStatus ...
func (s outboundRequestGlobalCare) GetORStatus(or mgwarehouse.OutboundRequest) (*model.SyncORStatusResponse, error) {
	res, err := s.client.GetOrder(or.Partner.Code)
	if err != nil {
		return nil, err
	}
	r := &model.SyncORStatusResponse{
		ORCode:         or.Partner.Code,
		OrderCode:      or.OrderCode,
		Status:         globalCareStatusMapping[res.Result.StatusId],
		DeliveryStatus: globalCareDeliveryStatusMapping[res.Result.StatusId],
		Data:           model.OrderORData{Link: res.Result.CertLink},
	}
	if r.Status == "" {
		r.Status = or.Status
	}
	if r.DeliveryStatus == "" {
		r.DeliveryStatus = or.DeliveryStatus
	}
	return r, nil
}

// Create ...
func (s outboundRequestGlobalCare) Create(p model.OutboundRequestPayload, whPartner responsemodel.ResponseWarehousePartnerConfig) (
	*model.OutboundRequestResponse, error) {
	if p.Insurance == nil {
		return nil, errors.New(errorcode.ORMissingInsuranceInfo)
	}
	r := &model.OutboundRequestResponse{
		OrderCode:    p.OrderCode,
		TrackingCode: p.TrackingCode,
		Status:       constant.WarehouseORStatusWaitingToConfirm,
	}
	insuranceInfo := p.Insurance
	customer := p.Customer
	vTypeID, err1 := format.StringToInt(insuranceInfo.VehicleTypeID)
	iTypeID, err2 := format.StringToInt(insuranceInfo.InsuranceTypeID)
	if err1 != nil || err2 != nil {
		return nil, errors.New(errorcode.ORMissingInsuranceInfo)
	}

	cfg := globalcare.MotorbikeConfig
	isCarInsurance := insuranceInfo.InsuranceType == constant.InsuranceType.Car
	if isCarInsurance {
		cfg = globalcare.CarConfig
	}
	payload := globalcare.CreateOrderPayload{
		ProductCode: cfg.ProductCode,
		ProviderID:  cfg.ProviderID,
		ProductID:   cfg.ProductID,
		PartnerID:   p.OrderCode,
		VehicleInfo: globalcare.VehicleInfo{
			TypeID:                       vTypeID,
			TypeName:                     insuranceInfo.VehicleTypeName,
			CarOccupantAccidentInsurance: iTypeID,
			License:                      insuranceInfo.License,
			Chassis:                      insuranceInfo.Chassis,
			Engine:                       insuranceInfo.Engine,
		},
		InsuredInfo: globalcare.InsuredInfo{
			BuyerName:        customer.Name,
			BuyerPhone:       customer.PhoneNumber,
			BuyerEmail:       customer.Email,
			BuyerAddress:     customer.Address.FullAddress,
			YearsOfInsurance: format.ToString(insuranceInfo.YearsOfInsurance),
			BeginDate:        insuranceInfo.BeginDate,
		},
	}
	// assign car insurance info
	if isCarInsurance {
		vi := payload.VehicleInfo
		buy := 1
		if insuranceInfo.NumberOfSeatsCarOccupantAccidentInsurance <= 0 {
			buy = 2
		}
		vi.CarOccupantAccidentInsurance = globalcare.CarOccupantAccidentInsuranceObj{
			NumberOfSeats: insuranceInfo.NumberOfSeatsCarOccupantAccidentInsurance,
			Buy:           buy,
		}

		if insuranceInfo.NumberOfSeats >= globalcare.NumOfSeatsMinValue {
			vi.NumberOfSeatsOver25 = insuranceInfo.NumberOfSeats
		}
		tonnageId, err := format.StringToInt(insuranceInfo.NumberOfSeatsOrTonnageId)
		if err != nil {
			return nil, errors.New(errorcode.ORMissingInsuranceInfo)
		}
		vi.NumberOfSeatsOrTonnage = tonnageId
		vi.NumberOfSeatsOrTonnageName = insuranceInfo.NumberOfSeatsOrTonnageName
		// https://airtable.com/shrDq2C069X6fEhTG/tblC4IFOYChKnalKN/viwYU1cA7BgEJwWfg/recXsORbi1qIW6EMP?blocks=hide
		// v2 = true if car type = 1 as default
		if vi.TypeID == 1 {
			vi.V2 = true
		}
		payload.VehicleInfo = vi
	}
	res, err := s.client.CreateOrder(payload)
	if err != nil {
		r.Status = constant.WarehouseORStatusFailed
		r.Reason = err.Error()
	} else {
		r.ORCode = res.Result.OrderCode
	}
	return r, err
}

// UpdateLogisticInfo ...
func (s outboundRequestGlobalCare) UpdateLogisticInfo(p model.UpdateOutboundRequestLogisticInfoPayload,
	or mgwarehouse.OutboundRequest) error {
	return nil
}

// Cancel ...
func (s outboundRequestGlobalCare) Cancel(or mgwarehouse.OutboundRequest, note string) error {
	return nil
}

// GetWebhookData ...
func (s outboundRequestGlobalCare) GetWebhookData(b []byte) (*responsemodel.OutboundRequestWebhookData, error) {
	var data globalcare.Webhook
	if err := pjson.Unmarshal(b, &data); err != nil {
		logger.Error("service.outboundRequestGlobalCare.GetWebhookData - parse webhook", logger.LogData{
			"data": string(b),
			"err":  err.Error(),
		})
		return nil, err
	}
	res := &responsemodel.OutboundRequestWebhookData{
		Status:         globalCareStatusMapping[data.Status],
		DeliveryStatus: globalCareDeliveryStatusMapping[data.Status],
		ORCode:         data.OrderCode,
		ORRequestID:    "",
		OrderCode:      data.PartnerOrderCode,
		UpdatedAt:      data.UpdatedAt,
		Reason:         data.Note,
		Link:           data.CertLink,
	}
	return res, nil
}

// GetIdentityCode ...
func (s outboundRequestGlobalCare) GetIdentityCode() string {
	return constant.WarehousePartnerCodeGlobalCare
}

func (s outboundRequestGlobalCare) getClient() (*globalcare.Client, error) {
	var info struct {
		PublicKey  string `json:"publicKey"`
		PrivateKey string `json:"privateKey"`
	}
	b := base64.Decode(s.auth)
	if err := pjson.Unmarshal(b, &info); err != nil {
		return nil, fmt.Errorf("outboundRequestGlobalCare.getClient: %v", err)
	}
	if info.PrivateKey == "" || info.PublicKey == "" {
		return nil, fmt.Errorf("outboundRequestGlobalCare.getClient - missing rsa key")
	}
	env := globalcare.EnvDev
	if config.IsEnvProduction() {
		env = globalcare.EnvProd
	} else if config.IsEnvStaging() {
		env = globalcare.EnvStaging
	}
	natsClient := natsio.GetServer()
	return globalcare.NewClient(env, info.PrivateKey, info.PublicKey, natsClient)
}

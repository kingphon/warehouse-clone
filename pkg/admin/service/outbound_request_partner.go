package service

import (
	"fmt"

	"git.selly.red/Selly-Modules/natsio/model"

	"git.selly.red/Selly-Server/warehouse/external/constant"
	mgwarehouse "git.selly.red/Selly-Server/warehouse/external/model/mg/warehouse"
	responsemodel "git.selly.red/Selly-Server/warehouse/pkg/admin/model/response"
)

// OutboundRequestPartner ...
type OutboundRequestPartner interface {
	Create(p model.OutboundRequestPayload, whPartner responsemodel.ResponseWarehousePartnerConfig) (*model.OutboundRequestResponse, error)
	UpdateLogisticInfo(p model.UpdateOutboundRequestLogisticInfoPayload, or mgwarehouse.OutboundRequest) error
	Cancel(or mgwarehouse.OutboundRequest, note string) error
	GetORStatus(or mgwarehouse.OutboundRequest) (*model.SyncORStatusResponse, error)
	GetWebhookData(b []byte) (*responsemodel.OutboundRequestWebhookData, error)
	GetIdentityCode() string
	HasUpdateORCode() bool
	UpdateDeliveryStatus(model.WarehouseORUpdateDeliveryStatus, mgwarehouse.OutboundRequest) error
}

// NewOutboundRequestService ...
func NewOutboundRequestService(partnerIdentityCode, auth string) (OutboundRequestPartner, error) {
	switch partnerIdentityCode {
	case constant.WarehousePartnerCodeTNC:
		return newORTNCService(auth)
	case constant.WarehousePartnerCodeGlobalCare:
		return newORGlobalCareService(auth)
	case constant.WarehousePartnerCodeOnPoint:
		return newOROnPointService(auth)
	}
	return nil, fmt.Errorf("service.NewOutboundRequestService: partner_code_not_supported %s", partnerIdentityCode)
}

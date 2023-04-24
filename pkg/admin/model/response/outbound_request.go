package responsemodel

import (
	"time"

	"git.selly.red/Selly-Server/warehouse/external/utils/ptime"
)

// OutboundRequestWebhookData ...
type OutboundRequestWebhookData struct {
	Status         string    `json:"status"`
	DeliveryStatus string    `json:"deliveryStatus"`
	ORCode         string    `json:"orCode"`
	ORRequestID    string    `json:"orRequestId"`
	OrderCode      string    `json:"orderCode"`
	UpdatedAt      time.Time `json:"updatedAt"`
	Reason         string    `json:"reason"`
	Link           string    `json:"link"`
}

// OutboundRequestStatus ...
type OutboundRequestStatus struct {
	ORCode         string    `json:"orCode"`
	Status         string    `json:"status"`
	DeliveryStatus string    `json:"deliveryStatus"`
	UpdatedAt      time.Time `json:"updatedAt"`
	Link           string    `json:"link"`
}

// OutboundRequestList ...
type OutboundRequestList struct {
	List []OutboundRequest `json:"list"`
}

// OutboundRequest ...
type OutboundRequest struct {
	ID           string                 `json:"_id"`
	Status       string                 `json:"status"`
	Partner      OutboundRequestPartner `json:"partner"`
	TrackingCode string                 `json:"trackingCode"`
	CreatedAt    *ptime.TimeResponse    `json:"createdAt"`
	UpdatedAt    *ptime.TimeResponse    `json:"updatedAt"`
}

// OutboundRequestPartner ...
type OutboundRequestPartner struct {
	IdentityCode string `json:"identityCode"`
	Code         string `json:"code"`
	RequestID    string `json:"requestId"`
}

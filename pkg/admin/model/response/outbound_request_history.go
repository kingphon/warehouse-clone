package responsemodel

import (
	"git.selly.red/Selly-Server/warehouse/external/utils/ptime"
)

// OutboundRequestHistory ...
type OutboundRequestHistory struct {
	ID        string              `json:"_id"`
	Status    string              `json:"status"`
	CreatedAt *ptime.TimeResponse `json:"createdAt"`
}

// OutboundRequestHistoryList ...
type OutboundRequestHistoryList struct {
	List []OutboundRequestHistory `json:"list"`
}

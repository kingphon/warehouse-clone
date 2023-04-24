package requestmodel

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"

	"git.selly.red/Selly-Server/warehouse/pkg/admin/errorcode"
)

// OutboundRequestHistoryQuery ...
type OutboundRequestHistoryQuery struct {
	Request string `json:"request" query:"request" example:"outbound_request_id"`
}

// Validate ...
func (q OutboundRequestHistoryQuery) Validate() error {
	return validation.ValidateStruct(
		&q,
		validation.Field(&q.Request, validation.Required.Error(errorcode.ORRequestRequired)),
	)
}

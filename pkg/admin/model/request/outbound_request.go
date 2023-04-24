package requestmodel

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"

	"git.selly.red/Selly-Server/warehouse/pkg/admin/errorcode"
)

// OutboundRequestQuery ...
type OutboundRequestQuery struct {
	Order string `json:"order" query:"order" example:"order_id"`
}

// Validate ...
func (q OutboundRequestQuery) Validate() error {
	return validation.ValidateStruct(
		&q,
		validation.Field(&q.Order, validation.Required.Error(errorcode.OROrderRequired)),
	)
}

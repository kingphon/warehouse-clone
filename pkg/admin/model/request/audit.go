package requestmodel

import (
	"time"

	"git.selly.red/Selly-Server/warehouse/external/response"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"

	"git.selly.red/Selly-Server/warehouse/pkg/admin/errorcode"
)

// Staff ...
type Staff struct {
	ID   string `json:"_id"`
	Name string ` json:"name"`
}

// AuditPayload ...
type AuditPayload struct {
	Action      string    `json:"action"`
	ActionAt    time.Time `json:"actionAt"`
	Checksum    string    `json:"checksum"`
	Data        string    `json:"data"`
	Message     string    `json:"message"`
	ServiceName string    `json:"serviceName"`
	Target      string    `json:"target"`
	Author      *Staff    `json:"author"`
}

type AllQuery struct {
	Target   string `query:"target"`
	TargetID string `query:"targetID"`

	// Pagination
	Page  int64 `query:"page"`
	Limit int64 `query:"limit"`
}

// Validate ...
func (m AllQuery) Validate() error {
	return validation.ValidateStruct(
		&m,
		validation.Field(&m.TargetID, validation.Required.Error(errorcode.AuditInvalidTargetID), is.MongoID.Error(errorcode.AuditInvalidTargetID)),
		validation.Field(&m.Page, validation.Min(0).Error(response.CommonInvalidPagination)),
		validation.Field(&m.Limit, validation.Min(0).Error(response.CommonInvalidPagination)),
	)
}

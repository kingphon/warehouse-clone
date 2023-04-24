package service

import (
	"context"

	"git.selly.red/Selly-Modules/audit"
	"git.selly.red/Selly-Server/warehouse/external/constant"
	"git.selly.red/Selly-Server/warehouse/external/utils/format"

	externalauth "git.selly.red/Selly-Server/warehouse/external/auth"
	"git.selly.red/Selly-Server/warehouse/external/utils/mgquery"
	responsemodel "git.selly.red/Selly-Server/warehouse/pkg/admin/model/response"
)

// AuditInterface ...
type AuditInterface interface {
	// All return audit
	All(ctx context.Context, q mgquery.AppQuery) responsemodel.ResponseAuditAll

	// Create ...
	Create(target, targetId, msg, action string, data interface{})
}

// auditImplement ...
type auditImplement struct {
	CurrentStaff externalauth.User
}

// Audit return audit service
func Audit(cs externalauth.User) AuditInterface {
	return auditImplement{
		CurrentStaff: cs,
	}
}

// Create ...
func (s auditImplement) Create(target, targetId, msg, action string, data interface{}) {
	p := audit.CreatePayload{
		Target:   target,
		TargetID: targetId,
		Action:   action,
		Data:     format.ToString(data),
		Message:  msg,
		Author: audit.CreatePayloadAuthor{
			ID:   s.CurrentStaff.ID,
			Name: s.CurrentStaff.Name,
			Type: constant.TypeStaffAudit,
		},
	}
	audit.GetInstance().Create(p)
}

// All return audit
func (auditImplement) All(ctx context.Context, q mgquery.AppQuery) responsemodel.ResponseAuditAll {
	query := audit.AllQuery{
		Target:   q.Audit.Target,
		TargetID: q.Audit.TargetID,
		Limit:    q.Limit,
		Page:     q.Page,
		Sort:     q.SortInterface,
	}

	result, total := audit.GetInstance().All(query)
	return responsemodel.ResponseAuditAll{
		List:  result,
		Total: total,
	}

}

package errorcode

import "git.selly.red/Selly-Server/warehouse/external/response"

const (
	AuditInvalidTargetID = "audit_invalid_targetID"
	AuditInvalidTarget   = "audit_invalid_target"
)

var audit = []response.Code{
	{
		Key:     AuditInvalidTargetID,
		Message: "targetID không hợp lệ",
		Code:    201,
	},
	{
		Key:     AuditInvalidTarget,
		Message: "target không hợp lệ",
		Code:    202,
	},
}

package responsemodel

import "git.selly.red/Selly-Modules/audit"

type (
	// ResponseAuditAll ...
	ResponseAuditAll struct {
		List  []audit.Audit `json:"list"`
		Total int64         `json:"total"`
	}
)

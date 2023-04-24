package responsemodel

import "git.selly.red/Selly-Modules/audit"

// StaffShort
type StaffShort struct {
	ID   string `json:"_id"`
	Name string `json:"name"`
}

type (
	// ResponseAuditAll ...
	ResponseAuditAll struct {
		List  []audit.Audit `json:"list"`
		Total int64         `json:"total"`
	}
)

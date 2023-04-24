package initialize

import (
	"git.selly.red/Selly-Modules/audit"
	"git.selly.red/Selly-Server/warehouse/external/constant"
	"git.selly.red/Selly-Server/warehouse/pkg/admin/config"
)

// InitAudit ...
func InitAudit() {
	var cfg = config.GetENV().MongoAudit

	// Init
	if err := audit.NewInstance(audit.Config{
		Targets: []string{
			constant.AuditTargetWarehouse,
		},
		MongoDB: cfg.GetConnectOptions(),
	}); err != nil {
		panic(err)
	}
}

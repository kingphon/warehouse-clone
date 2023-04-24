package initialize

import (
	"git.selly.red/Selly-Server/warehouse/pkg/admin/config"
	"git.selly.red/Selly-Server/warehouse/pkg/admin/errorcode"
	"git.selly.red/Selly-Server/warehouse/pkg/admin/schedule"
)

// Init ...
func Init() {
	config.Init()
	errorcode.Init()
	zookeeper()
	mongoDB()
	authentication()
	InitAudit()
	nats()
	redis()
	schedule.Init()
}

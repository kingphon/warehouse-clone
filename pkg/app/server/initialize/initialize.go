package initialize

import (
	"git.selly.red/Selly-Server/warehouse/pkg/app/config"
	"git.selly.red/Selly-Server/warehouse/pkg/app/errorcode"
)

// Init ...
func Init() {
	config.Init()
	errorcode.Init()
	zookeeper()
	mongoDB()
	InitAudit()
	nats()
}

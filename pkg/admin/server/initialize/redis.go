package initialize

import (
	"git.selly.red/Selly-Modules/redisdb"

	"git.selly.red/Selly-Server/warehouse/pkg/admin/config"
)

func redis() {
	cfg := config.GetENV().Redis
	if err := redisdb.Connect(cfg.URI, cfg.PWD); err != nil {
		panic(err)
	}
}

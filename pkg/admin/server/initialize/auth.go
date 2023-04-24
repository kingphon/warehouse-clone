package initialize

import (
	"git.selly.red/Selly-Server/warehouse/pkg/admin/auth"
	"git.selly.red/Selly-Server/warehouse/pkg/admin/config"
)

// authentication
func authentication() {
	env := config.GetENV()
	auth.InitAuthentication(env.Nats.APIKey, env.Nats)
}

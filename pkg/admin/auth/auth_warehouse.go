package auth

import (
	"time"

	"git.selly.red/Selly-Modules/authentication"
	"git.selly.red/Selly-Modules/natsio"

	externalauth "git.selly.red/Selly-Server/warehouse/external/auth"
	"git.selly.red/Selly-Server/warehouse/external/auth/permission"
	"git.selly.red/Selly-Server/warehouse/pkg/admin/config"
)

var client *authentication.Client

// InitAuthentication ...
func InitAuthentication(apiKey string, nats config.NatsConfig) {
	setSecretKey, err := externalauth.Init(apiKey, permission.WarehouseSource, natsio.Config{
		URL:            nats.URL,
		User:           nats.Username,
		Password:       nats.Password,
		RequestTimeout: 3 * time.Minute,
	})

	if err != nil {
		panic(err)
	}

	envVars := config.GetENV()
	envVars.SecretKey = setSecretKey

}

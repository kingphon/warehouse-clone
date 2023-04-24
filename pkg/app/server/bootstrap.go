package server

import (
	"git.selly.red/Selly-Modules/logger"
	"github.com/labstack/echo/v4"

	"git.selly.red/Selly-Server/warehouse/pkg/app/route"

	"git.selly.red/Selly-Server/warehouse/pkg/app/server/initialize"
)

// Bootstrap ...
func Bootstrap(e *echo.Echo) {
	logger.Init("selly", "warehouse-app")

	// Init modules
	initialize.Init()

	// Routes
	route.Init(e)
}

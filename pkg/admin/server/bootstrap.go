package server

import (
	"git.selly.red/Selly-Modules/logger"
	"git.selly.red/Selly-Server/warehouse/pkg/admin/route"
	"git.selly.red/Selly-Server/warehouse/pkg/admin/server/initialize"
	"github.com/labstack/echo/v4"
)

// Bootstrap ...
func Bootstrap(e *echo.Echo) {
	logger.Init("selly", "warehouse-admin")

	// Init modules
	initialize.Init()

	// Routes
	route.Init(e)
}

package route

import (
	"strings"

	"git.selly.red/Selly-Server/warehouse/external/utils/echocontext"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"git.selly.red/Selly-Server/warehouse/external/utils/routemiddleware"
	"git.selly.red/Selly-Server/warehouse/pkg/app/config"
)

// Init ...
func Init(e *echo.Echo) {
	var (
		envVars = config.GetENV()
	)

	// Middlewares ...
	e.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey: []byte(envVars.SecretKey),
		Skipper: func(c echo.Context) bool {
			token := echocontext.GetToken(c)
			return token == "" || strings.Contains(c.Path(), "supplier-holidays")
		},
	}))

	e.Use(routemiddleware.CORSConfig())
	e.Use(routemiddleware.Locale)

	r := e.Group("/app/warehouse")

	// Components
	common(r)
	warehouse(r)
	supplierHoliday(r)
}

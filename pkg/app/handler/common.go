package handler

import (
	"github.com/labstack/echo/v4"

	"git.selly.red/Selly-Server/warehouse/external/response"
)

// Common ...
type Common struct{}

// Ping godoc
// @tags Common
// @summary Ping
// @id ping
// @security ApiKeyAuth
// @accept json
// @produce json
// @router /ping [get]
func (Common) Ping(c echo.Context) error {
	return response.R200(c, echo.Map{}, "")
}

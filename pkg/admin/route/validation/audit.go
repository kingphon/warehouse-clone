package routevalidation

import (
	"git.selly.red/Selly-Server/warehouse/external/response"
	"git.selly.red/Selly-Server/warehouse/external/utils/echocontext"
	requestmodel "git.selly.red/Selly-Server/warehouse/pkg/admin/model/request"
	"github.com/labstack/echo/v4"
)

// Audit
type Audit struct{}

// GetAll ...
func (Audit) GetAll(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		var (
			query requestmodel.AllQuery
		)

		if err := c.Bind(&query); err != nil {
			return response.R400(c, nil, "")
		}

		if err := query.Validate(); err != nil {
			return response.RouteValidation(c, err)
		}

		echocontext.SetQuery(c, query)
		return next(c)
	}
}

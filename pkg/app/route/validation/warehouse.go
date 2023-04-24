package routevalidation

import (
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"git.selly.red/Selly-Server/warehouse/pkg/app/errorcode"

	"git.selly.red/Selly-Server/warehouse/external/response"
	"git.selly.red/Selly-Server/warehouse/external/utils/echocontext"
)

// Warehouse ...
type Warehouse struct{}

// Detail ...
func (Warehouse) Detail(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		var (
			id = c.Param("id")
		)

		if !primitive.IsValidObjectID(id) {
			return response.R404(c, nil, errorcode.WarehouseInvalidID)
		}

		objID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return response.R400(c, nil, "")
		}

		echocontext.SetParam(c, "id", objID)
		return next(c)
	}
}

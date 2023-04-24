package handler

import (
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"git.selly.red/Selly-Server/warehouse/pkg/app/service"

	"git.selly.red/Selly-Server/warehouse/external/response"
	"git.selly.red/Selly-Server/warehouse/external/utils/echocontext"
)

// Warehouse ...
type Warehouse struct{}

// Detail godoc
// @tags Warehouse
// @summary Detail
// @id app-warehouse-detail
// @security ApiKeyAuth
// @accept json
// @produce json
// @param device-id header string true "Device id"
// @Param  id path string true "Warehouse id"
// @success 200 {object} responsemodel.ResponseWarehouseDetail
// @router /warehouses/{id} [get]
func (Warehouse) Detail(c echo.Context) error {
	var (
		ctx = echocontext.GetContext(c)
		s   = service.Warehouse()
		id  = echocontext.GetParam(c, "id").(primitive.ObjectID)
	)

	result, err := s.Detail(ctx, id)
	if err != nil {
		return response.R404(c, echo.Map{}, err.Error())
	}
	return response.R200(c, echo.Map{
		"data": result,
	}, "")
}

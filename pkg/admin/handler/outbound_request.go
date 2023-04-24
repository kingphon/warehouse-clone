package handler

import (
	"github.com/labstack/echo/v4"

	"git.selly.red/Selly-Server/warehouse/external/response"
	"git.selly.red/Selly-Server/warehouse/external/utils/echocontext"
	requestmodel "git.selly.red/Selly-Server/warehouse/pkg/admin/model/request"
	responsemodel "git.selly.red/Selly-Server/warehouse/pkg/admin/model/response"
	"git.selly.red/Selly-Server/warehouse/pkg/admin/service"
)

// OutboundRequest ...
type OutboundRequest struct{}

// GetList godoc
// @tags OutboundRequest
// @summary Get List
// @id admin-outbound-request-list
// @security ApiKeyAuth
// @accept json
// @produce json
// @param device-id header string true "Device id"
// @Param payload query requestmodel.OutboundRequestQuery false "Query"
// @success 200 {object} responsemodel.OutboundRequestList
// @router /outbound-requests [get]
func (h OutboundRequest) GetList(c echo.Context) error {
	var (
		params = echocontext.GetQuery(c).(requestmodel.OutboundRequestQuery)
		s      = service.OutboundRequest()
		ctx    = echocontext.GetContext(c)
	)

	list, err := s.GetList(ctx, params)
	if err != nil {
		return response.R400(c, nil, err.Error())
	}
	return response.R200(c, responsemodel.OutboundRequestList{List: list}, "")
}

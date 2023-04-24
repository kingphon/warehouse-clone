package handler

import (
	"github.com/labstack/echo/v4"

	"git.selly.red/Selly-Server/warehouse/external/response"
	"git.selly.red/Selly-Server/warehouse/external/utils/echocontext"
	requestmodel "git.selly.red/Selly-Server/warehouse/pkg/admin/model/request"
	responsemodel "git.selly.red/Selly-Server/warehouse/pkg/admin/model/response"
	"git.selly.red/Selly-Server/warehouse/pkg/admin/service"
)

// OutboundRequestHistory ...
type OutboundRequestHistory struct{}

// GetList godoc
// @tags OutboundRequestHistory
// @summary Get List
// @id admin-outbound-request-history-list
// @security ApiKeyAuth
// @accept json
// @produce json
// @param device-id header string true "Device id"
// @Param payload query requestmodel.OutboundRequestHistoryQuery false "Query"
// @success 200 {object} responsemodel.OutboundRequestHistoryList
// @router /outbound-request-histories [get]
func (h OutboundRequestHistory) GetList(c echo.Context) error {
	var (
		params = echocontext.GetQuery(c).(requestmodel.OutboundRequestHistoryQuery)
		s      = service.OutboundRequestHistory()
		ctx    = echocontext.GetContext(c)
	)

	list, err := s.GetList(ctx, params)
	if err != nil {
		return response.R400(c, nil, err.Error())
	}
	return response.R200(c, responsemodel.OutboundRequestHistoryList{List: list}, "")
}

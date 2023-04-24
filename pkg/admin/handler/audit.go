package handler

import (
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"

	externalauth "git.selly.red/Selly-Server/warehouse/external/auth"
	"git.selly.red/Selly-Server/warehouse/external/response"
	"git.selly.red/Selly-Server/warehouse/external/utils/echocontext"
	"git.selly.red/Selly-Server/warehouse/external/utils/mgquery"
	requestmodel "git.selly.red/Selly-Server/warehouse/pkg/admin/model/request"
	"git.selly.red/Selly-Server/warehouse/pkg/admin/service"
)

// Audit
type Audit struct{}

// All godoc
// @tags Audit
// @summary All
// @id admin-audit-all
// @security ApiKeyAuth
// @accept json
// @produce json
// @param Device-ID header string true "DeviceID"
// @param payload query requestmodel.AllQuery true "Query"
// @success 200 {object} responsemodel.ResponseAuditAll
// @router /audits [get]
func (Audit) All(c echo.Context) error {
	var (
		ctx     = echocontext.GetContext(c)
		qParams = echocontext.GetQuery(c).(requestmodel.AllQuery)
		s       = service.Audit(externalauth.User{})
		q       = mgquery.AppQuery{
			Page:          qParams.Page,
			Limit:         qParams.Limit,
			SortInterface: bson.M{"createdAt": -1},
			Audit: mgquery.Audit{
				Target:   qParams.Target,
				TargetID: qParams.TargetID,
			},
		}
	)

	result := s.All(ctx, q)
	return response.R200(c, result, "")
}

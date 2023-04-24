package routeauth

import (
	"github.com/labstack/echo/v4"

	externalauth "git.selly.red/Selly-Server/warehouse/external/auth"
	"git.selly.red/Selly-Server/warehouse/external/response"
	"git.selly.red/Selly-Server/warehouse/external/utils/echocontext"
)

// RequiredLogin ...
func RequiredLogin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// check invalid token
		user := externalauth.GetCurrentUserByToken(c.Get("user"))
		if user == nil || user.ID == "" {
			return response.R403(c, echo.Map{}, response.CommonForbidden)
		}

		staff := externalauth.User{
			ID:   user.ID,
			Name: user.Name,
		}

		c.Set("staff", staff)
		return next(c)
	}
}

// CheckPermission ...
func CheckPermission(scopes []string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			user := externalauth.GetCurrentUserByToken(c.Get("user"))
			if user == nil || user.ID == "" {
				return response.R403(c, echo.Map{}, response.CommonForbidden)
			}

			err := externalauth.CheckPermission(scopes, externalauth.StaffCheckPermissionBody{
				StaffID:  user.ID,
				DeviceID: echocontext.GetDeviceID(c),
				Token:    echocontext.GetToken(c),
			})
			if err != nil {

				return response.R403(c, echo.Map{}, err.Error())
			}

			return next(c)
		}
	}
}

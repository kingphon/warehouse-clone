package routeauth

import (
	"fmt"
	externalauth "git.selly.red/Selly-Server/warehouse/external/auth"

	"git.selly.red/Selly-Server/warehouse/external/utils/echocontext"
	responsemodel "git.selly.red/Selly-Server/warehouse/pkg/app/model/response"
	"github.com/labstack/echo/v4"

	"git.selly.red/Selly-Server/warehouse/external/response"

	"git.selly.red/Selly-Server/warehouse/pkg/app/service"
)

// RequiredLogin ...
func RequiredLogin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// check invalid token
		if externalauth.GetCurrentUserByToken(c.Get("user")).ID == "" {
			return response.R403(c, echo.Map{}, response.CommonForbidden)
		}
		return next(c)
	}
}

// CheckTokenSupplier ...
func CheckTokenSupplier(permissions []string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			token := echocontext.GetToken(c)
			res, err := service.CheckTokenSupplierUser(token, permissions)
			if err != nil {
				return err
			}

			if res.IsValid == false {
				return response.R403(c, nil, "")

			}

			fmt.Println("user: ", res)

			c.Set("user", &responsemodel.ResponseUserInfo{
				ID:         res.User.ID,
				Name:       res.User.Name,
				SupplierID: res.User.SupplierID,
			})

			return next(c)
		}
	}
}

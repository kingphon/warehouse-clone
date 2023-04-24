package auth

import (
	"strings"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"git.selly.red/Selly-Server/warehouse/external/utils/routemiddleware"
	responsemodel "git.selly.red/Selly-Server/warehouse/pkg/app/model/response"
)

// GetCurrentUser ...
func GetCurrentUser(c echo.Context) (user responsemodel.User) {
	token := c.Get("user")
	if token == nil {
		return
	}

	data, ok := token.(*jwt.Token)
	if !ok {
		return
	}

	m, ok := data.Claims.(jwt.MapClaims)
	if ok && data.Valid && m["_id"] != "" {
		user.ID = GetAppIDFromHex(m["_id"].(string))
	}

	if ok && data.Valid && m["name"] != "" {
		user.Name = m["name"].(string)
	}
	return
}

// GetAppIDFromHex ...
func GetAppIDFromHex(s string) primitive.ObjectID {
	id, _ := primitive.ObjectIDFromHex(s)
	return id
}

// GetToken ...
func GetToken(c echo.Context) string {
	token := c.Request().Header.Get(echo.HeaderAuthorization)

	split := strings.Split(token, " ")
	if len(split) == 1 {
		return ""
	}

	return split[1]
}

// GetDeviceID ...
func GetDeviceID(c echo.Context) string {
	return c.Request().Header.Get(routemiddleware.HeaderDeviceID)
}

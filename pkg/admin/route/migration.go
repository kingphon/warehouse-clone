package route

import (
	"github.com/labstack/echo/v4"

	"git.selly.red/Selly-Server/warehouse/pkg/admin/migration"
)

func migrate(r *echo.Group) {
	g := r.Group("/migration")

	migration.Migrate(g)
}

package initialize

import (
	"git.selly.red/Selly-Server/warehouse/pkg/admin/database"
)

// mongoDB ...
func mongoDB() {
	database.ConnectMongoDBWarehouse()
}

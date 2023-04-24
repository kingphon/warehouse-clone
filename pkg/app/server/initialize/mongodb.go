package initialize

import (
	"git.selly.red/Selly-Server/warehouse/pkg/app/database"
)

// mongoDB ...
func mongoDB() {
	database.ConnectMongoDBWarehouse()
}
